package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/Shopify/sarama"
	"go-boilerplate/internal/models"
	"go.uber.org/zap"
)

type Consumer struct {
	consumer sarama.Consumer
	topic    string
	logger   *zap.Logger
}

func NewConsumer(brokers []string, topic string, logger *zap.Logger) (*Consumer, error) {
	// Create admin client first to create topic if it doesn't exist
	adminConfig := sarama.NewConfig()
	adminConfig.Version = sarama.V2_8_0_0 // Use a recent version that supports topic creation

	admin, err := sarama.NewClusterAdmin(brokers, adminConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster admin: %w", err)
	}
	defer admin.Close()

	// Try to create topic
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	err = admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		// Ignore error if topic already exists
		if !strings.Contains(err.Error(), "Topic with this name already exists") {
			return nil, fmt.Errorf("failed to create topic: %w", err)
		}
	}

	// Create consumer
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sarama.V2_8_0_0 // Match admin client version

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
		logger:   logger,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return err
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()

			for {
				select {
				case msg := <-pc.Messages():
					var message models.Message
					if err := json.Unmarshal(msg.Value, &message); err != nil {
						c.logger.Error("Failed to unmarshal message", zap.Error(err))
						continue
					}

					c.logger.Info("Received message",
						zap.String("id", message.ID.String()),
						zap.String("content", message.Content),
					)

				case <-ctx.Done():
					return
				}
			}
		}(pc)
	}

	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
