package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github/shipment-auth-notification/internal/config"
	"github/shipment-auth-notification/internal/model"
	"log"
	"math/rand"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

var statuses = []string{model.AUTHORIZED, model.CANCELLED}

func main() {
	cfg, err := config.Load("conf/local.properties")
	if err != nil {
		log.Fatalf("Error de configuración: %v", err)
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KafkaBrokers...),
		kgo.DefaultProduceTopic(cfg.KafkaTopic),
	)
	if err != nil {
		log.Fatalf("Error de conexión: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	fmt.Println("🚀 Productor iniciado. Enviando autorizaciones simuladas...")

	for i := 1; i <= 5; i++ {
		// Simular datos de negocio reales
		auth := model.Auth{
			ShipmentID:     uint64(rand.Intn(1000)),
			SenderID:       uint64(rand.Intn(1000)),
			DestinationID:  uint64(rand.Intn(1000)),
			FacilityID:     "MLA123",
			TrackingNumber: fmt.Sprintf("MLA-%d", rand.Intn(100000000)),
			DateCreated:    time.Now(),
			Status:         statuses[rand.Intn(len(statuses))],
		}

		payload, _ := json.Marshal(auth)

		err = client.ProduceSync(ctx, &kgo.Record{Value: payload}).FirstErr()
		if err != nil {
			log.Printf("❌ Error al enviar tracking number %s: %v", auth.TrackingNumber, err)
		} else {
			fmt.Printf("✅ Tracking number enviado: %s \n", auth.TrackingNumber)
		}
		time.Sleep(1 * time.Second)
	}
}
