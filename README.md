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
- 🐳 Optional Docker support for easy deployment with configurable volumes, user IDs (PUID/PGID), and environment variables

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
Then visit 👉 http://localhost:8888

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

## 📚 API Documentation (Swagger)

This project uses Swagger (OpenAPI) for automated API documentation, powered by [swaggo/swag](https://github.com/swaggo/swag).

### 📁 Local Swagger files
The dopcumentation is already generated at the
```
./docs/swagger.json
```
place

### 🔧 How to generate Swagger docs

If you make changes to your API routes or comments, regenerate the documentation with:
```
swag init
```
This will update the docs/ folder with the latest Swagger spec.
Make sure swag is installed:
```
go install github.com/swaggo/swag/cmd/swag@latest
```
###Where to access the docs online
Once the application is running, the Swagger UI is available at:
http://localhost:8888/swagger/index.html

## 🐳 Usage with Docker

### 📦 Build the image

Clone the repository and build the image using the Dockerfile located in the `docker/` directory:

```bash
git clone https://github.com/ignisVeneficus/library.git
cd library

# Build the Docker image
docker build -f docker/Dockerfile -t library .
```
To always fetch the latest version from GitHub (no build cache):
```bash
docker build --no-cache -f docker/Dockerfile -t library .
```
### 🚀 Run the container
```bash
docker run -d \
  --name library \
  -p 8880:8888 \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Europe/Budapest \
  -e LIBRARY_DB_USERNAME=dbuser \
  -e LIBRARY_DB_PASSWORD=mysecret \
  -e LIBRARY_DB_HOST=mydatabase \
  -e LIBRARY_DB_DATABASE=library \

  -v /host/books:/books \
  -v /host/covers:/covers \
  library
```
### 🧪 Using Docker Compose
```yaml
services:
  library:
    build:
      context: .
      dockerfile: docker/Dockerfile
    container_name: library
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Europe/Budapest
      - "LIBRARY_DB_USERNAME=dbuser"
      - "LIBRARY_DB_PASSWORD=mysecret"
      - "LIBRARY_DB_HOST=mydatabase"
      - "LIBRARY_DB_DATABASE=library"
    volumes:
      - /host/books:/books
      - /host/covers:/covers
    ports:
      - "8888:8888"
    restart: unless-stopped
```

## 🛣️ Roadmap
 - ☐ Edit Author
 - ☐ Edit Series
 - ☐ Edit Tag
 - ☐ Maintaince tasks (delete orphan authors, tags, series, etc )
 - ☐ UI redesign for smaller screens
 - ☐ Import 
 - ☐ Possibility to change UI to different webapp

## 📝 License
This project is licensed under the MIT License.
See the LICENSE file for details.

