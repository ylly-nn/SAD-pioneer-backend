package config

import (
	"fmt"
	"net/smtp"
	"os"
	"time"
)

// Отправка писем через SMTP
type SMTPEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// Создание нового экземпляра SMTPEmailService
func NewSMTPEmailService(host, port, username, password, from string) *SMTPEmailService {
	return &SMTPEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// Создает и настраивает SMTP сервис из env
func ConnectSMTP() (*SMTPEmailService, error) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	return &SMTPEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}, nil
}

// Интерфейс для отправки email
type EmailSender interface {
	SendVerificationCode(toEmail, code string) error
}

// Отправляние кода подтверждения
func (s *SMTPEmailService) SendVerificationCode(toEmail, code string) error {
	subject := "Код подтверждения регистрации"

	// Тело письма
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h2>Подтверждение регистрации</h2>
				</div>
				<div class="content">
					<p>Для завершения регистрации введите следующий код подтверждения:</p>
					<div class="code">%s</div>
					<p>Код действителен в течение 10 минут.</p>
					<p>Если вы не регистрировались, просто проигнорируйте это письмо.</p>
				</div>
				<div class="footer">
					<p>© %d Pioneer</p>
				</div>
			</div>
		</body>
		</html>
	`, code, time.Now().Year())

	// Заголовки
	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Сборка сообщения
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Аутентификация
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Адрес сервера
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	// Отправление с таймаутом
	errChan := make(chan error, 1)
	go func() {
		errChan <- smtp.SendMail(addr, auth, s.from, []string{toEmail}, []byte(message))
	}()

	// Ожидание ответа (15 секунд)
	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("ошибка отправки email: %w", err)
		}
		return nil
	case <-time.After(15 * time.Second):
		return fmt.Errorf("таймаут отправки email")
	}
}
