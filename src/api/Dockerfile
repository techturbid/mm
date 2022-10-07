FROM golang:1
ENV GOPATH /go
COPY main.go /go
WORKDIR /go
EXPOSE 8123
RUN go get github.com/gorilla/mux
RUN go get -u github.com/aws/aws-sdk-go/...
RUN go get github.com/Sirupsen/logrus
RUN go build -o main .
CMD ["/go/main"]
