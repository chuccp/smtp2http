package entity

type SendMail struct {
	SMTPId     uint     `json:"SMTPId"`
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Content    string   `json:"content"`
}

type SendMailApi struct {
	Token      string   `json:"token"`
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Content    string   `json:"content"`
}
