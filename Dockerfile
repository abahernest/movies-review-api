# syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

ARG app_env=prod

ENV APP_ENV=$app_env

RUN go build -o ./movies-review-api ./cmd/main.go

## Deploy
FROM alpine:3.11.3

WORKDIR /

COPY --from=builder /app/movies-review-api .

EXPOSE 6001

#USER nonroot:nonroot

CMD [ "./movies-review-api" ]