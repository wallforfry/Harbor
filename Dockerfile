FROM golang:1.8

WORKDIR /go/src/wallforfry.fr/harbor
COPY . .

RUN go get -d -v wallforfry.fr/harbor
RUN go install -v ./...

CMD ["/go/bin/harbor"]
