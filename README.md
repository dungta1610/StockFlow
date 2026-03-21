# StockFlow

## Overview

A mini e-commerce backend focused on **order lifecycle**, **inventory consistency**, and **modular backend design**:

- Manage core resources such as **users**, **products**, and **warehouses**
- Create and track **warehouse orders** with item-level data
- Handle **inventory adjustment** and inspect **inventory transaction history**
- Simulate a basic **payment flow** with checkout and callback endpoints
- Use **PostgreSQL transactions** for important write flows
- Add **Redis-based rate limiting middleware** to protect the API surface
- Organize the codebase using a **module-first architecture** with **handwritten SQL**

This project is built as a practical backend showcase for an **Intern Backend portfolio**. The goal is not to build a full production marketplace, but to model the backend foundation of a **mini e-commerce / warehouse order management system** with realistic service boundaries.

## Tech Stack

- **Go**
- **Gin** (HTTP framework)
- **PostgreSQL 16**
- **pgx / pgxpool**
- **Redis 7**
- **Docker / Docker Compose**
- **Handwritten SQL** (no ORM)

## Features

### 1) Product Management

- Create products
- Get product detail by ID
- List products with paging
- Product data is stored in **PostgreSQL** and exposed through thin **Gin handlers**

### 2) Warehouse Management

- Create warehouses
- Get warehouse detail by ID
- List warehouses with paging
- Warehouses act as **inventory locations** for stock operations and order fulfillment

### 3) User Management

- Create users
- Get user detail by ID
- Update user information
- List users with paging

### 4) Inventory Management

- Adjust stock quantities through a dedicated inventory endpoint
- Read inventory detail for a product in a warehouse
- Inspect inventory transaction history
- Inventory logic is separated into **entity models**, **use cases**, and **SQL storage**

### 5) Order Flow

- Create orders with order items
- Get order detail by ID
- List orders with paging/filter support in the module design
- Cancel an order
- Expire an order
- Important order writes are implemented through **transactional storage code**

### 6) Payment Flow

- Create checkout payment records
- Receive payment callback updates
- Get payment detail by ID
- List payments with paging
- Payment status transitions are modeled in a separate **module**

### 7) Redis Rate Limiting

- Global **Gin middleware** backed by **Redis**
- Helps mitigate **API abuse** during local testing and demo scenarios

### 8) Dockerized Local Environment

- Runs **PostgreSQL** and **Redis** with **Docker Compose**
- Application configuration is driven through **.env**

### 9) Experimental Outbox Module

- The repository contains an **outbox module** with enqueue/list/mark endpoints
- In the current source state, this module exists as infrastructure groundwork and is not yet a core integrated business flow

## Database Schema

![Database Schema](StockFlow.png)

## System Architecture

This backend can be viewed as a mini e-commerce architecture centered on **product catalog**, **warehouse stock**, **customer orders**, and **payment processing**.

```text
┌──────────────────────────────────────────────────────────────────┐
│                          Client / Postman                        │
└──────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                       Gin HTTP API Layer                         │
│  /users  /products  /warehouses  /inventories  /orders  /payments │
└──────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                     Module-First Application                     │
│                                                                  │
│  User Module        Product Module        Warehouse Module       │
│  Inventory Module   Order Module          Payment Module         │
│  Outbox Module (experimental)                                   │
└──────────────────────────────────────────────────────────────────┘
                               │
                ┌──────────────┴──────────────┐
                ▼                             ▼
┌──────────────────────────────┐   ┌──────────────────────────────┐
│ PostgreSQL                   │   │ Redis                        │
│ - users                      │   │ - rate limit counters        │
│ - products                   │   │ - short-lived cache/limits   │
│ - warehouses                 │   └──────────────────────────────┘
│ - inventory / transactions   │
│ - orders / order_items       │
│ - payments                   │
│ - outbox_events              │
└──────────────────────────────┘
```

## Domain View

- Users place or own **orders**
- Products are stored in **warehouses**
- Inventory tracks available stock per **warehouse/product** pair
- Orders contain one or more **order items**
- Payments track **checkout** and **callback** state for orders
- Outbox is intended for asynchronous event recording, but is currently infrastructure-level rather than a fully integrated runtime flow

## Request Flow

