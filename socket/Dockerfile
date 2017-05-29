FROM golang:1.7

COPY . /go/src/github.com/ivch/aws-zombie/socket/
WORKDIR /go/src/github.com/ivch/aws-zombie/socket/
RUN go get && go build

CMD ./socket