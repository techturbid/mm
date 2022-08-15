FROM golang:1-buster
ENV GOPATH /go
COPY main.go /go
WORKDIR /go
RUN go get github.com/gorilla/mux
RUN go get -u github.com/aws/aws-sdk-go/...
RUN go get github.com/Sirupsen/logrus
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go build -o main .
CMD ["/go/main"]
