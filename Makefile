all:
	go build -o bin/chat-demo ./main.go
	go build -o bin/hello ./examples/hello/main.go
	go build -o bin/prompt ./examples/prompt/main.go
	go build -o bin/stdin ./examples/stdin/main.go
	go build -o bin/chat ./examples/chat/main.go
