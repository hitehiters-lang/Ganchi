package models

type ContactMail struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone,omitempty"`
	Content string `json:"content,omitempty"`
}
