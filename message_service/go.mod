module github.com/ciscapello/message_service

go 1.22.3

require (
	github.com/ciscapello/lib v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/rabbitmq/amqp091-go v1.10.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

replace github.com/ciscapello/lib => ../lib
