# PostgreSQL Query Optimizer Toolkit (Phase 1) version 0.1.1

## Project Overview

The PostgreSQL Query Optimizer Toolkit is a CLI tool designed to help database administrators and developers analyze PostgreSQL query execution plans. This initial phase focuses on foundational query analysis capabilities using PostgreSQL's `EXPLAIN ANALYZE` functionality.

## Phase 1 Features

- **Query Execution Analysis**: Runs `EXPLAIN ANALYZE` on provided queries
- **Plan Parsing**: Converts JSON query plans into structured data
- **Basic Bottleneck Detection**: Identifies potentially slow operations
- **Configuration Management**: Secure handling of database credentials

## Technology Stack

- **Language**: Go 1.18+
- **Database Driver**: github.com/jackc/pgx/v5
- **CLI Framework**: github.com/spf13/cobra
- **Configuration**: github.com/spf13/viper

## Project Setup

### Prerequisites

1. Go 1.18 or later
2. PostgreSQL 12+ (with permissions to run EXPLAIN ANALYZE)
3. Git (for version control)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/pg-opt-toolkit.git
cd pg-opt-toolkit

# Install dependencies
go mod download

# Build the tool
go build -o pgopt cmd/pgopt/main.go

# Configuration in yaml
db:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "" # Can be set via PGPASSWORD environment variable
  database: "your_database"
  sslmode: "prefer"

# Alternatively, set environment variables
export PGOPT_DB_HOST=localhost
export PGOPT_DB_PORT=5432
export PGOPT_DB_USER=postgres
export PGOPT_DB_PASSWORD=yourpassword
export PGOPT_DB_DATABASE=your_database
```

### Usage

```bash
# Analyze a query
./pgopt analyze "SELECT * FROM users WHERE email = 'test@example.com'"

# Example output:
# Query Analysis:
# Execution Time: 15.42 ms
# Planning Time: 2.31 ms
#
# Node Type: Seq Scan
# Relation: users (alias: users)
# Cost: 0.00..25.88
# Actual Time: 0.04..15.32 ms, Rows: 1, Loops: 1
#
# Potential Bottlenecks:
# - Slow operation: Seq Scan (15.32 ms)
```

### Command line Options

```bash
# Show help
./pgopt --help

# Analyze command help
./pgopt analyze --help
```
