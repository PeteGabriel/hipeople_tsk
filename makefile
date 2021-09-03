build:	
		mkdir static	
		mkdir static/images	

run:
		go build -o bin/main cmd/api/main.go && ./bin/main -port=4000 -host=localhost

test:
		go test ./...