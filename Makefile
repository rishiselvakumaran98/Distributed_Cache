build:
	go build -o bin/ggcache ./main

run: build
	./bin/ggcache

runFollower: build
	./bin/ggcache --listenAddr :4000 --leaderAddr :3000