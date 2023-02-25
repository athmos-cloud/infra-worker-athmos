
FROM golang:1.19.3-alpine

ENV GO111MODULE=on
RUN apk add --no-cache -U bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf
RUN go install github.com/githubnemo/CompileDaemon@v1.4.0


WORKDIR /go/src/app

COPY . .

RUN go mod download

ENTRYPOINT CompileDaemon --build="go build -tags musl -o infra-worker main.go" --command=./infra-worker