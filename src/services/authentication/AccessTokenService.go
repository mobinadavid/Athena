package authentication

import (
	"athena/src/api/errs"
	"athena/src/database/scopes"
	"athena/src/hash"
	"athena/src/models"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/repositories"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TokenType string
type OwnerType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type IAccessTokenService interface {
	GetList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetActiveTokens(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetByUuid(ctx context.Context, accessTokenUuid *uuid.UUID) (*models.AccessTokenModel, error)
	Create(ctx context.Context, owner interface{}, dto *JwtDTO, ip, userAgent string) (*models.AccessTokenModel, error)
	UpdateLastUsedAt(ctx context.Context, accessToken *models.AccessTokenModel) (*models.AccessTokenModel, error)
	RefreshAccessTokens(ctx context.Context, refreshToken, ownerType string) (*JwtDTO, error)
	Validate(ctx context.Context, tokenString string, tokenType TokenType, ownerType string) (*models.AccessTokenModel, error)
	RevokeTokens(ctx context.Context, auth *utils.Auth) error
	RevokeTokenByUuid(ctx context.Context, accessTokenUuid *uuid.UUID, ownerID uint, ownerType string) error
}

type AccessTokenService struct {
	AccessTokenRepository repositories.IAccessTokenRepository
	JwtService            IJwtService
	UserRepository        repositories.IUserRepository
}

func (service *AccessTokenService) GetList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	res, err := service.AccessTokenRepository.GetList(builder)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get list access tokens", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AccessTokenService) GetActiveTokens(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	res, err := service.AccessTokenRepository.GetActiveTokens(builder)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get list access tokens", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AccessTokenService) Create(ctx context.Context, owner interface{}, dto *JwtDTO, ip, userAgent string) (*models.AccessTokenModel, error) {
	accessToken := &models.AccessTokenModel{
		Uuid:                  dto.Uuid,
		AccessToken:           []byte(dto.AccessTokenString),
		AccessTokenExpiresAt:  dto.AccessTokenExpiresAt,
		RefreshToken:          []byte(dto.RefreshTokenString),
		RefreshTokenExpiresAt: dto.RefreshTokenExpiresAt,
		IP:                    ip,
		UserAgent:             userAgent,
	}

	switch owner := owner.(type) {
	case *models.UserModel:
		accessToken.OwnerID = owner.ID
		accessToken.OwnerType = "user"
	case *models.AdminModel:
		accessToken.OwnerID = owner.ID
		accessToken.OwnerType = "admin"
	default:
		err := errors.New("unsupported owner type")
		logger.LogErrorWithFieldsV2(ctx, "unsupported owner type", service, err,
			zap.String("ip", ip),
			zap.String("user-agent", userAgent),
		)
		return nil, err
	}

	atOrm, err := service.AccessTokenRepository.Create(accessToken)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to create access token", service, err)
		return nil, err
	}

	return atOrm, nil
}

func (service *AccessTokenService) GetByUuid(ctx context.Context, accessTokenUuid *uuid.UUID) (*models.AccessTokenModel, error) {
	accessToken, err := service.AccessTokenRepository.GetByUuid(accessTokenUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get token by uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return accessToken, nil
}

func (service *AccessTokenService) UpdateLastUsedAt(ctx context.Context, accessToken *models.AccessTokenModel) (*models.AccessTokenModel, error) {
	res, err := service.AccessTokenRepository.UpdateLastUsedAt(accessToken, time.Now())
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update token", service, err)
		return nil, err
	}

	return res, nil
}

