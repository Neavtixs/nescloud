# API Specification

## Project

Cloud Storage Web Application

---

# Base URL

Development

```
http://localhost:8080/api
```

Production

```
https://api.example.com/api
```

---

# Authentication

Menggunakan

```
JWT Access Token
Refresh Token (HttpOnly Cookie)
```

Access Token dikirim melalui

```
Authorization: Bearer <token>
```

Refresh Token dikirim melalui HttpOnly Cookie

```
Set-Cookie: refresh_token=<uuid>; HttpOnly; Secure; SameSite=Lax; Path=/api/auth; MaxAge=604800
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
  "message": "Unauthorized"
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

## Register — ✅ Implemented

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

Response `201 Created`

```json
{
  "message": "register user success",
  "data": {
    "access_token": "<jwt>"
  }
}
```

Set-Cookie: `refresh_token=<uuid>; HttpOnly; Secure; SameSite=Lax; Path=/api/auth; MaxAge=604800`: `refresh_token=<uuid>; HttpOnly; Secure; SameSite=Lax; Path=/api/auth; MaxAge=604800`

---

## Login — ❌ Not implemented

POST

```
/auth/login
```

Request

```json
{
  "email": "lawrient@mail.com",
  "password": "password"
}
```

Response `200 OK`

```json
{
  "message": "login user success",
  "data": {
    "access_token": "<jwt>",
    "expires_in": 900
  }
}
```

Set-Cookie: `refresh_token=<uuid>; HttpOnly; Secure; SameSite=Lax; Path=/api/auth; MaxAge=604800`

---

## Refresh Token — ❌ Not implemented

POST

```
/auth/refresh
```

Request

Tidak ada request body. `refresh_token` dibaca dari cookie.

Response `200 OK`

```json
{
  "message": "token refreshed",
  "data": {
    "access_token": "<jwt>",
    "expires_in": 900
  }
}
```

Set-Cookie: `refresh_token=<new-uuid>; HttpOnly; Secure; SameSite=Lax; Path=/api/auth; MaxAge=604800` (rotate token)

---

## Logout — ❌ Not implemented

POST

```
/auth/logout
```

Response `200 OK`

```json
{
  "message": "logout success"
}
```

Menghapus Refresh Token dari Redis dan clear cookie: `Set-Cookie: refresh_token=; MaxAge=0; HttpOnly; Secure; SameSite=Lax; Path=/api/auth`

---

## Current User — ❌ Not implemented

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

# Folder — ❌ All not implemented

## Create Folder

POST

```
/folders
```

Request

```json
{
  "parent_folder_id": "",
  "name":[118;1:3u "Documents"
}
```

Response `201 Created`

```json
{
  "message": "folder created",
  "data": {
    "id": "uuid"
  }
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

Response `200 OK`

```json
{
  "message": "success",
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 0,
    "total_pages": 0
  }
}
```

---

## Rename Folder

PATCH

```
/folders/{id}
```

Request

```json
{
  "name": "New Name"
}
```

Response `200 OK`

```json
{
  "message": "folder renamed",
  "data": {
    "id": "uuid",
    "name": "New Name"
  }
}
```

---

## Delete Folder

DELETE

```
/folders/{id}
```

Response `200 OK`

```json
{
  "message": "folder moved to trash"
}
```

Masuk Trash.

---

## Restore Folder

POST

```
/trash/folders/{id}/restore
```

Response `200 OK`

```json
{
  "message": "folder restored"
}
```

---

# File Upload — ❌ All not implemented

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
  "folder_id": "uuid",
  "file_name": "photo.png",
  "mime_type": "image/png",
  "size": 102400
}
```

Response `201 Created`

```json
{
  "message": "upload initialized",
  "data": {
    "file_id": "uuid",
    "upload_url": "https://s3...",
    "expired_at": "2026-06-25T10:15:00Z"
  }
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

Response `200 OK`

```json
{
  "message": "upload completed",
  "data": {
    "file_id": "uuid"
  }
}
```

Action

- Validasi object exists
- Update status uploaded
- Tambah quota

---

# File — ❌ All not implemented

## Get File List

GET

```
/files
```

Query

```
?folder_id=uuid
&page=1
&limit=20
&search=holiday
```

Response `200 OK`

```json
{
  "message": "success",
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 0,
    "total_pages": 0
  }
}
```

---

## Get File Detail

GET

```
/files/{id}
```

Response `200 OK`

```json
{
  "message": "success",
  "data": {
    "id": "uuid",
    "original_name": "photo.png",
    "mime_type": "image/png",
    "size": 102400,
    "folder_id": "uuid",
    "created_at": "2026-06-25T10:00:00Z"
  }
}
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

Response `200 OK`

```json
{
  "message": "file renamed",
  "data": {
    "id": "uuid",
    "original_name": "holiday.png"
  }
}
```

Storage Key tidak berubah.

---

## Delete File

DELETE

```
/files/{id}
```

Response `200 OK`

```json
{
  "message": "file moved to trash"
}
```

Masuk Trash.

---

## Download File

GET

```
/files/{id}/download
```

Response `200 OK`

Backend membuat Presigned Download URL.

```json
{
  "message": "success",
  "data": {
    "download_url": "https://s3..."
  }
}
```

---

# Public Link — ❌ All not implemented

## Generate

POST

```
/files/{id}/public-link
```

Response `201 Created`

```json
{
  "message": "public link created",
  "data": {
    "url": "https://..."
  }
}
```

---

## Revoke

DELETE

```
/files/{id}/public-link
```

Response `200 OK`

```json
{
  "message": "public link revoked"
}
```

---

## Access Public File

GET

```
/public/{token}
```

Authentication

Tidak diperlukan.

---

# Trash — ❌ All not implemented

## Get Trash

GET

```
/trash
```

Response `200 OK`

```json
{
  "message": "success",
  "data": {
    "files": [],
    "folders": []
  }
}
```

---

## Restore File

POST

```
/trash/files/{id}/restore
```

Response `200 OK`

```json
{
  "message": "file restored"
}
```

---

## Restore Folder

POST

```
/trash/folders/{id}/restore
```

Response `200 OK`

```json
{
  "message": "folder restored"
}
```

---

## Permanent Delete File

DELETE

```
/trash/files/{id}
```

Response `200 OK`

```json
{
  "message": "file permanently deleted"
}
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

Response `200 OK`

```json
{
  "message": "folder permanently deleted"
}
```

Menghapus seluruh isi folder.

---

## Empty Trash

DELETE

```
/trash
```

Response `200 OK`

```json
{
  "message": "trash emptied"
}
```

Menghapus seluruh isi Trash.

---

# Search — ❌ Not implemented

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

# Pagination — ❌ Not implemented

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

# Scheduler — ❌ Not implemented

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
