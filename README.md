# Go Sales API

This is a robust backend API for a sales system, built with Go, Gin, and GORM, following clean architecture principles.

## Features

- **Layered Architecture**: Clear separation between handlers, services, and repositories.
- **RESTful API**: Well-defined REST endpoints for managing entities.
- **Database Migrations**: Schema management using `golang-migrate` for safe and versioned database changes.
- **Structured Logging**: JSON-formatted logs using `zerolog` for efficient monitoring.
- **Configuration Management**: Environment-based configuration using Viper.
- **Dependency Injection**: Loosely coupled components for better testability and maintainability.

---

##  Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.18 or higher.
- **PostgreSQL**: A running instance of PostgreSQL.
- **golang-migrate**: The migration tool. [Installation Instructions](#step-2-install-migration-tool).
- **Make**: The build automation tool (usually pre-installed on Linux/macOS).

---

## 🚀 Getting Started

Follow these steps to get your development environment set up and running.

### Step 1: Clone the Repository

```bash
git clone https://github.com/herculanocm/go-sales.git
cd go-sales
```

### Step 2: Install Migration Tool

If you haven't already, install `golang-migrate` which is essential for managing the database schema.

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

> **Note:** If you get a `command not found: migrate` error after installation, you need to add Go's binary directory to your system's PATH. Add the following line to your `~/.bashrc`, `~/.zshrc`, or equivalent shell profile file, then restart your terminal.
> ```bash
> export PATH=$PATH:$(go env GOPATH)/bin
> ```

### Step 3: Configure Environment Variables

The application uses an `.env` file for configuration.

1.  Copy the example file:
    ```bash
    cp .env.example .env
    ```
2.  Edit the `.env` file and update the `DB_URL` and other variables to match your local PostgreSQL setup.

    ```properties
    # .env
    APP_ENV=development
    APP_API_PREFIX=/api/v1
    APP_API_PORT=8081

    # Update with your PostgreSQL connection details
    DB_URL="postgres://YOUR_USER:YOUR_PASSWORD@localhost:5432/YOUR_DB_NAME?sslmode=disable&search_path=master"
    DB_NAME=YOUR_DB_NAME
    DEFAULT_SCHEMA=master
    ```

### Step 4: Set Up the Database

With the environment configured, you can now create the database schema by running the migrations.

The project includes a `Makefile` to simplify common commands.

```bash
make migrate-up
```

This command reads your `DB_URL` from the `.env` file and applies all pending migrations from the `internal/database/migrations` directory.

### Step 5: Install Dependencies and Run the Project

Finally, install the Go module dependencies and run the application.

1.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

2.  **Run the Project:**
    ```bash
    go run main.go
    ```

The API server should now be running on the port specified in your `.env` file (e.g., `http://localhost:8081`).

---

## Development Workflow

### Managing Database Migrations

Use the `Makefile` to manage your database schema.

-   **Create a new migration:**
    ```bash
    # Example: make migrate-create name=create_products_table
    make migrate-create name=<your_migration_name>
    ```
    This will create new `.up.sql` and `.down.sql` files in the migrations directory.

-   **Apply all pending migrations:**
    ```bash
    make migrate-up
    ```

-   **Revert the last applied migration:**
    ```bash
    make migrate-down
    ```

-   **Reset the database (DANGER):** This will drop all tables and re-apply all migrations from scratch. Use with caution.
    ```bash
    make migrate-reset
    ```
