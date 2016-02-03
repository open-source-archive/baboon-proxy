default: build.test

clean:
		rm -rf build

build.test:
		go test -v ./...

prepare:
		mkdir -p build/linux
		mkdir -p build/osx

build.local: prepare
		godep go build -o build/baboon-proxy -ldflags "-X main.BuildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.GitHash=`git rev-parse HEAD`" 

build.linux: prepare
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -o build/linux/baboon-proxy -ldflags "-X main.BuildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.GitHash=`git rev-parse HEAD`" 

build.osx: prepare
		GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 godep go build -o build/osx/baboon-proxy -ldflags "-X main.BuildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.GitHash=`git rev-parse HEAD`" 
