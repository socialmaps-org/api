.PHONY: bin/socialmaps-api internal/model

bin/socialmaps-api: internal/model
	go build -tags sqlite_math_functions -o bin/socialmaps-api cmd/socialmaps-api/main.go

internal/model:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@0e3f5404b0e69e3402bd81db0aaee2d7275a8785
	sqlc generate
