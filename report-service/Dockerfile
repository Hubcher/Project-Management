# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

WORKDIR /src

COPY report-service ./report-service
COPY contracts ./contracts

RUN go work init ./report-service ./contracts

WORKDIR /src/report-service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/report-service ./cmd

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/report-service /app/report-service

ENTRYPOINT ["/app/report-service"]