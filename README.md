# GoCloud: Microservices-based Cloud Application

## Overview

**GoCloud** is a scalable, microservices-based project built with **Go**. It follows **Clean Code Architecture** principles and supports **gRPC** and **message brokers** for service-to-service communication. The project includes **authentication & authorization**, **email and push notifications**, **SQL & NoSQL databases**, **caching**, **file storage**, and **distributed tracing**.

The project is about Delivery System which consists of four microservices:
- **API Gateway**
- **Orders Service**
- **Partners Service**
- **Notifications Service**

Each service is containerized using **Docker** and orchestrated via **Docker Compose** and **Kubernetes**. The project also features **automatic database schema migration**, **mock tests**, and **Swagger documentation** for APIs.

---

## Project Structure

```
gocloud/
│── cmd/                        # Main service binaries
│   ├── api/                    # API Gateway, authentication & authorization service
│   ├── notifications/          # Push Notifications service
│   ├── orders/                 # Orders service
│   ├── partners/               # Partners and products service
│── configs/                    # Configuration files
│── design/                     # System design diagrams
│── docs/                       # API documentation (Swagger)
│── internal/                   # Core business logic (private)
│── migrations/                 # Database migrations
│── pkg/                        # Shared utility packages
│── docker-compose.yml          # Docker Compose configuration
│── .dockerignore               # Files and directories to exclude when building Docker images
│── Makefile                    # Automation commands
│── go.mod, go.sum              # Go dependencies
│── README.md                   # Project documentation
```

---

## Services

### 1. API Gateway (`cmd/api/`)

The **API Gateway** is the central entry point for all **customer-facing HTTP requests** in the system. It handles **authentication, authorization, product retrieval, order management, and callback processing** from payment and partner systems. It supports **high-performance request routing, validation, caching, and distributed tracing** while acting as a bridge between clients and internal microservices. Provides **Swagger documentation** for easy API interaction.  

#### Key Responsibilities

1. **Handles HTTP Requests Only**  
   - Serves as the primary **HTTP API** for frontend and external consumers.  
   - Supports **RESTful routes** for **authentication, product availability, and order management**.  

2. **Implements Custom HTTP Router**  
   - Uses a **custom-built router** based on `http.ServeMux`.  
   - Supports common **HTTP methods (GET, POST, PUT, DELETE)**.  
   - Provides **parameter parsing and request routing** for different endpoints.  

3. **Authentication & Authorization**  
   - Uses **third-party authentication providers** for customer login and registration.  
   - Implements **middleware-based authorization** to enforce access control.  

4. **Routes & Endpoints**  
   - **Authentication & Authorization** – Manages customer login, registration, and session validation.  
   - **Product Availability** – Fetches currently available products from the **Partners Service**.  
   - **Order Management** – Handles order creation, updates, and retrieval.  
   - **Payment Callbacks** – Processes responses from **Payment API systems**.  
   - **Partner Callbacks** – Receives and validates status updates from **Partners API systems**.  

