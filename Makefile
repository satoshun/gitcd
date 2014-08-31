BINARY=pkg/darwin/gitcd

deploy: build-darwin
	git add $(BINARY)
	git commit -m "update: $(VERSION)"
	git tag $(VERSION)
	git push origin master
	git push origin $(VERSION)

build-darwin:
	GOARCH=amd64 GOOS=darwin go build -v -o $(BINARY)
