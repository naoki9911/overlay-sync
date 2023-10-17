LDFLAGS ?= -s -w -extldflags '-static'
overlay-sync: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=$(shell dpkg --print-architecture) go build -ldflags "${LDFLAGS}" -o overlay-sync main.go
