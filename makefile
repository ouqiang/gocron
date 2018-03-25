
default: build

.PHONY: build
build: gocron node

.PHONY: gocron
gocron:

	go build -o bin/gocron ./cmd/gocron


.PHONY: node
node:

	go build -o bin/gocron-node ./cmd/node

.PHONY: clean
clean:

	rm bin/gocron
	rm bin/gocron-node