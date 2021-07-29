FROM golang:1.16
WORKDIR /go/src/app
COPY . .
RUN go get -d -v github.com/gorilla/websocket
RUN go mod init
CMD ["go", "run", "main.go"]
