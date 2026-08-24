package services

import (
	"athena/src/api/errs"
	"athena/src/api/http/requests/Users/UserRequests"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/repositories"
	"context"

	"github.com/google/uuid"
)

type IUserService interface {
	GetList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	Create(ctx context.Context, request *UserRequests.CreateRequest) (*models.UserModel, error)
	Save(ctx context.Context, user *models.UserModel) (*models.UserModel, error)
	Update(ctx context.Context, user *models.UserModel) (*models.UserModel, error)
	GetByUuid(ctx context.Context, uuid *uuid.UUID) (*models.UserModel, error)
	GetById(ctx context.Context, userID uint, includeRelations ...string) (*models.UserModel, error)
	GetProfile(ctx context.Context, userID uint, includeRelations ...string) (*models.UserModel, error)
	GetByMobile(ctx context.Context, mobile string) (*models.UserModel, error)
	GetByNationalIdentityCode(ctx context.Context, nationalIdentityCode string, includeRelations ...string) (*models.UserModel, error)
	GetByMobileAndNationalIdentityCode(ctx context.Context, mobile, nationalIdentityCode string) (*models.UserModel, error)
	DeleteByUuid(ctx context.Context, userUuid *uuid.UUID) error
	FirstOrCreate(ctx context.Context, mobile string, nationalIdentityCode string) (*models.UserModel, error)
	UpdateOrCreate(ctx context.Context, search *models.UserModel, assign *models.UserModel) (*models.UserModel, error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (service *UserService) GetList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	res, err := service.UserRepository.GetList(builder)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get list of users", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) FirstOrCreate(ctx context.Context, mobile string, nationalIdentityCode string) (*models.UserModel, error) {
	res, err := service.UserRepository.FirstOrCreate(mobile, nationalIdentityCode)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get or create user", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) GetByMobileAndNationalIdentityCode(ctx context.Context, mobile, nationalIdentityCode string) (*models.UserModel, error) {
	res, err := service.UserRepository.GetByMobileAndNationalIdentityCode(mobile, nationalIdentityCode)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to user by mobile and national identity code", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) GetById(ctx context.Context, userID uint, includeRelations ...string) (*models.UserModel, error) {
	res, err := service.UserRepository.GetById(userID, includeRelations...)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by id", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) GetProfile(ctx context.Context, userID uint, includeRelations ...string) (*models.UserModel, error) {
	res, err := service.UserRepository.GetById(userID, includeRelations...)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by id", service, err)
		return nil, errs.SomeThingWentWrong
	}
	fullName := ""
	if res.FirstName != "" || res.LastName != "" {
		fullName = res.FirstName + " " + res.LastName
	}
	res.FullName = fullName
	return res, nil
}

func (service *UserService) GetByUuid(ctx context.Context, uuid *uuid.UUID) (*models.UserModel, error) {
	res, err := service.UserRepository.GetByUuid(uuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}
	return res, nil
}

func (service *UserService) GetByMobile(ctx context.Context, mobile string) (*models.UserModel, error) {
	res, err := service.UserRepository.GetByMobile(mobile)
	if err != nil && utils.CheckError(err, errs.RecordNotFound) == false {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by mobile", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) GetByNationalIdentityCode(ctx context.Context, nationalIdentityCode string, includeRelations ...string) (*models.UserModel, error) {
	res, err := service.UserRepository.GetByNationalIdentityCode(nationalIdentityCode, includeRelations...)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by notional identity code", service, err)
		return nil, err
	}

	return res, nil
}

func (service *UserService) Save(ctx context.Context, user *models.UserModel) (*models.UserModel, error) {
	res, err := service.UserRepository.Save(user)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update user", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) Update(ctx context.Context, user *models.UserModel) (*models.UserModel, error) {
	res, err := service.UserRepository.Update(user)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update user", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) Create(ctx context.Context, request *UserRequests.CreateRequest) (*models.UserModel, error) {
	user := &models.UserModel{
		FirstName:            request.FirstName,
		LastName:             request.LastName,
		NationalIdentityCode: request.NationalIdentityCode,
		Mobile:               request.Mobile,
		Email:                request.Email,
		Password:             []byte(request.Password),
	}

	userOrm, err := service.UserRepository.Create(user)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to create user", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return userOrm, nil
}

func (service *UserService) DeleteByUuid(ctx context.Context, userUuid *uuid.UUID) error {
	user, err := service.GetByUuid(ctx, userUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by uuid", service, err)
		return errs.SomeThingWentWrong
	}

	err = service.UserRepository.Delete(user)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to delete user", service, err)
		return errs.SomeThingWentWrong
	}
	return nil
}

func (service *UserService) UpdateOrCreate(ctx context.Context, search *models.UserModel, assign *models.UserModel) (*models.UserModel, error) {
	res, err := service.UserRepository.UpdateOrCreate(search, assign)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to create or update user", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}
