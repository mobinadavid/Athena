package drivers

import (
	"crypto/tls"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Username   string
	Password   string
	Host       string
	Port       string
	connection *amqp.Connection // Store connection as part of struct
}

// Connect establishes a new connection to RabbitMQ.
// It configures the connection with TLS settings for enhanced security.
func (rabbitmq *RabbitMQ) Connect() error {
	var err error
	rabbitmq.connection, err = amqp.DialTLS(
		fmt.Sprintf("amqps://%s:%s@%s:%s",
			rabbitmq.Username, rabbitmq.Password, rabbitmq.Host, rabbitmq.Port,
		),
		&tls.Config{MinVersion: tls.VersionTLS13},
	)
	return err // Directly return error without redundant check
}

// Close terminates the connection to RabbitMQ.
// It handles errors that may occur during the close operation.
func (rabbitmq *RabbitMQ) Close() error {
	if rabbitmq.connection == nil {
		return nil // Avoid attempting to close a nil connection
	}
	return rabbitmq.connection.Close()
}

// GetClient returns the current RabbitMQ connection.
func (rabbitmq *RabbitMQ) GetClient() *amqp.Connection {
	return rabbitmq.connection
}
