FROM golang:latest
MAINTAINER Logo_fox Edward <logo_fox@163.com>

WORKDIR $GOPATH/src/hellodocker
ADD . $GOPATH/src/hellodocker

RUN go build .

EXPOSE 12345

ENTRYPOINT ["./hellodocker"]
