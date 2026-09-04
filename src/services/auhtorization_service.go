package services

import (
	"athena/src/api/errs"
	"athena/src/api/http/dto"
	"athena/src/api/http/requests/AuthorizationRequests"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/repositories"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IAuthorizationService interface {
	GetPermissionsList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetPermissionGroupsList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	CreatePermissionGroup(ctx context.Context, request *AuthorizationRequests.CreatePermissionGroupRequest) (*models.PermissionGroupModel, error)
	GetPermissionGroupByUuid(ctx context.Context, uuid *uuid.UUID, includeRelations ...string) (*models.PermissionGroupModel, error)
	UpdatePermissionGroupByUuid(ctx context.Context, groupUuid *uuid.UUID, request *AuthorizationRequests.UpdatePermissionGroupRequest) (*models.PermissionGroupModel, error)
	DeletePermissionGroupByUuid(ctx context.Context, groupUuid *uuid.UUID) error
	GetRolesList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetRoleByUuid(ctx context.Context, roleUuid *uuid.UUID) (*models.RoleModel, error)
	CreateRole(ctx context.Context, request *AuthorizationRequests.CreateRoleRequest) (*models.RoleModel, error)
	UpdateRoleByUuid(ctx context.Context, roleUuid *uuid.UUID, request *AuthorizationRequests.UpdateRoleRequest) (*models.RoleModel, error)
	DeleteRoleByUuid(ctx context.Context, roleUuid *uuid.UUID) error
	GetAdminsForRole(ctx context.Context, roleUuid *uuid.UUID) ([]*models.AdminModel, error)
	GetUsersForRole(ctx context.Context, roleUuid *uuid.UUID) ([]*models.UserModel, error)
	IsAuthorized(ctx context.Context, userId uint, userType, permission string) (bool, error)
}

type AuthorizationService struct {
	RoleRepository            repositories.IRoleRepository
	PermissionRepository      repositories.IPermissionRepository
	PermissionGroupRepository repositories.IPermissionGroupRepository
	UserRepository            repositories.IUserRepository
	AdminRepository           repositories.IAdminRepository
}

func (service *AuthorizationService) GetPermissionsList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	res, err := service.PermissionRepository.GetList(builder)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get list of permission", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AuthorizationService) GetPermissionGroupsList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	res, err := service.PermissionGroupRepository.GetList(builder)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get list of permission group", service, err)
		return nil, errs.SomeThingWentWrong
	}
	// to return the dto
	groups, ok := res.Items.(*[]*models.PermissionGroupModel)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to cast res.Items to []*models.PermissionGroupModel", service, nil)
		return nil, errs.SomeThingWentWrong
	}
	var groupsDto []*dto.PermissionGroupListModel
	for _, group := range *groups {
		permissionGroup := &dto.PermissionGroupListModel{
			ID:               group.ID,
			Uuid:             group.Uuid,
			Name:             group.Name,
			Title:            group.Title,
			Permissions:      group.Permissions,
			PermissionsCount: len(group.Permissions),
			IsActive:         *group.IsActive,
			Description:      group.Description,
			UpdatedAt:        group.UpdatedAt,
			CreatedAt:        group.CreatedAt,
		}
		groupsDto = append(groupsDto, permissionGroup)
	}
	res.Items = groupsDto
	return res, nil
}

func (service *AuthorizationService) CreatePermissionGroup(ctx context.Context, request *AuthorizationRequests.CreatePermissionGroupRequest) (*models.PermissionGroupModel, error) {
	// Create a new RoleModel instance with the title from the request.
	group := &models.PermissionGroupModel{
		Name:        request.Name,
		Title:       request.Title,
		Description: request.Description,
	}

	// Loop through the permission UUIDs provided in the request.
	for _, requestPermUuid := range request.Permissions {
		// Parse the permission UUID.
		permissionUuid, err := uuid.Parse(requestPermUuid)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to parse permission uuid", service, err)
			return nil, errs.SomeThingWentWrong
		}
		// Retrieve the permission from the repository using the UUID.
		permission, err := service.PermissionRepository.GetByUuid(&permissionUuid)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get permission by uuid", service, err)
			return nil, errs.SomeThingWentWrong
		}
		// Append the permission to the group model's slice.
		group.Permissions = append(group.Permissions, permission)
	}

	// Create the permission group  in the repository.
	groupOrm, err := service.PermissionGroupRepository.Create(group)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to create permission group", service, err)
		if utils.CheckError(errs.DuplicatePermissionGroupName, err) {
			return nil, err
		}
		return nil, errs.SomeThingWentWrong
	}

	// Return the created group.
	return groupOrm, nil
}

