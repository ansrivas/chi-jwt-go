.DEFAULT_GOAL := help

DEP= $(shell command -v dep 2>/dev/null)
VERSION=$(shell git describe --always --long)
PROJECT_NAME := chi-jwt
CLONE_URL:=github.com/ansrivas/chi-jwt
IDENTIFIER= $(VERSION)-$(GOOS)-$(GOARCH)
BUILD_TIME=$(shell date -u +%!F(MISSING)T%!T(MISSING)%!z(MISSING))
LDFLAGS="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

help:          ## Show available options with this Makefile
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY : test crossbuild release
test:          ## Run all the tests
test:
	chmod +x ./test.sh && ./test.sh

clean:         ## Clean the application
clean:
	@go clean -i ./...
	@rm -rf ./{PROJECT_NAME}

build: vendor	clean
	go build -i -v -ldflags $(LDFLAGS) $(FLAGS) $(CLONE_URL)

dep:           ## Go get dep
dep:
	go get -u github.com/golang/dep/cmd/dep

ensure:        ## Run dep ensure.
ensure:
ifndef DEP
	make dep
endif
	dep ensure
	touch vendor

crossbuild: ensure
	mkdir -p build/${PROJECT_NAME}-$(IDENTIFIER)
	make build FLAGS="-o build/${PROJECT_NAME}-$(IDENTIFIER)/${PROJECT_NAME}"
	cd build \
	&& tar cvzf "${PROJECT_NAME}-$(IDENTIFIER).tgz" "${PROJECT_NAME}-$(IDENTIFIER)" \
	&& rm -rf "${PROJECT_NAME}-$(IDENTIFIER)"

release:       ## Create a release build.
release:
	make crossbuild GOOS=linux GOARCH=amd64
	make crossbuild GOOS=linux GOARCH=386
	make crossbuild GOOS=darwin GOARCH=amd64


bench:	       ## Benchmark the code.
bench:
	@go test -o bench.test -cpuprofile cpu.prof -memprofile mem.prof -bench .

prof:          ## Run the profiler.
prof:	bench
	@go tool pprof cpu.prof

prof_svg:      ## Run the profiler and generate image.
prof_svg:	clean	bench
	@echo "Do you have graphviz installed? sudo apt-get install graphviz."
	@go tool pprof -svg bench.test cpu.prof > cpu.svg