# Real-Time Analytics Dashboard

A real-time data processing pipeline and interactive dashboard for monitoring business metrics and generating insights, built with Go.

## 🚀 Features

- **Real-time data ingestion** via RESTful API
- **Streaming data processing** with metric calculations
- **WebSocket communication** for instant dashboard updates
- **Interactive visualizations** with Chart.js
- **Modular architecture** for easy extension

## 📋 Components

The application is structured into several key components:

- **Data Service**: Handles data storage and retrieval
- **Processing Engine**: Computes metrics and analyzes incoming data
- **Visualization Manager**: Manages WebSocket connections and real-time updates
- **Web Dashboard**: Interactive UI with real-time charts and tables

## 🔧 Technology Stack

- **Backend**: Go (Golang)
- **API Framework**: Gin Web Framework
- **Real-time Communication**: WebSockets (gorilla/websocket)
- **Frontend**: HTML, CSS, JavaScript
- **Visualization**: Chart.js
- **Data Processing**: Custom Go implementation

## 📥 Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Setup

1. Clone the repository:
```bash
git clone https://github.com/klejdi94/real-time-analytics.git
cd real-time-analytics
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run cmd/server/main.go
```

4. Open your browser and navigate to:
```
http://localhost:8080
```

## 🧪 Generating Test Data

The project includes a data generator tool to populate the dashboard with sample metrics:

```bash
# Generate 100 data points at 3-second intervals
go run cmd/generator/main.go --count 100 --interval 3

# Generate an infinite stream of data
go run cmd/generator/main.go

# Send data in batches of 5
go run cmd/generator/main.go --batch 5
```

## 📊 Dashboard Usage

The dashboard is divided into three main tabs:

1. **Data Engineering**: Shows total events, event distribution, and recent data table
2. **Stream Processing**: Displays time series data for sales and user metrics
3. **Visualization**: Features region distribution, metrics comparison, and real-time metrics

## 🏗️ Project Structure

```
├── cmd
│   ├── generator        # Test data generator tool
│   └── server           # Main application server
├── pkg
│   ├── data             # Data storage and management
│   ├── processing       # Data processing and analysis
│   └── visualization    # WebSocket and real-time updates
└── web
    ├── static           # Static assets (JS, CSS)
    └── templates        # HTML templates
```

## 🚧 Extending the Dashboard

### Adding a New Metric

1. Define the metric in `pkg/processing/processor.go`
2. Add visualization in the appropriate tab in `web/templates/index.html`
3. Update the JavaScript handler in `web/static/js/dashboard.js`

### Adding a New Data Type

1. Extend the data model in `pkg/data/model.go`
2. Add processing logic in `pkg/processing/processor.go`
3. Update the visualization components as needed

## 📄 API Reference

### Data Ingestion

```
POST /api/ingest
Content-Type: application/json

{
  "timestamp": "2023-04-01T12:00:00Z",
  "source": "web-app",
  "type": "sales",
  "values": {
    "amount": 150,
    "units": 5,
    "region": "North America"
  }
}
```

### Metrics Retrieval

```
GET /api/metrics
```

### WebSocket Connection

```
WS /api/ws
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details. 