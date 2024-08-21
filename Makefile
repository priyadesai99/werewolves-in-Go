build:
	go build -o bin/client client/main.go
	go build -o bin/server server/main.go

clean:
	rm bin/client
	rm bin/server