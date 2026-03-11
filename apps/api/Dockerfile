# syntax=docker/dockerfile:1.7

FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY apps/api/go.mod apps/api/go.sum* ./apps/api/
WORKDIR /src/apps/api
RUN go mod download

COPY apps/api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

FROM alpine:3.21
RUN adduser -D -H appuser
WORKDIR /app
COPY --from=builder /out/api /app/api

EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/api"]
