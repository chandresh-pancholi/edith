package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chandresh-pancholi/edith/constant"
	"github.com/chandresh-pancholi/edith/model"
	"io/ioutil"
	"net/http"
	"time"


	"go.uber.org/zap"

	"github.com/Shopify/sarama"
)

//HTTPResponse is a response from downstream rest call
type HTTPResponse struct {
	response *http.Response
	err      error
}

//Invoke hit consumer API and store the response into ES
func Invoke(msg *sarama.ConsumerMessage, consumerList []model.Consumer, logger *zap.Logger) {
	var event model.Event
	ch := make(chan HTTPResponse)

	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		logger.Error("message value deserialisation to event object failed", zap.String("value", string(msg.Value)), zap.Error(err))
	}

	for _, c := range consumerList {
		go handler(c.HTTPMethod, c.URL, event.Payload, logger, ch)
	}

	for _, c := range consumerList {
		// Use the response (<-ch).body
		r := <-ch
		response := r.response

		deliveredMessage := model.DeliveredMessage{}

		deliveredMessage.EventID = event.ID

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Error("response body read failed", zap.String("eventId", event.ID), zap.String("topic", c.Topic), zap.Error(err))
		}
		deliveredMessage.Response = string(body)
		deliveredMessage.Payload = string(msg.Value)
		deliveredMessage.CreatedAt = event.Time
		deliveredMessage.DeliveredAt = time.Now().Unix()

		err = r.err
		if err != nil {
			deliveredMessage.Status = constant.Failed
		}
		deliveredMessage.Status = constant.Success
		deliveredMessage.GroupID = c.GroupID

		fmt.Sprintf("%s", deliveredMessage)
	}
}

//Handler is
func handler(method, url, payload string, logger *zap.Logger, ch chan<- HTTPResponse) {
	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	request, err := http.NewRequest(method, url, bytes.NewBufferString(payload))
	if err != nil {
		logger.Error("http new request creation failed", zap.String("url", url), zap.String("method", method), zap.Error(err))
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		logger.Error("http request to upstream server failed", zap.String("url", url), zap.String("method", method), zap.Error(err))
	}
	//defer resp.Body.Close()

	ch <- HTTPResponse{resp, err}
}
