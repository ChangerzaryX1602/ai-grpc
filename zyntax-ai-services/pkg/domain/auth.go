package domain

import (
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/models"
)

type AuthRepository interface {
	Login(models.MainUser) (*models.MainUser, error)
}

type AuthService interface {
	Login(string) (*models.MainUser, error)
}