5. **Data Validation**  
   - Ensures incoming requests contain the necessary fields.  
   - Returns appropriate error messages for malformed requests.  
   - Utilizes standard validation tags along with custom validations.  
   - Uses [Package validator](https://github.com/go-playground/validator).

6. **Swagger API Documentation**  
   - Generates API docs for HTTP routes using [http-swagger](https://github.com/swaggo/http-swagger).  
   - Accessible at `http://delivery.local/docs/index.html`.  
   - Allows developers to **test endpoints interactively**.  

7. **Distributed Tracing**  
   - Uses **distributed tracing** to link request traces across microservices.  
   - Ensures **full visibility** of API calls, including performance metrics.  

8. **File Storage for Customer Profile Avatars**  
   - Stores **customer profile images** securely.  
   - Allows users to upload and retrieve their **profile avatars**.  

9. **Caching for Order Processing**  
   - Stores **validated and checked orders** in a cache layer (e.g., Redis).  
   - Reduces load on the database and service calls by avoiding redundant validations.  

10. **Communication with Other Microservices**  
    - Uses **gRPC** for high-performance service-to-service communication.  
    - Publishes **messages via message broker** for event-driven workflows.  

### 2. Push Notifications Service (`cmd/notifications/`)

The **Push Notifications Service** is responsible for handling **real-time push notifications** across different microservices. It processes **incoming events** via a **message broker** and delivers push messages through **third-party notification providers**. The service ensures **reliable delivery, priority-based handling, and historical tracking** of notifications.  

It integrates **SQL** and **NoSQL** databases for configuration and message storage while leveraging **distributed tracing** to maintain visibility into the request flow across multiple services.  

#### **Key Responsibilities**  

1. **Receives Events via Message Broker**  
   - Listens to **message broker** queues for notification-related events.  
   - Supports various event types (e.g., order updates, payment confirmations, system alerts).  

2. **Manages Push Notification Settings**  
   - Uses an **SQL database** to store service-specific **notification settings**.  
   - Configures **API keys, priorities, and delivery preferences** for different services (agents).  

3. **Stores Notification History**  
   - Uses a **NoSQL database** to store **sent notifications and delivery status**.  
   - Maintains a **history of notifications** for tracking and troubleshooting.  

4. **Integrates with Third-Party Push Providers**  
   - Sends push notifications via external services (e.g., **Firebase Cloud Messaging (FCM), Apple Push Notification Service (APNS), or other providers**).  
   - Supports **multi-provider configurations** based on service priority and fallback mechanisms.  

5. **Implements Distributed Tracing**  
   - Links with **previous service traces** to provide full observability.  
   - Ensures that each notification event is trackable in the request lifecycle.  

### 3. Orders and Payments Service (`cmd/orders/`)

The **Orders and Payments Service** is responsible for managing customer orders and processing payments. It communicates with other microservices via **gRPC** and **RabbitMQ (message broker)** and stores data in both **SQL (PostgreSQL)** and **NoSQL (MongoDB or another NoSQL database)** for optimized performance.  

This service ensures seamless order management, records transaction history, and logs errors using **distributed tracing**.  

#### Key Responsibilities  

1. **Accepts requests via gRPC and Message Broker**  
   - Exposes gRPC endpoints for order and payment management.  
   - Listens to **message broker** queues for asynchronous processing.  

2. **Uses SQL and NoSQL databases**  
   - **SQL Database (PostgreSQL):**  
      - Stores **banks** (payment providers that can process payments).  
      - Stores **orders** (order details, status, partner and customer info).  
   - **NoSQL Database (MongoDB or similar):**  
      - Stores **payments** (transaction details).  
      - Stores **order history** (audit logs, previous order statuses).  

3. **Manages Orders & Payments**  
   - Creates, updates, and retrieves orders.  
   - Logs all payment transactions in noSQL database for future reference.  

4. **Distributed Tracing & Error Logging**  
   - Uses **distributed tracing** to connect to chain traces of previous services.  

5. **Publishes Events to Message Broker**  
   - Sends messages to **message broker** when actions are required in other microservices.  
   - Notifies **inventory service**, **notifications service**, and other relevant systems.  

### 4. Partners and Products Service (`cmd/partners/`)

The **Partners and Products Service** is responsible for managing **partners, products, and their availability**. It plays a crucial role in ensuring that partner businesses have up-to-date product listings and availability statuses. This service also notifies partners when a customer **pays for an order**, ensuring that they prepare the necessary products for fulfillment.  

The service **accepts requests via gRPC and a message broker** and integrates with **distributed tracing** to link its operations with previous service traces.  

#### **Key Responsibilities**  

1. **Handles Partner and Product Management**  
   - Manages **partners** (registered businesses that supply products).  
   - Manages **products** offered by each partner.  
   - Tracks **availability** of products with pricing information.  

2. **Accepts Requests via gRPC and Message Broker**  
   - Provides **gRPC endpoints** for managing partners and products.  
   - Listens to **message broker** for updates and notifications from other services.  

3. **Uses Distributed Tracing**  
   - Links to traces from previous services (**Orders & Payments, API Gateway**).  
   - Provides full observability across the service chain.  

4. **Database Structure**  
   - Stores **partners** (businesses providing products).  
   - Stores **products** (items offered by partners).  
   - Stores **availability** (which partners have which products and at what price).  

5. **Notifies Partners When Orders are Paid**  
   - Sends requests to **partners' external APIs** when an order is paid.  
   - Ensures partners are aware of **which products to prepare and in what quantities** for customers.  

---

## Features

- **Clean Architecture**: Well-structured and maintainable.
- **Microservices Architecture**: Independent, scalable services.
- **gRPC & Message Broker**: Efficient inter-service communication.
- **Authentication & Authorization**: Secure access management.
- **Redis Caching**: Performance improvements with in-memory caching.
- **SQL & NoSQL Databases**: PostgreSQL and a NoSQL solution for optimized and flexible data storage.
- **Push Notifications**: Asynchronous messaging for push notifications.
- **File Storage**: Persistent storage solution.
- **Distributed Tracing**: OpenTelemetry for performance monitoring.
- **Swagger API Documentation**: Easy API consumption.
- **Comprehensive Testing**: Unit tests and integration tests.
- **Docker & Kubernetes Ready**: Containerized for scalability.
- **Custom HTTP Router**: Optimized request processing.

---

## Design of the Delivery System

Communication methods are color-coded in design diagrams:  
   - **Silver**: WebSocket  
   - **Green**: HTTP  
   - **Pink**: gRPC  
   - **Orange**: Message Broker  
   - **Violet**: Third-Party Integrations  
   - **Cyan**: Internal Calls  

![Design](./design/design-0-delivery.svg)

---

## Scenarios

### 1. Customer Registers in the System  

   1. **Front App** subscribes the client to a push notification provider and obtains a token: "Turn on notifications: Allow/Cancel."  
   2. **Front App** sends the token along with the customer's profile details to the **API Gateway** for registration.  
   3. **API Gateway** parses and validates the request, then transfers it to its service.  
   4. **Gateway Service** calls a third-party **Auth Service** to register the customer.  
   5. The customer logs into the system using their email and password and receives access and refresh tokens. (**2**->**3**->**4**)  
   6. **Front App** saves the tokens and uses them in subsequent requests.  

![1](./design/design-1-sing-up.svg)

### 2. Customer Lists Available Products  

   1. **Front App** sends a request to the **API Gateway**.  
   2. **API Gateway** verifies the access token, parses and validates the request, and transfers it to its service.  
   3. **Gateway Service** calls an appropriate method of the **Partners API** via gRPC.  
   4. **Partners API** accepts and transfers the request to its service.  
   5. **Partners Service** retrieves data from its database and returns it.  

![2](./design/design-2-list-products.svg)

### 3. Customer Leaves an Order: Check  

   1. The customer adds products to their cart.  
   2. **Front App** sends a check request to the **API Gateway**.  
   3. **API Gateway** verifies the access token, parses and validates the request, and transfers it to its service.  
   4. **Gateway Service** calls an appropriate method of the **Partners API** via gRPC.  
   5. **Partners API** accepts and transfers the request to its service.  
   6. **Partners Service** checks product and partner availability in its database.  
   7. If everything is valid, the **Gateway Service** saves the checked request in the cache for 10 minutes.  

![3](./design/design-3-check-order.svg)

### 4. Customer Leaves an Order: Confirm 

   1. **Front App** sends a confirmation request to the **API Gateway** if the check request was successful.  
   2. **API Gateway** verifies the access token, parses and validates the request, and transfers it to its service.  
   3. **Gateway Service** checks its cache to verify if the checked order exists.  
   4. If the order exists, the **Gateway Service** calls an appropriate method of the **Orders API** via gRPC.  
   5. **Orders API** accepts and transfers the request to its service.  
   6. **Orders Service** checks its database to confirm if the chosen payment system is registered and retrieves its web checkout information.  
   7. **Orders Service** saves the order in its database and returns payment and callback information.  

![4](./design/design-4-confirm-order.svg)

### 5. Customer Pays for the Order  

   1. The customer is redirected to the **Payment System's** web checkout page to complete the payment.  
   2. **Payment System** sends a callback request to the **API Gateway**.  
   3. **API Gateway** parses and validates the request and transfers it to its service.  
   4. **Gateway Service** calls an appropriate method of the **Orders API** via gRPC.  
   5. **Orders API** accepts and transfers the request to its service.  
   6. **Orders Service** saves the payment details in its database and updates the order's status.  
   7. **Orders Service** sends the order products to the **Partners API** via the message broker.  
   8. **Partners API** accepts the products and transfers them to its service.  
   9. **Partners Service** retrieves the **Partner System URL** from its database.  
   10. **Partners Service** sends a request to the **Partner System** with the products to be prepared and readiness callback information.  

![5](./design/design-5-pay-order.svg)

### 6. Order is Ready  

   1. **Partner System** sends a callback about the order readiness to the **API Gateway**.  
   2. **API Gateway** parses and validates the request and transfers it to its service.  
   3. **Gateway Service** publishes the request to the **Orders API** via the message broker.  
   4. **Orders API** accepts the request and transfers it to its service.  
   5. **Orders Service** updates the order status to "ready."  

![6](./design/design-6-pickup-order.svg)

### 7. Deliverer Chooses the Order  

   1. The deliverer selects the order and sends an assignment request to the **API Gateway**.  
   2. **API Gateway** verifies the deliverer's access token, parses and validates the request, and transfers it to its service.  
   3. **Gateway Service** calls an appropriate method of the **Orders API** via gRPC.  
   4. **Orders API** accepts the request and transfers it to its service.  
   5. **Orders Service** updates the order status in its database.  
   6. **Orders Service** sends a message to the **Notifications API** via the message broker.  
   7. **Notifications API** accepts the message and transfers it to its service.  
   8. **Notifications Service** checks the sender agent in its database and retrieves their priority and settings.  
   9. **Notifications Service** saves the notification message in its database.  
   10. **Notifications Service** sends a push notification to the customer via a third-party push notification provider.  
   11. The **third-party push notification provider** locates the customer's device using the provided token and sends the push message.  

![7](./design/design-7-assign-order.svg)

---

## Installation

### Prerequisites
- **Go**

### Clone the Repository
```sh
git clone https://github.com/shahzodshafizod/gocloud.git
cd gocloud
```

### Run Unit & Integration Tests (Mock Tests)
```sh
make tests-run
```

### Run Tests with Coverage
```sh
make tests-cover
make tests-clear  # Clears generated test-cover.out files
```

---

## On-Premises Implementations

In this branch, I have implemented a variety of on-premises solutions, utilizing the following technologies:  

- **Authentication**: [**Keycloak**](https://www.keycloak.org) for identity and access management.
- **Caching**: [**Redis**](https://redis.io) for in-memory caching to enhance performance.
- **Email Services**: [**MailHog**](https://github.com/mailhog/MailHog) for local email testing. It can be easily changed to any production solution.  
- **Push Notifications**: [**OneSignal**](https://onesignal.com) for cross-platform push notifications (requires API setup).
- **NoSQL Databases**: [**MongoDB**](https://www.mongodb.com) for handling unstructured/semi-structured data.
- **SQL Databases**: [**PostgreSQL**](https://www.postgresql.org) as the relational database.
- **Message Brokers**: [**RabbitMQ**](https://www.rabbitmq.com) and [**NATS JetStream**](https://nats.io) for messaging and event-driven communication.
- **Storage**: [**MinIO**](https://min.io), an S3-compatible object storage solution.
- **Distributed Tracing**: [**OpenTelemetry**](https://opentelemetry.io) for observability, with [**Jaeger**](https://www.jaegertracing.io) for trace visualization.

---

### Preconfiguration of Push Notifications
1. Create a Web App on [OneSignal](https://onesignal.com) and copy the App ID.
2. Set `ONESIGNAL_APP_ID` in `configs/config.env`.
3. In [OneSignal Dashboard](https://dashboard.onesignal.com/), go to "Settings > Keys & IDs", generate an API Key, and set `ONESIGNAL_REST_API_KEY` in `configs/config.env`.
4. Update `appId` in `pkg/onprem/onesignal/index.html`.

## Setting Up and Running the Application (on Docker)

### Prerequisites
- **Docker & Docker Compose**
- **Docker Images**: Ensure all required Docker images are available locally or accessible from a container registry.

### 1. Create a Docker Network
```sh
docker network create gocloud
```
This ensures communication between services running in different Docker Compose files.

### 2. Start Dependencies
```sh
docker compose -f pkg/onprem/docker-compose.yml up -d
```

### 3. Set Up PostgreSQL Databases
```sh
docker exec -it postgres psql -d postgres -U odmin
# CREATE DATABASE notificationsdb;
# CREATE DATABASE ordersdb;
# CREATE DATABASE partnersdb;
```

### 4. Mapping of the IP address 127.0.0.1 (localhost) to the hostname delivery.local
```sh
echo "127.0.0.1 delivery.local" | sudo tee -a /etc/hosts
```

### 5. Configure MinIO
1. Access the MinIO console at [`http://delivery.local:9090`](http://delivery.local:9090) (credentials in `pkg/onprem/docker-compose.yml`).  
2. Go to "Access Keys", create a new access key, and update:
   - `MINIO_ACCESS_KEY` in `configs/config.env`
   - `MINIO_SECRET_KEY` in `configs/config.env`

### 6. Build Docker Images
```sh
make docker-build

# or build a separate image
docker build -t delivery/api -f cmd/api/Dockerfile .
```

### 7. Start Application Services
```sh
docker compose up -d

# or start a separate service:
docker compose up -d deliveryapi
```

### Access Application Components
- **API Documentation (Swagger):** http://delivery.local/docs/
- **Push Notification Demo:** http://delivery.local:4444/
  - Find "Subscription ID" at: `https://dashboard.onesignal.com/apps/{APP_ID}/subscriptions`
- **Email Sandbox (MailHog):** http://delivery.local:8025/
- **Keycloak Admin Panel:** http://delivery.local:8080/admin/delivery/console/
- **Tracing Console (Jaeger):** http://delivery.local:16686/
- **MinIO Console Panel:** http://delivery.local:9090/

---

## Running on Kubernetes (Minikube)
This guide provides step-by-step instructions to deploy the application on a local Kubernetes cluster using Minikube.

### Prerequisites

- **Minikube**: Ensure Minikube is installed and running.
- **kubectl**: Ensure `kubectl` is installed and configured to interact with your Minikube cluster.
- **Docker Images**: Ensure all required Docker images are available locally or accessible from a container registry.

### 1. Start Minikube and Enable Required Add-ons
```sh
minikube start
minikube addons enable ingress
minikube addons enable ingress-dns
```

### 2. Load Docker Images into Minikube Registry (or they will be downloaded when you run deployment manifests)
```sh
minikube image load quay.io/keycloak/keycloak:latest
minikube image load quay.io/minio/minio:latest
minikube image load postgres:latest
minikube image load rabbitmq:latest
minikube image load nats:latest
minikube image load jaegertracing/all-in-one:latest
minikube image load redis:latest
minikube image load mongo:latest
minikube image load nginx:latest
minikube image load mailhog/mailhog:latest

minikube image load delivery/api:latest
minikube image load delivery/notifications:latest
minikube image load delivery/partners:latest
minikube image load delivery/orders:latest
```

### 3. Create Namespace and ConfigMaps
```sh
kubectl create namespace gocloud

kubectl -n gocloud create configmap keycloakimport --from-file=pkg/onprem/keycloak/import/import.json
kubectl -n gocloud create configmap onesignal --from-file=pkg/onprem/onesignal/

kubectl -n gocloud create configmap configs --from-env-file=configs/config.env
kubectl -n gocloud create configmap migrationnotifications --from-file=migrations/notifications/
kubectl -n gocloud create configmap migrationorders --from-file=migrations/orders/
kubectl -n gocloud create configmap migrationpartners --from-file=migrations/partners/
```

### 4. Apply Kubernetes Dependencies Manifests
```sh
kubectl apply -f pkg/onprem/k8s/
```

### 5. Set Up PostgreSQL Databases
```sh
kubectl -n gocloud exec -it <POSTGRES_POD_NAME> -- psql -d postgres -U odmin
# CREATE DATABASE notificationsdb;
# CREATE DATABASE ordersdb;
# CREATE DATABASE partnersdb;
```

### 6. Map Minikube IP to Ingress Hostnames
```sh
echo "$(minikube ip) delivery.local" | sudo tee -a /etc/hosts
echo "$(minikube ip) jaeger.delivery.local" | sudo tee -a /etc/hosts
echo "$(minikube ip) keycloak.delivery.local" | sudo tee -a /etc/hosts
echo "$(minikube ip) mailhog.delivery.local" | sudo tee -a /etc/hosts
echo "$(minikube ip) minio.delivery.local" | sudo tee -a /etc/hosts
echo "$(minikube ip) push.delivery.local" | sudo tee -a /etc/hosts
```

### 7. Configure Keycloak Frontend URL

Set the Frontend URL in Keycloak to ensure correct confirmation callback URLs in emails:

1. Visit: http://keycloak.delivery.local/admin/delivery/console/
2. Go to **Realm Settings > General**.
3. Set **Frontend URL** to: `http://keycloak.delivery.local/`

### 8. Apply Kubernetes Application Manifests
```sh
kubectl apply -f k8s/
```

### Access Application Components
- **API Documentation (Swagger):** http://delivery.local/docs/
- **Push Notification Demo:** http://push.delivery.local/
  - Find "Subscription ID" at: `https://dashboard.onesignal.com/apps/{APP_ID}/subscriptions`
- **Email Sandbox (MailHog):** http://mailhog.delivery.local/
- **Keycloak Admin Panel:** http://keycloak.delivery.local/admin/delivery/console/
- **Tracing Console (Jaeger):** http://jaeger.delivery.local/
- **MinIO Console Panel:** http://minio.delivery.local/
