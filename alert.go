// Erik Oomen, jul 2019
//
// Reads kafka business topic and filter out the severe errors
//
// Warning, not written for stability or performance

package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "encoding/json"
    kafka "github.com/segmentio/kafka-go"
)

const (
	KAFKA_SERVER    = "192.168.1.106:9092"
	SENSOR_TOPIC    = "sensors"
	BUSINESS_TOPIC  = "business"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 51, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait: 100 * time.Millisecond,
	})
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
        return kafka.NewWriter(kafka.WriterConfig{
                Brokers:  []string{kafkaURL},
                Topic:    topic,
                Balancer: &kafka.LeastBytes{},
                BatchTimeout: 50 * time.Millisecond,
        })
}

func main() {
    kafkaReader := getKafkaReader(KAFKA_SERVER, BUSINESS_TOPIC, "")
    defer kafkaReader.Close()

	fmt.Println("start consuming ... !!")
	for {
                var result map[string]interface{}
		m, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
                json.Unmarshal([]byte(m.Value),  &result)
		if result["Severity"] != nil && result["Severity"].(string) == "1" {
			fmt.Printf("ALERT @ %v\n", m.Time.Format("Jan _2 15:04:05 2006"));
		}
	}
}

