package client

import (
	"github.com/chandresh-pancholi/edith/model"
	"github.com/chandresh-pancholi/edith/pkg/mysql"
)

//Create insert client to MySQL DB
func Create(db *mysql.DB, consumer model.Consumer) error {
	sqlStatement := `
	INSERT INTO consumers (consumer_id, name, topic, partition, url, http_method, description, total_consumer, group_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	insert, err := db.Query(sqlStatement, consumer.ConsumerID, consumer.Name, consumer.Topic, consumer.Partition, consumer.URL,
		consumer.HTTPMethod, consumer.Description, consumer.TotalConsumer, consumer.GroupID)
	if err != nil {
		return err
	}

	defer insert.Close()

	return nil
}

//Stop to stop consumer
func Stop(db *mysql.DB, consumerID, topic string) error {
	sqlStatement := `UPDATE consumers SET status = 0 where consumer_id=$1 and topic = $2`

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(consumerID, topic)
	if err != nil {
		return err
	}

	return nil
}
