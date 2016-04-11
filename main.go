package main

import (
	"os"
	"strconv"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/hyperboloide/dispatch"
)

var (
	// MainMailer connect to the queue and handle messages
	MainMailer *Mailer
)

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

	tmpls, err := template.ParseGlob(os.Getenv("TEMPLATES"))
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
}

func main() {
	Configure()

	if err := MainMailer.Queue.ListenBytes(MainMailer.Listenner); err != nil {
		log.Fatal(err)
	}
}
