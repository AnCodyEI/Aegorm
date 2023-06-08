package repositories

import (
	"errors"
	"github.com/AnCodyEI/Aegorm/ressources/models"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetRoleByID(ID uint) (role models.Role, err error)
	GetRoleByName(name string) (role models.Role, err error)
	GetRoles() (role models.Roles, err error)
	CreateRole(name string, description string) (ID uint, err error)
	DeleteRole(role models.Role) (err error)

	AddUserRole(UserID uint, r interface{}) error
	RemoveUserRole(UserID uint, r interface{}) error
	UserHasRole(UserID uint, r interface{}) (bool, error)
}

type SRoleRepository struct {
	Database *gorm.DB
}

func (repository *SRoleRepository) GetRoleByID(ID uint) (role models.Role, err error) {
	err = repository.Database.First(&role, "roles.id = ?", ID).Error
	return
}

func (repository *SRoleRepository) GetRoleByName(name string) (role models.Role, err error) {
	err = repository.Database.First(&role, "roles.guard_name = ?", slug.Make(name)).Error
	return
}

func (repository *SRoleRepository) GetRoles() (roles models.Roles, err error) {
	err = repository.Database.Find(&roles).Error
	return
}

func (repository *SRoleRepository) CreateRole(name string, description string) (ID uint, err error) {
	err = repository.Database.Create(&models.Role{Name: name, Description: description}).Error
	return
}

func (repository *SRoleRepository) DeleteRole(role models.Role) (err error) {
	tx := repository.Database.Begin()

	if err = tx.Delete(&models.UserRoles{}, "user_roles.role_id = ?", role.ID).Error; err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Delete(&role).Error; err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

// AddUserRole adds a user role to the repository based on the provided UserID and role parameter.
// The role parameter can be either a role ID (uint), role name (string), or a role struct (models.Role).
func (repository *SRoleRepository) AddUserRole(UserID uint, r interface{}) error {
	if isUint(r) {
		role, err := repository.GetRoleByID(r.(uint))
		if err != nil {
			return err
		}
		return repository.Database.Create(&models.UserRole{UserID: UserID, RoleID: role.ID}).Error
	} else if isString(r) {
		role, err := repository.GetRoleByName(r.(string))
		if err != nil {
			return err
		}
		return repository.Database.Create(&models.UserRole{UserID: UserID, RoleID: role.ID}).Error
	} else if isRole(r) {
		return repository.Database.Create(&models.UserRole{UserID: UserID, RoleID: r.(models.Role).ID}).Error
	} else {
		return errors.New("bad parameters")
	}
}

// RemoveUserRole remove a user role to the repository based on the provided UserID and role parameter.
// The role parameter can be either a role ID (uint), role name (string), or a role struct (models.Role).
func (repository *SRoleRepository) RemoveUserRole(UserID uint, r interface{}) error {
	if isUint(r) {
		role, err := repository.GetRoleByID(r.(uint))
		if err != nil {
			return err
		}
		return repository.Database.Delete(&models.UserRole{}, "user_id = ? AND role_id = ?", UserID, role.ID).Error
	} else if isString(r) {
		role, err := repository.GetRoleByName(r.(string))
		if err != nil {
			return err
		}
		return repository.Database.Delete(&models.UserRole{}, "user_id = ? AND role_id = ?", UserID, role.ID).Error
	} else if isRole(r) {
		return repository.Database.Delete(&models.UserRole{}, "user_id = ? AND role_id = ?", UserID, r.(models.Role).ID).Error
	} else {
		return errors.New("bad parameters")
	}
}

// UserHasRole verifies a user's role in the repository on the basis of the identifier and role supplied.
// The role parameter can be either a role ID (uint), role name (string), or a role struct (models.Role).
func (repository *SRoleRepository) UserHasRole(UserID uint, r interface{}) (bool, error) {
	var role models.Role
	var count int64
	var err error

	if isUint(r) {
		role, err = repository.GetRoleByID(r.(uint))
		if err != nil {
			return false, err
		}
	} else if isString(r) {
		role, err = repository.GetRoleByName(r.(string))
		if err != nil {
			return false, err
		}
	} else if isRole(r) {
		role = r.(models.Role)
	} else {
		return false, errors.New("bad parameters")
	}

	repository.Database.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", UserID, role.ID).Count(&count)
	return count > 0, nil
}
