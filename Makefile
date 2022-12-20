.PHONY: faktory
faktory:
	mkdir -p data
	docker run --rm -it -v $(shell pwd)/data:/var/lib/faktory/db -e "FAKTORY_PASSWORD=password" -p 127.0.0.1:7419:7419 -p 127.0.0.1:7420:7420 contribsys/faktory:latest /faktory -b :7419 -w :7420 -e production

.PHONY: faktory-dev
faktory-dev:
	mkdir -p data
	docker run --rm -it -v $(shell pwd)/data:/var/lib/faktory/db -p 127.0.0.1:7419:7419 -p 127.0.0.1:7420:7420 contribsys/faktory:latest /faktory -b :7419 -w :7420

.PHONY: test
test:
	go test -tags=docker -v ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	test -z $$(gofmt -l .) # This will return non-0 if unsuccessful  run `go fmt ./...` to fix
