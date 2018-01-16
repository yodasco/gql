
test: dep
	go test -v ./...

dep:
	dep ensure