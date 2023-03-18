FROM golang:1.19.3-alpine

RUN apk add -U build-base git curl\
    make openssh-client

WORKDIR /go/src/app

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

ADD . .

ENTRYPOINT ["air"]
CMD  ["-c", ".air.toml"]