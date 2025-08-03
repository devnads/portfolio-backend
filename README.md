## Technologies

- **Backend**: Go 1.24.5
- **Router**: Chi v5
- **Database**: PostgreSQL with pgx/v5
- **Config**: godotenv

## Installation

### 1. Prerequisites
```bash
# Go 1.24.5+
# PostgreSQL 13+
# Git
```

### 2. Clone the project
```bash
git clone https://github.com/Veenoway/monad-portfolio-backend.git
cd monad-portfolio-backend
```

### 3. Install dependencies
```bash
go mod download
```

### 4. Database Setup

#### Option A: Local PostgreSQL Installation

**Install PostgreSQL:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS (with Homebrew)
brew install postgresql
brew services start postgresql

# Windows - Download from: https://www.postgresql.org/download/windows/
```

**Setup Database:**
```bash
# Switch to postgres user and open psql
sudo -u postgres psql

# Or on macOS/Windows, simply:
psql postgres
```

**Create database and user:**
```sql
-- Create database
CREATE DATABASE monad_dev;

-- Create user with password
CREATE USER monad_user WITH PASSWORD 'your_secure_password';

-- Grant all privileges
GRANT ALL PRIVILEGES ON DATABASE monad_dev TO monad_user;

-- Exit psql
\q
```

#### Option B: Cloud Database (Recommended for beginners)

**Free PostgreSQL hosting options:**
- [Supabase](https://supabase.com) - Free tier with 500MB
- [Neon](https://neon.tech) - Serverless PostgreSQL
- [ElephantSQL](https://www.elephantsql.com) - Free plan with 20MB

1. Sign up for any of these services
2. Create a new PostgreSQL database
3. Copy the connection string provided

### 5. Configuration
Create a `.env` file at the root:

```env
# Database (Local setup)
DB_URL=postgres://monad_user:your_secure_password@localhost:5432/monad_portfolio?sslmode=disable

# Database (Cloud setup example)
# DB_URL=postgres://username:password@your-host.com:5432/dbname?sslmode=require

# Server
SERVER_PORT=8080

# Logs
LOG_LEVEL=info
```

**⚠️ Important**: 
- Replace `your_secure_password` with your actual password
- For cloud databases, use the connection string provided by your service
- Keep your `.env` file secret and never commit it to git

### 6. Verify Database Connection

**Test your connection:**
```bash
# Using psql (replace with your credentials)
psql "postgres://monad_user:your_secure_password@localhost:5432/monad_portfolio"

# Should connect successfully and show:
# monad_dev=#
```

**Common connection issues:**
```bash
# If you get "connection refused":
sudo systemctl start postgresql  # Linux
brew services start postgresql   # macOS

# If you get "authentication failed":
# Check your username/password in .env file

# If you get "database does not exist":
# Make sure you created the database in step 4
```

### 7. Run
```bash
# Development
go run cmd/server/main.go

# Should see:
# ✅ Connected to database successfully
# ✅ Schema dropped and recreated successfully
# Server running on port :8080

# Production
go build -o bin/server cmd/server/main.go
./bin/server
```

### 8. Verify Everything Works
```bash
# Test the API
curl http://localhost:8080
# Should return: "Welcome to Monad Dev Portfolio API"

# Test database endpoints
curl http://localhost:8080/devs
# Should return: [] (empty array)
```

## 🛠️ Troubleshooting Database Issues

### Connection Problems
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql  # Linux
brew services list | grep postgres  # macOS

# Check if port 5432 is open
netstat -an | grep 5432
```


### Quick Setup Script
Create a `setup.sh` file:
```bash
#!/bin/bash
echo "🚀 Setting up Monad Portfolio Backend..."

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found!"
    echo "📝 Please create .env file with your database credentials"
    echo "📖 See README.md for example configuration"
    exit 1
fi

# Run the application
echo "🎯 Starting server..."
go run cmd/server/main.go
```

Make it executable and run:
```bash
chmod +x setup.sh
./setup.sh
```

## 🌐 API Endpoints

### 👨‍💻 Developers

