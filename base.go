package aegorm

import (
	"github.com/AnCodyEI/Aegorm/ressources/models"
	"github.com/AnCodyEI/Aegorm/ressources/repositories"
	"gorm.io/gorm"
)

type Aegorm struct {
	repositories.IRoleRepository
}

func New(db *gorm.DB) (a *Aegorm, err error) {
	roleRepository := &repositories.SRoleRepository{Database: db}

	a = &Aegorm{roleRepository}

	err = db.AutoMigrate(&models.Role{}, &models.UserRole{})
	return
}
