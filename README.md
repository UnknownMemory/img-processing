# Image Processing Service
## Getting Started
### Prerequisites
- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Goose](https://github.com/pressly/goose)
- [sqlc](https://github.com/sqlc-dev/sqlc)
- [Air](https://github.com/air-verse/air) (for live reloading)
- [Bruno](https://www.usebruno.com/) (optional)

### Local Development Setup
```shell
git clone https://github.com/UnknownMemory/img-processing
cd img-processing

# Install dependencies
go mod download
```
Configure the ``garage.toml`` file using ``garage.toml.example`` as an example.
``rpc_secret`` is a random 32-byte hex-encoded string.

Configure the ``.env`` file using ``.env.example`` as an example.

```shell
# Run the services
docker-compose up

# Setup goose environment variables (use goose.env.ps1 on Windows)
./goose.env.sh

# Apply migrations
goose up

# Run the service
air -c .\.air.toml

# Run the worker
air -c .\.air.worker.toml
```

### Build
```shell
docker-compose -f compose.prod.yaml up -d --build --no-deps --force-recreate
```