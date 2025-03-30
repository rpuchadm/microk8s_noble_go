# Build stage (usando imagen ligera de Go)
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server

# Runtime stage (imagen ultraligera)
FROM alpine:3.21
WORKDIR /app

# Copia el binario desde el builder
COPY --from=builder /app/server /app/server

# Puerto expuesto y variable de entorno (sobrescribible en k8s)
EXPOSE 8080
ENV NOMBRE="ValorPorDefecto"

# Usuario no-root para seguridad
RUN adduser -D appuser
USER appuser

CMD ["/app/server"]