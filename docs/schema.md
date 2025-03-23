# Database Schema Documentation

## Table List

- [Database Schema Documentation](#database-schema-documentation)
  - [Table List](#table-list)
  - [Table Details](#table-details)
    - [Users](#users)
    - [Activation Tokens](#activation-tokens)

## Table Details

### Users

**Table Name:** `users`

**Description:** Stores user account information including authentication details.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|--------|-----------|----------|---------|-------------|
| id | SERIAL | no | auto_increment | Primary Key |
| email | VARCHAR(255) | no | - | User's email address (unique) |
| fullname | VARCHAR(255) | no | - | User's full name |
| password_hash | VARCHAR(255) | no | - | Hashed user password |
| created_at | TIMESTAMP | no | CURRENT_TIMESTAMP | Record creation time |
| updated_at | TIMESTAMP | yes | NULL | Last updated time |

**Index:**
- PRIMARY KEY (`id`)
- UNIQUE INDEX (`email`)

**Migration History:**
- `20250320000001_create_users_table.sql` - Initial table creation

### Activation Tokens

**Table Name:** `activation_tokens`

**Description:** Manages tokens for various user account operations like email verification and password reset.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|--------|-----------|----------|---------|-------------|
| id | SERIAL | no | auto_increment | Primary Key |
| email | VARCHAR(255) | no | - | Associated user email |
| token_hash | VARCHAR(255) | no | - | Hashed token value |
| type | VARCHAR(50) | no | - | Token type (e.g., activation, reset) |
| expires_at | TIMESTAMP | no | - | Token expiration timestamp |
| created_at | TIMESTAMP | no | CURRENT_TIMESTAMP | Record creation time |

**Index:**
- PRIMARY KEY (`id`)
- UNIQUE INDEX (`email`, `type`)

**Migration History:**
- `20250320000002_create_activation_tokens_table.sql` - Initial table creation