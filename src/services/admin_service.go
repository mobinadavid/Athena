package services

import (
	"athena/src/api/errs"
	"athena/src/api/http/dto"
	"athena/src/api/http/requests/AdminRequests"
	"athena/src/cache"
	"athena/src/config"
	"athena/src/database"
	"athena/src/database/scopes"
	"athena/src/hash"
	"athena/src/models"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/repositories"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IAdminService interface {
	GetList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	Create(ctx context.Context, request *AdminRequests.CreateAdminRequest, creatorAdminID uint) (*models.AdminModel, error)
	GetByUuid(ctx context.Context, uuid *uuid.UUID) (*models.AdminModel, error)
	GetProfile(ctx context.Context, userId uint) (*models.AdminModel, error)
	GetById(ctx context.Context, adminID uint) (*models.AdminModel, error)
	GetByMobile(ctx context.Context, mobile string) (*models.AdminModel, error)
	GetByNationalIdentityCode(ctx context.Context, nationalIdentityCode string) (*models.AdminModel, error)
	GetByUsername(ctx context.Context, username string) (*models.AdminModel, error)
	UpdateByUuid(ctx context.Context, adminUuid *uuid.UUID, request *AdminRequests.UpdateAdminRequest) (*models.AdminModel, error)
	DeleteByUuid(ctx context.Context, adminUuid *uuid.UUID) error
	Update(ctx context.Context, admin *models.AdminModel) (*models.AdminModel, error)
	Save(ctx context.Context, admin *models.AdminModel) (*models.AdminModel, error)
	ChangePasswordSendOTP(ctx context.Context) (string, error)
	ChangePassword(ctx context.Context) error
	ChangePasswordResendOTP(ctx context.Context) error
}

// AdminService struct implements the AdminService interface.
type AdminService struct {
	AdminRepository repositories.IAdminRepository
	RoleRepository  repositories.IRoleRepository
	//CDNService      ICDNService
	OTPService IOTPService
}

func (service *AdminService) GetList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	builder.Relations = append(builder.Relations, "Roles")
	res, err := service.AdminRepository.GetList(builder)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get list of admins", service, err)
		return nil, errs.SomeThingWentWrong
	}
	var adminList []dto.AdminListModel
	admins, ok := res.Items.(*[]*models.AdminModel)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to cast res.Items to []*models.AdminModel", service, nil)
		return nil, errs.SomeThingWentWrong
	}
	for _, admin := range *admins {

		//var name string
		//var profilePresigned *dto.PreSignedURLResponse
		//var cdnClient = cdn.NewCdnClient()
		//if admin.AdminImageUuid != "" {
		//	profilePresigned, err = cdnClient.GetPreSigned(config.GetInstance().Get("CDN_BUCKET_NAME"), admin.AdminImageUuid)
		//	name = admin.AdminImageUuid
		//} else {
		//	profilePresigned, err = cdnClient.GetPreSigned(config.GetInstance().Get("CDN_BUCKET_NAME"), config.GetInstance().Get("IMAGE_UUID_DEFAULT"))
		//	name = config.GetInstance().Get("IMAGE_UUID_DEFAULT")
		//}
		//if err != nil {
		//	logger.LogErrorWithFieldsV2(ctx, "failed to get presigned url", service, err)
		//	return nil, errs.SomeThingWentWrong
		//}
		//attachment := models.Attachment{
		//	Name: name,
		//	URL:  base64.StdEncoding.EncodeToString([]byte(profilePresigned.Data.URL)),
		//}
		//admin.AdminImage = attachment

		adminModel := dto.AdminListModel{
			Uuid:     admin.Uuid,
			Name:     admin.FirstName,
			LastName: admin.LastName,
			Mobile:   admin.Mobile,
			//ProfileImage: models.Attachment{
			//	Name: admin.AdminImageUuid,
			//	URL:  attachment.URL,
			//},
			Username:    admin.Username,
			IsActive:    admin.IsActive,
			CreatedDate: admin.CreatedAt,
			UpdatedDate: admin.UpdatedAt,
		}
		var adminRoles []dto.AdminRole
		for _, role := range admin.Roles {
			adminRoles = append(adminRoles, dto.AdminRole{
				Title: role.Title,
			})
		}
		adminModel.Roles = adminRoles
		adminList = append(adminList, adminModel)
	}
	res.Items = adminList

	return res, nil
}

func (service *AdminService) GetById(ctx context.Context, adminID uint) (*models.AdminModel, error) {
	res, err := service.AdminRepository.GetById(adminID)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AdminService) GetByUuid(ctx context.Context, uuid *uuid.UUID) (*models.AdminModel, error) {
	res, err := service.AdminRepository.GetByUuid(uuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}
	return res, nil
}

func (service *AdminService) GetProfile(ctx context.Context, userId uint) (*models.AdminModel, error) {
	res, err := service.AdminRepository.GetById(userId)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
		return nil, errs.SomeThingWentWrong
	}

	//if res.AdminImageUuid != "" {
	//	res.AdminImage = GetImageAttachment(res.AdminImageUuid)
	//} else {
	//	res.AdminImage = GetImageAttachment(config.GetInstance().Get("IMAGE_UUID_DEFAULT"))
	//}
	return res, nil
}

