.PHONY: bin/socialmaps-api bin/socialmaps-auth

all: bin/socialmaps-api bin/socialmaps-auth

bin/socialmaps-api:
	go build -tags sqlite_math_functions -o bin/socialmaps-api cmd/socialmaps-api/main.go

bin/socialmaps-auth:
	go build -tags sqlite_math_functions -o bin/socialmaps-auth cmd/socialmaps-auth/main.go

