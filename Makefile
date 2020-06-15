compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/persist-post cmd/persist-post/persist-post.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/get-posts-for-user cmd/get-posts-for-user/get-posts-for-user.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/delete-post cmd/delete-post/delete-post.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/update-post cmd/update-post/update-post.go