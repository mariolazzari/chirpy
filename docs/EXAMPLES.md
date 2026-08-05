# API Usage Examples

Practical examples for using the Chirpy API with curl and other tools.

## Prerequisites

- Chirpy server running on `http://localhost:8080`
- `curl` or similar HTTP client
- `jq` (optional, for JSON formatting)

---

## User Registration & Authentication Flow

### 1. Create a New User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "mysecurepassword"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z",
  "is_chirpy_red": false
}
```

### 2. Login to Get Tokens

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "mysecurepassword"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1NTBlODQwMC1lMjliLTQxZDQtYTcxNi00NDY2NTU0NDAwMDAiLCJleHAiOjE3MDUzMjc4MDB9.abc123...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "is_chirpy_red": false
}
```

Save the tokens for subsequent requests:
```bash
export ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
export REFRESH_TOKEN="a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
```

---

## Working with Chirps

### 3. Create a Chirp

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "body": "Hello Chirpy! This is my first post."
  }'
```

**Response:**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "body": "Hello Chirpy! This is my first post.",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-01-15T12:00:00Z",
  "updated_at": "2025-01-15T12:00:00Z"
}
```

### 4. Create a Chirp with Profanity (Gets Filtered)

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "body": "What a kerfuffle this sharbert situation is!"
  }'
```

**Response:**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440002",
  "body": "What a **** this **** situation is!",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-01-15T12:05:00Z",
  "updated_at": "2025-01-15T12:05:00Z"
}
```

### 5. Get All Chirps (Ascending Order)

```bash
curl http://localhost:8080/api/chirps
```

**Response:**
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "body": "Hello Chirpy! This is my first post.",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-01-15T12:00:00Z",
    "updated_at": "2025-01-15T12:00:00Z"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440002",
    "body": "What a **** this **** situation is!",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-01-15T12:05:00Z",
    "updated_at": "2025-01-15T12:05:00Z"
  }
]
```

### 6. Get All Chirps (Descending Order)

```bash
curl "http://localhost:8080/api/chirps?sort=desc"
```

### 7. Get Chirps by Specific Author

```bash
curl "http://localhost:8080/api/chirps?author_id=550e8400-e29b-41d4-a716-446655440000"
```

### 8. Get Chirps by Author with Descending Sort

```bash
curl "http://localhost:8080/api/chirps?author_id=550e8400-e29b-41d4-a716-446655440000&sort=desc"
```

### 9. Get a Specific Chirp by ID

```bash
curl http://localhost:8080/api/chirps/660e8400-e29b-41d4-a716-446655440001
```

### 10. Delete Your Own Chirp

```bash
curl -X DELETE http://localhost:8080/api/chirps/660e8400-e29b-41d4-a716-446655440001 \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Response:** `204 No Content`

### 11. Try to Delete Someone Else's Chirp (Fails)

```bash
curl -X DELETE http://localhost:8080/api/chirps/770e8400-e29b-41d4-a716-446655440099 \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Response:** `403 Forbidden`
```json
{
  "error": "forbidden"
}
```

---

## User Profile Management

### 12. Update User Information

```bash
curl -X PUT http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "email": "newemail@example.com",
    "password": "newsecurepassword"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "newemail@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T14:20:00Z",
  "is_chirpy_red": false
}
```

---

## Token Management

### 13. Refresh Access Token

When your access token expires (after 1 hour), use the refresh token:

```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer $REFRESH_TOKEN"
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.NEW_TOKEN_HERE..."
}
```

Update your access token:
```bash
export ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.NEW_TOKEN_HERE..."
```

### 14. Revoke Refresh Token (Logout)

```bash
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer $REFRESH_TOKEN"
```

**Response:** `204 No Content`

After revoking, the refresh token can no longer be used.

---

## Webhooks

### 15. Polka Webhook (Upgrade to Premium)

This endpoint is typically called by the Polka payment service:

```bash
curl -X POST http://localhost:8080/api/polka/webhooks \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey your-polka-api-key-here" \
  -d '{
    "event": "user.upgraded",
    "data": {
      "user_id": "550e8400-e29b-41d4-a716-446655440000"
    }
  }'
```

**Response:** `204 No Content`

The user's `is_chirpy_red` status is now `true`.

---

## Health Check & Metrics

### 16. Health Check

```bash
curl http://localhost:8080/api/healthz
```

**Response:**
```json
{
  "status": "ok"
}
```

### 17. View Metrics (Admin)

```bash
curl http://localhost:8080/admin/metrics
```

**Response:**
```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 42 times!</p>
  </body>
</html>
```

### 18. Reset Metrics (Admin)

```bash
curl -X POST http://localhost:8080/admin/reset
```

---

## Error Handling Examples

### Invalid Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "wrongpassword"
  }'
```

**Response:** `401 Unauthorized`
```json
{
  "error": "Incorrect email or password"
}
```

### Chirp Too Long

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "body": "This chirp is way too long and exceeds the 140 character limit that we have set for all chirps in this application to keep things short and sweet like Twitter used to be."
  }'
```

**Response:** `400 Bad Request`

### Missing Authentication

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -d '{
    "body": "Trying to chirp without auth"
  }'
```

**Response:** `401 Unauthorized`
```json
{
  "error": "unauthorized"
}
```

---

## Using with JavaScript (Fetch API)

```javascript
// Create a user
const createUser = async () => {
  const response = await fetch('http://localhost:8080/api/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: 'user@example.com',
      password: 'securepassword123'
    })
  });

  return await response.json();
};

// Login
const login = async (email, password) => {
  const response = await fetch('http://localhost:8080/api/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password })
  });

  const data = await response.json();
  localStorage.setItem('access_token', data.token);
  localStorage.setItem('refresh_token', data.refresh_token);

  return data;
};

// Create a chirp
const createChirp = async (body) => {
  const token = localStorage.getItem('access_token');

  const response = await fetch('http://localhost:8080/api/chirps', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({ body })
  });

  return await response.json();
};

// Get all chirps
const getChirps = async (authorId = null, sort = 'asc') => {
  let url = 'http://localhost:8080/api/chirps?sort=' + sort;
  if (authorId) {
    url += '&author_id=' + authorId;
  }

  const response = await fetch(url);
  return await response.json();
};
```

---

## Using with Python (Requests)

```python
import requests

BASE_URL = "http://localhost:8080"

# Create a user
def create_user(email, password):
    response = requests.post(
        f"{BASE_URL}/api/users",
        json={"email": email, "password": password}
    )
    return response.json()

# Login
def login(email, password):
    response = requests.post(
        f"{BASE_URL}/api/login",
        json={"email": email, "password": password}
    )
    data = response.json()
    return data['token'], data['refresh_token']

# Create a chirp
def create_chirp(token, body):
    response = requests.post(
        f"{BASE_URL}/api/chirps",
        headers={"Authorization": f"Bearer {token}"},
        json={"body": body}
    )
    return response.json()

# Get all chirps
def get_chirps(author_id=None, sort='asc'):
    params = {'sort': sort}
    if author_id:
        params['author_id'] = author_id

    response = requests.get(f"{BASE_URL}/api/chirps", params=params)
    return response.json()

# Example usage
if __name__ == "__main__":
    # Create user and login
    create_user("test@example.com", "password123")
    access_token, refresh_token = login("test@example.com", "password123")

    # Create a chirp
    chirp = create_chirp(access_token, "Hello from Python!")
    print(f"Created chirp: {chirp['id']}")

    # Get all chirps
    chirps = get_chirps(sort='desc')
    print(f"Total chirps: {len(chirps)}")
```
