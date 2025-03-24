# Real-time Analytics Dashboard

A real-time data processing pipeline and interactive dashboard for monitoring business metrics and generating insights, built with Go.

## Features

- **Data Engineering**: Collect and store event data from various sources
- **Stream Processing**: Process incoming data in real-time to generate metrics
- **Visualization**: Interactive dashboard with charts and tables
- **Real-time Updates**: WebSocket connection for live data updates

## Project Structure

```
analytic-dashboard/
├── cmd/
│   └── server/         # Main application entry point
├── pkg/
│   ├── data/           # Data models and storage
│   ├── processing/     # Data processing and metrics calculation
│   └── visualization/  # WebSocket and visualization handling
├── web/
│   ├── static/         # Static assets (CSS, JS)
│   └── templates/      # HTML templates
└── go.mod              # Go module file
```

## Requirements

- Go 1.21 or higher
- Modern web browser with JavaScript enabled

## Installation

1. Install Go from [https://golang.org/dl/](https://golang.org/dl/)

2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/analytic-dashboard.git
   cd analytic-dashboard
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Application

1. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

## API Endpoints

- `GET /`: Main dashboard interface
- `POST /api/ingest`: Data ingestion endpoint
- `GET /api/metrics`: Get current metrics
- `GET /api/ws`: WebSocket connection for real-time updates

## Data Ingestion

To send data to the dashboard, make a POST request to the `/api/ingest` endpoint with JSON data in the following format:

```json
{
  "timestamp": "2023-03-15T12:34:56Z",
  "source": "sales-system",
  "type": "sales",
  "values": {
    "amount": 125,
    "region": "Europe"
  }
}
```

The `timestamp` field is optional. If not provided, the current time will be used.

## Development

### Adding New Data Types

1. Define the data structure in `pkg/data/data.go`
2. Add processing logic in `pkg/processing/processor.go`
3. Update the visualization in the frontend JavaScript

### Customizing the Dashboard

The dashboard UI can be customized by modifying the HTML, CSS, and JavaScript files in the `web` directory:

- `web/templates/index.html`: Main dashboard layout
- `web/static/css/styles.css`: Styling and layout
- `web/static/js/dashboard.js`: Dashboard functionality and charts

## License

MIT 