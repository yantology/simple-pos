# simple-pos

simple-pos backend with golang and gin

support Swagger API

## Environment Configuration

Copy the `.env.example` file to `.env` and customize the values:

```bash
cp .env.example .env
```

### Available Environment Variables

#### App Configuration
- `APP_PORT`: Server port (default: 8080)

#### Database Configuration
- `DB_HOST`: Database host (default: 127.0.0.1)
- `DB_PORT`: Database port (default: 3306)
- `DB_NAME`: Database name
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_DRIVER`: Database driver (mysql/postgres)

#### JWT Configuration
- `JWT_ACCESS_SECRET`: Secret key for access tokens
- `JWT_REFRESH_SECRET`: Secret key for refresh tokens
- `JWT_ACCESS_DURATION_MINUTES`: Access token duration in minutes (default: 15)
- `JWT_REFRESH_DURATION_DAYS`: Refresh token duration in days (default: 7)
- `JWT_ISSUER`: Token issuer name (default: retail-pro)

#### CORS Configuration
- `CORS_ALLOW_ORIGINS`: Comma-separated list of allowed origins

#### Resend Email Configuration
- `RESEND_API_KEY`: Resend API key
- `RESEND_DOMAIN`: Email domain
- `RESEND_NAME`: Sender name
