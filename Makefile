build:
	go build -o bin/goweb *.go

run-prod:
	./bin/goweb --prod=true
