# Login API & JWT Middleware Implementation

Plan ini akan mengimplementasikan sistem Login tanpa password (hanya menggunakan Nomor Telepon) beserta pengecekan IP (lock 1 device) dan Middleware JWT.

## Proposed Changes

### 1. Dependencies & Environment
- Install package `github.com/golang-jwt/jwt/v5` untuk men-generate JWT.
- Menambahkan key rahasia `JWT_SECRET` pada file `.env`.

### 2. Controllers
#### [NEW] controllers/auth_controller.go
Membuat fungsi `Login`:
1. Menerima payload JSON berisi `phone`.
2. Mengecek apakah nomor telepon tersebut ada di database (sudah didaftarkan oleh admin). Jika tidak, tolak.
3. **Mengecek IP (Lock Device)**: 
   - Mendapatkan IP request client saat ini.
   - Jika field `ip_address` milik user di database masih `null`, maka simpan IP client ini ke database (First time login).
   - Jika field `ip_address` sudah terisi, cocokkan dengan IP client saat ini. Jika berbeda, tolak login dengan pesan device tidak dikenali.
4. Jika sukses, buat token JWT (berisi `id` dan `role`) lalu kirimkan token tersebut sebagai response.

### 3. Middlewares
#### [NEW] middlewares/auth_middleware.go
Membuat 2 buah fungsi middleware:
1. `RequireAuth`: Mengekstrak Bearer Token dari header HTTP, memvalidasinya dengan `JWT_SECRET`, lalu menyisipkan data klaim user ke dalam *Gin Context*.
2. `RequireRole(role)`: Middleware lanjutan untuk memvalidasi apakah user memiliki role tertentu (contoh: mengecek apakah dia Admin).

### 4. Routes & Main App
#### [NEW] routes/auth_routes.go
- Mendaftarkan route publik `POST /api/auth/login`.

#### [MODIFY] routes/user_routes.go
- Mengaplikasikan middleware `RequireAuth` dan `RequireRole("admin")` ke seluruh rute `/api/admin/users` agar hanya Admin dengan token valid yang bisa mengakses.

#### [MODIFY] main.go
- Mendaftarkan rute otentikasi.

## Open Questions
> [!NOTE]
> Jika ada customer/employee yang tidak sengaja ter-lock di device lama (misal: HP hilang/rusak atau ganti HP), apakah kita perlu membuatkan fitur/endpoint khusus untuk Admin agar bisa melakukan **Reset IP Device** milik user tersebut? Jika iya, akan saya tambahkan fitur tersebut di `user_controller.go` yang sebelumnya.

## Verification Plan
1. Mencoba login menggunakan nomor telepon yang belum terdaftar (harus gagal).
2. Mencoba login dengan nomor telepon terdaftar (sukses, token JWT didapat, IP tersimpan di DB).
3. Mencoba login dengan nomor terdaftar namun menggunakan IP/Device berbeda (harus diblokir).
4. Mengakses CRUD API `/api/admin/users` tanpa token JWT (harus diblokir, *Unauthorized*).
5. Mengakses CRUD API menggunakan token JWT Admin (harus sukses).
