FROM golang:1.18-alpine as BUILD

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV PASSWORD=pg_password

ENV JWT_SECRET=jwt_secret

RUN go build -o exec .

FROM alpine:latest

WORKDIR /newapi

COPY --from=BUILD /api/exec /newapi/

RUN apk update && apk add --no-cache nats-server

EXPOSE 4222 8080

CMD ["sh", "-c", "nats-server & sleep 3 && ./exec"]

