# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

WORKDIR /src

COPY payment-calendar-service ./payment-calendar-service
COPY contracts ./contracts

RUN go work init ./payment-calendar-service ./contracts

WORKDIR /src/payment-calendar-service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/payment-calendar-service ./cmd

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/payment-calendar-service /app/payment-calendar-service

ENTRYPOINT ["/app/payment-calendar-service"]
