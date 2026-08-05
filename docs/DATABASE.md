# Database Documentation

Complete database schema and migration guide for Chirpy.

## Database System

Chirpy uses **PostgreSQL** as its database system with the following tools:

- **sqlc**: Generate type-safe Go code from SQL queries
- **Goose**: Database migration management

---

## Database Schema

### Tables Overview

1. **users** - User accounts and authentication
2. **chirps** - Short message posts
3. **refresh_tokens** - Token management for authentication

---

## Schema Details

### users

Stores user account information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Unique user identifier |
| `email` | TEXT | UNIQUE, NOT NULL | User's email address |
| `hashed_password` | TEXT | NOT NULL | Argon2id hashed password |
| `is_chirpy_red` | BOOLEAN | DEFAULT false | Premium membership status |
| `created_at` | TIMESTAMP | NOT NULL | Account creation timestamp |
| `updated_at` | TIMESTAMP | NOT NULL | Last update timestamp |

**Indexes:**
- Primary key on `id`
- Unique constraint on `email`

**Migration File:** `001_users.sql`, `003_add_column_hashed_password.sql`, `005_alter_table_users_add_is_chirpy_red.sql`

---

### chirps

Stores user posts (chirps).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PRIMARY KEY | Unique chirp identifier |
| `body` | VARCHAR(140) | NOT NULL | Chirp content (max 140 chars) |
| `user_id` | UUID | FOREIGN KEY, NOT NULL | References `users(id)` |
| `created_at` | TIMESTAMP | NOT NULL | Chirp creation timestamp |
| `updated_at` | TIMESTAMP | NOT NULL | Last update timestamp |

**Foreign Keys:**
- `user_id` REFERENCES `users(id)` ON DELETE CASCADE

**Indexes:**
- Primary key on `id`
- Foreign key index on `user_id`

**Cascade Behavior:**
- When a user is deleted, all their chirps are automatically deleted

**Migration File:** `002_chirps.sql`

---

### refresh_tokens

Manages authentication refresh tokens.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `token` | TEXT | PRIMARY KEY | The refresh token string |
| `user_id` | UUID | FOREIGN KEY, NOT NULL | References `users(id)` |
| `expires_at` | TIMESTAMP | NOT NULL | Token expiration time (60 days) |
| `revoked_at` | TIMESTAMP | NULL | When token was revoked (if applicable) |
| `created_at` | TIMESTAMP | NOT NULL | Token creation timestamp |
| `updated_at` | TIMESTAMP | NOT NULL | Last update timestamp |

**Foreign Keys:**
- `user_id` REFERENCES `users(id)` ON DELETE CASCADE

**Indexes:**
- Primary key on `token`
- Foreign key index on `user_id`

**Token Lifecycle:**
- Tokens expire after 60 days
- Tokens can be manually revoked (logout)
- Expired or revoked tokens cannot be used

**Migration File:** `004_create_refresh_tokens.sql`

---

## Entity Relationships

```
users (1) ──< (many) chirps
  │
  └──< (many) refresh_tokens
```

- One user can have many chirps
- One user can have many refresh tokens
- Deleting a user cascades to delete all their chirps and tokens

---

## Migrations

### Migration Files Location

```
sql/schema/
├── 001_users.sql
├── 002_chirps.sql
├── 003_add_column_hashed_password.sql
├── 004_create_refresh_tokens.sql
└── 005_alter_table_users_add_is_chirpy_red.sql
```

### Running Migrations

#### Install Goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

#### Apply All Migrations (Up)

```bash
goose -dir sql/schema postgres "postgresql://username:password@localhost:5432/chirpy?sslmode=disable" up
```

#### Rollback Last Migration (Down)

```bash
goose -dir sql/schema postgres "postgresql://username:password@localhost:5432/chirpy?sslmode=disable" down
```

#### Check Migration Status

```bash
goose -dir sql/schema postgres "postgresql://username:password@localhost:5432/chirpy?sslmode=disable" status
```

#### Reset Database (Down All, Then Up All)

