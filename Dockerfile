FROM golang:1.23-alpine3.22 AS builder

ENV GO111MODULE=on
WORKDIR /app

RUN apk add --no-cache git tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app .

FROM alpine:3.22

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/app .
COPY .env .

EXPOSE 8080

CMD ["./app"]