func (service *AdminService) GetByMobile(ctx context.Context, mobile string) (*models.AdminModel, error) {
	res, err := service.AdminRepository.GetByMobile(mobile)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by mobile", service, err)
		return nil, err
	}

	return res, nil
}

func (service *AdminService) GetByNationalIdentityCode(ctx context.Context, nationalIdentityCode string) (*models.AdminModel, error) {
	res, err := service.AdminRepository.GetByNationalIdentityCode(nationalIdentityCode)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by notional identity code", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AdminService) GetByUsername(ctx context.Context, username string) (*models.AdminModel, error) {
	res, err := service.AdminRepository.GetByUsername(username)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by username", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AdminService) Create(ctx context.Context, request *AdminRequests.CreateAdminRequest, creatorAdminID uint) (*models.AdminModel, error) {
	if _, err := service.AdminRepository.GetById(creatorAdminID); err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
		return nil, errs.SomeThingWentWrong
	}
	username := strings.ToLower(request.Username)
	re := regexp.MustCompile(`[^a-z0-9_]`)
	username = re.ReplaceAllString(username, "")
	request.Username = username
	admin := &models.AdminModel{
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Mobile:       request.Mobile,
		Username:     request.Username,
		Password:     []byte(request.Password),
		TwoFaEnabled: false,
		IsActive:     request.IsActive,
	}
	//if request.ProfileImage != "" {
	//	admin.AdminImageUuid = request.ProfileImage
	//} else {
	//	url, err := cdn.NewCdnClient().GetPreSigned(config.GetInstance().Get("CDN_BUCKET_NAME"), config.GetInstance().Get("IMAGE_UUID_DEFAULT"))
	//	if err != nil {
	//		logger.LogErrorWithFieldsV2(ctx, "failed to get presigned url", service, err)
	//		return nil, errs.SomeThingWentWrong
	//	}
	//	admin.AdminImageUuid = url.Data.FileName
	//}
	// Loop through the role UUIDs provided in the request.
	exists := false
	for _, requestRoleUuid := range request.Roles {
		// Parse the role UUID.
		RoleUuid, err := uuid.Parse(requestRoleUuid)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to parse uuid", service, err)
			return nil, errs.SomeThingWentWrong
		}

		// Retrieve the role from the service using the UUID.
		role, err := service.RoleRepository.GetByUuid(&RoleUuid)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get role by uuid", service, err)
			return nil, err
		}
		if role.Name == models.DefaultRole {
			exists = true
		}
		// Append the role to slice.
		admin.Roles = append(admin.Roles, role)
	}
	if !exists {
		// Retrieve the role from the service using the name.
		role, err := service.RoleRepository.GetByName(models.DefaultRole)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get role by Name", service, err)
			return nil, err
		}
		admin.Roles = append(admin.Roles, role)
	}

	adminOrm, err := service.AdminRepository.Create(admin)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to create admin", service, err)
		if utils.CheckError(err, errs.DuplicateUsername) {
			return nil, err
		}
		return nil, errs.SomeThingWentWrong
	}

	return adminOrm, nil
}

func (service *AdminService) UpdateByUuid(ctx context.Context, adminUuid *uuid.UUID, request *AdminRequests.UpdateAdminRequest) (*models.AdminModel, error) {
	admin, err := service.GetByUuid(ctx, adminUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}

	if request.FirstName != "" {
		admin.FirstName = request.FirstName
	}
	if request.LastName != "" {
		admin.LastName = request.LastName
	}
	if request.Mobile != "" {
		admin.Mobile = request.Mobile
	}
	if request.Username != "" {
		admin.Username = request.Username
	}
	if request.ProfileImage != "" {
		admin.AdminImageUuid = request.ProfileImage
	}
	if request.IsActive != nil {
		admin.IsActive = request.IsActive
	}
	if request.Password != "" {
		admin.Password = []byte(request.Password)
	}
	if len(request.Roles) > 0 {
		var roles []*models.RoleModel
		for _, requestRoleUuid := range request.Roles {
			RoleUuid, err := uuid.Parse(requestRoleUuid)
			if err != nil {
				logger.LogErrorWithFieldsV2(ctx, "failed to parse uuid", service, err)
				return nil, errs.SomeThingWentWrong
			}
			role, err := service.RoleRepository.GetByUuid(&RoleUuid)
			if err != nil {
				logger.LogErrorWithFieldsV2(ctx, "failed to get role by uuid", service, err)
				return nil, err
			}
			roles = append(roles, role)
		}

		db := database.GetInstance().GetClient()
		err = db.Model(admin).Association("Roles").Replace(roles)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to replace roles", service, err)
			return nil, errs.SomeThingWentWrong
		}
		admin.Roles = roles
	}
	var IPs string
	if len(request.AllowedIPs) != 0 {
		IPs = request.AllowedIPs[0]
		for i := 1; i < len(request.AllowedIPs); i++ {
			IPs = IPs + ";" + request.AllowedIPs[i]
		}
	}

	res, err := service.AdminRepository.Update(admin)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update admin", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AdminService) DeleteByUuid(ctx context.Context, adminUuid *uuid.UUID) error {
	admin, err := service.GetByUuid(ctx, adminUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin uuid", service, err)
		return errs.SomeThingWentWrong
	}

	err = service.AdminRepository.Delete(admin)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to delete admin", service, err)
		return errs.SomeThingWentWrong
	}
	return nil
}

