# roe

Creating a local RSS reader from scratch.

<!-- TOC -->
* [roe](#roe)
  * [Prerequisites](#prerequisites)
    * [1. Install runtimes](#1-install-runtimes)
    * [Create a `.env` file with the following environment variables](#create-a-env-file-with-the-following-environment-variables)
    * [2. Install Go dependencies](#2-install-go-dependencies)
    * [3. Install Node dependencies](#3-install-node-dependencies)
    * [4. Create a PostgreSQL local database using Docker](#4-create-a-postgresql-local-database-using-docker)
    * [5. Create a database schema](#5-create-a-database-schema)
  * [Available Tasks](#available-tasks)
    * [🚀 Run the application in dev mode](#-run-the-application-in-dev-mode)
    * [🧪 Run tests](#-run-tests)
    * [🎨 Format the code](#-format-the-code)
    * [🔍 Lint the codebase](#-lint-the-codebase)
    * [⬆️ Update and tidy dependencies](#-update-and-tidy-dependencies)
    * [🔧 Build the Go binary](#-build-the-go-binary)
    * [🚀 Build and run the application](#-build-and-run-the-application)
    * [🧹 Clean build artifacts](#-clean-build-artifacts)
<!-- TOC -->

## Prerequisites

- [Go 1.25](https://go.dev)
- [golangci-lint](https://golangci-lint.run/) for linting
- [Task](https://taskfile.dev/) for running tasks
- [Node 22](https://nodejs.org/en)
- [pnpm](https://pnpm.io/)
- [Docker](https://www.docker.com/)

### 1. Install runtimes

**Using MacOS and `brew`**

```bash
brew install go@1.25 golangci-lint go-task node@22 pnpm
brew install --cask docker-desktop
```

### Create a `.env` file with the following environment variables

```dotenv
DATABASE_USER=user
DATABASE_PASSWORD=password
DATABASE_HOST=localhost
DATABASE_NAME=testdb

# CORS - optional, use the URL and port of the frontend
CORS_ALLOWED_ORIGINS=http://localhost:3000
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

Use the file `internal/sql/schema.sql` to create a schema in your PostgreSQL database.

## Available Tasks

Below are the available tasks defined in [`Taskfile.yml`](./Taskfile.yml):

### 🚀 Run the application in dev mode

```sh
task dev
```

### 🧪 Run tests

```sh
task test
```

### 🎨 Format the code

```sh
task fmt
```

### 🔍 Lint the codebase

```sh
task lint
```

### ⬆️ Update and tidy dependencies

```sh
task update-deps
```

### 🔧 Build the Go binary

```sh
task build
```

### 🚀 Build and run the application

```sh
task run
```

### 🧹 Clean build artifacts

```sh
task clean
```