```bash
goose -dir sql/schema postgres "postgresql://username:password@localhost:5432/chirpy?sslmode=disable" reset
```

---

## SQL Queries

All SQL queries are located in `sql/queries/` and are compiled to type-safe Go code using **sqlc**.

### Query Files

```
sql/queries/
├── createUser.sql          # Create new user
├── getUserByEmail.sql      # Get user by email (for login)
├── updateUserInfo.sql      # Update user email/password
├── updateIsChirpyRed.sql   # Upgrade user to premium
├── deleteAllUsers.sql      # Delete all users (admin)
├── createChirp.sql         # Create new chirp
├── getAllChirps.sql        # Get all chirps (ordered by created_at ASC)
├── getChirpByID.sql        # Get specific chirp
├── getChirpsByUserId.sql   # Get chirps by author
├── deleteChirp.sql         # Delete chirp
├── createRefreshToken.sql  # Create refresh token
├── getRefreshToken.sql     # Get refresh token details
└── revokeRefreshToken.sql  # Revoke refresh token
```

### Generating Go Code from SQL

After modifying SQL query files, regenerate Go code:

```bash
sqlc generate
```

This creates type-safe Go functions in the `internal/database` package.

---

## Example SQL Queries

### Create a User

```sql
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;
```

### Get All Chirps (Ordered)

```sql
-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;
```

### Get Chirps by User ID

```sql
-- name: GetChirpsByUserId :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;
```

### Delete Chirp

```sql
-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;
```

---

## Database Setup

### Create Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE chirpy;

# Create user (optional)
CREATE USER chirpy_user WITH PASSWORD 'your_password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE chirpy TO chirpy_user;
```

### Environment Configuration

Set your database URL in `.env`:

```bash
DB_URL=postgresql://username:password@localhost:5432/chirpy?sslmode=disable
```

### Initialize Database

```bash
# Run all migrations
goose -dir sql/schema postgres "$DB_URL" up

# Verify tables were created
psql $DB_URL -c "\dt"
```

---

## Data Types

### UUIDs

All primary keys use UUID v4 for:
- Unpredictability (security)
- Global uniqueness
- No sequential guessing

Generated using PostgreSQL's `gen_random_uuid()` function.

### Timestamps

All timestamps are stored in UTC using PostgreSQL's `TIMESTAMP` type (without timezone).

- `created_at`: Set once on creation
- `updated_at`: Updated on every modification

### Password Hashing

Passwords are hashed using **Argon2id** algorithm with:
- Memory: 64 MB
- Iterations: 3
- Parallelism: 2
- Salt length: 16 bytes
- Key length: 32 bytes

---

## Backup and Restore

### Backup Database

```bash
pg_dump -U postgres -d chirpy -F c -f chirpy_backup.dump
```

### Restore Database

```bash
pg_restore -U postgres -d chirpy -c chirpy_backup.dump
```

---

## Performance Considerations

### Indexes

- Primary keys (`id` columns) are automatically indexed
- Foreign keys (`user_id` columns) have indexes for JOIN performance
- `email` has unique constraint (implicitly indexed) for fast lookups

### Query Optimization

- All queries retrieving multiple chirps are ordered by `created_at`
- User lookup by email is optimized via unique index
- Cascade deletes prevent orphaned records

### Connection Pooling

The application uses Go's `database/sql` package which provides built-in connection pooling.

---

## Development Tips

### View Table Structure

```sql
\d users
\d chirps
\d refresh_tokens
```

### Count Records

```sql
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM chirps;
SELECT COUNT(*) FROM refresh_tokens;
```

### Find Expired Tokens

```sql
SELECT * FROM refresh_tokens
WHERE expires_at < NOW()
   OR revoked_at IS NOT NULL;
```

### Get User with Their Chirp Count

```sql
SELECT u.email, COUNT(c.id) as chirp_count
FROM users u
LEFT JOIN chirps c ON u.id = c.user_id
GROUP BY u.id, u.email
ORDER BY chirp_count DESC;
```

### Find Premium Users

```sql
SELECT id, email, created_at
FROM users
WHERE is_chirpy_red = true;
```
