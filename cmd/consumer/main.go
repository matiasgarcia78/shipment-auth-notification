package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github/shipment-auth-notification/internal/model"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	client, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		kgo.ConsumeTopics("authorizations"),
		kgo.ConsumerGroup("notification-service-group"),
	)
	if err != nil {
		log.Fatalf("Error de conexión: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	fmt.Println("🕵️‍♂️ Consumidor listo. Monitoreando autorizaciones..")

	for {
		fetches := client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Fatalf("Error en Kafka: %v", errs)
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			var auth model.Auth
			if err := json.Unmarshal(record.Value, &auth); err != nil {
				log.Printf("Error al decodificar: %v", err)
				continue
			}

			// Aplicar lógica de negocio
			fmt.Printf("Shipment %d %s!\n", auth.ShipmentID, auth.Status)
		}
	}
}
