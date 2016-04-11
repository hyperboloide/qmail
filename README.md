# qmail

[![Build Status](https://travis-ci.org/hyperboloide/qmail.svg?branch=master)](https://travis-ci.org/hyperboloide/qmail)
[![GoDoc](https://godoc.org/github.com/hyperboloide/qmail?status.svg)](https://godoc.org/github.com/hyperboloide/qmail)

A mailer that reads from a RabbitMQ queue and generates messages from Markdown templates.

This project has a client and a server.

## Client
To send email you need to create and instance of the `client.Mailer`.

Start by importing the client package:

```
import "github.com/hyperboloide/qmail/client"
```

Then connect to the queue and send an email:

```
mailer, err := client.New(
	os.Getenv("QUEUE_NAME"),
	os.Getenv("QUEUE_HOST"))

if err != nil {
	log.Fatal(err)
}

email := client.Mail{
	Dests:    []string{"dest@example.com"},
	Subject:  "test",
	Template: "example_template",
	Data:     map[string]string{"User": "test user"},
	Files:    []string{"./some_file.txt"},
}

if err := mailer.Send(email); err != nil {
	log.Fatal(err)
}

```

## Server

The server is available as a Docker container

```
docker pull hyperboloide/qmail
```

All configuration options are passed as environement variables:

```
docker run \
	-v ~/templates:/templates \
	-link rabbitmq:rabbitmq \
	-e TEMPLATES=/templates/*.md \
	-e QUEUE_NAME=mails \
	-e QUEUE_HOST=rabbitmq \
	-e SMTP_HOST=smtp.example.com \
	-e SMTP_PORT=465 \
	-e SMTP_USER=user@example.com \
	-e SMTP_PASSWORD=password \
	-e SENDER="Example User <user@example.com>" \
	hyperboloide/qmail
```