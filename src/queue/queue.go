package queue

import (
	"athena/src/config"
	queueDrivers "athena/src/queue/drivers"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

var (
	connectOnce     sync.Once // Ensures queue connection is established only once.
	getInstanceOnce sync.Once // Ensures a single instance of Queue is created.
	instance        *Queue    // Holds the singleton instance of Queue.
)

// IQueueDriver defines the interface for queue drivers.
// It specifies the methods required for a queue driver to be compatible with the Queue struct.
type IQueueDriver interface {
	Connect() error              // Connect establishes a connection to the queue.
	Close() error                // Close terminates the connection to the queue.
	GetClient() *amqp.Connection // GetConnection returns the underlying connection for direct queue operations.
}

// Queue encapsulates the queue operations and driver.
// It serves as a central point for queue interactions, leveraging a driver that implements the IQueueDriver interface.
type Queue struct {
	driver IQueueDriver // The queue driver, implementing IQueueDriver for queue operations.
}

// Init initializes the queue by establishing a connection.
// It retrieves the singleton instance of the Queue and calls Connect on it.
func Init() (err error) {
	return GetInstance().Connect()
}

// Connect establishes a connection to the queue if not already connected.
// It uses connectOnce to ensure that the queue connection is established only once,
// preventing multiple connections in a concurrent environment.
func (queue *Queue) Connect() (err error) {
	connectOnce.Do(func() {
		configs := config.GetInstance() // Retrieve configurations
		// Initialize the driver with configuration values
		queue.driver = &queueDrivers.RabbitMQ{
			Username: configs.Get("RABBITMQ_USERNAME"),
			Password: configs.Get("RABBITMQ_PASSWORD"),
			Host:     configs.Get("RABBITMQ_HOST"),
			Port:     configs.Get("RABBITMQ_PORT"),
		}

		if err = queue.driver.Connect(); err != nil {
			log.Fatalln("Failed to connect to Queue Service:", err) // Log and halt on error
		}
	})

	return
}

// Close terminates the queue connection.
// It delegates the close operation to the queue driver and logs the closure.
func (queue *Queue) Close() (err error) {
	if err = queue.driver.Close(); err != nil {
		log.Println("Error closing queue connection:", err)
	} else {
		log.Println("Queue Service: Disconnected Successfully.")
	}
	return
}

// GetClient retrieves the amqp.Connection from the queue driver.
// It allows for direct queue operations using the connection.
func (queue *Queue) GetClient() *amqp.Connection {
	return queue.driver.GetClient()
}

// GetInstance returns the singleton instance of the Queue.
// It ensures that only one instance of Queue is created and used throughout the application,
// leveraging getInstanceOnce to enforce this constraint.
func GetInstance() *Queue {
	getInstanceOnce.Do(func() {
		instance = &Queue{} // Initialize the singleton instance if not already created
	})
	return instance
}
