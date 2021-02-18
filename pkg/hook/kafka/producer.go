/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 16/02/2021
*/
package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	kafkaGo "github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer Creation of an emitter and send any type of value to its topic in kafka
func (k *Kafka) Producer(messageKey string, messageValue interface{}, opts ...Option) error {
	log := logger.WithName("Producer")

	options := NewOptions(opts...)

	producerCM := &kafkaGo.ConfigMap{
		"bootstrap.servers": k.BootstrapServer,
	}
	if options.AuthKafka {
		if err := producerCM.SetKey("sasl.mechanisms", options.Mechanisms); err != nil {
			return err
		}
		if err := producerCM.SetKey("security.protocol", options.Protocol); err != nil {
			return err
		}
		if err := producerCM.SetKey("bootstrap.servers", k.BootstrapServer); err != nil {
			return err
		}
		if err := producerCM.SetKey("sasl.username", k.Username); err != nil {
			return err
		}
		if err := producerCM.SetKey("sasl.password", k.Password); err != nil {
			return err
		}
	}

	p, err := kafkaGo.NewProducer(producerCM)
	if err != nil {
		return fmt.Errorf("failed to create producer: %s", err)
	}
	defer p.Close()
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafkaGo.Message:
				if ev.TopicPartition.Error != nil {
					log.Error(ev.TopicPartition.Error, "Failed to deliver",
						"stacktrace", ev.TopicPartition)
				} else {
					log.Info("Successfully produced record",
						"topic", *ev.TopicPartition.Topic,
						"partition", ev.TopicPartition.Partition,
						"offset", ev.TopicPartition.Offset)
				}
			}
		}
	}()

	message, err := json.Marshal(&messageValue)
	if err != nil {
		log.Error(err, "marshal failed")
		return err
	}
	if err := p.Produce(&kafkaGo.Message{
		TopicPartition: kafkaGo.TopicPartition{Topic: &k.Topic, Partition: kafkaGo.PartitionAny},
		Key:            []byte(messageKey),
		Value:          message,
		Timestamp:      time.Now(),
	}, nil); err != nil {
		log.Error(err, "produce failed")
		return err
	}
	p.Flush(int(options.WriteTimeout / time.Millisecond))
	return nil
}
