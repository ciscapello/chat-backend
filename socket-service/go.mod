module github.com/ciscapello/socket-service

go 1.22.3

require (
	github.com/ciscapello/chat-lib v0.0.0-20240719151044-765b77492ffa
	github.com/gorilla/websocket v1.5.1
	github.com/joho/godotenv v1.5.1
	go.uber.org/zap v1.27.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.25.0 // indirect
)

replace github.com/ciscapello/lib => ../lib
