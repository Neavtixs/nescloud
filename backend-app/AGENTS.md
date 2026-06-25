# Backend Rules

Panduan ini dibuat dari struktur dan pola kode yang ada di repo. Gunakan sebagai aturan utama saat membaca, menulis, atau mengubah kode di proyek ini.

## Ringkasan Proyek

- Backend API ditulis dengan Go, Gin, PostgreSQL, Redis, JWT cookie auth, OAuth Google, dan Logrus.
- Entry point ada di `main.go`.
- Konfigurasi aplikasi ada di `configs/`.
- Wiring dependency manual ada di `internal/apps.go`.
- Routing terpusat di `internal/route/route.go`.
- Domain model dan database access ada di `internal/apps/domain/`.
- Feature utama ada di `internal/apps/feature/<feature>/` dengan pola `handler.go` dan `service.go`.
- DTO request, response, input service, dan result service ada di `internal/dto/`.
- Error domain reusable ada di `internal/errs/custom.go`.
- Helper umum ada di `internal/helper/`.
- Database migration ada di `db/migrations/`.

## Config DB Rules

- Konfigurasi koneksi database ada di `configs/db.go` lewat fungsi `GetConnection()`.
- `GetConnection()` harus tetap mengembalikan `*sql.DB` dan memakai driver PostgreSQL `lib/pq`.
- Ambil konfigurasi database dari environment variable yang sudah ada: `DB_USERNAME`, `DB_PASSWORD`, `DB_DATABASE`, `DB_HOST`, dan `DB_PORT`.
- Jangan hardcode username, password, host, port, database name, atau DSN lengkap di kode.
- Build DSN memakai `net/url.URL` dan `url.UserPassword` agar credential dengan karakter khusus tetap aman.
- Parameter PostgreSQL seperti `sslmode` ditambahkan lewat query parameter URL, bukan string concatenation manual.
- Jika mengubah pool connection, pertahankan konfigurasi eksplisit untuk `SetMaxIdleConns`, `SetMaxOpenConns`, dan `SetConnMaxLifetime`.
- Jangan membuka koneksi database baru di feature, handler, service, atau repository. Gunakan dependency `*sql.DB` dari wiring `internal/apps.go`.
- Jangan menjalankan migration atau query aplikasi dari `configs/db.go`; file ini hanya untuk membuat dan mengatur koneksi database.
- Jangan log DSN penuh atau password database.

## Struktur Feature

Ikuti pola ini saat menambah atau mengubah feature:

```text
internal/apps/feature/<feature>/
  handler.go
  service.go
```

- `handler.go` hanya menangani HTTP concern: bind request, ambil path/query/cookie/context value, validasi DTO, panggil service, mapping result ke response JSON.
- `service.go` menangani business logic, transaksi database, pemanggilan repository, Redis, UUID generation, dan mapping entity ke result service.
- Jangan letakkan query SQL di handler atau service. Query SQL harus ada di repository.
- Jangan akses `gin.Context` di repository.
- Jangan membuat route langsung di feature. Tambahkan route di `internal/route/route.go`.
- Tambahkan dependency feature baru lewat `internal/apps.go`, bukan lewat global variable.

## Dependency Injection

- Constructor handler memakai pola:

```go
func NewHandler(service *Service, validate *validator.Validate, log *logrus.Logger) *Handler
```

- Constructor service menerima dependency eksplisit seperti `*sql.DB`, `*redis.Client`, dan repository yang dibutuhkan.
- Buat repository lewat constructor `NewXRepo()` di `internal/apps.go`.
- Hindari singleton/global state baru. Dependency harus terlihat dari struct dan constructor.

## Handler Rules

