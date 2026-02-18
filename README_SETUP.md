# MIC Website Backend - Setup Instructions

## Prerequisites

- Go 1.25.5 or higher
- PostgreSQL 12 or higher
- psql (PostgreSQL command-line tool)

## Database Setup

### 1. Create a PostgreSQL database

You have two options:

#### Option A: Using psql command line
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE startups;

# Exit psql
\q
```

#### Option B: Using the setup script
```bash
# First, create a .env file (see below)
# Then run the setup script
chmod +x setup_db.sh
./setup_db.sh
```

### 2. Configure Environment Variables

Create a `.env` file in the root directory:

```bash
cp configs/config.example.env .env
```

Edit `.env` and update the `DATABASE_URL` with your PostgreSQL credentials:

```
DATABASE_URL=postgres://username:password@localhost:5432/startups?sslmode=disable
```

Replace:
- `username` with your PostgreSQL username
- `password` with your PostgreSQL password
- `localhost:5432` with your PostgreSQL host and port (if different)
- `startups` with your database name (if different)

### 3. Run Migrations

If you haven't used the setup script, run the migrations manually:

```bash
# Source your .env file
source .env  # or export DATABASE_URL="your_connection_string"

# Run migrations
psql $DATABASE_URL -f migrations/0001_startups.up.sql
psql $DATABASE_URL -f migrations/0002_users.up.sql
```

### 4. Install Go Dependencies

```bash
go mod download
```

### 5. Run the Server

```bash
# Make sure DATABASE_URL is set
export DATABASE_URL="postgres://user:pass@localhost:5432/startups?sslmode=disable"

# Run the server
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## Testing the Integration

1. Open your browser and navigate to `http://localhost:8080/login.html`
2. Try signing up with an email ending in:
   - `@learner.manipal.edu` (will create a STUDENT role)
   - `@manipal.edu` (will create a FACULTY role)
3. After signup, try logging in with the same credentials

## API Endpoints

- `POST /api/signup` - Create a new user account
  - Body: `{"name": "Full Name", "email": "user@learner.manipal.edu", "password": "password123"}`
  
- `POST /api/login` - Login and get JWT token
  - Body: `{"email": "user@learner.manipal.edu", "password": "password123"}`
  - Returns: `{"token": "jwt_token_here"}`

## Troubleshooting

### Database connection errors
- Verify PostgreSQL is running: `pg_isready`
- Check your DATABASE_URL is correct
- Ensure the database exists: `psql -l | grep startups`

### Migration errors
- Make sure you run migrations in order (0001, then 0002)
- Check if tables already exist: `psql $DATABASE_URL -c "\dt"`

### Port already in use
- Change the port in `cmd/server/main.go` or stop the process using port 8080

