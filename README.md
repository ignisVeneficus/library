# 📚 Library

**Library** is a lightweight 📖 web application built *exclusively* for cataloging and organizing your local ebook collection. It is **not** an online ebook reader — instead, it is designed as a central repository from which various ebook readers can download ebooks.

Think of it as a personal ebook server: it scans your local collection, extracts metadata, and makes it accessible through a clean, searchable web interface.

---

## ✨ Features

- 🔍 Scans directories for ebooks (`.epub`, `.mobi`)
- 🏷️ Extracts and categorizes metadata (title, author, etc.)
- 🌐 Downloads missing metadata from external sources (currently [moly.hu](https://moly.hu))
- 🧭 Searchable web interface
- 💾 Stores data in a database (local, or remote) (mysql or mariadb)
- 📤 Designed to serve ebooks to external ebook readers (e.g. via download links)

## ⚠️ Disclaimer

This application does **not implement any user management, authentication, or security features**.  
It is designed to run:

- locally on a trusted machine, **or**
- behind a secure reverse proxy with authentication (e.g., [Authelia](https://www.authelia.com/), OAuth2 proxy, etc.)

**Do not expose this application directly to the internet without proper protection.**

## 🚀 Getting Started

### Requirements

- ✅ Go 1.21+
- ✅ Mysql/mariadb database
- ✅ Local folder with ebooks
- ✅ (Optional) SQLite database

### Installation

```bash
git clone https://github.com/ignisVeneficus/library.git
cd library
go build -o library
```
Set up the environment variables
than
```bash
./library
```
Then visit 👉 http://localhost:8080

⚙️ Configuration

Set environment variables:

The application uses the following environment variables for configuration:

| Variable                   | Description                          | Required |
|----------------------------|--------------------------------------|----------|
| `LIBRARY_DB_USERNAME`      | Database username                    | ✅       |
| `LIBRARY_DB_PASSWORD`      | Database password                    | ✅       |
| `LIBRARY_DB_HOST`          | Database host URL                    | ✅       |
| `LIBRARY_DB_DATABASE`      | Database name                        | ✅       |
| `LIBRARY_BOOKS`            | Path to the directory with ebooks    | ✅       |
| `LIBRARY_COVERS`           | Path to the directory for cover images, the covers will extracted | ✅     |

Example usage in shell:

```bash
export LIBRARY_DB_USERNAME=myuser
export LIBRARY_DB_PASSWORD=secret
export LIBRARY_DB_HOST=localhost:5432
export LIBRARY_DB_DATABASE=library
export LIBRARY_BOOKS=/path/to/ebooks
export LIBRARY_COVERS=/path/to/covers
./library
```
### Command-line Flags

| Flag                  | Attribute | Description                                          |
|-----------------------|-----------|------------------------------------------------------|
| `--forceUpdate`, `-fu`| *none* | Force update: force re-read all book, extract covers | *none* |
| `--noServer`, `-ns`   | *none* |Not start the server (for maintainne)       |
| `--noCheck`,`-nc`     | *none* |No eBooks check at start                    |
| `--export`, `-e`      | *export file name* | Export database in json format |

## 🪵 Logging

This application uses [zerolog](https://github.com/rs/zerolog) for structured, high-performance logging.

Configuration is managed via [zeroconfig](https://pkg.go.dev/go.mau.fi/zeroconfig), allowing runtime log level and format customization through a config file.

### 🔧 Configuration

Logging settings are read from a file named `log.config` located in the working directory.  
It supports both JSON and plain text output, log levels, and other zerolog features.

Example `log.config` (already part of the code):

```yaml
min_level: info
#min_level: trace
max_level: fatal
caller: false
metadata: null
writers:
- type: file
  filename: server.log
  max_size: 10
  max_age: 0
  max_backups: 1
  local_time: false
  compress: false
- type: stdout
  format: pretty-colored
  time_format: 2006-01-02 15:04:05
  min_level: trace
  max_level: fatal
```
For full configuration options, see the [zeroconfig documentation](https://github.com/tulir/zeroconfig).

## 🛣️ Roadmap
 - ☐ Edit Author
 - ☐ Edit Series
 - ☐ Edit Tag
 - ☐ Maintaince tasks (delete orphan authors, tags, series, etc )
 - ☐ UI redesign for smaller screens
 - ☐ Import 
 - ☐ Docker container
 - ☐ Possibility to change UI to different webapp

## 📝 License
This project is licensed under the MIT License.
See the LICENSE file for details.

