.PHONEY: clean get

VERSION=`git describe --tags`
BUILD=`git rev-parse HEAD`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

build:
	 if [ ! -d "./bin/" ]; then mkdir -p ./bin/server && mkdir -p ./bin/client; fi
	 env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -v -o ./bin/server/426c-server ./server/
	 env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -v -o ./bin/client/426c ./client/
get:
	 go get -d ./client/
	 go get -d ./server/
clean:
	go clean
install:
ifneq ($(shell uname),Linux)
	echo "Install only available on Linux"
	exit 1
endif
	cp ./bin/server/server /usr/local/sbin/
	if [ ! -d "/etc/426c/" ]; then mkdir /etc/426c/; fi
	cp 426c.service /etc/systemd/system/
	systemctl daemon-reload