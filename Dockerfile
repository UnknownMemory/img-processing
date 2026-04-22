FROM golang:1.26-alpine AS build

RUN apk add --no-cache build-base pkgconfig git vips-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /build/api ./cmd/api
RUN CGO_ENABLED=1 GOOS=linux go build -o /build/worker ./cmd/worker

FROM alpine:3.23 AS api

RUN apk add --no-cache ca-certificates vips

WORKDIR /app

COPY --from=build /build/api /app/api
COPY .env.prod /app/.env

EXPOSE 8080

ENTRYPOINT ["/app/api"]

FROM alpine:3.23 AS worker

RUN apk add --no-cache ca-certificates vips

WORKDIR /app

COPY --from=build /build/worker /app/worker
COPY .env.prod /app/.env

ENTRYPOINT ["/app/worker"]