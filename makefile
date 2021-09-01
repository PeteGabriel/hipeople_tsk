run:
		mkdir static	
		mkdir static/images	
		go build -o bin/main cmd/api/main.go && ./bin/main -port=4000

test:
		go test ./...