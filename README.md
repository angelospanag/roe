# roe

Creating a local RSS reader from scratch.

- [roe](#roe)
  - [Prerequisites](#prerequisites)
    - [1. Install runtimes](#1-install-runtimes)
    - [2. Install Go dependencies](#2-install-go-dependencies)
    - [3. Install Node dependencies](#3-install-node-dependencies)
    - [4. Create a PostgreSQL local database using Docker](#4-create-a-postgresql-local-database-using-docker)
    - [5. Create a database schema](#5-create-a-database-schema)
  - [Running](#running)
  - [Run backend development server](#run-backend-development-server)
  - [Run frontend development server](#run-frontend-development-server)
  - [Build binary](#build-binary)

## Prerequisites

- [Go 1.24](https://go.dev/)
- [Node 22](https://nodejs.org/en)
- [pnpm](https://pnpm.io/)

### 1. Install runtimes

**Using MacOS and `brew`**

```bash
brew install go@1.24 node@22 pnpm
```

### 2. Install Go dependencies

```bash
go mod tidy
```

### 3. Install Node dependencies

```bash
cd frontend
pnpm install
```

### 4. Create a PostgreSQL local database using Docker

```bash
docker compose up -d db
```

### 5. Create a database schema

Use the file `sql/schema.sql` to create a schema in your PostgreSQL database.

## Running

## Run backend development server

```bash
go run main.go
```

The server will run on http://127.0.0.1:8000

## Run frontend development server

```bash
pnpm run dev
```

The server will run on http://127.0.0.1:5173

## Build binary

```bash
go build
```
