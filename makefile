

.PHONY: build
build: gocron node

.PHONY: build-race
build-race: enable-race build

.PHONY: run
run: build kill
	./bin/gocron-node &
	./bin/gocron web

.PHONY: run-race
run-race: enable-race run

.PHONY: kill
kill:
	-killall gocron-node

.PHONY: gocron
gocron:
	go build $(RACE) -o bin/gocron ./cmd/gocron

.PHONY: node
node:
	go build $(RACE) -o bin/gocron-node ./cmd/node

.PHONY: test
test:
	go test $(RACE) ./...

.PHONY: test-race
test-race: enable-race test

.PHONY: enable-race
enable-race:
	$(eval RACE = -race)

.PHONY: clean
clean:
	rm bin/gocron
	rm bin/gocron-node
