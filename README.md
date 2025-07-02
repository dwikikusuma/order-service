# Order Service

This is an Order Service built using **Go**, leveraging **PostgreSQL** for data storage, **Redis** for caching, **GORM** as the ORM for PostgreSQL, and **Kafka** for messaging. The service is secured with **JWT authentication** to protect API endpoints.

## Tech Stack

- 🦸 **Go**: The main programming language for building the service.
- 🐘 **PostgreSQL**: Used for storing order data.
- 🔥 **Redis**: Utilized for caching and improving performance.
- 🦦 **GORM**: An ORM for Go to interact with PostgreSQL.
- 🐐 **Kafka**: Implements event-driven architecture for messaging and decoupling services.
- 🔑 **JWT Authentication**: Provides secure authentication for the API.

## Overview

The Order Service manages the creation and retrieval of orders while integrating with Kafka for asynchronous processing and Redis for caching frequently accessed data. The service ensures that sensitive API endpoints are protected with JWT-based authentication, allowing only authorized users to access certain operations.

This service can be scaled to fit more complex architectures, utilizing Kafka to handle high throughput and asynchronous tasks. With the flexibility of GORM and Redis, this service ensures both reliability and performance.
