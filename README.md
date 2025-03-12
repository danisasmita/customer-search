# Customer Search API

API untuk mencari data pelanggan dengan autentikasi JWT.

## Fitur
- Register dan Login dengan JWT.
- Pencarian pelanggan berdasarkan nama (dilindungi oleh JWT).
- Terintegrasi dengan PostgreSQL menggunakan GORM.
- Docker support untuk development dan staging.
- Dokumentasi API otomatis menggunakan Swagger UI.

## Cara Menjalankan
1. Clone repository ini.
2. Jalankan Docker Compose:
   ```bash
   docker-compose up --build