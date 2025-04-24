# Database Migration Guide

## Overview

This guide covers how to work with database migrations in our Go project. We use SQL migrations to manage database schema changes in a version-controlled and organized manner.

## Migration Structure

All migrations are stored in the `migrations` directory. Each migration consists of two files following this naming convention:

```
YYYYMMDDHHMMSS_description.up.sql   # For applying changes
YYYYMMDDHHMMSS_description.down.sql # For rolling back changes
```

## Creating New Migrations

### File Naming

The timestamp format should be: `YYYYMMDDHHMMSS` (Year, Month, Day, Hour, Minute, Second)
Example: `20240115143022_create_users_table.up.sql`

### File Structure

1. Create two files in the `migrations` directory:

   ```
   migrations/
   ├── YYYYMMDDHHMMSS_description.up.sql
   └── YYYYMMDDHHMMSS_description.down.sql
   ```

2. Example Migration Content:

   ```sql
   -- YYYYMMDDHHMMSS_create_users_table.up.sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       username VARCHAR(255) NOT NULL,
       email VARCHAR(255) UNIQUE NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   -- YYYYMMDDHHMMSS_create_users_table.down.sql
   DROP TABLE IF EXISTS users;
   ```

## Creating or updating documentation Schema

 <!-- for guide please look #file:tabels-doc-prompt.md for intructuion and create on ../docs/shcema.md -->
