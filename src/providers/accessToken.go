package providers

import (
	adminAuthcontrollers "athena/src/api/http/controllers/admins"
	userAuthcontrollers "athena/src/api/http/controllers/users/authentication"
	"athena/src/database"
	"athena/src/repositories"
	"athena/src/services/authentication"
)

func ProvideAccessTokenService(accessTokenRepository *repositories.AccessTokenRepository,
	jwtService *authentication.JwtService,
	UserRepository *repositories.UserRepository) *authentication.AccessTokenService {
	return &authentication.AccessTokenService{
		AccessTokenRepository: accessTokenRepository,
		JwtService:            jwtService,
		UserRepository:        UserRepository,
	}
}

func ProvideAccessTokenRepository(db *database.Database) *repositories.AccessTokenRepository {
	return &repositories.AccessTokenRepository{
		DatabaseHandler: db,
	}
}

func ProvideUserAccessTokenController(accessTokenService *authentication.AccessTokenService) *userAuthcontrollers.AccessTokenController {
	return &userAuthcontrollers.AccessTokenController{
		AccessTokenService: accessTokenService,
	}
}

func ProvideAdminAccessTokenController(accessTokenService *authentication.AccessTokenService) *adminAuthcontrollers.AccessTokenController {
	return &adminAuthcontrollers.AccessTokenController{
		AccessTokenService: accessTokenService,
	}
}
