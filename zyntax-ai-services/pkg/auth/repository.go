package auth

import (
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/domain"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type authRepository struct {
	*gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{DB: db}
}
func (r *authRepository) Login(user models.MainUser) (*models.MainUser, error) {
	if r == nil {
		err := fiber.NewError(fiber.StatusServiceUnavailable, "Database server has gone away")
		return nil, err
	}
	err := r.First(&user, "user_id = ?", user.ID)
	if err.Error != nil {
		r.Create(&user)
		return &user, nil
	}
	return &user, nil
}
