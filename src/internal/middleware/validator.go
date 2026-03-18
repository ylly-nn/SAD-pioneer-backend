package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("email", validateEmail)
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	isValid, _ := validateEmailDetail(email)
	return isValid
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
	HasSpace    bool
	LetterCount int
	IsLatin     bool
	Length      int
}

// SendValidationError отправляет структурированную ошибку валидации
func SendValidationError(w http.ResponseWriter, err error) {

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

// validateEmailDetail проверяет email на допустимые символы
func validateEmailDetail(email string) (isValid bool, invalidChars string) {
	// Проверка на пробелы
	if strings.Contains(email, " ") {
		return false, "email не должен содержать пробелы"
	}

	// Проверка email на допустимые символы
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_@"
	for _, char := range email {
		if char > 127 {
			return false, "email содержит недопустимые символы. Разрешены только латинские буквы, цифры и символы . - _ @"
		}
		charStr := string(char)
		if !strings.Contains(allowedChars, charStr) {
			return false, "email содержит недопустимый символ '" + charStr + "'. Разрешены: латиница, цифры, '.', '-', '_', '@'"
		}
	}

	// Проверка на количество @
	atCount := strings.Count(email, "@")
	if atCount != 1 {
		return false, "email должен содержать ровно один символ '@'"
	}

	parts := strings.Split(email, "@")
	name := parts[0]
	domain := parts[1]

	// Проверка имени (до @)
	if len(name) == 0 {
		return false, "отсутствует часть перед @"
	}

	// Проверка домена (после @)
	if len(domain) == 0 {
		return false, "отсутствует домен после @"
	}

	// Проверка на точку в конце
	if strings.HasSuffix(domain, ".") {
		return false, "домен не может заканчиваться точкой"
	}

	// Проверка, что домен не начинается с точки
	if strings.HasPrefix(domain, ".") {
		return false, "домен не может начинаться с точки"
	}
	return true, ""
}

// getEmailErrorMessage возвращает понятное сообщение для ошибок email
func getEmailErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "Email обязателен для заполнения"
	case "email":
		email := e.Value().(string)
		if isValid, message := validateEmailDetail(email); !isValid {
			return message
		}
		return "Некорректный формат email"
	case "min":
		return "Email должен содержать минимум 6 символов"
	case "max":
		return "Email должен содержать максимум 64 символа"
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
	if len(password) > 24 {
		messages = append(messages, "максимум 24 символа")
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
		messages = append(messages, "хотя бы один спецсимвол")
	}
	if result.LetterCount < 4 {
		messages = append(messages, "минимум 4 буквы")
	}
	if !result.IsLatin {
		messages = append(messages, "только латинские буквы")
	}
	if result.HasSpace {
		messages = append(messages, "не должен содержать пробелы")
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

	specialChars := "~!?@#$%^&*_-+()[]{}/'.,:;"

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
		case strings.ContainsRune(specialChars, char):
			result.HasSpecial = true
		case char == ' ':
			result.HasSpace = true
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
		len(password) <= 24
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
			result += " и " + msg
		} else if i == 0 {
			result += msg
		} else {
			result += ", " + msg
		}
	}
	return result
}
