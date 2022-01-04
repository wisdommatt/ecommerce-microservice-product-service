# Ecommerce Product Service

This is the GRPC product service that handles all activities in the ecommerce application that has to do with products (adding products, updating products, deleting products etc).

## Applications

* ## [Jaeger](https://www.jaegertracing.io/)

  * Jaeger is an **open source software for tracing transactions between distributed services**.
* ## [Nats](https://nats.io)

  * NATS is **an open-source messaging system** (sometimes called message-oriented middleware).
  * NATS is used in the product service for communicating with the notification service when a new product is added.
* ## [MySQL](https://www.mysql.com/)

  * MySQL is an open-source relational database management system.
  * The product service stores products in MySQL.

### Usage

To install / run the user microservice run the command below:

```bash
docker-compose up
```

## Requirements

The application requires the following:

* Go (v 1.5+)
* Docker (v3+)
* Docker Compose

## Other Micro-Services / Resources

* #### [User Service](https://github.com/wisdommatt/ecommerce-microservice-user-service)
* #### [Notification Service](https://github.com/wisdommatt/ecommerce-microservice-notification-service)
* #### [Cart Service](https://github.com/wisdommatt/ecommerce-microservice-cart-service)
* #### [Shared](https://github.com/wisdommatt/ecommerce-microservice-shared)

## Public API

The public graphql API that interacts with the microservices internally can be found in [https://github.com/wisdommatt/ecommerce-microservice-public-api](https://github.com/wisdommatt/ecommerce-microservice-public-api).
