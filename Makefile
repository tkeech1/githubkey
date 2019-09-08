init:
	go mod init github.com/tkeech1/githubkey

tidy:
	go mod tidy

clean-testcache:
	go clean -testcache github.com/tkeech1/githubkey

test: clean-testcache	
	go test -race -v -covermode=atomic -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt

test-circleci: 
	go test -race -covermode=atomic -coverprofile=coverage.txt ./...

deps: 
	go get -v -t -d ./...
