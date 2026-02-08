# Stage 1: Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Kompilasi aplikasi menjadi binary bernama 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Production
FROM alpine:latest
WORKDIR /app
# Ambil hanya file binary hasil build
COPY --from=builder /app/main .
# Jika ada folder static atau .env, copy juga di sini
# COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]