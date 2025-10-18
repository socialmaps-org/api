.PHONY: bin/socialmaps-api

bin/socialmaps-api:
	go build -tags sqlite_math_functions -o bin/socialmaps-api cmd/socialmaps-api/main.go
