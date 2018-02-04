
all: build

build: gocron node

gocron:

	go build -o bin/gocron ./cmd/gocron


node:

	go build -o bin/gocron-node ./cmd/node

clean:

	rm bin/gocron
	rm bin/gocron-node