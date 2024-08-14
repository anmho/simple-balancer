

all: build

.PHONY: build
build:
	@echo "building load balancer"
	@go build -o ./bin/balancer ./cmd/balancer

.PHONY: clean
clean:
	@rm ./bin/*