package main

import (
	"os"
	"strconv"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/hyperboloide/dispatch"
)

var (
	mailer *Mailer
)

func configure() {
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

	mailer = &Mailer{
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
	configure()

	if err := mailer.Queue.ListenBytes(mailer.Listenner); err != nil {
		log.Fatal(err)
	}
}
