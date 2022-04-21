.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: postgres
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=Qwerty123 -d postgres

.PHONY: redis
redis: 
	docker run --name redis -p 6379:6379 -d redis


.PHONY: docker
docker:
	docker image build . -t image
	docker container run -p 9090:8080 -d --name rest-api image

clean:
	rm apiserver
.DEFAULT_GOAL := build