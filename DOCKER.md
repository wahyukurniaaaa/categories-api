# Docker Build & Deployment Guide

## 📦 Build Docker Image

### Menggunakan Script (Recommended)

```bash
# Build dengan tag latest
./build-docker.sh

# Build dengan tag spesifik
./build-docker.sh v1.0.0
```

### Manual Build

```bash
# Build image
docker build -t ghcr.io/wahyukurniaaaa/category-api-golang:latest .

# Build dengan tag spesifik
docker build -t ghcr.io/wahyukurniaaaa/category-api-golang:v1.0.0 .
```

## 🚀 Push ke GitHub Container Registry

### Menggunakan Script (Recommended)

```bash
# Push tag latest
./push-docker.sh

# Push tag spesifik
./push-docker.sh v1.0.0
```

### Manual Push

```bash
# Login ke GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u wahyukurniaaaa --password-stdin

# Push image
docker push ghcr.io/wahyukurniaaaa/category-api-golang:latest
```

## 🔄 GitHub Actions (Otomatis)

GitHub Actions akan otomatis build dan push image saat:

1. **Push ke branch `main`** - akan membuat tag `latest`
2. **Membuat tag baru** (contoh: `v1.0.0`) - akan membuat tag sesuai versi
3. **Pull Request** - hanya build, tidak push

### Cara Membuat Release dengan Tag

```bash
# Commit perubahan
git add .
git commit -m "Release v1.0.0"

# Buat tag
git tag v1.0.0

# Push dengan tag
git push origin main --tags
```

GitHub Actions akan otomatis:
- Build Docker image
- Push ke `ghcr.io/wahyukurniaaaa/category-api-golang:v1.0.0`
- Push ke `ghcr.io/wahyukurniaaaa/category-api-golang:latest`

## 🏃 Menjalankan Container

### Dengan .env file

```bash
docker run -p 8080:8080 --env-file .env ghcr.io/wahyukurniaaaa/category-api-golang:latest
```

### Dengan environment variables

```bash
docker run -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PORT=5432 \
  -e DB_USER=your-user \
  -e DB_PASSWORD=your-password \
  -e DB_NAME=your-database \
  ghcr.io/wahyukurniaaaa/category-api-golang:latest
```

### Menggunakan Docker Compose

```bash
docker-compose up -d
```

## 📋 Environment Variables yang Diperlukan

Pastikan file `.env` atau environment variables berikut sudah diset:

- `DB_HOST` - Database host
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name

## 🔍 Melihat Image yang Tersedia

```bash
# List local images
docker images | grep category-api-golang

# Pull dari registry
docker pull ghcr.io/wahyukurniaaaa/category-api-golang:latest
```

## 🛠️ Troubleshooting

### Permission Denied saat Push

Pastikan Anda sudah login ke GitHub Container Registry:

```bash
# Buat Personal Access Token di GitHub dengan scope: write:packages
echo $GITHUB_TOKEN | docker login ghcr.io -u wahyukurniaaaa --password-stdin
```

### Build Gagal

Pastikan semua dependencies sudah ada di `go.mod` dan `go.sum`:

```bash
go mod tidy
```

### Container Tidak Bisa Connect ke Database

Pastikan environment variables sudah benar dan database bisa diakses dari container.

## 📝 Multi-Platform Build

GitHub Actions sudah dikonfigurasi untuk build multi-platform (amd64 dan arm64).

Untuk build manual multi-platform:

```bash
docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 \
  -t ghcr.io/wahyukurniaaaa/category-api-golang:latest \
  --push .
```
