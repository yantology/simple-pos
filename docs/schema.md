# Database Schema Documentation

## Table List

- [Database Schema Documentation](#database-schema-documentation)
  - [Table List](#table-list)
  - [Table Details](#table-details)
    - [Users](#users)
    - [Activation Tokens](#activation-tokens)
    - [Tenant](#tenant)
    - [Roles](#roles)
    - [Tenant Users](#tenant-users)
    - [Products](#products)
    - [Invoices](#invoices)
    - [Invoice Items](#invoice-items)

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

### Tenant

**Table Name:** `tenant`

**Description:** Stores tenant information linked to users.

**Structure:**

| Column      | Data Type      | Nullable | Default           | Description                     |
|-------------|----------------|----------|-------------------|---------------------------------|
| id          | SERIAL         | no       | auto_increment    | Primary Key                    |
| username    | VARCHAR(255)   | no       | -                 | Tenant username (unique)       |
| description | TEXT           | yes      | NULL              | Description of the tenant      |
| user_id     | INT            | no       | -                 | Foreign key referencing `users`|
| created_at  | TIMESTAMP      | no       | CURRENT_TIMESTAMP | Record creation time           |

**Index:**
- PRIMARY KEY (`id`)
- UNIQUE INDEX (`username`)
- INDEX (`user_id`)

**Relations:**
- `user_id` references `users(id)` with `ON DELETE CASCADE`

**Migration History:**
- `20250320000002_create_tenant_table.up.sql` - Initial table creation

### Roles

**Table Name:** `roles`

**Description:** Simple role definitions for tenant user management.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|--------|-----------|----------|---------|-------------|
| id | SERIAL | no | auto_increment | Primary Key |
| name | VARCHAR(50) | no | - | Role name (unique) |

**Index:**
- PRIMARY KEY (`id`)
- UNIQUE INDEX (`name`)

**Migration History:**
- `20250404053711_create_roles_table.up.sql` - Initial table creation with two roles (creator, user)

### Tenant Users

**Table Name:** `tenant_users`

**Description:** Links users to tenants with specific roles. Creators can add users to their tenant.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|--------|-----------|----------|---------|-------------|
| id | SERIAL | no | auto_increment | Primary Key |
| tenant_id | INT | no | - | Foreign key to `tenant` |
| user_id | INT | no | - | Foreign key to `users` |
| role_id | INT | no | - | Foreign key to `roles` |

**Index:**
- PRIMARY KEY (`id`)
- UNIQUE INDEX (`tenant_id`, `user_id`)

**Relations:**
- `tenant_id` references `tenant(id)` with `ON DELETE CASCADE`
- `user_id` references `users(id)` with `ON DELETE CASCADE`
- `role_id` references `roles(id)`

**Migration History:**
- `20250404053712_create_tenant_users_table.up.sql` - Initial table creation

### Products

**Table Name:** `products`

**Description:** Stores product information for each tenant with unique name constraint within each tenant.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|--------|-----------|----------|---------|-------------|
| id | SERIAL | no | auto_increment | Primary Key |
| tenant_id | INT | no | - | Foreign key to `tenant` |
| name | VARCHAR(255) | no | - | Product name (unique per tenant) |
| price | DECIMAL(10,2) | no | - | Product price |
| created_at | TIMESTAMP | no | CURRENT_TIMESTAMP | Record creation time |
| updated_at | TIMESTAMP | yes | NULL | Last updated time |

**Index:**
- PRIMARY KEY (`id`)
- UNIQUE INDEX (`tenant_id`, `name`)
- INDEX (`name`)
- INDEX (`price`)
- GIN INDEX (`name` using pg_trgm) for text search

**Relations:**
- `tenant_id` references `tenant(id)` with `ON DELETE CASCADE`

**Migration History:**
- `20250405000001_create_products_table.up.sql` - Initial table creation

### Invoices

**Table Name:** `invoices`

**Description:** Stores invoice records for each tenant's transactions.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|-------|-----------|----------|---------|------------|
| id | SERIAL | no | auto_increment | Primary Key |
| transaction_date | timestamp | no | CURRENT_TIMESTAMP | Date of transaction |
| total_amount | decimal(10,2) | no | - | Total amount of invoice (must be >= 0) |
| tenant_id | int | no | - | Reference to tenant table |
| user_id | int | no | - | Reference to users table |
| created_at | timestamp | no | CURRENT_TIMESTAMP | Record creation time |

**Index:**
- PRIMARY KEY (`id`)
- INDEX `idx_invoices_tenant_id` (`tenant_id`)
- INDEX `idx_invoices_user_id` (`user_id`)

**Relations:**
- `tenant_id` references `tenant(id)` ON DELETE CASCADE
- `user_id` references `users(id)` ON DELETE CASCADE

**Migration History:**
- `20250405000002_create_invoices_table.up.sql` - Initial table creation with foreign key indexes

### Invoice Items

**Table Name:** `invoice_items`

**Description:** Stores individual line items for each invoice.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|-------|-----------|----------|---------|------------|
| id | SERIAL | no | auto_increment | Primary Key |
| invoice_id | int | no | - | Reference to invoices table |
| product_name | varchar(255) | no | - | Name of the product |
| quantity | int | no | - | Quantity purchased (must be > 0) |
| unit_price | decimal(10,2) | no | - | Price per unit (must be >= 0) |
| total_price | decimal(10,2) | no | - | Total price for this item (must be >= 0) |
| created_at | timestamp | no | CURRENT_TIMESTAMP | Record creation time |

**Index:**
- PRIMARY KEY (`id`)
- INDEX `idx_invoice_items_invoice_id` (`invoice_id`)

**Relations:**
- `invoice_id` references `invoices(id)` ON DELETE CASCADE

**Migration History:**
- `20250405000003_create_invoice_items_table.up.sql` - Initial table creation with index on invoice_id