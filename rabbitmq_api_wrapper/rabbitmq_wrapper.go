package rabbitmqapiwrapper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type QueueInfo struct {
	Name          string `json:"name"`
	Messages      int    `json:"messages"`
	MessagesReady int    `json:"messages_ready"`
	MessagesUnack int    `json:"messages_unacknowledged"`
	Vhost         string `json:"vhost"`
}

/*
* Queries RabbitMQ API to get info on all queues, extracts the information and parses the response
 */
func GetAllQueueMetrics(host string) (err error, metrics []QueueInfo) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s:/api/queues", host)
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		// handle error
		log.Fatal("Err Creating Req => ", err)
	}

	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")

	authHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password))))

	req.Header.Add("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		log.Fatal("Err sending req => ", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return err, nil
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Non-200 resonse: %d\n", res.StatusCode)
		fmt.Println("resonse body:", string(body))
		return
	}

	// Parse the json response
	var queuesDet []QueueInfo
	errr := json.Unmarshal(body, &queuesDet)

	if errr != nil {
		log.Fatal("Problem encountered parsing json...")
	}

	return nil, queuesDet
}
