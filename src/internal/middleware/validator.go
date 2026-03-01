package middleware

import (
	"encoding/json"
	"net/http"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("password", validatePassword)
}

// ValidateStruct валидирует структуру и возвращает ошибку
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// PasswordValidationResult детализированный результат проверки пароля
type PasswordValidationResult struct {
	HasUpper    bool
	HasLower    bool
	HasNumber   bool
	HasSpecial  bool
	LetterCount int
	IsLatin     bool
	Length      int
}

// SendValidationError отправляет структурированную ошибку валидации
func SendValidationError(w http.ResponseWriter, err error) {
	validationErrors := err.(validator.ValidationErrors)
	fields := make(map[string]string)

	for _, e := range validationErrors {
		field := e.Field()
		switch field {
		case "Email":
			fields["email"] = getEmailErrorMessage(e)
		case "Password":
			fields["password"] = getPasswordErrorMessage(e)
		case "ConfirmPassword":
			fields["confirm_password"] = "Пароли не совпадают"
		default:
			fields[field] = "Некорректное значение"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := map[string]interface{}{
		"error":  "Ошибка валидации",
		"fields": fields,
	}

	json.NewEncoder(w).Encode(response)
}

// getEmailErrorMessage возвращает понятное сообщение для ошибок email
func getEmailErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "Email обязателен для заполнения"
	case "email":
		return "Некорректный формат email"
	case "min":
		return "Email должен содержать минимум 6 символов"
	case "max":
		return "Email должен содержать максимум 100 символов"
	default:
		return "Некорректный email"
	}
}

// getPasswordErrorMessage возвращает понятное сообщение для ошибок пароля
func getPasswordErrorMessage(e validator.FieldError) string {
	password := e.Value().(string)

	var messages []string

	if len(password) < 8 {
		messages = append(messages, "минимум 8 символов")
	}
	if len(password) > 20 {
		messages = append(messages, "максимум 20 символов")
	}

	result := validatePasswordDetail(password)

	if !result.HasUpper {
		messages = append(messages, "хотя бы одна заглавную букву")
	}
	if !result.HasLower {
		messages = append(messages, "хотя бы одна строчную букву")
	}
	if !result.HasNumber {
		messages = append(messages, "хотя бы одну цифру")
	}
	if !result.HasSpecial {
		messages = append(messages, "хотя бы один спецсимвол (!@#$%^)")
	}
	if result.LetterCount < 4 {
		messages = append(messages, "минимум 4 буквы")
	}
	if !result.IsLatin {
		messages = append(messages, "только латинские буквы")
	}

	if len(messages) == 0 {
		return "Некорректный пароль"
	}

	return "Пароль должен содержать: " + joinMessages(messages)
}

// validatePasswordDetail возвращает детальную информацию о пароле
func validatePasswordDetail(password string) PasswordValidationResult {
	result := PasswordValidationResult{
		Length: len(password),
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			result.HasUpper = true
			result.LetterCount++
		case unicode.IsLower(char):
			result.HasLower = true
			result.LetterCount++
		case unicode.IsNumber(char):
			result.HasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			result.HasSpecial = true
		}
	}

	// Проверка на латиницу
	result.IsLatin = true
	for _, char := range password {
		if unicode.IsLetter(char) && !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z') {
			result.IsLatin = false
			break
		}
	}

	return result
}

// validatePassword проверяет сложность пароля
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	result := validatePasswordDetail(password)

	return result.HasUpper &&
		result.HasLower &&
		result.HasNumber &&
		result.HasSpecial &&
		result.LetterCount >= 4 &&
		result.IsLatin &&
		len(password) >= 8 &&
		len(password) <= 20
}

// joinMessages объединяет сообщения в строку
func joinMessages(messages []string) string {
	if len(messages) == 0 {
		return ""
	}
	if len(messages) == 1 {
		return messages[0]
	}

	result := ""
	for i, msg := range messages {
		if i == len(messages)-1 {
			result += "и " + msg
		} else if i == 0 {
			result += msg
		} else {
			result += ", " + msg
		}
	}
	return result
}
