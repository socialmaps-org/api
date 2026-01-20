.PHONY: bin/socialmaps-api internal/model

bin/socialmaps-api: internal/model
	go build -tags sqlite_math_functions -o bin/socialmaps-api cmd/socialmaps-api/main.go

internal/model:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
	sqlc generate
