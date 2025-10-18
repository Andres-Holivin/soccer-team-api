# Go REST API with Gin + GORM + PostgreSQL

A modern REST API boilerplate built with Go, featuring Gin web framework, GORM ORM, and PostgreSQL database.

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin) - Fast HTTP web framework
- **ORM**: [GORM](https://gorm.io/) - Developer-friendly ORM library
- **Database**: PostgreSQL 12+ - Production-ready relational database
- **Authentication**: JWT (JSON Web Tokens)
- **Security**: bcrypt for password hashing

## ✨ Features

- **RESTful API Architecture** - Clean and organized REST endpoints
- **JWT Authentication** - Secure token-based authentication
- **Database Migrations** - Automatic schema migrations with GORM
- **Soft Deletes** - Safe data deletion with recovery options
- **Request Validation** - Input validation using Gin binding
- **Error Handling** - Consistent error responses
- **Pagination Support** - Built-in pagination for list endpoints
- **Environment Configuration** - .env file support
- **Database Seeding** - CLI tool for populating test data
- **CORS Support** - Cross-origin resource sharing
- **Timezone Aware** - PostgreSQL timestamptz support

## 📋 Prerequisites

- **Go** 1.21 or higher
- **PostgreSQL** 12 or higher
- **Git** (for version control)

## 🛠️ Installation & Setup

### 1. Clone Repository

```bash
git clone https://github.com/Andres-Holivin/soccer-team-api
cd soccer-team-api
```

### 2. Install Dependencies

```bash
go mod download
```

Or initialize a new Go module:

```bash
go mod init <your-module-name>
go mod tidy
```

### 3. Configure Environment Variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` file with your configuration:

```env
GIN_MODE=debug

APP_PORT=3000
APP_ENV=development

ALLOW_ORIGINS=*

DATABASE_URL='postgresql://local_user:local_password@localhost:5432/localdb?sslmode=disable'

JWT_SECRET=your_jwt_secret_key
JWT_EXPIRES_IN=1h

CLOUDINARY_URL=cloudinary://123456789012345:example_key@example_cloud
CLOUDINARY_ROOT_FOLDER=soccer-team-api
```

### 4. Run the Application

#### Development Mode

```bash
# Run directly
go run main.go

# Or with hot reload (requires air)
air
```

The server will start at `http://localhost:3000`

### 5. Verify Installation

Test if the server is running:

```bash
curl http://localhost:3000/health
```

## 📚 API Documentation

### Base URL
```
http://localhost:3000/api/v1
```

### Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Standard Response Format

#### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

#### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error information"
}
```

#### Paginated Response
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with cost 14 for secure password storage
- **JWT Authentication**: Token-based authentication with configurable expiration
- **Protected Routes**: Middleware-based route protection
- **Soft Deletes**: Safe data deletion with recovery capability
- **Input Validation**: Request validation using Gin binding tags
- **SQL Injection Protection**: Parameterized queries via GORM ORM
- **CORS Support**: Configurable cross-origin resource sharing
- **Environment Variables**: Sensitive data stored in .env files
## 📖 Additional Resources

### Recommended Tools

- **[Postman](https://www.postman.com/)** - API testing
- **[DBeaver](https://dbeaver.io/)** - Database management
- **[Air](https://github.com/cosmtrek/air)** - Hot reload for Go
- **[golangci-lint](https://golangci-lint.run/)** - Go linter
