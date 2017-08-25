FROM golang:latest
MAINTAINER Logo_fox Edward <logo_fox@163.com>

WORKDIR $GOPATH/src/logo_fox/autocall
ADD . $GOPATH/src/logo_fox/autocall

RUN go build .

EXPOSE 12345

ENTRYPOINT ["./autocall"]
