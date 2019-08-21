# Edith
Reliable messaging among multiple microservices using message brokers like Kafka, RabbitMQ, NATS Streaming, Kinesis, Google pub sub, Azure event hubs. Currently edith supports Kafka. Future releases will support mentioned message brokers.

Modern day architecture build upon Microservice Architecture to build & ship product very fast.

Microservice Architecture helps in achieving loose coupling and high cohesion. and to make system loosely coupled, team are heaviliy using async publisher subscriber technlogies like Kafka, RabbitMq, NATS streaming, AWS Kinesis and many more.

Edith currently supports Apache Kafka for Async reliable communication among microservices.

### Technology Stack
##### Core Development 
     1. Golang
     2. gRPC
     3. Protobuf
     4. Apache Kafka
     5. Elasticsearch
##### Monitoring and Alerting
    1. Open Census
    2. Jaeger
    3. Prometheus
##### Container Orchestrator stack
    1. Docker
    2. Kubernetes
    3. Istio
 
#### Regular Microservice Architecture
 ![Alt text](docs/micro-architecture.png?raw=true "Title")  
   
#### Microservice Architecture with Edith
![Alt text](docs/edith-01.png?raw=true "Title")
#### Road map

##### MVP - V1
* gRPC client support in language Java, Go
* gRPC server support 
* TLS encryption support between grpc client -- server
* consumer disable notification
* end to end monitoring on Prometheus
* Kubernetes deployment
* daily roll over index on ES
* Archival support of older ES data
