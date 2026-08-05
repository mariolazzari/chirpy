# Chirpy API Documentation

Complete API reference for the Chirpy social media platform.

## Base URL

```
http://localhost:8080
```

## Authentication

Chirpy uses JWT (JSON Web Token) for authentication. Protected endpoints require a Bearer token in the Authorization header.

### Header Format

```
Authorization: Bearer <your-jwt-token>
```

### Token Types

- **Access Token**: Short-lived (1 hour) JWT for API access
- **Refresh Token**: Long-lived token to obtain new access tokens

---

## Table of Contents

1. [Health & Metrics](#health--metrics)
2. [User Management](#user-management)
3. [Authentication](#authentication-endpoints)
4. [Chirps](#chirps)
5. [Webhooks](#webhooks)
6. [Error Responses](#error-responses)

---

## Health & Metrics

### Check Health Status

```http
GET /api/healthz
```

Returns server health status.

**Response: 200 OK**
```json
{
  "status": "ok"
}
```

### Get Admin Metrics

```http
GET /admin/metrics
```

View server metrics including request count.

**Response: 200 OK**
```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited X times!</p>
  </body>
</html>
```

### Reset Metrics

```http
POST /admin/reset
```

Reset server metrics (development only).

**Response: 200 OK**
```json
{
  "message": "Metrics reset"
}
```

---

## User Management

### Create User

```http
POST /api/users
```

Register a new user account.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response: 201 Created**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z",
  "is_chirpy_red": false
}
```

**Error Responses**
- `400 Bad Request`: Invalid email or password format
- `500 Internal Server Error`: User creation failed

### Update User

```http
PUT /api/users
```

Update user email and/or password. Requires authentication.

**Headers**
```
Authorization: Bearer <access-token>
```

**Request Body**
```json
{
  "email": "newemail@example.com",
  "password": "newpassword123"
}
```

**Response: 200 OK**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "newemail@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T11:45:00Z",
  "is_chirpy_red": false
}
```

**Error Responses**
- `401 Unauthorized`: Invalid or missing token
- `400 Bad Request`: Invalid request body
- `500 Internal Server Error`: Update failed

---

## Authentication Endpoints

### Login

```http
POST /api/login
```

Authenticate user and receive access and refresh tokens.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response: 200 OK**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0",
  "is_chirpy_red": false
}
```

**Error Responses**
- `400 Bad Request`: Missing email or password
- `401 Unauthorized`: Incorrect email or password

### Refresh Access Token

```http
POST /api/refresh
```

Get a new access token using a refresh token.

**Headers**
```
Authorization: Bearer <refresh-token>
```

**Response: 200 OK**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses**
- `401 Unauthorized`: Invalid, expired, or revoked refresh token

### Revoke Refresh Token

```http
POST /api/revoke
```

Revoke a refresh token to log out.

**Headers**
```
Authorization: Bearer <refresh-token>
```

**Response: 204 No Content**

**Error Responses**
- `401 Unauthorized`: Missing authorization header
- `404 Not Found`: Token not found
- `500 Internal Server Error`: Revocation failed

---

## Chirps

### Create Chirp

```http
POST /api/chirps
```

Create a new chirp. Requires authentication.

**Headers**
```
Authorization: Bearer <access-token>
```

**Request Body**
```json
{
  "body": "This is my first chirp!"
}
```

**Response: 201 Created**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "body": "This is my first chirp!",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-01-15T12:00:00Z",
  "updated_at": "2025-01-15T12:00:00Z"
}
```

**Profanity Filter**

The following words are automatically replaced with `****`:
- kerfuffle
- sharbert
- fornax

**Validation Rules**
- Maximum 140 characters
- Body cannot be empty

**Error Responses**
- `401 Unauthorized`: Invalid or missing token
- `400 Bad Request`: Chirp exceeds 140 characters or invalid format

### Get All Chirps

```http
GET /api/chirps
```

Retrieve all chirps, with optional filtering and sorting.

**Query Parameters**
- `author_id` (optional): Filter chirps by user UUID
- `sort` (optional): Sort order - `asc` (default) or `desc`

**Examples**
```http
GET /api/chirps
GET /api/chirps?sort=desc
GET /api/chirps?author_id=550e8400-e29b-41d4-a716-446655440000
GET /api/chirps?author_id=550e8400-e29b-41d4-a716-446655440000&sort=desc
```

**Response: 200 OK**
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "body": "This is my first chirp!",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-01-15T12:00:00Z",
    "updated_at": "2025-01-15T12:00:00Z"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440002",
    "body": "Another chirp here",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-01-15T13:00:00Z",
    "updated_at": "2025-01-15T13:00:00Z"
  }
]
```

**Error Responses**
- `400 Bad Request`: Invalid author_id UUID or invalid sort parameter
- `404 Not Found`: Author ID not found

### Get Chirp by ID

```http
GET /api/chirps/{chirpID}
```

Retrieve a specific chirp by its UUID.

**Response: 200 OK**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "body": "This is my first chirp!",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-01-15T12:00:00Z",
  "updated_at": "2025-01-15T12:00:00Z"
}
```

**Error Responses**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Chirp not found

### Delete Chirp

```http
DELETE /api/chirps/{chirpID}
```

Delete a chirp. Requires authentication. Users can only delete their own chirps.

**Headers**
```
Authorization: Bearer <access-token>
```

**Response: 204 No Content**

**Error Responses**
- `401 Unauthorized`: Invalid or missing token
- `403 Forbidden`: User does not own this chirp
- `404 Not Found`: Chirp not found
- `400 Bad Request`: Invalid chirp ID format

---

## Webhooks

### Polka Webhook

```http
POST /api/polka/webhooks
```

Webhook endpoint for Polka payment service to upgrade users to Chirpy Red (premium).

**Headers**
```
Authorization: ApiKey <polka-api-key>
```

**Request Body**
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Response: 204 No Content**

**Behavior**
- Only processes `user.upgraded` events
- Other event types return 204 but take no action
- Upgrades user's `is_chirpy_red` status to `true`

**Error Responses**
- `401 Unauthorized`: Invalid or missing API key
- `400 Bad Request`: Invalid user ID format
- `404 Not Found`: User not found

---

## Error Responses

### Standard Error Format

```json
{
  "error": "error description"
}
```

### HTTP Status Codes

- `200 OK`: Request succeeded
- `201 Created`: Resource created successfully
- `204 No Content`: Request succeeded with no response body
- `400 Bad Request`: Invalid request format or parameters
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Authenticated but not authorized for this action
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error occurred

---

## Rate Limiting

Currently, no rate limiting is implemented. This is a development/learning project.

## CORS

CORS headers are not configured by default. Configure in production as needed.

## Content Type

All request and response bodies use `application/json` unless otherwise specified.
