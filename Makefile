.PHONEY: clean get

VERSION=`git describe --tags`
BUILD=`git rev-parse HEAD`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

build: linux windows darwin
linux:
	 env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -v -o ./bin/server-linux64/426c-server ./server/
	 env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -v -o ./bin/client-linux64/426c ./client/
windows:
	 env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -v -o ./bin/client-win64/426c.exe ./client/
darwin:
	 env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -v -o ./bin/client-osx/426c ./client/
get:
	 go get -d ./client/
	 go get -d ./server/
clean:
	go clean

