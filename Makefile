NAME = qmail
DOCKERID = hyperboloide

all: container push

container:
	GOOS=linux GOARCH=amd64 go build -a
	docker build -t $(DOCKERID)/$(NAME) .

push:
	docker push     $(DOCKERID)/$(NAME)

.PHONY: all container push
