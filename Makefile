NAME = aheadaviation/users
DBNAME = aheadaviation/users-db
INSTANCE = users
TESTDB = aheadaviationtestusersdb
GROUP = aheadaviationdemos
TEST?=$$(go list ./... |grep -v 'vendor')

default: docker

pre:
	go get -v github.com/golang/dep/cmd/dep

dep: pre
		 dep ensure -v

rm-dep:
	rm -rf vendor

test:
	@docker build -t $(INSTANCE)-test -f ./Dockerfile-test .
	@docker run --rm -it $(INSTANCE)-test /bin/sh -c 'go test -i $(TEST)'

cover:
	go test -v -covermode=count
