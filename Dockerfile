FROM golang:1.19.3-alpine

RUN apk add -U build-base git curl\
    make openssh-client

WORKDIR /go/src/app

ADD . .

ENTRYPOINT ["go", "run", "main.go"]