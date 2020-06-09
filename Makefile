build:
	go build -o bin/persist-post cmd/persist-post.go

run:
	go run cmd/persist-post.go

compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/persist-post cmd/persist-post.go