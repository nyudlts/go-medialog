Go-Medialog
===========
A digital media logging and reporting application

## Running the Application

### Prerequisites

- A MySQL/MariaDB database configured and accessible
- A YAML configuration file (see `go-medialog.yml` for an example)

### Configuration

The config file defines environments (e.g., `dev`, `prod`) with database connection details, log location, and admin email:

```yaml
dev:
  log: medialog_dev.log
  port: 8080
  database:
    username: <db_user>
    password: <db_password>
    url: <db_host>
    port: 3306
    database_name: <db_name>
  admin_email: admin@example.com
```

### CLI Flags

| Flag | Type | Description |
|------|------|-------------|
| `--config` | string | Path to the YAML configuration file |
| `--environment` | string | Environment key to use from the config file (e.g., `dev`) |
| `--version` | bool | Print the application version and exit |
| `--prod` | bool | Run in production mode (logs to file instead of stdout) |
| `--migrate` | bool | Run database migrations and exit |
| `--rollback` | bool | Roll back database migrations and exit |
| `--automigrate` | bool | Auto-migrate the database schema and exit |
| `--create-admin` | bool | Create the admin user (email from config) and exit |
| `--create-json` | bool | Export database to JSON and exit |
| `--gorm-debug` | bool | Enable GORM debug logging |

### Common Commands

**Start the server (development):**
```sh
./medialog --config go-medialog.yml --environment dev
```

**Start the server (production):**
```sh
./medialog --config go-medialog.yml --environment prod --prod
```

**Run database migrations:**
```sh
./medialog --config go-medialog.yml --environment dev --migrate
```

**Roll back database migrations:**
```sh
./medialog --config go-medialog.yml --environment dev --rollback
```

**Auto-migrate the database schema:**
```sh
./medialog --config go-medialog.yml --environment dev --automigrate
```

**Create the admin user:**
```sh
./medialog --config go-medialog.yml --environment dev --create-admin
```

**Print the version:**
```sh
./medialog --version
```

The server listens on port **8080** by default.
