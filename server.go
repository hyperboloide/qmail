package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	html "html/template"
	text "text/template"

	log "github.com/Sirupsen/logrus"
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
	Queue     dispatch.PersistantQueue
	Sender    string
	Templates *text.Template
	Body      *html.Template
}

func doError(err error, email *client.Mail) error {
	if email != nil {
		log.WithFields(map[string]interface{}{
			"template": email.Template,
			"dests":    email.Dests,
			"files":    email.Files,
		}).Error(err)
	}
	return err
}

// Listenner is called by the Listen polling function of the queue when a
// message is available
func (m Mailer) Listenner(buff []byte) error {
	email := &client.Mail{}
	if err := json.Unmarshal(buff, email); err != nil {
		return doError(err, nil)
	}

	textBuff := &bytes.Buffer{}
	tmplName := fmt.Sprintf("%s.md", email.Template)

	if t := m.Templates.Lookup(tmplName); t == nil {
		return doError(fmt.Errorf("template '%s' not found", tmplName), email)
	} else if err := t.Execute(textBuff, email.Data); err != nil {
		return doError(err, email)
	}

	md := blackfriday.MarkdownCommon(textBuff.Bytes())
	htmlBuff := &bytes.Buffer{}
	err := m.Body.Execute(
		htmlBuff,
		html.HTML(string(md[:])))
	if err != nil {
		return doError(err, email)
	}

	messages := []*gomail.Message{}
	for _, dest := range email.Dests {
		gm := gomail.NewMessage()
		gm.SetHeader("From", m.Sender)
		gm.SetHeader("To", dest)
		gm.SetHeader("Subject", email.Subject)
		gm.SetBody("text/plain", string(textBuff.Bytes()[:]))
		gm.AddAlternative("text/html", string(htmlBuff.Bytes()[:]))
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

	if err := d.DialAndSend(messages...); err != nil {
		return doError(err, email)
	}

	log.WithFields(map[string]interface{}{
		"template": email.Template,
	}).Infof("Message sent to '%d' recipients", len(email.Dests))

	return nil
}
