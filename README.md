<div align="center">
<h1>Go-Microservices-Email</h1>
<p>
HTTP / gRPC / RabbitMQ Golang email microservice
</p>
</div>


# 📖 About
This is an email service for sending emails through grpc / http ([grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)) and rabbitmq interfaces. With full cycle of DevOps methodologies ([Ops repository](https://google.com))

**⚙️ Technologies used:**
* [GRPC](https://grpc.io/) - gRPC
* [GRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway) - gRPC gateway for accessing either http and gRPC interfaces
* [Swagger](https://github.com/go-swagger/go-swagger) - Swagger docs
* [RabbitMQ](https://github.com/streadway/amqp) - RabbitMQ
* [Gorm](https://github.com/go-gorm/gorm) - ORM library for Golang
* [Echo](https://github.com/labstack/echo) - High performance, minimalist Go web framework
* [Logrus](https://github.com/sirupsen/logrus) - Structured, pluggable logger for Golang
* [viper](https://github.com/spf13/viper) - Go configuration with fangs
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [Docker](https://www.docker.com/) - Docker
* [Prometheus](https://prometheus.io/) - Prometheus
* [Gomail](https://github.com/go-gomail/gomail/tree/v2) - Simple and efficient package to send emails
* [Prometheus-go-client](https://github.com/prometheus/client_golang) - Prometheus instrumentation library for Go applications
* [otel](https://github.com/open-telemetry/opentelemetry-go) - otel client for Golang
* [Jaeger](https://www.jaegertracing.io) - otel Jaeger exporter

# ✅ Setup
First you need SMTP credentials for [gomail](https://github.com/go-gomail/gomail), you can obtain it from various providers, for example [Yahoo](https://help.yahoo.com/kb/SLN4724.html). Provider used as email delivery vendor for sending and receive email messages througth [SMTP Protocol](https://en.wikipedia.org/wiki/Simple_Mail_Transfer_Protocol)

## Example config
 
```yaml
smtp:
  Host: ""
  Port: 0
  User: ""
  Password: ""

logger:
  Mode: "prod"

metrics:
  Port: 8242

gateway:
  Port: 8192

jaeger:
  Host: localhost:14268

grpc:
  Port: 50001

rabbit:
  Host: "localhost"
  Port: 5672
  User: "guest"
  Password: "guest"
  QueueName: "test"
  ConsumePool: 2

database:
  Host: "localhost"
  Port: 5432
  User: "root"
  Password: "postgres"
  DbName: "mails_db"
  SslMode: "disable"
 ```

# 🏃 Running

For running you need a CONFIG variable setted in env, for test purposes you can edit values in config/config-local.yml and then export it with following command
```bash
export CONFIG=$(cat config/config-local.yml)
```

## Local run
```bash
make compose_up
make run
```

# Endpoints

*Note: these endpoints valid with default `config/config-local.yml` file*

### gRPC gateway (REST)
http://localhost:8192/

### gRPC Server
http://localhost:50001/

### Metrics and Swagger
http://localhost:8242/metrics

### Swagger
For running swagger execute the following command

```bash
make swagger
```

# Credits

- [Ignat Belousov](https://github.com/fayvori) (Author)

# License

```
MIT License

Copyright (c) 2023 Ignat Belousov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
