# YAD Backend API Implementation Guide

## Overview

This document provides a comprehensive guide to the YAD backend API implementation, covering all components, architecture, and API endpoints.

## Technology Stack

- **Language**: Go 1.25
- **Web Framework**: Fiber v2 (Fast HTTP framework)
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT (JSON Web Tokens)
- **Security**: Helmet middleware
- **Configuration**: Viper + godotenv
- **ID Generation**: Google UUID

## Clean Architecture Layers

The backend follows clean architecture principles with strict separation of concerns:

### 1. **Domain Layer** (`internal/domain/`)
Contains pure business logic and entity definitions.

**Models implemented:**
- `User` - User account and authentication
- `Minyan` - Prayer group gatherings
- `Notification` - User notifications
- `Place` - Saved venues/locations by users
- `Location` - Generic location references
- `UserSettings` - User preferences and settings
- `AuthToken` - JWT token management

Each domain model includes:
- Entity struct with GORM tags
- Constructor function (`NewEntity`)
- Repository interface defining data access contracts

### 2. **Repository Layer** (`internal/repository/`)
Handles all database operations with GORM.

**Implemented repositories:**
- `UserRepository` - User CRUD operations
- `AuthTokenRepository` - Token management
- `MinyanRepository` - Minyan CRUD with geospatial queries
- `NotificationRepository` - Notification CRUD
- `PlaceRepository` - Saved places management
- `LocationRepository` - Location queries
- `UserSettingsRepository` - Settings management

**Key patterns:**
- Repository interfaces defined in domain models
- Factory functions create repository instances
- All error handling with detailed messages
- Pagination support where applicable

### 3. **Service Layer** (`internal/service/`)
Implements business logic and validation rules.

**Implemented services:**
- `AuthService` - Authentication and JWT token generation
- `UserService` - User management and profile updates
- `MinyanService` - Minyan lifecycle management
- `NotificationService` - Notification creation and management
- `PlaceService` - Saved place management
- `LocationService` - Location search and retrieval
- `UserSettingsService` - Settings management

**Responsibilities:**
- Business logic validation
- Service composition
- Error handling and messaging
- Cross-cutting concerns

### 4. **Handler Layer** (`internal/handler/`)
HTTP request/response handling for REST API.

**Implemented handlers:**
- `AuthHandler` - Login, register, token refresh
- `UserHandler` - User profile endpoints
- `MinyanHandler` - Minyan CRUD endpoints
- `NotificationHandler` - Notification endpoints
- `PlaceHandler` - Saved places endpoints
- `LocationHandler` - Location search endpoints
- `UserSettingsHandler` - Settings management endpoints

**Key features:**
- Request/response DTOs with validation tags
- Standardized error responses
- Pagination parameter parsing
- Context-based request data (userId from middleware)

### 5. **Configuration** (`internal/config/`)
Environment and application configuration management.

**Configuration sections:**
- `Port` - Server port (default: 8080)
- `Env` - Environment (development/production)
- `Database` - MySQL connection details
- `JWT` - Token secret and duration settings
- `Datadog` - Optional observability configuration

### 6. **Middleware** (`internal/middleware/`)
Cross-cutting concerns and request processing.

**Implemented middleware:**
- `AuthMiddleware` - JWT token verification and user context extraction
- `CORSMiddleware` - Cross-origin resource sharing
- Helmet security headers

### 7. **Database** (`internal/database/`)
Database initialization and migrations.

**Features:**
- MySQL connection pool management
- Automatic schema migration with GORM
- Support for all domain models

## API Endpoints

### Authentication Endpoints

#### Register User
```
POST /api/v1/users/register
Content-Type: application/json

{
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "password": "SecurePass123"
}

Response (201):
{
  "id": "uuid-string",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "createdAt": "2024-12-20T10:30:00Z"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123"
}

Response (200):
{
  "accessToken": "eyJhbGciOiJIUzI1NiJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiJ9...",
  "user": {
    "id": "uuid-string",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }
}
```

#### Refresh Token
```
POST /api/v1/auth/refresh
Content-Type: application/json
Authorization: Bearer {refreshToken}

{}

Response (200):
{
  "accessToken": "new-jwt-token",
  "refreshToken": "new-refresh-token"
}
```

### Minyan Endpoints (Protected)

