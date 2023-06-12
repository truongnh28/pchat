package service

import (
	"bytes"
	"chat-app/config"
	gomail "gopkg.in/gomail.v2"
	"log"
	"text/template"
)

type MailService interface {
	SendOpt(receiver []string, opt string) error
}

type mailImpl struct {
	mailConfig       *config.MailConfig
	mailTemplatePath string
}

type Otp struct {
	Otp string `json:"Otp"`
}

func NewMailService(config *config.MailConfig, mailTemplatePath string) MailService {
	return &mailImpl{
		mailConfig:       config,
		mailTemplatePath: mailTemplatePath,
	}
}

func (m *mailImpl) SendOpt(receiver []string, opt string) error {
	var (
		smtpHost     = m.mailConfig.SmtpHost
		smtpPort     = m.mailConfig.SmtpPort
		mailSender   = m.mailConfig.MailSender
		password     = m.mailConfig.Password
		templatePath = m.mailTemplatePath
		err          = error(nil)
		templateName = "template.html"
		t            = template.New(templateName)
		//tpl          bytes.Buffer
		body bytes.Buffer
	)
	t, err = t.ParseFiles(templatePath)
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(&body, templateName, &Otp{Otp: opt})
	if err != nil {
		log.Println(err)
		return err
	}
	var mail = gomail.NewMessage()
	mail.SetHeader("From", mailSender)
	mail.SetHeader("To", receiver...)
	// m.SetAddressHeader("Cc", "<RECIPIENT CC>", "<RECIPIENT CC NAME>")
	mail.SetHeader("Subject", "PChat OTP")
	mail.SetBody("text/html", body.String())
	// m.Attach(t) // attach whatever you want

	d := gomail.NewDialer(smtpHost, smtpPort, mailSender, password)

	return d.DialAndSend(mail)
}
