FROM golang:latest
ENV GOPATH /go
COPY main.go /go
WORKDIR /go
EXPOSE 8123
RUN ["go", "get", "github.com/gorilla/mux"]
RUN go build -o main .
CMD ["/go/main"]
