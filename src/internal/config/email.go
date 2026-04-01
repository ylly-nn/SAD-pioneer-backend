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
	SendVerificationResetCode(toEmail, code string) error
}

// Отправление кода подтверждения для регистрации
func (s *SMTPEmailService) SendVerificationCode(toEmail, code string) error {
	subject := "Код подтверждения регистрации"

	// Тело письма
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
  </head>

  <body style="margin:0; padding:0; background:#0f1c1e; font-family:Arial, sans-serif; color:white;">

    <table width="100%%" cellpadding="0" cellspacing="0" style="padding:40px 0;">
      <tr>
        <td align="center">

          <table width="520" cellpadding="0" cellspacing="0" style="background:#1c2a2c; border-radius:16px; overflow:hidden;">

            <tr>
              <td style="padding:30px; text-align:center; background:rgb(22,71,71);">
                <div style="font-size:26px; font-weight:bold; letter-spacing:2px;">
                  PIONEER
                </div>
              </td>
            </tr>

            <tr>
              <td style="padding:35px 40px; text-align:left;">

                <h2 style="margin-top:0; color:white;">
                  Код подтверждения
                </h2>

                <p style="color:#cfd8dc; line-height:1.6;">
                  Чтобы завершить регистрацию, используйте следующий код подтверждения:
                </p>

                <div style="
                  font-size:40px;
                  font-weight:bold;
                  letter-spacing:8px;
                  margin:30px 0;
                  padding:18px 25px;
                  background:rgb(22,71,71);
                  border-radius:10px;
                  text-align:center;
                ">
                  %s
                </div>

                <p style="color:#cfd8dc;">
                  Код действителен в течение <b>10 минут</b>.
                </p>

                <p style="color:#cfd8dc; margin-top:30px;">
                  С уважением,<br>
                  <b>Pioneer</b>
                </p>

              </td>
            </tr>

            <tr>
              <td style="padding:20px; text-align:center; font-size:12px; color:#9ea7aa; background:#182426;">
                Если вы не регистрировались, просто проигнорируйте это письмо.
              </td>
            </tr>

            <tr>
              <td style="padding:20px; text-align:center; font-size:12px; color:#9ea7aa;">
                © %d Pioneer
              </td>
            </tr>

          </table>

        </td>
      </tr>
    </table>

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

// Отправление кода подтверждения для восстановления пароля
func (s *SMTPEmailService) SendVerificationResetCode(toEmail, code string) error {
	subject := "Код подтверждения для восстановления пароля"

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
					<h2>Восстановление пароля</h2>
				</div>
				<div class="content">
					<p>Для восстановления пароля введите следующий код подтверждения:</p>
					<div class="code">%s</div>
					<p>Код действителен в течение 10 минут.</p>
					<p>Если запрашивали код не вы, просто проигнорируйте это письмо.</p>
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