| Method | Endpoint | Description | Auth |
|---------|----------|-------------|------|
| `GET` | `/devs` | List all developers | ❌ |
| `GET` | `/devs?search=john` | Search by name | ❌ |
| `GET` | `/devs?roles=admin` | Filter by role | ❌ |
| `GET` | `/devs?sort_by=username&sort_dir=asc` | Sort | ❌ |
| `GET` | `/devs?limit=10&offset=20` | Pagination | ❌ |
| `POST` | `/devs` | Create a developer | ✅ Admin |

### 📂 Projects

| Method | Endpoint | Description | Auth |
|---------|----------|-------------|------|
| `GET` | `/projects` | List all projects | ❌ |
| `GET` | `/projects?creator_id=123` | Projects by developer | ❌ |
| `GET` | `/projects?categories=DeFi` | Filter by category | ❌ |
| `POST` | `/projects` | Create a project | ✅ Admin |

### 📊 Request Examples

```bash
# Get all developers
curl "http://localhost:8080/devs"

# Search for a developer
curl "http://localhost:8080/devs?search=alice&limit=5"

# Get projects by developer
curl "http://localhost:8080/projects?creator_id=dev-123"

# Create a project (requires admin)
curl -X POST "http://localhost:8080/projects" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MyDapp",
    "description": "Innovative DeFi application",
    "categories": ["DeFi", "Web3"],
    "creator_id": "dev-123"
  }'
```

## 🗄️ Database Schema

### Table `devs`
```sql
CREATE TABLE devs (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  roles TEXT[],                 -- ['admin', 'developer']
  profile_image TEXT,
  address TEXT UNIQUE NOT NULL, -- Wallet address
  twitter TEXT,
  discord TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Table `projects`
```sql
CREATE TABLE projects (
  id TEXT PRIMARY KEY,
  dev_id TEXT REFERENCES devs(id),
  mission_id TEXT REFERENCES missions(id),
  name TEXT NOT NULL,
  image TEXT,
  categories TEXT[],            -- ['DeFi', 'NFT', 'Gaming']
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Table `missions`
```sql
CREATE TABLE missions (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  start_time TIMESTAMPTZ,
  end_time TIMESTAMPTZ,
  round INT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## 🧪 Testing

```bash
# Unit tests
go test ./...

# Tests with coverage
go test -cover ./...

# Integration tests
go test -tags=integration ./...
```

## 🚀 Features

- ✅ **RESTful API** with Chi router
- ✅ **PostgreSQL** integration with pgx
- ✅ **Search & Filtering** with advanced query parameters
- ✅ **Pagination** with limit/offset
- ✅ **Sorting** by multiple fields
- ✅ **Admin Authentication** middleware
- ✅ **Environment Configuration** with .env
- ✅ **Database Migrations** automated
- ✅ **Error Handling** comprehensive

## 🔧 API Features

### Advanced Querying
```bash
# Complex search with multiple filters
GET /devs?search=alice&roles=developer&sort_by=created_at&sort_dir=desc&limit=10&offset=0

# Project filtering by multiple categories
GET /projects?categories=DeFi&creator_id=dev-123&sort_by=name
```

### Response Format
```json
{
  "id": "dev-123",
  "username": "alice_dev",
  "profile_image": "https://...",
  "roles": ["developer", "admin"],
  "address": "0x742d35Cc6665C0532846c88e0662E04a6dd5f81E",
  "discord": "alice#1234",
  "twitter": "@alice_dev",
  "created_at": "2024-01-15T10:30:00Z"
}
```

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go conventions
- Add tests for new features
- Update documentation
- Use meaningful commit messages

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Novee** - *Lead Developer* - [@Veenoway](https://github.com/Veenoway)

## 🙏 Acknowledgments

- [Monad Labs](https://monad.xyz/) for blockchain infrastructure
- Go community for excellent libraries
- Project contributors

## 📞 Support

- Create an [Issue](https://github.com/your-username/monad-portfolio-backend/issues)
- Follow on [Twitter](https://twitter.com/monad_xyz)

---

**Built with ❤️ for the Monad ecosystem** 
