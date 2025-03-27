package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Zentrix-Software-Hive/zyntax-ai-services/internal/handlers"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/domain"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	helpers "github.com/zercle/gofiber-helpers"
)

type AuthHandler struct {
	authService domain.AuthService
	*models.JwtResources
	*handlers.RouterResources
}

func NewAuthHandler(authRoute fiber.Router, authService domain.AuthService, jwt *models.JwtResources) {
	handler := &AuthHandler{
		JwtResources: jwt,
		authService:  authService,
	}
	authRoute.Post("/", handler.Login())
}

// @Summary Login
// @Description https://oauth.kku.ac.th/authorize?response_type=code&client_id=e8fdb4894be17a3a&redirect_uri=http://localhost:8080/api/v1/swagger/index.html
// @Tags Auth
// @Accept json
// @Produce json
// @Param code body models.Oauth true "Code from oauth"
// @Router /api/v1/auth/ [post]
func (h *AuthHandler) Login() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		responseForm := helpers.ResponseForm{}
		request := models.Oauth{}
		err = c.BodyParser(&request)
		if err != nil {
			responseForm.Errors = []helpers.ResponseError{
				{
					Code:    http.StatusBadRequest,
					Source:  helpers.WhereAmI(),
					Message: err.Error(),
				},
			}
			return c.Status(http.StatusBadRequest).JSON(responseForm)
		}
		login, err := h.authService.Login(request.Code)
		if err != nil {
			responseForm.Errors = []helpers.ResponseError{
				{
					Code:    http.StatusUnauthorized,
					Source:  helpers.WhereAmI(),
					Message: err.Error(),
				},
			}
			return c.Status(http.StatusUnauthorized).JSON(responseForm)
		}
		fmt.Println("login", login)
		token := jwt.NewWithClaims(h.JwtSigningMethod, &jwt.RegisteredClaims{})
		claims := token.Claims.(*jwt.RegisteredClaims)
		claims.Subject = login.ID
		claims.Issuer = c.Hostname()
		claims.Audience = []string{fmt.Sprintf("%v:%v", "highest permission", "5")}
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
		signToken, err := token.SignedString(h.JwtSignKey)
		if err != nil {
			responseForm.Errors = []helpers.ResponseError{
				{
					Code:    http.StatusUnauthorized,
					Source:  helpers.WhereAmI(),
					Message: err.Error(),
				},
			}
			return c.Status(http.StatusUnauthorized).JSON(responseForm)
		}
		responseForm.Result = fiber.Map{
			"token": signToken,
		}
		responseForm.Messages = []string{login.NameTh + " have been login successfully."}
		responseForm.Success = true
		return c.Status(http.StatusOK).JSON(responseForm)
	}
}
