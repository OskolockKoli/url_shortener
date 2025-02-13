build:
    go build -o bin/app ./cmd

run-mem:
    PORT=50051 DB_TYPE=memory go run ./cmd/main.go

run-pgsql:
    PORT=50051 DB_TYPE=postgres go run ./cmd/main.go

test:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

docker-build:
    docker build -t url-shortener .

docker-run-mem:
    docker run -it -p 50051:50051 url-shortener -p 50051 -d memory

docker-run-pgsql:
    docker run -it -e POSTGRES_URL="postgres://user:password@host:5432/database" -p 50051:50051 url-shortener -p 50051 -d postgres