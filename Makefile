clean:
	@ rm -dr .dist || true

create-dist:
	@ mkdir .dist

run-server: clean create-dist
	$(info Running Server on PORT 8080)
	@ PORT=8080; go run main.go

gofmt:
	@ go fmt .
