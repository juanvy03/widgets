package redisclient

import (
	"context"
	"encoding/json"
	"log"
	"spike/internal/model"
	"spike/internal/mongoclient"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func InitRedisClient() RedisClient {
	opt, err := redis.ParseURL("redis://default:password@localhost:6379/0")
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	log.Printf("Redis Client initialized with options: Addr=%s, DB=%d", opt.Addr, opt.DB)
	return RedisClient{client: rdb}
}

func (r RedisClient) AddRecord(key string, value string) error {
	ctx := context.Background()
	log.Printf("Attempting to add record to Redis: Key=%s, Value=%s", key, value)
	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("Error adding record to Redis: %v", err)
		return err
	}

	log.Println("Successfully added record to Redis.")
	return nil
}

func (r RedisClient) AddNewWidget(widget model.WidgetProperties) error {
	ctx := context.Background()

	log.Println("Attempting to add widget to Redis:", widget)

	widgetProperties := map[string]interface{}{
		"name":          widget.Name,
		"serial_number": widget.SerialNumber,
		"port_type_p":   widget.PortTypeP,
		"port_type_r":   widget.PortTypeR,
		"port_type_q":   widget.PortTypeQ,
		"timestamp":     widget.Timestamp,
	}

	err := r.client.HSet(ctx, widget.Name, widgetProperties).Err()
	if err != nil {
		log.Println("Error adding record to Redis: ", err)
		return err
	}

	log.Println("Successfully added record to Redis.")
	return nil
}

func (r RedisClient) DeleteTempWidgetHash(widget model.WidgetProperties) error {

	ctx := context.Background()
	err := r.client.HDel(ctx, widget.Name)
	if err != nil {
		log.Println("Error deleting widget: ", err)
	}
	log.Println("Successfully deleted record to Redis.")
	return nil
}

func (r RedisClient) StreamProducer(event_type string, widget interface{}) error {
	ctx := context.Background()
	log.Println("Publishing new event")

	jsonData, erro := json.Marshal(widget)
	if erro != nil {
		log.Println("Marshall Error: ", erro.Error())
		return erro
	}

	stringEvent := string(jsonData)

	newWidgetEvent := map[string]interface{}{
		"event_type": event_type,
		"event":      stringEvent,
		"timestamp":  time.Now().Unix(),
	}

	log.Println(newWidgetEvent)

	err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "events",
		Values: newWidgetEvent,
	}).Err()

	if err != nil {
		log.Printf("Error pushing Event to Redis: %v", err)
		return err
	}

	log.Println("Event pushed successfully!")
	return nil
}

func (r RedisClient) StreamConsumer(stream string) error {
	ctx := context.Background()
	log.Println("Redis StreamConsumer :: Started")

	consumerGroup := "events-consumer"
	consumerName := "processor"

	err := r.client.XGroupCreateMkStream(ctx, stream, consumerGroup, "0").Err()
	if err != nil {
		log.Println(err)
	}

	for {
		entries, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: consumerName,
			Streams:  []string{stream, ">"},
			Count:    2,
			Block:    0,
			NoAck:    false,
		}).Result()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Events backlog:", len(entries[0].Messages))
		for i := 0; i < len(entries[0].Messages); i++ {
			log.Println(entries[0].Messages[i])
			event := entries[0].Messages[i]

			eventType, ok := event.Values["event_type"].(string)
			if !ok {
				log.Println("event_type field is missing or not a string")
				continue
			}

			var eventMap map[string]interface{}
			err = json.Unmarshal([]byte(event.Values["event"].(string)), &eventMap)
			if err != nil {
				log.Fatalln("Error decoding JSON string:", err)
			}

			/*
				Save to MongoDB
			*/
			switch eventType {
			case "registration":
				log.Println(">>> Registering:", eventType, eventMap)
				mongoclient.MongoProcessRegistration(eventType, eventMap)

				/*
					ACK the message in Streams (This is a copy paste)
				*/
				_, err = r.client.XAck(ctx, stream, consumerGroup, event.ID).Result()
				if err != nil {
					log.Printf("Error acknowledging message: %v", err)
				} else {
					log.Println("Message acknowledged:", event.ID)
				}
			case "deletion":
				log.Println("deregistration")
				mongoclient.MongoProcessDeregistration(eventMap)
			}
		}
	}

}
