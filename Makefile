REGISTRY=ghcr.io/elwin
REPOSITORY=wiki-index
VERSION=latest

build:
	docker build . -t $(REPOSITORY)

run: build
	docker run -p 8080:80 $(REPOSITORY)

publish: build
	docker tag $(REPOSITORY) $(REGISTRY)/$(REPOSITORY):$(VERSION)
	docker push $(REGISTRY)/$(REPOSITORY):$(VERSION)

