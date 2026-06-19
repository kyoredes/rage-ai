FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY proto/ proto/
COPY auth-service/ auth-service/

WORKDIR /build/auth-service

ENV GOTOOLCHAIN=auto
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/auth-service ./cmd/server
RUN CGO_ENABLED=0 go build -o /bin/auth-migrate ./cmd/migration

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/auth-service /bin/auth-service
COPY --from=builder /bin/auth-migrate /bin/auth-migrate
COPY auth-service/docker-entrypoint.sh /docker-entrypoint.sh

RUN chmod +x /docker-entrypoint.sh

EXPOSE 50051

ENTRYPOINT ["/docker-entrypoint.sh"]
