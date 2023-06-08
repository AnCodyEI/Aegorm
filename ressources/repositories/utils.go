package repositories

import (
	"github.com/AnCodyEI/Aegorm/ressources/models"
)

func isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func isUint(value interface{}) bool {
	_, ok := value.(uint)
	return ok
}

func isRole(value interface{}) bool {
	_, ok := value.(models.Role)
	return ok
}
