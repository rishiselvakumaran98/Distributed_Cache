build:
	go build -o bin/ggcache ./main

run: build
	./bin/ggcache