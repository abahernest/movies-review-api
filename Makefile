run-dev:
	export APP_ENV="dev" && \
 	cd ./cmd && \
 	go run -race main.go

run-prod:
	export APP_ENV="prod" && \
	cd ./cmd && \
 	go run -race main.go

build-dev:
	docker build -f Dockerfile-dev --tag abahernest/movies-review-api:dev-$(version) --platform=linux/amd64 .

build-prod:
	docker build --tag abahernest/movies-review-api:$(version) --platform=linux/amd64 .

test:
	go test -failfast -race -v ./...

test-all-pkg:
	cd ./pkg && go test -failfast -race -v ./...

# test a single package
test-pkg:
	cd $(pkg) && go test -failfast -race -v ./...
	# usage: make test-pkg pkg=./pkg/logger

# test single function within a package
test-pkg-fxn:
	cd $(pkg) && go test -failfast -race -v -run $(fxn) ./...
    # usage: make test-pkg-fxn pkg=./pkg/logger fxn=TestInitLogger