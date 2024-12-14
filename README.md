# WebSocket Relay Server for Binance Futures Data

## Project Overview

This project implements a WebSocket relay server to collect real-time data from Binance Futures WebSocket streams for the top 150 `USDT perpetual` pairs based on 24-hour trading volume. The server efficiently relays this data to multiple clients in real time, allowing clients to subscribe to specific coin symbols.

The solution is designed with a focus on minimal latency, high throughput, dynamic subscription handling, and robust data integrity.

---

## Features

- **Real-time Data Relay:** Fetches and relays live cryptocurrency price updates for top `USDT perpetual` pairs (e.g., BTCUSDT, ETHUSDT).
- **Dynamic Client Subscriptions:** Clients can subscribe to specific coin symbols and only receive relevant updates.
- **High Throughput & Low Latency:** Optimized to ensure minimal delays between the Binance WebSocket and end-users.
- **Scalable Architecture:** Designed to handle large volumes of data and multiple simultaneous client connections.
- **Performance Benchmarking Client:** Includes a WebSocket client for comparing the server's performance with other implementations.

---

## How It Works

1. **Data Collection:**
   - Connects to Binance Futures WebSocket streams to fetch the top 150 `USDT perpetual` pairs based on 24-hour trading volume.
   - Processes and parses the incoming data efficiently.

2. **Data Relay:**
   - Acts as a relay server to push real-time data to connected clients.
   - Filters data based on client subscriptions to specific coin symbols.

3. **Dynamic Subscriptions:**
   - Clients can dynamically subscribe or unsubscribe to coin symbols without disrupting the server.

4. **Benchmarking Client:**
   - A WebSocket client is included to benchmark the server’s performance against other WebSocket implementations, comparing metrics such as latency and throughput.

---

## Installation

### Prerequisites

- Go should be installed.

### Steps

1. Clone this repository:

   ```bash
   git clone https://github.com/ahmetG23/binance-relay-server.git
   cd websocket-relay-server
   ```

2. Run the server:

   Run the following commands and wait until see the message "Server started at localhost:8080".

   ```bash
   cd server
   go install
   go run main.go
   ```

3. Run the client:

   On a new terminal:

   ```bash
   cd client
   go install
   go run main.go
   ```

---

## Usage

### Server

1. Start the server to establish a connection with Binance Futures WebSocket and prepare for client connections.
2. The server automatically retrieves the top 150 `USDT perpetual` pairs based on 24-hour trading volume.

### Clients

1. Connect to the WebSocket relay server as described above.
2. Receive real-time updates for the subscribed coin symbols.

### Testing

To test the first data receive time elapse, you can run the client with the following commands:

```bash
cd client
go test -bench=.
```

An example output is as follows (showing that the project receives data faster than getting data directly from Binance):

```txt
ORIGINAL first receive time elapse: 1.443661958s
SERVER first receive time elapse: 1.3844005s
```

## Acknowledgments

- Binance for providing the Futures WebSocket API.

## Notes

Please refer to the comments in the codes for more detailed information about the implementation.
