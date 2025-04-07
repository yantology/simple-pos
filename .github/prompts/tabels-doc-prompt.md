# Database Schema Documentation

This document lists all the database tables in the project and their migration history. Use this documentation as a reference when creating new tables or modifying existing ones.

## Table List

- [Table 1](#table-1)
- [Table 2](#table-2)
- [Tenant](#tenant)

## Table Details

### Table 1

**Table Name:** `table_name_1`

**Description:** A brief description of the purpose of this table.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|-------|-----------|----------|---------|------------|
| id | int | no | auto_increment | Primary Key |
| name | varchar(255) | no | - | Username |
| email | varchar(255) | no | - | Email must be unique |
| created_at | timestamp | no | CURRENT_TIMESTAMP | Record creation time |
| updated_at | timestamp | yes | NULL | Last updated time |

**Index:**

- PRIMARY KEY (`id`)
- UNIQUE INDEX `email_unique` (`email`)

**Relationship:**

- `user_id` references `users(id)`

**Migration History:**

- `YYYYMMDDHHMMSS_create_table_1.php` - Initial table creation
- `YYYYMMDDHHMMSS_add_column_to_table_1.php` - Adding a new column

### Table 2

**Table Name:** `table_name_2`

**Description:** A brief description of the purpose of this table.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|-------|-----------|----------|---------|------------|
| id | int | no | auto_increment | Primary Key |
| title | varchar(100) | no | - | Item title |
| description | text | yes | NULL | Full description |
| status | enum('active','inactive') | no | 'active' | Item status |
| created_at | timestamp | no | CURRENT_TIMESTAMP | Record creation time |
| updated_at | timestamp | yes | NULL | Last updated time |

**Index:**

- PRIMARY KEY (`id`)
- INDEX `status_index` (`status`)

**Relations:**

- `category_id` references `categories(id)`

**Migration History:**

- `YYYYMMDDHHMMSS_create_table_2.php` - Initial table creation
- `YYYYMMDDHHMMSS_status_column.php` - Adding status column

### Tenant

**Table Name:** `tenant`

**Description:** Stores tenant information linked to users.

**Structure:**

| Column      | Data Type      | Nullable | Default           | Description                     |
|-------------|----------------|----------|-------------------|---------------------------------|
| id          | int            | no       | auto_increment    | Primary Key                    |
| username    | varchar(255)   | no       | -                 | Tenant username                |
| description | text           | yes      | NULL              | Description of the tenant      |
| user_id     | int            | no       | -                 | Foreign key referencing `users`|
| created_at  | timestamp      | no       | CURRENT_TIMESTAMP | Record creation time           |

**Index:**

- PRIMARY KEY (`id`)

**Relations:**

- `user_id` references `users(id)` with `ON DELETE CASCADE`

**Migration History:**

- `20250320000002_create_tenant_table.up.sql` - Initial table creation

<!--
## Template for New Table

### Table Name

**Table Name:** `new_table_name`

**Description:** A brief description of the purpose of this table.

**Structure:**

| Column | Data Type | Nullable | Default | Description |
|-------|-----------|----------|---------|------------|
| id | int | no | auto_increment | Primary Key |
| ... | ... | ... | ... | ... |

**Index:**
- PRIMARY KEY (`id`)
- ...

**Relations:**
- ...

**Migration History:**
- `YYYY_MM_DD_create_new_table.php` - Initial table creation
-->

## Updating Document Guidelines

1. When creating a new table, add new entries using the same format as the existing table
2. When modifying an existing table, update the table structure and add new entries to the migration history
3. Always make sure to record the migration file name along with the date

## Important Notes

- Do not delete old migration history, this is important for tracking database changes
- Make sure all relationships between tables are properly documented
- If there are significant table structure changes, add comments explaining the reason for the change
