FROM golang:1.15-alpine3.12 AS builder

COPY . /github.com/olegvolkov91/pocketer-bot/
WORKDIR /github.com/olegvolkov91/pocketer-bot/

RUN go mod download
RUN go build -o ./bin/ cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/olegvolkov91/pocketer-bot/bin/bot .
COPY --from=0 /github.com/olegvolkov91/pocketer-bot/configs configs/

EXPOSE 80

CMD ["./bot"]