# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

WORKDIR /src

COPY project-service ./project-service
COPY contracts ./contracts

RUN go work init ./project-service ./contracts

WORKDIR /src/project-service

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/project-service ./cmd

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/project-service /app/project-service


ENTRYPOINT ["/app/project-service"]