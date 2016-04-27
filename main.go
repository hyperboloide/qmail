package main

import (
	html "html/template"
	"io/ioutil"
	"os"
	"strconv"
	text "text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/hyperboloide/dispatch"
)

var (
	// MainMailer connect to the queue and handle messages
	MainMailer *Mailer
)

const defaultBody = `
<!doctype html>
<html>
  <body>{{ . }}</body>
</html>`

// Configure the application from environement
func Configure() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})

	queue, err := dispatch.NewAMQPQueue(
		os.Getenv("QUEUE_NAME"),
		os.Getenv("QUEUE_HOST"))
	if err != nil {
		log.Fatal(err)
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	log.WithField("path", os.Getenv("TEMPLATES")).Info("Loading templates")
	tmpls, err := text.ParseGlob(os.Getenv("TEMPLATES"))
	if err != nil {
		log.Fatal(err)
	}

	MainMailer = &Mailer{
		SMTP: SMTPConf{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     smtpPort,
			User:     os.Getenv("SMTP_USER"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
		Queue:     queue,
		Sender:    os.Getenv("SENDER"),
		Templates: tmpls,
	}

	if os.Getenv("HTML_BODY") == "" {
		log.Info("No html body set, using default body")
		if MainMailer.Body, err = html.New("body").Parse(defaultBody); err != nil {
			log.Fatal(err)
		}
	} else if content, err := ioutil.ReadFile(os.Getenv("HTML_BODY")); err != nil {
		log.Fatal(err)
	} else if MainMailer.Body, err = html.New("body").Parse(string(content[:])); err != nil {
		log.Fatal(err)
	} else {
		log.WithField("path", os.Getenv("HTML_BODY")).Info("HTML body provided")
	}
}

func main() {
	Configure()

	log.WithField("queue", os.Getenv("QUEUE_NAME")).Info("qmail started")

	if err := MainMailer.Queue.ListenBytes(MainMailer.Listenner); err != nil {
		log.Fatal("Program failed, exiting.")
	}
}