#### Create Minyan
```
POST /api/v1/minyans
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "prayerType": "Shacharit",
  "date": "2024-12-20",
  "time": "06:30",
  "locationName": "Central Synagogue",
  "latitude": 40.7615,
  "longitude": -73.9776,
  "notes": "Everyone welcome"
}

Response (201):
{
  "id": "uuid-string",
  "userId": "user-uuid",
  "prayerType": "Shacharit",
  "date": "2024-12-20",
  "time": "06:30",
  "locationName": "Central Synagogue",
  "latitude": 40.7615,
  "longitude": -73.9776,
  "notes": "Everyone welcome",
  "status": "draft",
  "participantCount": 1,
  "createdAt": "2024-12-20T10:30:00Z"
}
```

#### Get User's Minyans
```
GET /api/v1/minyans/my?limit=20&offset=0
Authorization: Bearer {accessToken}

Response (200):
{
  "data": [...],
  "count": 5,
  "limit": 20,
  "offset": 0
}
```

#### Find Nearby Minyans
```
GET /api/v1/minyans/nearby?latitude=40.7128&longitude=-74.0060&radiusKm=5&limit=20
Authorization: Bearer {accessToken}

Response (200):
{
  "data": [
    {
      "id": "uuid",
      "prayerType": "Shacharit",
      "date": "2024-12-20",
      "locationName": "Central Synagogue",
      "latitude": 40.7615,
      "longitude": -73.9776,
      "distance": 3.2,
      "participantCount": 3
    }
  ],
  "count": 1
}
```

#### Publish Minyan
```
POST /api/v1/minyans/{id}/publish
Authorization: Bearer {accessToken}

Response (200):
{
  "status": "published",
  "updatedAt": "2024-12-20T11:00:00Z"
}
```

#### Join Minyan
```
POST /api/v1/minyans/{id}/join
Authorization: Bearer {accessToken}

Response (200):
{
  "message": "Successfully joined minyan",
  "participantCount": 4
}
```

### Notification Endpoints (Protected)

#### Create Notification
```
POST /api/v1/notifications
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "title": "Minyan Available",
  "message": "New minyan created near you",
  "type": "minyan_alert",
  "relatedMinyanId": "minyan-uuid"
}

Response (201):
{
  "id": "uuid",
  "title": "Minyan Available",
  "message": "New minyan created near you",
  "type": "minyan_alert",
  "isRead": false,
  "createdAt": "2024-12-20T10:30:00Z"
}
```

#### Get Notifications
```
GET /api/v1/notifications?limit=20&offset=0
Authorization: Bearer {accessToken}

Response (200):
{
  "data": [...],
  "count": 10,
  "limit": 20,
  "offset": 0
}
```

#### Get Unread Count
```
GET /api/v1/notifications/unread-count
Authorization: Bearer {accessToken}

Response (200):
{
  "unreadCount": 3
}
```

#### Mark as Read
```
PUT /api/v1/notifications/{id}/read
Authorization: Bearer {accessToken}

Response (200): OK
```

#### Mark All as Read
```
PUT /api/v1/notifications/read-all
Authorization: Bearer {accessToken}

Response (200): OK
```

### Place Endpoints (Protected)

#### Save Place
```
POST /api/v1/places/saved
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "placeId": "google-place-id",
  "name": "Central Synagogue",
  "address": "123 E 55th St, New York, NY 10022",
  "latitude": 40.7615,
  "longitude": -73.9776,
  "placeType": "Synagogue",
  "phoneNumber": "+1 (212) 838-5122",
  "website": "https://www.centralsynagogue.org"
}

Response (201):
{
  "id": "uuid",
  "placeId": "google-place-id",
  "name": "Central Synagogue",
  "address": "123 E 55th St, New York, NY 10022",
  "latitude": 40.7615,
  "longitude": -73.9776,
  "savedAt": "2024-12-20T10:30:00Z"
}
```

#### Get Saved Places
```
GET /api/v1/places/saved?limit=20&offset=0
Authorization: Bearer {accessToken}

Response (200):
{
  "data": [...],
  "count": 5,
  "limit": 20,
  "offset": 0
}
```

#### Check if Place is Saved
```
GET /api/v1/places/saved/{placeId}
Authorization: Bearer {accessToken}

Response (200):
{
  "placeId": "google-place-id",
  "isSaved": true
}
```

