package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"spike/internal/model"
	"spike/internal/redisclient"
	"time"
)

func RemoveWidget(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Fprintf(w, "This is a GET method")
	case "POST":
		fmt.Fprintf(w, "POST Request")
		log.Println("Got Request")

		var widget model.WidgetDeregistration
		err := json.NewDecoder(req.Body).Decode(&widget)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		widget.IsActive = false
		widget.ValidUntil = time.Now().UnixNano()
		log.Println("Received Payload: ", widget)

		//redisclient.InitRedisClient().AddNewWidget(widget)
		redisclient.InitRedisClient().StreamProducer("deletion", widget)
		/* Once consumed from Redis streams update on ColumnDB(TBD) and delete from Redis */

	default:
		fmt.Fprintf(w, "No method allowed")

	}
}
