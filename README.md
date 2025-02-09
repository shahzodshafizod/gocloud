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

<details>

<summary><b style="font-size:20px">1. API Gateway (`cmd/api/`)</b></summary>

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

</details>
<details>

<summary><b style="font-size:20px">2. Push Notifications Service (`cmd/notifications/`)</b></summary>

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

</details>
<details>

<summary><b style="font-size:20px">3. Orders and Payments Service (`cmd/orders/`)</b></summary>

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

</details>
<details>

<summary><b style="font-size:20px">4. Partners and Products Service (`cmd/partners/`)</b></summary>

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

</details>

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

## Next Steps

- **Implement Authentication** (e.g., Keycloak, Amazon Cognito, Firebase Authentication, Azure AD)
- **Implement Caching** (e.g., Redis, Memcached, Amazon ElastiCache, Memorystore, Azure Cache for Redis)
- **Implement Email Services** (e.g., Gmail, Amazon SES, SendGrid/Mailgun/Mailjet on Google Cloud, Azure Communication Services)
- **Implement Push Notifications** (e.g., OneSignal, AWS SNS, FCM, Azure Notification Hubs)
- **Implement NoSQL Databases** (e.g., MongoDB, Amazon DynamoDB, Cloud Firestore, Cloud Datastore, Azure Cosmos DB)
- **Implement SQL Databases** (e.g., Amazon RDS, Amazon Aurora, Cloud SQL, Azure SQL Database, Azure Database for PostgreSQL)
- **Implement Message Brokers** (e.g., RabbitMQ, JetStream, Amazon SQS, Amazon MSK (Managed Streaming for Apache Kafka), Amazon MQ, Cloud Pub/Sub, Azure Service Bus)
- **Implement Storage** (e.g., MinIO, Amazon S3, Google Cloud Storage, Azure Blob Storage)
- **Implement Tracing** (e.g., OpenTelemetry, AWS X-Ray, Cloud Trace, Azure Monitor Application Insights)

---
