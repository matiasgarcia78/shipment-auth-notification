# Shipment Auth Notification

Servicio de ejemplo en Go para publicar y consumir autorizaciones de envíos mediante Apache Kafka. El productor genera cinco autorizaciones simuladas y las publica en el tópico `authorizations`; el consumidor lee ese tópico e informa el shipment y su estado.

## Requisitos

- Go 1.26.5 o compatible con la versión indicada en `go.mod`.
- Docker y Docker Compose.
- El puerto `9092` disponible en `localhost`.

## Levantar Kafka

Desde la raíz del proyecto:

```bash
docker compose up -d
docker compose ps
```

El broker queda disponible en `localhost:9092`. Si ya existe otro contenedor Kafka usando ese puerto, detenelo antes de levantar este proyecto:

```bash
docker ps
docker stop <contenedor-kafka-anterior>
```

Creá el tópico que utilizan el productor y el consumidor:

```bash
docker exec portfolio-kafka \
  /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --if-not-exists \
  --topic authorizations \
  --partitions 1 \
  --replication-factor 1
```

Podés comprobarlo con:

```bash
docker exec portfolio-kafka \
  /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --list
```

## Instalar dependencias

```bash
go mod download
go mod tidy
```

La aplicación usa [`franz-go`](https://github.com/twmb/franz-go) como cliente Kafka.

## Ejecutar

En una terminal, iniciá el consumidor:

```bash
go run ./cmd/consumer
```

En otra terminal, ejecutá el productor:

```bash
go run ./cmd/producer
```

El productor enviará cinco mensajes JSON al tópico `authorizations`. El consumidor permanecerá ejecutándose y mostrará cada autorización recibida.

También podés compilar ambos comandos:

```bash
go build ./cmd/producer
go build ./cmd/consumer
```

## Detener Kafka

```bash
docker compose down
```

## Solución de problemas

- `UNKNOWN_TOPIC_OR_PARTITION`: verificá que Kafka esté activo y que exista el tópico `authorizations`.
- `connection refused`: comprobá que el contenedor esté `Up` y que `localhost:9092` no esté ocupado por otro Kafka.
- Si hay dos contenedores Kafka, el productor puede conectarse al broker equivocado. Dejá un único contenedor exponiendo `0.0.0.0:9092->9092/tcp`.
