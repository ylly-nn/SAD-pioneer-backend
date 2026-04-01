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
	validate.RegisterValidation("inn", validateINN)
	validate.RegisterValidation("kpp", validateKPP)
	validate.RegisterValidation("ogrn", validateOGRN)
	validate.RegisterValidation("org_name", validateOrgName)
	validate.RegisterValidation("org_short_name", validateOrgShortName)
	validate.RegisterValidation("person_name", validatePersonName)
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("address", validateAddress)
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
		case "Password", "NewPassword":
			fields["password"] = getPasswordErrorMessage(e)
		case "ConfirmPassword":
			fields["confirm_password"] = "Пароли не совпадают"
		case "INN":
			fields["inn"] = "ИНН должен содержать 10 или 12 цифр"
		case "KPP":
			fields["kpp"] = "КПП должен содержать 9 цифр"
		case "OGRN":
			fields["ogrn"] = "ОГРН должен содержать 13 цифр"
		case "OrgName":
			fields["org_name"] = "Название организации должно содержать от 3 до 100 символов (латиница, кириллица, пробел, кавычки)"
		case "OrgShortName":
			fields["org_short_name"] = "Короткое название организации должно содержать от 3 до 50 символов (латиница, кириллица, пробел, кавычки)"
		case "Name", "Surname", "Patronymic":
			fields[strings.ToLower(field)] = "Поле должно содержать только буквы русского алфавита и тире (от 2 до 100 символов)"
		case "Phone":
			fields["phone"] = "Телефон должен содержать 10 цифр"
		case "Address":
			fields["address"] = "Адрес должен содержать от 10 до 50 символов (кириллица, цифры, пробел, точка, запятая)"
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

	// Имя не может начинаться с точки
	if strings.HasPrefix(name, ".") {
		return false, "имя почты не может начинаться с точки"
	}

	// Имя не может заканчиваться точкой
	if strings.HasSuffix(name, ".") {
		return false, "имя почты не может заканчиваться точкой"
	}

	// Имя не может содержать две точки подряд
	if strings.Contains(name, "..") {
		return false, "имя почты не может содержать две точки подряд"
	}

	// Проверка домена (после @)
	if len(domain) == 0 {
		return false, "отсутствует домен после @"
	}

	// Домен не может начинаться с точки
	if strings.HasPrefix(domain, ".") {
		return false, "домен не может начинаться с точки"
	}

	// Проверка на точку в конце
	if strings.HasSuffix(domain, ".") {
		return false, "домен не может заканчиваться точкой"
	}

	// Домен не может содержать две точки подряд
	if strings.Contains(domain, "..") {
		return false, "домен не может содержать две точки подряд"
	}

	// После последней точки должно быть минимум 2 символа
	domainParts := strings.Split(domain, ".")
	lastPart := domainParts[len(domainParts)-1]
	if len(lastPart) < 2 {
		return false, "домен после последней точки должен содержать минимум 2 символа (например, .ru, .com)"
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
		return "Пароль не должен содержать пробелы"
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
		!result.HasSpace &&
		result.LetterCount >= 4 &&
		result.IsLatin &&
		len(password) >= 8 &&
		len(password) <= 24
}

// validateINN проверяет ИНН
func validateINN(fl validator.FieldLevel) bool {
	inn := fl.Field().String()
	if len(inn) != 10 && len(inn) != 12 {
		return false
	}
	for _, ch := range inn {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// validateKPP проверяет КПП
func validateKPP(fl validator.FieldLevel) bool {
	kpp := fl.Field().String()
	if len(kpp) != 9 {
		return false
	}
	for _, ch := range kpp {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// validateOGRN проверяет ОГРН
func validateOGRN(fl validator.FieldLevel) bool {
	ogrn := fl.Field().String()
	if len(ogrn) != 13 {
		return false
	}
	for _, ch := range ogrn {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// validateOrgName проверяет название организации (латиница, кириллица, кавычки, пробелы)
func validateOrgName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	length := 0
	for range name {
		length++
	}

	if length < 3 || length > 100 {
		return false
	}

	// Разрешены: латиница, кириллица, пробелы, кавычки (")
	for _, ch := range name {
		// Латиница
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			continue
		}
		// Кириллица
		if (ch >= 'а' && ch <= 'я') || (ch >= 'А' && ch <= 'Я') || ch == 'ё' || ch == 'Ё' {
			continue
		}
		// Пробел и кавычки
		if ch == ' ' || ch == '"' {
			continue
		}
		return false
	}
	return true
}

// validateOrgShortName проверяет короткое название организации (латиница, кириллица, кавычки, пробелы)
func validateOrgShortName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	length := 0
	for range name {
		length++
	}

	if length < 3 || length > 50 {
		return false
	}

	// Разрешены: латиница, кириллица, пробелы, кавычки (")
	for _, ch := range name {
		// Латиница
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			continue
		}
		// Кириллица
		if (ch >= 'а' && ch <= 'я') || (ch >= 'А' && ch <= 'Я') || ch == 'ё' || ch == 'Ё' {
			continue
		}
		// Пробел и кавычки
		if ch == ' ' || ch == '"' {
			continue
		}
		return false
	}
	return true
}

// validatePersonName проверяет имя/фамилию/отчество (только кириллица и тире)
func validatePersonName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	length := 0
	for range name {
		length++
	}

	if length < 2 || length > 100 {
		return false
	}

	for _, ch := range name {
		// Кириллица
		if (ch >= 'а' && ch <= 'я') || (ch >= 'А' && ch <= 'Я') || ch == 'ё' || ch == 'Ё' {
			continue
		}
		// Тире
		if ch == '-' {
			continue
		}
		return false
	}
	return true
}

// validatePhone проверяет телефон
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) != 10 {
		return false
	}
	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// validateAddress проверяет адрес (кириллица, цифры, пробел, точка, запятая)
func validateAddress(fl validator.FieldLevel) bool {
	address := fl.Field().String()
	length := 0
	for range address {
		length++
	}

	if length < 10 || length > 50 {
		return false
	}

	for _, ch := range address {
		// Кириллица
		if (ch >= 'а' && ch <= 'я') || (ch >= 'А' && ch <= 'Я') || ch == 'ё' || ch == 'Ё' {
			continue
		}
		// Цифры
		if ch >= '0' && ch <= '9' {
			continue
		}
		// Пробел, точка, запятая
		if ch == ' ' || ch == '.' || ch == ',' {
			continue
		}
		return false
	}
	return true
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
