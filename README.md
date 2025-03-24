# AnalyticsAI User Service

A microservice for managing user files and logs in the AnalyticsAI platform. This service handles file uploads, storage, and management using Google Cloud Storage and MongoDB.

## Features

- File upload from local storage
- File upload from URL
- File download
- File deletion (soft delete)
- File hiding
- List user files
- Google Cloud Storage integration
- MongoDB for metadata storage

## Prerequisites

- Go 1.21 or higher
- MongoDB 6.0 or higher
- Google Cloud Platform account with Cloud Storage enabled
- Google Cloud credentials

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
PORT=8080
ENV=development

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=analyticsai

# Google Cloud Storage Configuration
GCS_BUCKET_NAME=your-bucket-name
GOOGLE_APPLICATION_CREDENTIALS=path/to/your/credentials.json
```

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/analyticsai-user-service.git
cd analyticsai-user-service
```

2. Install dependencies:

```bash
go mod download
```

3. Set up your environment variables in `.env`

4. Run the service:

```bash
go run cmd/main.go
```

## API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

All endpoints require authentication. Include the JWT token in the Authorization header:

```http
Authorization: Bearer <your_jwt_token>
```

### Common Response Format

#### Success Response

```json
{
  "status": "success",
  "data": {
    // Response data specific to the endpoint
  }
}
```

#### Error Response

```json
{
  "status": "error",
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {} // Optional additional error details
  }
}
```

### Common HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

### Endpoints

#### 1. Upload File from URL

Upload a file by providing its URL.

```http
POST /files/url
Content-Type: application/json
Authorization: Bearer <token>

{
    "name": "example.log",
    "url": "https://example.com/logs/example.log"
}
```

##### Response (201 Created)

```json
{
  "status": "success",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "user_id": 123,
    "name": "example.log",
    "original_url": "https://example.com/logs/example.log",
    "storage_key": "files/123/example.log",
    "size": 1024,
    "mime_type": "text/plain",
    "status": "active",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
  }
}
```

##### Error Response (400 Bad Request)

```json
{
  "status": "error",
  "error": {
    "code": "INVALID_URL",
    "message": "Failed to download file from URL",
    "details": {
      "status_code": 404
    }
  }
}
```

#### 2. Upload Local File

Upload a file directly from your local system.

```http
POST /files
Content-Type: multipart/form-data
Authorization: Bearer <token>

file: <file>
```

##### Response (201 Created)

```json
{
  "status": "success",
  "data": {
    "id": "507f1f77bcf86cd799439012",
    "user_id": 123,
    "name": "uploaded_file.log",
    "storage_key": "files/123/uploaded_file.log",
    "size": 2048,
    "mime_type": "text/plain",
    "status": "active",
    "created_at": "2024-03-20T10:05:00Z",
    "updated_at": "2024-03-20T10:05:00Z"
  }
}
```

##### Error Response (400 Bad Request)

```json
{
  "status": "error",
  "error": {
    "code": "INVALID_FILE",
    "message": "Invalid file format or size",
    "details": {
      "max_size": 10485760, // 10MB in bytes
      "allowed_types": ["text/plain", "application/json", "application/xml"]
    }
  }
}
```

#### 3. List User Files

Get a list of all files for the authenticated user.

```http
GET /files
Authorization: Bearer <token>
```

##### Query Parameters

| Parameter | Type    | Description                   | Default    |
| --------- | ------- | ----------------------------- | ---------- |
| page      | integer | Page number for pagination    | 1          |
| per_page  | integer | Number of items per page      | 10         |
| status    | string  | Filter by file status         | all        |
| sort_by   | string  | Sort field (created_at, name) | created_at |
| order     | string  | Sort order (asc, desc)        | desc       |

##### Response (200 OK)

```json
{
  "status": "success",
  "data": {
    "files": [
      {
        "id": "507f1f77bcf86cd799439011",
        "name": "example.log",
        "size": 1024,
        "mime_type": "text/plain",
        "status": "active",
        "created_at": "2024-03-20T10:00:00Z",
        "download_url": "http://localhost:8080/api/v1/files/507f1f77bcf86cd799439011/download"
      }
    ],
    "pagination": {
      "total": 1,
      "page": 1,
      "per_page": 10,
      "total_pages": 1
    }
  }
}
```

#### 4. Download File

Download a specific file by its ID.

```http
GET /files/{id}/download
Authorization: Bearer <token>
```

##### Response (200 OK)

- Content-Type: Based on file's MIME type
- Content-Disposition: attachment; filename=<filename>
- Body: File content

##### Error Response (404 Not Found)

```json
{
  "status": "error",
  "error": {
    "code": "FILE_NOT_FOUND",
    "message": "File not found"
  }
}
```

#### 5. Delete File

Soft delete a file (marks it as deleted but keeps the record).

```http
DELETE /files/{id}
Authorization: Bearer <token>
```

##### Response (204 No Content)

- Empty response with status code 204

##### Error Response (404 Not Found)

```json
{
  "status": "error",
  "error": {
    "code": "FILE_NOT_FOUND",
    "message": "File not found"
  }
}
```

#### 6. Hide File

Hide a file from the user's file list.

```http
PUT /files/{id}/hide
Authorization: Bearer <token>
```

##### Response (204 No Content)

- Empty response with status code 204

##### Error Response (404 Not Found)

```json
{
  "status": "error",
  "error": {
    "code": "FILE_NOT_FOUND",
    "message": "File not found"
  }
}
```

#### 7. Get File Details

Get detailed information about a specific file.

```http
GET /files/{id}
Authorization: Bearer <token>
```

##### Response (200 OK)

```json
{
  "status": "success",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "user_id": 123,
    "name": "example.log",
    "original_url": "https://example.com/logs/example.log",
    "storage_key": "files/123/example.log",
    "size": 1024,
    "mime_type": "text/plain",
    "status": "active",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z",
    "download_url": "http://localhost:8080/api/v1/files/507f1f77bcf86cd799439011/download"
  }
}
```

### File Status Types

| Status    | Description                              |
| --------- | ---------------------------------------- |
| active    | File is visible and accessible           |
| hidden    | File is hidden from the user's file list |
| deleted   | File is marked as deleted (soft delete)  |
| analyzing | File is currently being processed        |

### File Size Limits

- Maximum file size: 10MB
- Supported file types:
  - Text files (.txt, .log)
  - JSON files (.json)
  - XML files (.xml)
  - CSV files (.csv)

### Rate Limiting

- 100 requests per minute per user
- 1000 requests per hour per user
- 10MB per file upload

### Error Codes

| Code           | Description                                    |
| -------------- | ---------------------------------------------- |
| INVALID_URL    | The provided URL is invalid or inaccessible    |
| INVALID_FILE   | The uploaded file is invalid or exceeds limits |
| FILE_NOT_FOUND | The requested file does not exist              |
| UNAUTHORIZED   | Missing or invalid authentication token        |
| FORBIDDEN      | User does not have permission to access file   |
| STORAGE_ERROR  | Error occurred while accessing storage         |
| DATABASE_ERROR | Error occurred while accessing database        |

## Project Structure

```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── handlers/
│   │   └── file_handler.go
│   ├── models/
│   │   └── file.go
│   ├── repository/
│   │   └── file_repository.go
│   └── service/
│       └── file_service.go
├── pkg/
│   └── storage/
│       └── gcs.go
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o user-service cmd/main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
