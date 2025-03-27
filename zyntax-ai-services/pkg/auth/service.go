package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/domain"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type authService struct {
	authRepository domain.AuthRepository
}

func NewAuthService(r domain.AuthRepository) domain.AuthService {
	return &authService{r}
}
func GetUserDetails(accessToken string) (map[string]interface{}, error) {
	userURL := fmt.Sprintf(viper.GetString("kku.oauth.host") + "/api/v1/user")
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(respBody))
	}

	var userDetails map[string]interface{}
	err = json.Unmarshal(respBody, &userDetails)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return userDetails, nil
}
func GetAccessToken(code string) (string, error) {
	tokenURL := fmt.Sprintf(viper.GetString("kku.oauth.host") + "/token")
	bodyParams := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {viper.GetString("kku.oauth.client_id")},
		"client_secret": {viper.GetString("kku.oauth.client_secret")},
		"code":          {code},
	}
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(bodyParams.Encode()))
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return "", fiber.NewError(fiber.StatusInternalServerError, string(respBody))
	}
	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	accessToken, ok := response["access_token"].(string)
	if !ok {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Access token not found in the response")
	}
	return accessToken, nil
}

func (s *authService) Login(code string) (*models.MainUser, error) {
	fmt.Println("code", code)
	accessToken, err := GetAccessToken(code)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	fmt.Println("accessToken", accessToken)
	userDetails, err := GetUserDetails(accessToken)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	fmt.Println("userDetails", userDetails)
	userID := strings.Split(userDetails["sso_mail"].(string), "@")[0]
	fullname := userDetails["sso_firstname"].(string) + " " + userDetails["sso_lastname"].(string)
	fullnameEng := userDetails["sso_firstname_eng"].(string) + " " + userDetails["sso_lastname_eng"].(string)
	user := models.MainUser{
		ID:     userID,
		NameTh: fullname,
		NameEn: fullnameEng,
		Email:  userDetails["sso_mail"].(string),
	}
	login, err := s.authRepository.Login(user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return &models.MainUser{
		ID:     login.ID,
		NameTh: login.NameTh,
		NameEn: login.NameEn,
		Email:  login.Email,
	}, nil
}
