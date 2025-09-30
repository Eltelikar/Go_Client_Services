FROM golang:1.24.3 as builder
WORKDIR /app
COPY . .
RUN go build -o app ./cmd/client-services

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/configs/config.yaml ./configs/config.yaml
EXPOSE 8080
CMD ["./app"]