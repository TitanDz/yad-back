#!/bin/bash

# YAD Backend Database Setup Script
# This script sets up MySQL for local development

set -e

echo "================================"
echo "YAD Backend - Database Setup"
echo "================================"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if MySQL is installed
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}✗ MySQL is not installed${NC}"
    echo ""
    echo "Please install MySQL first:"
    echo "  - macOS (Homebrew): brew install mysql"
    echo "  - Linux (Ubuntu/Debian): sudo apt-get install mysql-server"
    echo "  - Windows: Download from https://dev.mysql.com/downloads/mysql/"
    exit 1
fi

# Check if MySQL is running
if ! mysqladmin ping -h localhost &> /dev/null; then
    echo -e "${YELLOW}⚠ MySQL is not running${NC}"
    echo ""
    echo "Starting MySQL..."
    
    # Try to start MySQL (different methods for different systems)
    if command -v brew &> /dev/null; then
        brew services start mysql
    elif command -v systemctl &> /dev/null; then
        sudo systemctl start mysql
    else
        echo -e "${RED}Could not automatically start MySQL${NC}"
        echo "Please start MySQL manually and run this script again"
        exit 1
    fi
    
    # Wait for MySQL to be ready
    sleep 3
fi

echo -e "${GREEN}✓ MySQL is running${NC}"

# Create database and user
echo ""
echo "Setting up database and user..."

# Read credentials from .env or use defaults
DB_USER=${DB_USER:-root}
DB_PASSWORD=${DB_PASSWORD:-yad_dev_password}
DB_NAME=${DB_NAME:-yad_db}

# Create database and user
mysql -u root <<EOF
-- Create database
CREATE DATABASE IF NOT EXISTS ${DB_NAME};

-- Create user with password
CREATE USER IF NOT EXISTS '${DB_USER}'@'localhost' IDENTIFIED BY '${DB_PASSWORD}';

-- Grant privileges
GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'localhost';
FLUSH PRIVILEGES;

-- Verify setup
SELECT 'Database setup complete!' as Status;
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Database '${DB_NAME}' created successfully${NC}"
    echo -e "${GREEN}✓ User '${DB_USER}' created with password${NC}"
else
    echo -e "${RED}✗ Failed to setup database${NC}"
    exit 1
fi

echo ""
echo "================================"
echo -e "${GREEN}✓ Setup Complete!${NC}"
echo "================================"
echo ""
echo "You can now run the backend:"
echo "  cd yad-back"
echo "  go run ./cmd/main.go"
echo ""
echo "Or build and run:"
echo "  go build -o bin/yad-back ./cmd/main.go"
echo "  ./bin/yad-back"
echo ""