1. Client calls a Gin endpoint
2. Handler validates/binds request and delegates to **biz/**
3. **biz/** executes the use case and calls **storage/**
4. **storage/** runs handwritten SQL against **PostgreSQL**
5. **Redis middleware** applies request throttling before protected routes are processed

## Mini E-commerce Scope

This project models the backend foundation of a simplified commerce system:

- **Catalog**: products
- **Operations**: warehouses and inventory control
- **Sales flow**: order creation and status changes
- **Payment flow**: checkout and callback
- **Platform concerns**: rate limiting, SQL transactions, modular structure

It is best described as a warehouse-oriented mini e-commerce backend, not a full marketplace. The system emphasizes backend engineering concerns more than front-end customer features.

## Main API Groups

### Users

- POST /users
- GET /users
- GET /users/:id
- PUT /users/:id

### Products

- POST /products
- GET /products
- GET /products/:id

### Warehouses

- POST /warehouses
- GET /warehouses
- GET /warehouses/:id

### Inventories

- POST /inventories/adjust
- GET /inventories/detail
- GET /inventories/transactions

### Orders

- POST /orders
- GET /orders
- GET /orders/:id
- POST /orders/:id/cancel
- POST /orders/:id/expire

### Payments

- POST /payments/checkout
- POST /payments/callback
- GET /payments
- GET /payments/:id

### Outbox (experimental)

- POST /outbox/events
- GET /outbox/events
- POST /outbox/events/:id/processed
- POST /outbox/events/:id/failed

## Project Structure

```text
stockflow/
├── main.go                              # Composition root: init Postgres, Redis, Gin, middleware, routes
├── docker-compose.yml                   # PostgreSQL + Redis orchestration
├── .env                                 # Local runtime configuration
├── go.mod / go.sum                      # Go module files
├── StockFlow.png                        # Project image / schema-related asset
│
├── component/
│   ├── postgres/
│   │   └── postgres.go                  # PostgreSQL connection helper (pgxpool)
│   ├── redis/
│   │   └── redis.go                     # Redis connection helper
│   └── ratelimit/
│       └── limiter.go                   # Redis-based rate limiter
│
├── middleware/
│   └── ratelimit.go                     # Gin middleware wrapper for limiter
│
├── db/
│   └── init/                            # Database init folder (currently no SQL files in repository)
│
└── module/
    ├── user/
    │   ├── model/
    │   ├── biz/
    │   ├── storage/
    │   └── transport/gin/
    │
    ├── product/
    │   ├── model/
    │   ├── biz/
    │   ├── storage/
    │   └── transport/gin/
    │
    ├── warehouse/
    │   ├── model/
    │   ├── biz/
    │   ├── storage/
    │   └── transport/gin/
    │
    ├── inventory/
    │   ├── model/
    │   ├── biz/
    │   ├── storage/
    │   └── transport/gin/
    │
    ├── order/
    │   ├── model/
    │   ├── biz/
    │   ├── storage/
    │   └── transport/gin/
    │
    ├── payment/
    │   ├── model/
    │   ├── biz/
    │   ├── storage/
    │   └── transport/gin/
    │
    └── outbox/
        ├── model/
        ├── biz/
        ├── storage/
        └── transport/gin/
```

## Module-First Layering

Each module follows a consistent internal structure:

**transport/gin**  -> HTTP handlers and route registration
**storage**        -> handwritten SQL and DB access
**biz**            -> use-case orchestration and business rules
**model**          -> entities, DTOs, filters, paging, errors

This keeps handlers thin, avoids **ORM coupling**, and makes it easier to reason about each domain separately.

## Clean Architecture View

```text
┌─────────────────────────────────────────────┐
│         Transport (Driver)                  │
│      module/*/transport/gin                 │
└─────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────┐
│         Storage (Adapters)                  │
│           module/*/storage                  │
└─────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────┐
│          Biz (Use-Case)                     │
│             module/*/biz                    │
└─────────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────┐
│         Models (Entities)                   │
│            module/*/model                   │
└─────────────────────────────────────────────┘
```

## Running Locally

### 1) Start infrastructure

```bash
docker compose up -d
```

### 2) Set environment variables

Default .env in the repository:

```env
DB_DSN=postgres://postgres:postgres@localhost:5432/stockflow?sslmode=disable
REDIS_ADDR=127.0.0.1:6379
PORT=8080
```

### 3) Run the application

```bash
go run main.go
```

### 4) Health check

```http
GET /health
```

## Notes

- The project intentionally uses **handwritten SQL** instead of an ORM
- **PostgreSQL** is the source of truth for persistence design
- **Redis** is currently used for rate limiting
- Some repository parts, especially **database bootstrap/migrations**, are still incomplete and may require manual schema setup before all endpoints work end-to-end
- The **outbox module** exists in the source tree, but its business integration level is not yet as complete as the core modules