func (service *AuthorizationService) GetPermissionGroupByUuid(ctx context.Context, uuid *uuid.UUID, includeRelations ...string) (*models.PermissionGroupModel, error) {
	res, err := service.PermissionGroupRepository.GetByUuid(uuid, includeRelations...)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get permission group by uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AuthorizationService) UpdatePermissionGroupByUuid(ctx context.Context, groupUuid *uuid.UUID, request *AuthorizationRequests.UpdatePermissionGroupRequest) (*models.PermissionGroupModel, error) {
	// Get existing permission group
	group, err := service.GetPermissionGroupByUuid(ctx, groupUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get permission group by uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}

	// Overwrite all fields (PUT semantics)
	group.Name = request.Name
	group.Title = request.Title
	if request.Description != nil {
		group.Description = *request.Description
	}
	group.IsActive = &request.IsActive

	// Build new permission list
	var newPermissions []*models.PermissionModel
	for _, permUuidStr := range request.Permissions {
		permUuid, err := uuid.Parse(permUuidStr)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "invalid permission uuid in request", service, err)
			return nil, errs.SomeThingWentWrong
		}

		perm, err := service.PermissionRepository.GetByUuid(&permUuid)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get permission by uuid", service, err)
			return nil, errs.SomeThingWentWrong
		}

		newPermissions = append(newPermissions, perm)
	}

	group.Permissions = newPermissions

	// Save updated group
	updatedGroup, err := service.PermissionGroupRepository.Update(group)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update permission group", service, err)
		if utils.CheckError(errs.DuplicatePermissionGroupName, err) {
			return nil, err
		}
		return nil, errs.SomeThingWentWrong
	}

	return updatedGroup, nil
}

func (service *AuthorizationService) DeletePermissionGroupByUuid(ctx context.Context, groupUuid *uuid.UUID) error {
	group, err := service.GetPermissionGroupByUuid(ctx, groupUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get permission group by uuid", service, err)
		return errs.SomeThingWentWrong
	}

	err = service.PermissionGroupRepository.Delete(group)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to delete permission group", service, err)
		if utils.CheckError(err, errs.CantDeletePermissionGroup) {
			return err
		}
		return errs.SomeThingWentWrong
	}

	return nil
}

func (service *AuthorizationService) GetRolesList(ctx context.Context, builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	builder.PrioritySorts = []scopes.PrioritySort{
		{
			Cases: []scopes.PriorityCase{
				{
					Conditions: []scopes.Condition{
						{
							Operand:  "name",
							Value:    models.DefaultRole,
							Operator: scopes.EqualOperator,
						},
					},
					Priority: 1,
				},
			},
		},
	}
	res, err := service.RoleRepository.GetList(builder)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get list of role", service, err)
		return nil, errs.SomeThingWentWrong
	}
	roles, ok := res.Items.(*[]*models.RoleModel)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to cast res.Items to []*models.RoleModel", service, nil)
		return nil, errs.SomeThingWentWrong
	}
	for _, role := range *roles {
		role.PermissionGroupCount = len(role.PermissionGroups)
	}
	res.Items = roles
	return res, nil
}

func (service *AuthorizationService) CreateRole(ctx context.Context, request *AuthorizationRequests.CreateRoleRequest) (*models.RoleModel, error) {
	// Create a new RoleModel instance with the title from the request.
	roleModel := &models.RoleModel{
		Title:       request.Title,
		Name:        request.Name,
		Description: request.Description,
	}

	// Loop through the permission group UUIDs provided in the request.
	for _, requestPermUuid := range request.PermissionGroups {
		// Parse the permission UUID.
		permissionGroupUuid, err := uuid.Parse(requestPermUuid)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to parse permission group uuid", service, err)
			return nil, errs.SomeThingWentWrong
		}
		// Retrieve the permission group from the repository using the UUID.
		group, err := service.PermissionGroupRepository.GetByUuid(&permissionGroupUuid)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get permission by uuid", service, err)
			return nil, errs.SomeThingWentWrong
		}
		// Append the permission group to the role model's permission groups slice.
		roleModel.PermissionGroups = append(roleModel.PermissionGroups, group)
	}

	// Create the role in the repository.
	roleOrm, err := service.RoleRepository.Create(roleModel)
	if err != nil {
		if utils.CheckError(errs.DuplicateRoleName, err) {
			return nil, err
		}
		logger.LogErrorWithFieldsV2(ctx, "failed to create role", service, err)
		return nil, errs.SomeThingWentWrong
	}

	// Return the created role.
	return roleOrm, nil
}

func (service *AuthorizationService) GetRoleByUuid(ctx context.Context, roleUuid *uuid.UUID) (*models.RoleModel, error) {
	res, err := service.RoleRepository.GetByUuid(roleUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get role by uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}
	res.PermissionGroupCount = len(res.PermissionGroups)
	return res, nil
}

