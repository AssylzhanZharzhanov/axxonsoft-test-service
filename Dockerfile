FROM golang:1.17-alpine as base
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY deployment .

FROM base AS builder
ARG APP_VERSION
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.AppVersion=${APP_VERSION}" -o ./.bin/app ./cmd/main.go

FROM alpine:3
WORKDIR /app/
COPY --from=builder ./app/.bin/app .
COPY --from=builder ./app/config ./config