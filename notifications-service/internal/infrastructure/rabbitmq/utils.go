package rabbitmq

type UserCreatedMessage struct {
	Email    string
	Username string
	Code     string
}

const (
	UserCreatedTopic = "user.created"
)
