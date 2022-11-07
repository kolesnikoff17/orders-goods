package kafkaproducer

import "github.com/Shopify/sarama"

// Conn keeps kafka producer
type Conn struct {
	Producer sarama.SyncProducer
}

// New is a constructor for Conn
func New(host, port string) (*Conn, error) {
	brokers := []string{host + ":" + port}
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Conn{Producer: producer}, nil
}

// Close disconnect connection to kafka
func (c *Conn) Close() {
	_ = c.Producer.Close()
}
