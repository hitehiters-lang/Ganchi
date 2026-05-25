package contacts

import (
	"encoding/json"
	"fmt"
	"ganchi_app/additional"
	"ganchi_app/config"
	"ganchi_app/models"
	"net/http"
	"net/mail"
	"strings"

	"gopkg.in/gomail.v2"
)

// @Tags Связь
// @Summary Отправить письмо
// @Description Отправляет письмо на почту компании
// @ID contact_us
// @Accept json
// @Produce json
// @Param contacter body models.ContactMail true "Поля ввода"
// @Success 200 {object} map[string]string "{"message":"mail sent"}"
// @Failure 400 {object} map[string]string "{"error": "bad request"} {"error": "name required"} {"error": "bad email format"} {"error": "content required"}"
// @Failure 500 {object} map[string]string "{"error": "failed to send mail"}"
// @Router /api/contact/send_mail [post]
func ContactUs(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	var input models.ContactMail
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		additional.PrintError(path, err)
		return
	}
	if input.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "email required"})
		additional.PrintError(path, fmt.Errorf("Не был введён Email"))
		return
	}
	if input.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "name required"})
		additional.PrintError(path, fmt.Errorf("Не было введно имя"))
		return
	}
	if input.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "content required"})
		additional.PrintError(path, fmt.Errorf("Не было введно содержание обращения"))
		return
	}
	if input.Phone == "" {
		input.Phone = "Не указан"
	}
	if errx := ValidateEmail(input.Email); errx != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad email format"})
		additional.PrintError(path, errx)
		return
	}

	email_from, email_to, pass := config.GetEmailData()
	m := gomail.NewMessage()
	m.SetHeader("From", email_from)
	m.SetHeader("To", email_to)
	m.SetHeader("Subject", "Обращение к Hebei Gangchi от "+input.Name)
	m.SetBody("text/plain", fmt.Sprintf(`
	Обращение к Hebei Gangchi
	Имя: %v
	Электронная почта: %v
	Номер телефона: %v
	Содержание обращения:
	%v
	`, input.Name, input.Email, input.Phone, input.Content))

	d := gomail.NewDialer("smtp.gmail.com", 587, email_from, pass)

	if err := d.DialAndSend(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to send mail"})
		additional.PrintServerError(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "mail sent"})
	additional.PrintSuccess(path, "Обращение от пользователя "+input.Name+" доставлено")
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("Invalid email format")
	}
	if len(email) > 254 {
		return fmt.Errorf("Email too long")
	}
	return nil
}
