package main_test

import (
	"os"

	. "github.com/hyperboloide/qmail"
	"github.com/hyperboloide/qmail/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {

	var cli *client.Mailer

	It("should configure", func() {
		Configure()
		Ω(MainMailer).ToNot(BeNil())
	})

	It("should purge the queue", func() {
		Ω(MainMailer.Queue.Purge()).To(BeNil())
	})

	It("should create a client", func() {
		c, err := client.New(os.Getenv("QUEUE_NAME"), os.Getenv("QUEUE_HOST"))
		Ω(err).To(BeNil())
		cli = c
	})

	It("should send some mails and crash", func() {
		errChan := make(chan error)
		go func() {
			errChan <- MainMailer.Queue.ListenBytes(MainMailer.Listenner)
		}()

		mOk := client.Mail{
			Dests:    []string{"opensource@hyperboloide.com"},
			Subject:  "test",
			Template: "example",
			Data:     map[string]string{"User": "test user"},
			Files:    []string{"./example.md"},
		}
		Ω(cli.Send(mOk)).To(BeNil())

		mFail := client.Mail{
			Dests:    []string{"opensource@hyperboloide.com"},
			Subject:  "test",
			Template: "not_found",
			Data:     map[string]string{"User": "test user"},
		}
		Ω(cli.Send(mFail)).To(BeNil())

		Ω(<-errChan).ToNot(BeNil())
	})

})