#### Remove Saved Place
```
DELETE /api/v1/places/saved/{id}
Authorization: Bearer {accessToken}

Response (200): OK
```

#### Search Places
```
GET /api/v1/places/search?q=synagogue&limit=20&offset=0
Authorization: Bearer {accessToken}

Response (200):
{
  "data": [...],
  "count": 2,
  "query": "synagogue",
  "limit": 20,
  "offset": 0
}
```

## Environment Configuration

Create a `.env` file in the project root:

```env
PORT=8080
ENV=development

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=yad_db

JWT_SECRET=your-super-secret-key-change-in-production-minimum-32-chars
JWT_ACCESS_DURATION=15
JWT_REFRESH_DURATION=168

DATADOG_ENABLED=false
DATADOG_API_KEY=
DATADOG_APP_KEY=
DATADOG_ENV=development
DATADOG_SERVICE=yad-backend
```

## Database Setup

### Create Database
```sql
CREATE DATABASE yad_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### Automatic Migration
The application automatically creates all tables on startup using GORM's AutoMigrate feature.

### Manual Schema Creation (Optional)
See the SQL schema in the `infrastructure/` directory for manual table creation.

## Running the Server

### Development
```bash
cd yad-back
go run ./cmd/main.go
```

### Production Build
```bash
go build -o bin/yad-back ./cmd/main.go
./bin/yad-back
```

### Docker
```bash
docker build -t yad-backend .
docker run -p 8080:8080 --env-file .env yad-backend
```

### Docker Compose
```bash
docker-compose up -d yad-backend
```

## Error Handling

All API responses follow a consistent error format:

```json
{
  "error": "Description of what went wrong"
}
```

### HTTP Status Codes

- **200** - Success
- **201** - Created
- **400** - Bad Request (validation error)
- **401** - Unauthorized (missing or invalid token)
- **404** - Not Found
- **500** - Internal Server Error

## Security Considerations

1. **Password Hashing**
   - All passwords are hashed using bcrypt
   - Never stored or returned in API responses

2. **JWT Tokens**
   - Access token valid for 15 minutes
   - Refresh token valid for 7 days
   - Tokens must be sent in Authorization header: `Bearer {token}`

3. **CORS**
   - Configured to allow frontend origin
   - Update in production with actual domain

4. **SQL Injection**
   - GORM prevents SQL injection through parameterized queries

5. **Rate Limiting**
   - Consider implementing rate limiting in production
   - Use middleware for protection against brute force attacks

## Testing API Endpoints

### Using cURL

```bash
# Register
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "password": "SecurePass123"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'

# Create Minyan (with token)
curl -X POST http://localhost:8080/api/v1/minyans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{
    "prayerType": "Shacharit",
    "date": "2024-12-20",
    "time": "06:30",
    "locationName": "Central Synagogue",
    "latitude": 40.7615,
    "longitude": -73.9776,
    "notes": "Everyone welcome"
  }'
```

### Using Postman

Import the API collection into Postman:
1. Create a collection "YAD API"
2. Create environment variables for `baseUrl` and `token`
3. Add requests for each endpoint
4. Use pre-request scripts to update token after login

## Integration with Frontend

The Flutter frontend expects:
- Base URL: `http://localhost:8080/api/v1`
- Response format: Standard JSON with data models
- Authentication: Bearer JWT tokens in Authorization header
- Error format: `{"error": "error message"}`

Update `AppConfig.useMockApi` to `false` in Flutter app when backend is running.

## Logging and Monitoring

The application logs important events:
- Database connection status
- Authentication attempts
- API request errors
- Server startup/shutdown

### Enabling Datadog (Optional)
Set `DATADOG_ENABLED=true` and provide API/App keys for cloud monitoring.

## Future Enhancements

1. **Rate Limiting** - Prevent brute force and abuse
2. **Request Validation** - Add comprehensive input validation
3. **Database Indexes** - Optimize geospatial queries
4. **Caching** - Redis for frequently accessed data
5. **WebSockets** - Real-time notifications
6. **File Upload** - User profile pictures
7. **Email Service** - Notification emails
8. **SMS Notifications** - Text message alerts
9. **Analytics** - User behavior tracking
10. **API Documentation** - Swagger/OpenAPI

