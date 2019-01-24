MODULE=qiita_tile38

all: build 

build: ./*.go
	go build -o $(MODULE)
	
test:
	go test . -v

