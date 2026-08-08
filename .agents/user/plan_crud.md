# Admin User CRUD Implementation

Membuat fitur CRUD (Create, Read, Update, Delete) untuk entitas `User`. Endpoint ini nantinya ditujukan khusus untuk **Admin** agar bisa mengelola data `employee` dan `customer`.

## Proposed Changes

### Controllers

#### [NEW] controllers/user_controller.go
Membuat controller yang berisi logika untuk CRUD User:
1. `CreateUser`: Menambahkan user baru (hanya Name, Phone, dan Role).
2. `GetUsers`: Menampilkan seluruh data user.
3. `GetUserByID`: Menampilkan spesifik user berdasarkan ID.
4. `UpdateUser`: Memperbarui data user (seperti mengubah Role, Nomor HP, atau Upload Foto).
5. `DeleteUser`: Menghapus data user.

### Routes

#### [NEW] routes/user_routes.go
Membuat file untuk mendefinisikan rute API.
- `POST /api/admin/users` -> Create
- `GET /api/admin/users` -> Read All
- `GET /api/admin/users/:id` -> Read One
- `PUT /api/admin/users/:id` -> Update
- `DELETE /api/admin/users/:id` -> Delete

### Main App

#### [MODIFY] main.go
Mendaftarkan rute yang dibuat di `routes/user_routes.go` ke dalam framework Gin agar bisa diakses.

## Open Questions
> [!NOTE]
> **Autentikasi & Middleware**: Saat ini endpoint CRUD akan dibuat secara publik terlebih dahulu untuk memudahkan Anda melakukan testing via Postman/Insomnia. Apakah Anda setuju jika middleware untuk **memastikan hanya Admin yang bisa mengakses route ini (misal via Token/JWT)** kita kerjakan di tahap selanjutnya setelah fitur Login selesai?

## Verification Plan
1. Menyimpan dan merestart server backend.
2. Melakukan serangkaian HTTP Requests via cURL atau Postman untuk mencoba menambah, melihat, mengedit, dan menghapus data user.
