FROM golang:1.24
LABEL maintainer="Étienne Michon <etienne@scalingo.com>"

RUN go install github.com/cespare/reflex@latest

WORKDIR $GOPATH/src/github.com/Scalingo/go-plugins-helpers

CMD $GOPATH/bin/go-plugins-helpers
