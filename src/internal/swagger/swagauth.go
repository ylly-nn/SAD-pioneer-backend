package swagger

import (
	"src/internal/auth"
)

type RegisterResponse struct {
	Message string `json:"message" example:"Verification code sent to email"`
	Email   string `json:"email" example:"user@example.com"`
}

// Register обрабатывает POST /auth/register, выполняет регистрацию
// @Summary      Регистрация нового пользователя
// @Description  Создаёт нового пользователя и отправляет код подтверждения на email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.RegisterRequest true "Данные для регистрации"
// @Success      202  {object}  RegisterResponse  "Сообщение об отправке кода и email пользователя"
// @Failure      400  {string}  string  "Invalid request body, missing fields, or validation error"
// @Failure      405  {string}  string  "Method not allowed"
// @Failure      409  {string}  string  "user already exists"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /auth/register [post]
func postAuthRegister() {
	var _ = auth.RefreshRequest{}
}

type VerifyCodeResponse struct {
	Message string `json:"message" example:"User successfully registered"`
}

// VerifyCode обрабатывает POST /auth/verify, отправляет код подтверждения
// @Summary      Подтверждение регистрации по коду
// @Description  Проверяет код подтверждения, отправленный на email, и завершает регистрацию пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.VerifyCodeRequest true "Email и код подтверждения"
// @Success      201  {object}  VerifyCodeResponse  "Сообщение об успешной регистрации"
// @Failure      400  {string}  string  "Invalid request body, missing email/code, invalid code, or code expired"
// @Failure      405  {string}  string  "Method not allowed"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /auth/verify [post]
func postAuthVerify() {
	var _ = auth.VerifyCodeRequest{}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// Login обрабатывает POST /auth/login, выполняет вход в аккаунт
// @Summary      Вход пользователя
// @Description  Аутентифицирует пользователя по email и паролю, возвращает access и refresh токены
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.LoginRequest true "Учётные данные пользователя"
// @Success      200  {object}  TokenResponse  "Токены доступа и обновления"
// @Failure      400  {string}  string  "Invalid request body or missing email/password"
// @Failure      401  {string}  string  "Invalid credentials"
// @Failure      405  {string}  string  "Method not allowed"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /auth/login [post]
func postAuthLogin() {
	var _ = auth.LoginRequest{}
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type LogoutResponse struct {
	Message string `json:"message" example:"Successfully logged out"`
}

// postAuthLogout Post /auth/logout
// @Summary      Выход пользователя
// @Description  Инвалидирует refresh token, завершая сессию пользователя.
// @Tags 		 auth
// @Accept       json
// @Produce      json
// @Param        request body LogoutRequest true "Refresh токен"
// @Success      200  {object}  LogoutResponse  "Успешный выход"
// @Failure      400  {string}  string  "Invalid request body or refresh token is required"
// @Failure      405  {string}  string  "Method not allowed"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /auth/logout [post]
func postAuthLogout() {
	var _ = LogoutRequest{}
	var _ = LogoutResponse{}
}

// postAuthRefresh обновляет access токен
// @Summary      Обновление токенов доступа
// @Description  Принимает refresh токен и выдаёт новую пару access/refresh токенов.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.RefreshRequest true "Refresh токен"
// @Success      200  {object}  auth.TokenResponse  "Новые токены"
// @Failure      400  {string}  string  "Invalid request body or missing refresh token"
// @Failure      401  {string}  string  "Invalid or expired refresh token"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /auth/refresh [post]
func postAuthRefresh() {
	var _ = auth.RefreshRequest{}
	var _ = auth.TokenResponse{}
}
