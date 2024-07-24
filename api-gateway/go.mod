module github.com/ciscapello/api-gateway

go 1.22.3

require (
	github.com/ciscapello/chat-lib v0.0.0-20240719151044-765b77492ffa
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.3
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/swaggo/http-swagger v1.3.4
	github.com/swaggo/swag v1.16.3
	go.uber.org/zap v1.27.0
)

require go.uber.org/mock v0.4.0

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/swaggo/files v1.0.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ciscapello/lib => ../lib
