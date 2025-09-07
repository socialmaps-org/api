.PHONY: bin/socialmaps-auth

bin/socialmaps-auth:
	go build -o bin/socialmaps-auth cmd/socialmaps-auth/main.go