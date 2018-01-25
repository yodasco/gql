
lint:
	go list ./... | grep -v /vendor/ |  xargs -L1 golint

test: dep
	go test `go list ./... | grep -v /vendor/`

dep:
	dep ensure