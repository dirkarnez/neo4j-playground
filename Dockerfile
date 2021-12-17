FROM golang:1.17.0-alpine

COPY . /go/src/github.com/dirkarnez/golang-hello-world
WORKDIR /go/src/github.com/dirkarnez/golang-hello-world

RUN go build -o app

ENTRYPOINT [ "./app" ]

EXPOSE 5000