- Nama handler HTTP berakhiran `Handler`, contoh `LoginHandler`, `AdminCreateProductHandler`.
- Ambil user login dari `c.Get("user_id")` untuk route authenticated.
- Validasi `user_id` wajib memastikan value ada, bertipe `string`, dan tidak kosong.
- Gunakan `c.ShouldBindJSON(&req)` untuk JSON request.
- Untuk multipart, parse field manual seperti pola product handler.
- Setelah bind, jalankan `h.Validate.Struct(req)` untuk DTO yang punya validation tag.
- Response harus memakai `dto.ResponseWeb[T]`.
- Gunakan status code sesuai kondisi:
  - `200 OK` untuk read/update/toggle sukses.
  - `201 Created`[118;1:3u untuk create sukses.
  - `400 Bad Request` untuk invalid request, validation, atau input domain yang salah.
  - `401 Unauthorized` untuk access token atau refresh token invalid.
  - `403 Forbidden` untuk user valid tetapi tidak boleh mengakses resource.
  - `404 Not Found` untuk data publik/detail yang tidak ada.
  - `500 Internal Server Error` untuk error internal.
- Jangan bocorkan error internal ke client. Pakai `errs.ErrInternal.Error()`.
- Untuk validation response, pakai:

```go
c.JSON(http.StatusBadRequest, dto.ResponseWeb[map[string]string]{
	Message: "validation failed",
	Data:    helper.ValidationMsg(err),
})
```

- Pakai `errors.Is(err, errs.SomeError)` untuk membandingkan error domain.
- Hindari typo pada response message baru. Jangan meniru typo lama seperti `"validation vailed"`.

## Service Rules

- Service method untuk operation kompleks menerima pointer DTO input:

```go
func (s *Service) SomeAction(input *dto.InputSomeAction) (*dto.ResultSomeAction, error)
```

- Context berasal dari `input.Ctx` atau parameter `ctx context.Context`.
- Service mengontrol transaksi database:

```go
tx, err := s.DB.BeginTx(input.Ctx, nil)
if err != nil {
	return nil, err
}
defer tx.Rollback()

// repository calls

if err := tx.Commit(); err != nil {
	return nil, err
}
```

- Repository harus dipanggil dengan `*sql.Tx`, bukan `*sql.DB`.
- Commit transaksi sebelum return sukses.
- Jangan return `nil, nil` saat commit gagal. Return error.
- Generate UUID dengan `uuid.NewString()`.
- Role check admin memakai `strings.EqualFold(strings.TrimSpace(role), "admin")`.
- Untuk authorization admin, gunakan error domain yang jelas bila menambah kode baru. Pola lama sering mengembalikan `ErrDataNotFound` untuk non-admin; jangan perluas pola ini bila bisa memakai error yang lebih eksplisit.
- Redis key harus diberi namespace, contoh `refresh_token:<token>`.
- Operasi cleanup best-effort boleh mengabaikan error dengan `_ = ...` hanya jika tidak memengaruhi hasil utama dan alasannya jelas.

## Repository Rules

- Repository ada di `internal/apps/domain/repository/`.
- Nama struct repository: `XRepo`.
- Constructor: `NewXRepo() *XRepo`.
- Method repository menerima `*sql.Tx` dan `context.Context`.
- Query SQL ditulis sebagai raw SQL multiline string.
- Gunakan parameter placeholder PostgreSQL `$1`, `$2`, dan seterusnya. Jangan interpolasi string user input ke SQL.
- Mapping `sql.ErrNoRows` harus menjadi `errs.ErrDataNotFound`.
- Mapping PostgreSQL error khusus memakai `*pq.Error` bila perlu, contoh:
  - `23505` untuk unique violation.
  - `22P02` untuk invalid UUID/input syntax.
- Setelah `ExecContext`, cek `RowsAffected()` untuk create/update/delete yang harus mengubah data.
- Method repository harus mengembalikan entity atau data query tanpa HTTP concern.

## Entity Rules

- Entity ada di `internal/apps/domain/entity/`.
- Entity merepresentasikan table database.
- Field entity memakai PascalCase tanpa JSON tag, kecuali ada kebutuhan khusus.
- Sertakan `CreatedAt time.Time` dan `UpdatedAt time.Time` bila table punya kolom tersebut.
- Entity tidak berisi validation tag dan tidak berisi business logic HTTP.

## DTO Rules

- Request/response HTTP ada di `internal/dto/data_handler.go`.
- Input/result service ada di `internal/dto/data_service.go`.
- Response wrapper tetap:

```go
type ResponseWeb[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
```

- Request DTO berakhiran `Req`.
- Response DTO berakhiran `Res`.
- Service input berawalan `Input`.
- Service result berawalan `Result`.
- JSON tag memakai snake_case.
- Validation tag ditulis di request DTO, bukan entity.
- Untuk optional boolean request, pakai `*bool` agar bisa membedakan `false` dari tidak dikirim.
- Jangan membuat DTO baru jika DTO existing masih tepat secara semantik.
- Saat menambah DTO baru, perbaiki ejaan pada nama baru. Jangan meniru typo existing seperti `InputSawriaWebhook`, `InpuShowDetail`, atau `ResultLeaderboradSaweria`.

## Error Handling

- Error domain reusable harus ditambahkan di `internal/errs/custom.go`.
- Gunakan `errors.Is` untuk mengecek error dari service/repository.
- Handler bertugas menerjemahkan error domain menjadi HTTP status dan response message.
- Service dan repository tidak boleh membuat response JSON.
- Jangan mengembalikan pesan error database langsung ke client.
- Bila menambah error baru, gunakan pesan singkat, lowercase, dan jelas.

## Logging

- Handler yang perlu logging memakai:

```go
log := helper.NewLog(h.Log, c)
log.WithField("layer", "<layer_name>").Info("<message>")
```

- `helper.NewLog` sudah menambahkan `request_id`, path, dan method.
- Jangan log password, token, cookie value, API key, OAuth code, atau secret.
- Hindari `log.Println` pada kode baru. Pakai Logrus via dependency handler atau logger yang sudah ada.

## Auth dan Security

- Authenticated route memakai `middleware.Authorization()`.
- Middleware membaca `access_token` dari cookie.
- Refresh token disimpan di Redis dengan TTL.
- Cookie auth harus tetap `HttpOnly`.
- Jangan expose JWT atau refresh token di JSON response.
- Jangan hardcode secret, credential, atau URL sensitif. Ambil dari environment variable.
- Jangan commit `.env` baru atau perubahan credential.

## Route Rules

- Semua route ditambahkan di `internal/route/route.go`.
- Group utama memakai `/api`.
- Public route berada di group `public`.
- Authenticated user route berada di group `user := api.Group("", middleware.Authorization())`.
- Admin route berada di group `/admin`.
- Gunakan path yang sudah konsisten:
  - singular resource pada route existing seperti `/product`, `/experience`.
  - action path existing seperti `/experience/activated/:id`.
- Jangan ubah route existing tanpa kebutuhan eksplisit karena itu mengubah kontrak API.

## Migration Rules

- Migration ada di `db/migrations/` dengan format:

```text
00000N_description.up.sql
00000N_description.down.sql
```

- Gunakan PostgreSQL syntax.
- Table name existing memakai singular snake_case, contoh `users`, `user_profile`, `roblox_experience`, `donation_history`, `product`.
- Kolom memakai snake_case.
- UUID primary key ditulis `UUID PRIMARY KEY NOT NULL`.
- Timestamp standar:

```sql
created_at TIMESTAMPTZ DEFAULT NOW(),
updated_at TIMESTAMPTZ DEFAULT NOW()
```

- Foreign key diberi nama constraint `fk_<child>_<parent>` dan sebaiknya punya `ON DELETE CASCADE` bila child harus ikut hilang.
- Down migration harus drop table dengan `DROP TABLE IF EXISTS <table> CASCADE;`.
- Untuk migration baru, pertahankan urutan nomor dan jangan mengubah migration lama kecuali sedang memperbaiki baseline secara eksplisit.

## Test Rules

- Test repository ada di `internal/apps/domain/repository/*_test.go`.
- Test database dibuat lewat `SetupTestDB(t)`.
- Gunakan `testify/assert`.
- Setiap test yang membuka database harus `defer db.Close()`.
- Setiap test repository sebaiknya membuka transaction dan `defer tx.Rollback()` agar data test tidak menetap.
- Tambahkan test untuk:
  - create success.
  - duplicate/constraint error.
  - get success.
  - not found.
  - invalid UUID/type jika query menerima UUID.
- Untuk perubahan service/handler yang berisiko, tambahkan test sesuai boundary yang disentuh.

## Formatting dan Gaya Go

- Jalankan `gofmt` setelah mengubah file Go.
- Import harus dibiarkan diatur oleh `go fmt`/`goimports` bila tersedia.
- Gunakan nama yang jelas dan idiomatis.
- Hindari singkatan berlebihan pada nama baru.
- Gunakan `http.Status...` daripada angka status literal di kode baru.
- Gunakan `any` untuk response generic kosong bila tidak ada data spesifik.
- Jangan menambahkan komentar yang menjelaskan hal obvious. Komentar cukup untuk konteks bisnis atau bagian yang rawan salah.
- Jaga perubahan tetap kecil dan sesuai feature yang disentuh.

## API Response Style

- Message response memakai bahasa Inggris sederhana.
- Format sukses:
  - `"register user success"`
  - `"login user success"`
  - `"success get product categories"`
  - `"product category created successfully"`
- Format error umum:
  - `"validation failed"`
  - `"invalid request format"`
  - `"invalid access token"`
  - `"internal server error"`
- Jangan mencampur struktur response baru selain `ResponseWeb[T]`.
- Bila endpoint mengembalikan list, `Data` harus berupa slice kosong `[]` saat tidak ada data, bukan `null`, jika memungkinkan.

## Hal yang Harus Dihindari

- Jangan menambah query SQL di handler/service.
- Jangan membuat global DB/Redis/logger baru.
- Jangan menggunakan string interpolation untuk SQL.
- Jangan mengubah route, JSON field, atau response shape existing tanpa permintaan jelas.
- Jangan memperbanyak typo nama type/method dari kode lama.
- Jangan log data rahasia.
- Jangan melakukan refactor besar di luar scope perubahan.

## Checklist Saat Menambah Feature

1. Tambahkan migration bila perlu table/kolom baru.
2. Tambahkan entity sesuai table.
3. Tambahkan repository dan test repository.
4. Tambahkan DTO request/response/input/result.
5. Tambahkan service dengan transaksi dan business logic.
6. Tambahkan handler dengan bind, validasi, auth context, dan response mapping.
7. Wire repository, service, dan handler di `internal/apps.go`.
8. Tambahkan route di `inter[27;1:3unal/route/route.go`.
9. Jalankan `gofmt`.
10. Jalankan test yang relevan.
