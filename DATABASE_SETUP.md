# YAD Backend - Database Setup Guide

## Overview

The YAD backend is configured to automatically handle database connectivity with a smart fallback mechanism:

1. **Primary**: Attempts to connect to MySQL
2. **Fallback**: Uses SQLite for local development (if MySQL is unavailable)

This means you can start developing immediately without setting up MySQL, and can easily switch to MySQL for production.

---

## Quick Start (Easiest Option - SQLite Fallback)

The application will automatically fall back to SQLite if MySQL is not available. Just run:

```bash
cd yad-back
go run ./cmd/main.go
```

This will:
- Attempt to connect to MySQL (and succeed if available)
- Fall back to SQLite and create `yad.db` automatically
- Start the API server on `http://localhost:8080`
- Create all necessary database tables automatically

**No manual setup required!**

---

## Setup Options

### Option 1: Local Development (Recommended for Quick Start)

Simply run the application. It will use SQLite automatically:

```bash
cd yad-back
go run ./cmd/main.go
```

The application will:
- Create `yad.db` in the `yad-back` directory
- Set up all tables automatically via migrations
- Be ready to use immediately

### Option 2: Local MySQL Setup

If you want to use MySQL locally:

#### Step 1: Install MySQL

**macOS (Homebrew):**
```bash
brew install mysql
brew services start mysql
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install mysql-server
sudo systemctl start mysql
```

**Windows:**
- Download from: https://dev.mysql.com/downloads/mysql/
- Run the installer and follow the wizard

#### Step 2: Create Database and User (Automated)

Use the provided setup script:

```bash
chmod +x setup-db.sh
./setup-db.sh
```

The script will:
- Check if MySQL is installed and running
- Create the `yad_db` database
- Create the `root` user with password `yad_dev_password`
- Grant all necessary privileges

#### Step 3: Update .env File

Create a `.env` file (copy from `.env.example`):

```bash
cp .env.example .env
```

Ensure these settings:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yad_dev_password
DB_NAME=yad_db
```

#### Step 4: Run the Application

```bash
go run ./cmd/main.go
```

---

### Option 3: Docker Setup (Production-like)

Use Docker Compose to run both MySQL and the backend:

```bash
# From the project root
docker-compose up
```

This will:
- Start MySQL container
- Start the backend container
- Automatically initialize the database
- Expose the API on `http://localhost:8080`

---

## Database Configuration

The application reads database configuration from environment variables (in order of preference):

1. Environment variables (highest priority)
2. `.env` file
3. Defaults (hardcoded)

### Environment Variables

```bash
# Database Configuration
DB_HOST=localhost        # MySQL host
DB_PORT=3306            # MySQL port
DB_USER=root            # MySQL username
DB_PASSWORD=password    # MySQL password
DB_NAME=yad_db          # Database name

# Server Configuration
PORT=8080               # API server port
ENV=development         # Environment (development/production)

# JWT Configuration
JWT_SECRET=your-secret-key
```

### Default Configuration

If no `.env` file or environment variables are set, defaults are:

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=(empty)
DB_NAME=yad_db
PORT=8080
ENV=development
```

---

## SQLite vs MySQL

### SQLite (Default/Fallback)

**Pros:**
- No installation required
- Perfect for local development
- Single file database (`yad.db`)
- Automatic fallback if MySQL is unavailable

**Cons:**
- Not suitable for production with multiple servers
- Limited concurrent access

**When to use:**
- Local development
- Quick testing
- Prototyping

### MySQL

**Pros:**
- Production-grade
- Supports multiple concurrent connections
- Enterprise-standard relational database
- Excellent for scaling

**Cons:**
- Requires separate installation/service
- More setup overhead

**When to use:**
- Production environments
- Team development
- Docker deployment
- When you need multiple concurrent connections

---

## Troubleshooting

### "Connection refused" Error

This error means the application tried to connect to MySQL and failed. Two solutions:

1. **Use SQLite (easiest)**: Just let it fall back automatically - the app will work fine
2. **Set up MySQL**: Follow Option 2 above to install MySQL locally

### Database Already Exists

If you get an error that the database already exists, that's fine. The migrations are designed to be idempotent.

### Permission Denied Error

This might happen on Linux when creating the database. Use `sudo`:

```bash
sudo ./setup-db.sh
```

Or set up MySQL manually:

```bash
sudo mysql
```

### Can't Connect to MySQL Socket

If MySQL is running but you get socket errors:

1. Check MySQL is running:
   ```bash
   mysqladmin ping -u root
   ```

2. Restart MySQL:
   ```bash
   brew services restart mysql  # macOS
   sudo systemctl restart mysql # Linux
   ```

### Port Already in Use

If port 3306 is already in use by another application:

1. Change the port in `.env`:
   ```env
   DB_PORT=3307
   ```

2. Update your MySQL connection to use the new port

---

## Verifying the Setup

### Test SQLite Setup

```bash
cd yad-back
go run ./cmd/main.go
```

You should see:
```
⚠ MySQL connection failed: dial tcp [::1]:3306: connect: connection refused
ℹ Falling back to SQLite for local development...
✓ Connected to SQLite database (dev mode)
✓ Starting YAD API server on port 8080
```

### Test MySQL Setup

```bash
# Check MySQL is running
mysqladmin ping -u root -p

# Check database exists
mysql -u root -p -e "SHOW DATABASES LIKE 'yad_db';"

# Run application
go run ./cmd/main.go
```

You should see:
```
✓ Connected to MySQL database
✓ Starting YAD API server on port 8080
```

### Test API

In another terminal:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "time": "2025-11-10T15:40:17.891808588Z"
}
```

---

## Database Migrations

All database schema migrations are handled automatically when the application starts.

### Manual Migrations

If needed, you can manually trigger migrations:

```bash
# Via Go code
go run ./cmd/main.go
```

The `RunMigrations` function in `internal/database/migrations.go` will:
- Create all tables if they don't exist
- Update schema if models change
- Handle relationships and constraints

### Viewing Tables

**SQLite:**
```bash
sqlite3 yad.db ".tables"
```

**MySQL:**
```bash
mysql -u root -p yad_db -e "SHOW TABLES;"
```

---

## Production Deployment

For production:

1. **Use MySQL** (not SQLite)
2. **Use Docker**: See Option 3 above
3. **Set environment variables**:
   ```bash
   export DB_HOST=your-db-host
   export DB_PORT=3306
   export DB_USER=your-username
   export DB_PASSWORD=your-password
   export DB_NAME=yad_db
   export ENV=production
   export JWT_SECRET=your-secure-secret
   ```

4. **Run with:**
   ```bash
   go build -o bin/yad-back ./cmd/main.go
   ./bin/yad-back
   ```

---

## Quick Reference

| Task | Command |
|------|---------|
| Run with SQLite (no setup) | `go run ./cmd/main.go` |
| Set up MySQL locally | `./setup-db.sh` |
| Run with local MySQL | `go run ./cmd/main.go` |
| Run with Docker | `docker-compose up` |
| Test API | `curl http://localhost:8080/health` |
| View SQLite tables | `sqlite3 yad.db ".tables"` |
| View MySQL tables | `mysql -u root yad_db -e "SHOW TABLES;"` |
| Clean database (SQLite) | `rm yad.db` |

---

## Need Help?

1. **Quick start issues**: Check that MySQL is not running and let SQLite handle it
2. **MySQL connection issues**: Verify MySQL is installed and running
3. **Port conflicts**: Change the port in `.env`
4. **Permissions issues**: Run with `sudo` or adjust file permissions

For more details, see the main README.md in the yad-back directory.
