FROM golang:1.16-alpine

WORKDIR /chat-manager

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /build/chat-manager.go src/main.go

EXPOSE 8080

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x /wait

CMD [ "/build/chat-manager.go" ]