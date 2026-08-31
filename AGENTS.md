# api-wifi

WiFi billing + employee attendance API. Go + Gin + GORM + PostgreSQL. Deploy via Docker.

## PENTING: jangan edit folder ini kalau ini server live

- Ngoding di `~/dev/api-wifi` (clone dari GitHub).
- `/www/wwwroot/api-wifi` = folder LIVE (checkout runner). Deploy via GitHub Actions self-hosted runner.
- Alur: edit di `~/dev` → commit → push → runner `docker compose build` + `up -d --force-recreate`.

## Stack

- Go 1.25, Gin, GORM, PostgreSQL, godotenv
- Container: `wifi-api-prod` (port 8080) dari branch `main`, `wifi-api-dev` (port 8081) dari branch `development`
- `network_mode: host`

## Branch → environment

- `main` → prod, `.env.prod`, port 8080
- `development` → dev, `.env.dev`, port 8081

## Perintah

```bash
go run main.go
go build -o wifi-api .
docker compose build wifi-api-dev
docker compose up -d --force-recreate wifi-api-dev
```

## Struktur

```
main.go        entry
config/        koneksi DB + env
controllers/   handler HTTP
services/      logika bisnis
models/        GORM model
middlewares/   auth/RBAC
routes/        definisi route
helpers/ utils/
```

## Tak masuk git

`.env`, `.env.prod`, `.env.dev` (di-`.gitignore`). File env tersimpan di folder live server, di-copy saat deploy.