func (service *AdminService) Update(ctx context.Context, admin *models.AdminModel) (*models.AdminModel, error) {
	res, err := service.AdminRepository.Update(admin)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update admin", service, err)
		return nil, errs.SomeThingWentWrong
	}
	return res, nil
}

func (service *AdminService) Save(ctx context.Context, admin *models.AdminModel) (*models.AdminModel, error) {
	res, err := service.AdminRepository.Save(admin)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save admin", service, err)
		return nil, errs.SomeThingWentWrong
	}
	return res, nil
}

func (service *AdminService) ChangePasswordSendOTP(ctx context.Context) (string, error) {
	req, ok := ctx.Value("req").(AdminRequests.ChangePasswordSendOtpRequest)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, nil)
		return "", errs.ChangePasswordFailed
	}
	// fetch data
	adminID := ctx.Value("admin_id").(uint)
	admin, err := service.AdminRepository.GetById(adminID)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
		return "", errs.SomeThingWentWrong
	}
	if req.CurrentPassword == req.NewPassword {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify password", service, nil)
		return "", errs.PasswordShouldBeNew
	}

	if req.NewPassword != req.NewPasswordConfirmation {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify password", service, nil)
		return "", errs.PasswordNotMatch
	}

	// Validate Given Current Password
	passwordHashCheck, err := hash.VerifyStoredHash(admin.Password, req.CurrentPassword)
	if err != nil || !passwordHashCheck {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify password", service, err)
		return "", errs.ErrAuthenticationFailed
	}

	// marshal the req to save in redis
	req.Mobile = admin.Mobile
	reqData, err := json.Marshal(req)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to marshal request data", service, err)
		return "", errs.ErrAuthenticationFailed
	}

	// get expire time
	expiration, err := strconv.Atoi(config.GetInstance().Get("CHANGE_PASSWORD_LIFETIME"))
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to convert CHANGE_PASSWORD_LIFETIME", service, err)
		expiration = 300
	}

	// Save the request data in Redis
	key := uuid.NewString()
	err = cache.GetInstance().GetClient().Set(ctx, key, reqData, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, err)
		return "", errs.ErrAuthenticationFailed
	}
	err = service.OTPService.RequestOTP(ctx, admin.Mobile)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to request otp", service, err,
			zap.String("mobile", admin.Mobile),
		)
		return "", err
	}

	return key, nil
}

func (service *AdminService) ChangePassword(ctx context.Context) error {
	req, ok := ctx.Value("req").(AdminRequests.ChangePasswordVerifyOtpRequest)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, nil)
		return errs.ChangePasswordFailed
	}
	adminID := ctx.Value("admin_id").(uint)
	admin, err := service.AdminRepository.GetById(adminID)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
		return errs.SomeThingWentWrong
	}

	//verify otp
	var otpIsValid bool
	otpIsValid, err = service.OTPService.VerifyOTP(ctx, admin.Mobile, req.OTP)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify change password otp", service, err,
			zap.String("mobile", admin.Mobile),
			zap.String("otp", req.OTP),
		)
		return errs.ChangePasswordFailed
	}

	if !otpIsValid {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify change password otp", service, nil,
			zap.String("mobile", admin.Mobile),
			zap.String("otp", req.OTP),
		)
		return errs.OTPIsNotValid
	}
	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.Key).Result()
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, err)
		return errs.SomeThingWentWrong
	}

	var resp AdminRequests.ChangePasswordSendOtpRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to unmarshal response data", service, err)
		return errs.SomeThingWentWrong
	}
	// Update Password.
	admin.Password = []byte(resp.NewPassword)
	_, err = service.AdminRepository.Update(admin)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update admin", service, err)
		return errs.SomeThingWentWrong
	}

	return nil
}

func (service *AdminService) ChangePasswordResendOTP(ctx context.Context) error {
	req, ok := ctx.Value("req").(*AdminRequests.ChangePasswordResendOtpRequest)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, nil)
		return errs.SomeThingWentWrong
	}

	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.Key).Result()
	if err != nil {
		if err == redis.Nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, err)
			return errs.ErrChangePasswordTimeOut
		}

		logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, err)
		return errs.ChangePasswordFailed

	}

	var resp AdminRequests.ChangePasswordSendOtpRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to unmarshal response data", service, err)
		return errs.SomeThingWentWrong
	}

	err = service.OTPService.RequestOTP(ctx, resp.Mobile)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to request otp", service, err,
			zap.String("mobile", resp.Mobile),
		)
		return err
	}

	return nil
}
