FROM golang:1.7-alpine

ADD . /go/src/github.com/tomhco/randomImage
RUN go install github.com/tomhco/randomImage

VOLUME /images

ENTRYPOINT /go/bin/randomImage -listen :8080 -directory /images

EXPOSE 8080
