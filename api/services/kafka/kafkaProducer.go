package kafka

import (
	"encoding/json"
	"fmt"
	webappConfig "webapp/config"
	"github.com/IBM/sarama"
)

var (
	kafkaTopic   = webappConfig.GetEnvConfig().KAFKA_TRACE_TOPIC
	kafkaBrokers = webappConfig.GetEnvConfig().KAFKA_BROKER
	producer sarama.SyncProducer
)

// TraceMetadata defines the structure of metadata sent to Kafka
type TraceMetadata struct {
	CourseId     int    `json:"courseId"`
	TermId       int    `json:"termId"`
	InstructorId int    `json:"instructorId"`
	ReportType   string `json:"reportType"`
	FolderPath   string `json:"folderPath"`
	Filename     string `json:"filename"`
}


// InitializeKafkaProducer sets up the Kafka producer
func InitializeKafkaProducer() error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3

	var err error
	producer, err = sarama.NewSyncProducer([]string{kafkaBrokers}, config)
	if err != nil {
		return fmt.Errorf("failed to initialize Kafka producer: %v", err)
	}
	return nil
}

// CloseKafkaProducer closes the Kafka producer
func CloseKafkaProducer() {
	if producer != nil {
		producer.Close()
	}
}

// sendToKafka sends metadata to the specified Kafka topic
func SendToKafka(metadata TraceMetadata) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %v", err)
	}

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.StringEncoder(data),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}
	return nil
}

