package model

import (
	"github.com/chandresh-pancholi/edith/constant"
	"io"
	"net/http"


	"github.com/Shopify/sarama"
)

//Event to ingest Kafka
type Event struct {
	ID          string                `json:"event_id"`
	Key         string                `json:"event_key"`
	Name        string                `json:"name"`
	Payload     string                `json:"payload"`
	Time        int64                 `json:"event_time"`
	Topic       string                `json:"topic"`
	ServiceName string                `json:"service_name"`
	Metadata    interface{}           `json:"params"`
	Headers     []sarama.RecordHeader `json:"headers"`
}

//Client to submit to MSG
type Client struct {
	ClientID      string `json:"client_id"`
	Name          string `json:"name"`
	Topic         string `json:"topic"`
	Partition     int    `json:"partition"`
	URL           string `json:"url"`
	HTTPMethod    string `json:"http_method"`
	Header        bool   `json:"header"`
	Description   string `json:"description"`
	TotalConsumer int    `json:"total_consumer"`
	GroupID       string `json:"group_id"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

//Consumer to submit to MSG
type Consumer struct {
	ConsumerID    string `json:"client_id"`
	Name          string `json:"name"`
	Topic         string `json:"topic"`
	Partition     int    `json:"partition"`
	URL           string `json:"url"`
	HTTPMethod    string `json:"http_method"`
	DynamicURL    bool   `json:"dynamic_url"`
	Status        bool   `json:"status"`
	MaxRetry      int8   `json:"max_retry"`
	Description   string `json:"description"`
	TotalConsumer int    `json:"total_consumer"`
	GroupID       string `json:"group_id"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

//Inbound to store message inbound information
type Inbound struct {
	EventID    string
	EventTopic string
	Delivered  bool
	CreatedAt  int64
	UpdatedAt  int64
}

//Outbound after message is deliver to Kafka consumer
type Outbound struct {
	EventID    string
	EventTopic string
	Delivered  string
	StatusCode int
	Payload    interface{}
	Retry      int8
	Exception  string
	CreatedAt  int64
	UpdatedAt  int64
}

//DeliveredMessage when message is delivered to Kafka
type DeliveredMessage struct {
	EventID     string
	GroupID     string
	Client      string
	Payload     string
	Response    string
	URL         string
	Status      constant.DeliveredMessageStatus
	Headers     map[string]string
	CreatedAt   int64
	DeliveredAt int64
}

//Response is api response
type Response struct {
	StatusCode int           `json:"status_code"`
	Header     http.Header   `json:"header"`
	Body       io.ReadCloser `json:"body"`
}
