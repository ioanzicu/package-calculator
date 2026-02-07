# 📦 Package Calculator

### Features
- **Optimal Distribution**: Uses a recursive exact-match algorithm to find the perfect distribution.
- **Persistence**: Automatically stores every calculation in a PostgreSQL database.
- **History View**: Dedicated UI section to browse past calculations, powered by HTMX triggers.
- **Dynamic UI**: Seamless experience without page reloads, using custom HTMX events for real-time history updates.
- **Testable**: Includes comprehensive unit tests, benchmark-ready algorithms, and E2E integration tests.

### How to Use
1.  **Start the database**: Run `make compose-up` (requires Docker) to start PostgreSQL and migrations.
2.  **Start the application**: Run `make run` or `go run cmd/main.go`.
3.  **Navigate to the calculator**: Open your browser at `http://localhost:8080`.
4.  **Enter Pack Sizes**: Provide comma-separated values (e.g., `23, 31, 53`).
5.  **Enter Amount**: Specify the total number of items you need to pack (e.g., `500000`).
6.  **Click Calculate**: See the results instantly!

### Screenshot
![Package Calculator](./main_page.png)

- **Input**: Pack Sizes [23, 31, 53], Amount 500,000
- **Optimal Result** (Sorted Descending):
    - 53: 9,429 packs
    - 31: 7 packs
    - 23: 2 packs
- **Total**: 500,000 (Exact Match!)


## 🛠️ Tech Stack & Dependencies

| Component | Technology | Version | Role | Source |
| :--- | :--- | :--- | :--- | :--- |
| **Language** | [Go](https://go.dev/) | 1.25.4 | Core programming language | [Official Site](https://go.dev/) |
| **Database** | [PostgreSQL](https://www.postgresql.org/) | 18.1 (Alpine) | Relational database system | [Official Site](https://www.postgresql.org/) |
| **Driver** | [pgx](https://github.com/jackc/pgx) | v5.8.0 | PostgreSQL driver and toolkit | [GitHub](https://github.com/jackc/pgx) |
| **ORM/Gen** | [sqlc](https://sqlc.dev/) | v1.30.0 | Type-safe SQL compiler | [GitHub](https://github.com/sqlc-dev/sqlc) |
| **Config** | [godotenv](https://github.com/joho/godotenv) | v1.5.1 | Environment variable loader | [GitHub](https://github.com/joho/godotenv) |
| **Migration** | [Goose](https://pressly.github.io/goose/) | v3.26.0 | Database migration tool | [GitHub](https://github.com/pressly/goose) |

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

-   **Go**: Version 1.25.4 or higher
-   **Docker**: For containerized execution
-   **Make**: For running project commands

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone git@github.com:ioanzicu/package-calculator.git
cd package-calculator
```

### 2. Configuration

The application uses environment variables for configuration.

**Local Development**: Create a `.env` file in the root directory:

```env
HTTP_HOST=127.0.0.1
HTTP_PORT=8080
PG_DSN=postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable
```

**Docker**: Variables are automatically handled in `docker-compose.yml`.

### 3. Running with Docker (Recommended)

Start the application and database with a single command:

```bash
make compose-up
```

This will: 
- Start Postgres 18.1
- Build and start the `app` service
- Expose the API on port `8080`

To stop the services:
```bash
make compose-down
```

To view logs:
```bash
make compose-logs
```

### 4. Running Locally

Install dependencies and tools:
```bash
make deps
```

Generate database code (if queries change):
```bash
make generate
```

Run the application:
```bash
make run
```

## 🛠️ Make Commands

The project includes a `Makefile` to simplify common tasks:

| Command | Description |
| :--- | :--- |
| `make help` | Display available commands |
| `make deps` | Install Go dependencies and tools (`sqlc`, `goose`) |
| `make run` | Run the application locally |
| `make build` | Build the application binary |
| `make test` | Run unit tests |
| `make generate` | Generate Go code from SQL using `sqlc` |
| `make migration-create` | Create a new database migration |
| `make migration-up` | Apply database migrations |
| `make compose-up` | Start Docker environment |
| `make compose-down` | Stop Docker environment |

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

Built with [🔥 Ignis](https://github.com/ioanzicu/ignis) ❤️