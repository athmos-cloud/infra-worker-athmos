FROM golang:1.20-bullseye

RUN apt-get update && apt-get install -y gcc openssl git curl\
    make openssh-client

WORKDIR /go/src/app

ADD . .

ENTRYPOINT ["go", "run", "main.go"]