func (service *AuthorizationService) UpdateRoleByUuid(ctx context.Context, roleUuid *uuid.UUID, request *AuthorizationRequests.UpdateRoleRequest) (*models.RoleModel, error) {
	role, err := service.GetRoleByUuid(ctx, roleUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get role by uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}

	if role.Name == models.SuperAdminRole {
		logger.LogErrorWithFieldsV2(ctx, "failed to update role by uuid", service, nil)
		return nil, errs.CantChangeSuperAdmin
	}

	role.Name = request.Name
	role.Title = request.Title
	role.IsActive = &request.IsActive
	if request.Description != nil {
		role.Description = *request.Description
	}

	// Build new permission list
	var newPermissionGroups []*models.PermissionGroupModel
	for _, permUuidStr := range request.PermissionGroups {
		groupUuid, err := uuid.Parse(permUuidStr)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "invalid permission group uuid in request", service, err)
			return nil, errs.SomeThingWentWrong
		}

		group, err := service.PermissionGroupRepository.GetByUuid(&groupUuid)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get permission group by uuid", service, err)
			return nil, errs.SomeThingWentWrong
		}

		newPermissionGroups = append(newPermissionGroups, group)
	}

	role.PermissionGroups = newPermissionGroups

	// Save updated group
	updatedRole, err := service.RoleRepository.Update(role)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update role", service, err)
		if utils.CheckError(errs.DuplicateRoleName, err) {
			return nil, err
		}
		return nil, errs.SomeThingWentWrong
	}

	return updatedRole, nil
}

func (service *AuthorizationService) DeleteRoleByUuid(ctx context.Context, roleUuid *uuid.UUID) error {
	role, err := service.GetRoleByUuid(ctx, roleUuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get role by uuid", service, err)
		return errs.SomeThingWentWrong
	}

	err = service.RoleRepository.Delete(role)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to delete role", service, err)
		if utils.CheckError(err, errs.CantDeleteThisRole) {
			return err
		}
		return errs.SomeThingWentWrong
	}

	return nil
}

func (service *AuthorizationService) GetUsersForRole(ctx context.Context, uuid *uuid.UUID) ([]*models.UserModel, error) {
	role, err := service.GetRoleByUuid(ctx, uuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get role by uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}

	res, err := service.RoleRepository.GetAssociateUser(role)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get associate users for role", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AuthorizationService) GetAdminsForRole(ctx context.Context, uuid *uuid.UUID) ([]*models.AdminModel, error) {
	role, err := service.GetRoleByUuid(ctx, uuid)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get role by uuid", service, err)
		return nil, errs.SomeThingWentWrong
	}

	res, err := service.RoleRepository.GetAssociateAdmin(role)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get associate admins for role", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AuthorizationService) IsAuthorized(ctx context.Context, userId uint, userType, permission string) (bool, error) {
	var owner interface{}
	var err error
	switch userType {
	case "user":
		owner, err = service.UserRepository.GetById(userId)
	case "admin":
		owner, err = service.AdminRepository.GetById(userId)
	default:
		logger.LogErrorWithFieldsV2(ctx, "unsupported owner type", service, nil,
			zap.String("given-user-type", userType),
		)
		return false, errs.SomeThingWentWrong
	}

	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get owner", service, err)
		return false, errs.SomeThingWentWrong
	}

	// Type assert owner to the appropriate type
	var roles []*models.RoleModel
	switch o := owner.(type) {
	case *models.UserModel:
		roles = o.Roles
	case *models.AdminModel:
		roles = o.Roles
	default:
		logger.LogErrorWithFieldsV2(ctx, "unsupported owner type", service, nil,
			zap.Any("given-owner-type", o),
		)
		return false, errs.SomeThingWentWrong
	}

	// Iterate over user's roles
	for _, role := range roles {
		if role.Name == models.SuperAdminRole {
			return true, nil
		}
		if role.IsActive != nil && !*role.IsActive {
			continue
		}
		if roleHasPermission(role, permission) {
			return true, nil
		}
	}
	// If none of the user's roles have the permission, user is not authorized
	logger.LogErrorWithFieldsV2(ctx, "permission failed", service, nil)
	return false, nil
}

// Helper function to check if a role has a specific permission
func roleHasPermission(role *models.RoleModel, permission string) bool {
	for _, permissionGroup := range role.PermissionGroups {
		if permissionGroup.IsActive != nil && !*permissionGroup.IsActive {
			continue
		}
		for _, perm := range permissionGroup.Permissions {
			if perm.Name == permission {
				return true
			}
		}
	}

	return false
}
