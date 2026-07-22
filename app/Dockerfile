FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates ffmpeg tzdata

RUN adduser -D -h /app app
WORKDIR /app

COPY --from=builder /build/server .
COPY --from=builder /build/migrations ./migrations

USER app

EXPOSE 3000

ENTRYPOINT ["./server"]
