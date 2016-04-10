package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/go-gomail/gomail"
	"github.com/hyperboloide/dispatch"
	"github.com/hyperboloide/qmail/client"
	"github.com/russross/blackfriday"
)

// SMTPConf contains connection informations
type SMTPConf struct {
	Host     string
	Port     int
	User     string
	Password string
}

// Mailer reads mails from the queue and send message to smtp
type Mailer struct {
	SMTP      SMTPConf
	Queue     dispatch.Queue
	Sender    string
	Templates *template.Template
}

// Listenner is call by the Listen polling function of the queue when a
// message is available
func (m Mailer) Listenner(buff []byte) error {
	email := &client.Mail{}
	if err := json.Unmarshal(buff, email); err != nil {
		return err
	}

	textBuff := &bytes.Buffer{}

	tmplName := fmt.Sprintf("%s.md", email.Template)
	if t := m.Templates.Lookup(tmplName); t == nil {
		return fmt.Errorf("template '%s' not found", tmplName)
	} else if err := t.Execute(textBuff, email.Data); err != nil {
		return err
	}

	html := blackfriday.MarkdownCommon(textBuff.Bytes())

	messages := []*gomail.Message{}
	for _, dest := range email.Dests {
		gm := gomail.NewMessage()
		gm.SetHeader("From", m.Sender)
		gm.SetHeader("To", dest)
		gm.SetHeader("Subject", email.Subject)
		gm.SetBody("text/plain", string(textBuff.Bytes()[:]))
		gm.AddAlternative("text/html", string(html[:]))
		for _, f := range email.Files {
			gm.Attach(f)
		}
		messages = append(messages, gm)
	}

	d := gomail.NewDialer(
		m.SMTP.Host,
		m.SMTP.Port,
		m.SMTP.User,
		m.SMTP.Password)

	return d.DialAndSend(messages...)
}