func (service *AccessTokenService) RefreshAccessTokens(ctx context.Context, refreshToken, ownerType string) (*JwtDTO, error) {
	//validate token
	token, err := service.Validate(ctx, refreshToken, RefreshToken, ownerType)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to validate refresh token", service, err)
		return nil, errs.ErrInvalidRefreshToken
	}

	// generate new jwt
	jwtDto, err := service.JwtService.Generate(ctx, ownerType)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to generate new token", service, err)
		return nil, errs.SomeThingWentWrong
	}

	// update user tokens
	_, err = service.AccessTokenRepository.RefreshAccessTokens(
		token,
		jwtDto.Uuid,
		[]byte(jwtDto.AccessTokenString),
		[]byte(jwtDto.RefreshTokenString),
		jwtDto.AccessTokenExpiresAt,
		jwtDto.RefreshTokenExpiresAt,
	)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update user tokens", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return jwtDto, nil
}

func (service *AccessTokenService) Validate(ctx context.Context, tokenString string, tokenType TokenType, ownerType string) (*models.AccessTokenModel, error) {
	// Validate the extracted JWT token and retrieve the user claims.
	userClaimed, err := service.JwtService.Validate(ctx, tokenString)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed validate jwt token", service, err,
			zap.String("token-type", string(tokenType)),
		)
		return nil, err
	}

	// Additionally, validate the token against the database and check for its expiry.
	claimedUuid, err := uuid.Parse(userClaimed.ID)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to parse uuid", service, err,
			zap.String("token-type", string(tokenType)),
		)
		return nil, err
	}

	// Retrieve the access token from the database using the parsed UUID.
	token, err := service.AccessTokenRepository.GetByUuid(&claimedUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get token by uuid", service, err,
			zap.String("token-type", string(tokenType)),
		)
		return nil, errs.RecordNotFound
	}

	if token.OwnerType != ownerType {
		logger.LogErrorWithFieldsV2(ctx, "owner type is not valid", service, nil,
			zap.String("token-type", string(tokenType)),
		)
		return nil, errs.ErrAuthenticationFailed
	}
	// Verify the hash of the stored token against the provided token to ensure they match.
	var storedHash []byte
	var tokenExpiresAt time.Time

	switch tokenType {
	case AccessToken:
		storedHash = token.AccessToken
		tokenExpiresAt = token.AccessTokenExpiresAt

	case RefreshToken:
		storedHash = token.RefreshToken
		tokenExpiresAt = token.RefreshTokenExpiresAt
	}

	hashCheck, err := hash.VerifyStoredHash(storedHash, tokenString)
	if err != nil || !hashCheck {
		logger.LogErrorWithFieldsV2(ctx, "the password is not valid", service, err,
			zap.String("token-type", string(tokenType)),
		)
		return nil, errs.ErrAuthenticationFailed
	}

	// Check if the token has expired by comparing its expiry timestamp against the current time.
	if tokenExpiresAt.Before(time.Now()) {
		logger.LogErrorWithFieldsV2(ctx, "token is expire", service, nil)
		return nil, errs.ErrTokenExpired
	}
	return token, nil
}

func (service *AccessTokenService) RevokeTokens(ctx context.Context, auth *utils.Auth) error {
	// get list of tokens
	accessTokens, err := service.AccessTokenRepository.GetAll(auth.OwnerId, auth.OwnerType)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user tokens", service, err)
		return errs.SomeThingWentWrong
	}

	// delete the tokens
	err = service.AccessTokenRepository.DeleteMany(accessTokens)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to delete user tokens", service, err)
		return errs.SomeThingWentWrong
	}
	return nil
}

func (service *AccessTokenService) RevokeTokenByUuid(ctx context.Context, accessTokenUuid *uuid.UUID, ownerID uint, ownerType string) error {
	accessToken, err := service.AccessTokenRepository.GetByUuid(accessTokenUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get token by uuid", service, err)
		return errs.InvalidRevokeTokenFlow
	}

	if accessToken.OwnerID != ownerID || accessToken.OwnerType != ownerType {
		logger.LogErrorWithFieldsV2(ctx, "revoke token failed - owner id or type is not equal to the ones in the token", service, nil)
		return errs.InvalidRevokeTokenFlow
	}

	err = service.AccessTokenRepository.Delete(accessToken)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to delete user token", service, err)
		return errs.InvalidRevokeTokenFlow
	}
	return nil
}
