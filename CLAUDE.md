# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

QuickQuery is a lightweight Go library for testing PostgreSQL database interactions. It wraps pgx/v4 and testify to provide convenient helper methods for executing SQL queries and asserting results during tests.

## Build and Test Commands

```bash
go build ./...      # Build the package
go test ./...       # Run tests
go mod tidy         # Clean up dependencies
```

## Architecture

The library has two main types:

### TestDB (`test_db.go`)
Database pool wrapper that validates test database naming (requires `_test` suffix) and manages pgxpool.Pool connections. Factory for creating QuickQuery instances.

### QuickQuery (`quick_query.go`)
Test helper that provides:
- Automatic transaction creation with rollback on `Done()` for test isolation
- Context with 10-second default timeout
- Integrated testify assertions (`A` for assert, `R` for require)
- Type-safe query methods: `Bool()`, `Int()`, `Int64()`, `String()`, `Time()`, `UUID()`, `IsNull()`
- `Exec()` for non-query statements
- `Throws(sql, pgErrCode)` for PostgreSQL SQLSTATE error validation

### Usage Pattern

```go
db := quickquery.NewTestDB(t, connString)  // Validates _test suffix
defer db.TearDown()

qq := db.QQ(t)
defer qq.Done()  // Rollback + context cancel

result := qq.String("SELECT name FROM users WHERE id = $1", 1)
```

## Key Dependencies

- `github.com/jackc/pgx/v4` - PostgreSQL driver
- `github.com/stretchr/testify` - Assertions
