# API Specification

## Project

Cloud Storage Web Application

---

# Base URL

Development

```
http://localhost:8080/api/v1
```

Production

```
https://api.example.com/api/v1
```

---

# Authentication

Menggunakan

```
JWT Access Token
Refresh Token
```

Access Token dikirim melalui

```
Authorization: Bearer <token>
```

---

# Response Format

Success

```json
{
  "message": "Success",
  "data": {}
}
```

Error

```json
{
  "message": "Unauthorized",
  "errors": {}
}
```

---

# HTTP Status

| Code | Description           |
| ---- | --------------------- |
| 200  | OK                    |
| 201  | Created               |
| 204  | No Content            |
| 400  | Bad Request           |
| 401  | Unauthorized          |
| 403  | Forbidden             |
| 404  | Not Found             |
| 409  | Conflict              |
| 413  | Payload Too Large     |
| 422  | Validation Error      |
| 429  | Too Many Requests     |
| 500  | Internal Server Error |

---

# Authentication

## Register

POST

```
/auth/register
```

Request

```json
{
  "name": "Lawrient",
  "email": "lawrient@mail.com",
  "password": "password"
}
```

Response

```json
{
  "message": "register user success"
}
```

---

## Login

POST

```
/auth/login
```

Response

```json
{
  "access_token": "",
  "refresh_token": "",
  "expires_in": 900
}
```

---

## Refresh Token

POST

```
/auth/refresh
```

Request

```json
{
  "refresh_token": "..."
}
```

---

## Logout

POST

```
/auth/logout
```

Menghapus Refresh Token.

---

## Current User

GET

```
/auth/me
```

Authentication

Required.

Response

```json
{
  "message": "Success",
  "data": {
    "id": "",
    "name": "Lawrient",
    "email": "lawrient@mail.com",
    "created_at": "2026-06-25T10:00:00Z"
  }
}
```

---

# Folder

## Create Folder

POST

```
/folders
```

Request

```json
{
  "parent_folder_id": "",
  "name": "Documents"
}
```

---

## Get Folder List

GET

```
/folders
```

Query

```
parent_folder_id=
```

---

## Rename Folder

PATCH

```
/folders/{id}
```

---

## Delete Folder

DELETE

```
/folders/{id}
```

Masuk Trash.

---

## Restore Folder

POST

```
/trash/folders/{id}/restore
```

---

# File Upload

Upload menggunakan Presigned URL.

Flow

```
Init Upload

↓

Upload Object Storage

↓

Complete Upload
```

---

## Init Upload

POST

```
/files/init-upload
```

Request

```json
{
  "folder_id": "",
  "file_name": "photo.png",
  "mime_type": "image/png",
  "size": 102400
}
```

Response

```json
{
  "file_id": "",
  "upload_url": "",
  "expired_at": "..."
}
```

Validation

- Max 100 MB
- Quota tersedia
- Nama file valid

---

## Complete Upload

POST

```
/files/complete
```

Request

```json
{
  "file_id": ""
}
```

Action

- Validasi object exists
- Update status uploaded
- Tambah quota

---

# File

## Get File List

GET

```
/files
```

Query

```
folder_id
page
limit
search
```

---

## Get File Detail

GET

```
/files/{id}
```

---

## Rename File

PATCH

```
/files/{id}
```

Request

```json
{
  "name": "holiday.png"
}
```

Storage Key tidak berubah.

---

## Delete File

DELETE

```
/files/{id}
```

Masuk Trash.

---

## Download File

GET

```
/files/{id}/download
```

Response

Backend membuat Presigned Download URL.

```json
{
  "download_url": "..."
}
```

---

# Public Link

## Generate

POST

```
/files/{id}/public-link
```

Response

```json
{
  "url": "https://..."
}
```

---

## Revoke

DELETE

```
/files/{id}/public-link
```

---

## [118;1:3uAccess Public File

GET

```
/public/{token}
```

Authentication

Tidak diperlukan.

---

# Trash

## Get Trash

GET

```
/trash
```

---

## Restore File

POST

```
/trash/files/{id}/restore
```

---

## Restore Folder

POST

```
/trash/folders/{id}/restore
```

---

## Permanent Delete File

DELETE

```
/trash/files/{id}
```

Action

- Delete Object Storage
- Delete Metadata
- Release Quota

---

## Permanent Delete Folder

DELETE

```
/trash/folders/{id}
```

Menghapus seluruh isi folder.

---

## Empty Trash

DELETE

```
/trash
```

Menghapus seluruh isi Trash.

---

# Search

GET

```
/files
```

Query

```
search=holiday
```

Search berdasarkan

```
original_name
```

---

# Pagination

Request

```
?page=1
&limit=20
```

Response

```json
{
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 120,
    "total_pages": 6
  }
}
```

---

# Scheduler

## Auto Trash

Interval

```
1 Hour
```

Action

Cari file

```
created_at <= NOW()-24h
```

↓

Move To Trash

---

## Auto Permanent Delete

Interval

```
1 Hour
```

Cari

```
deleted_at <= NOW()-24h
```

↓

Delete Object Storage

↓

Delete Metadata

↓

Release Quota

---

# Validation

File Size

```
Max 100 MB
```

Quota

```
1 GB
```

Allowed Owner

```
Owner Only
```

---

# Error Code

FILE_NOT_FOUND

```
404
```

INVALID_FILE_SIZE

```
413
```

QUOTA_EXCEEDED

```
409
```

UPLOAD_FAILED

```
500
```

PUBLIC_LINK_NOT_FOUND

```
404
```

UNAUTHORIZED

```
401
```

FORBIDDEN

```
403
```

---

# API Flow

Upload

```
Init Upload

↓

Presigned URL

↓

Upload Object Storage

↓

Complete Upload

↓

Done
```

Delete

```
Delete

↓

Trash

↓

Permanent Delete
```

Restore

```
Trash

↓

Restore

↓

Active
```
