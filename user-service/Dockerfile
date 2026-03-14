# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

WORKDIR /src

COPY user-service ./user-service
COPY contracts ./contracts

RUN go work init ./user-service ./contracts

WORKDIR /src/user-service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/user-service ./cmd

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/user-service /app/user-service


ENTRYPOINT ["/app/user-service"]