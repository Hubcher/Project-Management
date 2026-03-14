# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

WORKDIR /src

COPY auth-service ./auth-service
COPY contracts ./contracts

RUN go work init ./auth-service ./contracts

WORKDIR /src/auth-service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/gateway ./cmd

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/auth-service /app/auth-service

ENTRYPOINT ["/app/auth-service"]