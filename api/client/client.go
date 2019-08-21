package client

import (
	"encoding/json"
	"fmt"
	"github.com/chandresh-pancholi/edith/model"
	"github.com/chandresh-pancholi/edith/pkg/elasticsearch"
	es "github.com/chandresh-pancholi/edith/pkg/elasticsearch/config"
	"github.com/chandresh-pancholi/edith/pkg/kafka/admin"
	"github.com/chandresh-pancholi/edith/pkg/mysql"
	"github.com/chandresh-pancholi/edith/pkg/mysql/client"
	"net/http"
	"time"


	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//Handler is
type Handler struct {
	esClient    *es.ESClient
	db          *mysql.DB
	logger      *zap.Logger
	adminClient *admin.AdminClient
}

//NewClientHandler to create client api object
func NewClientHandler(db *mysql.DB, logger *zap.Logger, adminClient *admin.AdminClient, esClient *es.ESClient) *Handler {
	h := new(Handler)

	h.db = db
	h.logger = logger
	h.adminClient = adminClient
	h.esClient = esClient
	return h
}

//CreateClient creates new client in Edith
func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var client model.Consumer
	err := decoder.Decode(&client)
	if err != nil {
		h.logger.Error("client creation input reuest is not correct",
			zap.Error(err))

		w.WriteHeader(http.StatusBadRequest)
	}

	client.ConsumerID = uuid.NewV4().String()
	client.CreatedAt = time.Now().Unix()
	client.UpdatedAt = time.Now().Unix()

	err = elasticsearch.IndexDocument(h.esClient, viper.GetString("elasticsearch.client.index"), viper.GetString("elasticsearch.client.index.type"), client.ConsumerID, client)
	if err != nil {
		h.logger.Error("document indexing failed in elasticsearch", zap.Error(err))
	}
	err = h.adminClient.CreateTopic(client.Topic, client.Partition, viper.GetInt("kafka.replication"), h.logger)
	if err != nil {
		h.logger.Error("Kafka topic creation failed",
			zap.String("client_id", client.ConsumerID),
			zap.String("name", client.Name),
			zap.String("topic", client.Topic),
			zap.Int("partition", client.Partition),
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

//GetClient fetches consumer from ES
func (h *Handler) GetClient(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("getting consumer.......")
	//consumerName := r.URL.Query().Get("cname")
	//elasticsearch.Search(h.esClient, viper.GetString("elasticsearch.client.index"))
}

// StopClient to stop consumer
func (h *Handler) StopClient(w http.ResponseWriter, r *http.Request) {
	consumerID := r.URL.Query().Get("cid")
	topic := r.URL.Query().Get("topic")

	if consumerID == "" || topic == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	err := client.Stop(h.db, consumerID, topic)
	if err != nil {
		h.logger.Error("Consumer client stop command to database failed",
			zap.String("consumer_id", consumerID),
			zap.String("topic", topic),
			zap.Error(err))
	}

	//TODO Figure out a way to unsubscribe kafka topic

	fmt.Println("CID ===> " + consumerID)
}
