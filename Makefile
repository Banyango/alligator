clean:
	@ rm -dr .dist || true

create-dist:
	@ mkdir .dist

run-server: clean create-dist build
	$(info Running Server on PORT 8080)
	@ cd ./.dist; addr:8080; ./alligator

gofmt:
	@ go fmt .

gotest:
	@ go test ./...

build: clean create-dist
	$(info Building...)
	@ cp alligator.toml.sample ./.dist/alligator.toml
	@ go build -o ./.dist/alligator main.go


build-docker: gotest
	test $(v)
	@ docker build -t alligator:$(v) .

build-linux: clean create-dist
	@ touch ./.dist/alligator.toml
	@ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./.dist/alligator .