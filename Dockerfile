FROM golang:1.16.2-alpine3.13
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go get github.com/gorilla/websocket
WORKDIR /app/src
RUN go build -o exec
CMD ["/app/src/exec"]