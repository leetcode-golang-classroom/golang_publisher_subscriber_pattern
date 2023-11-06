.PHONY=build_pub_sub

build_pub_sub:
	@go build -o bin/main pub-sub/main.go

run_pub_sub: build_pub_sub
	@./bin/main