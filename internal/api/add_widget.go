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

func AddWidget(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Fprintf(w, "GET Request")
		log.Println("GET Request, nothing to do.")
	case "POST":
		fmt.Fprintf(w, "POST Request")
		log.Println("Got Request")

		var widget model.WidgetProperties
		err := json.NewDecoder(req.Body).Decode(&widget)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		widget.Timestamp = time.Now().UnixNano()
		widget.IsActive = true
		widget.ValidUntil = nil
		log.Println("Received Payload: ", widget)

		//redisclient.InitRedisClient().AddNewWidget(widget)
		redisclient.InitRedisClient().StreamProducer("registration", widget)
		/* Once consumed from Redis streams update on ColumnDB(TBD) and delete from Redis */
	default:
		fmt.Fprintf(w, "No method allowed")

	}
}
