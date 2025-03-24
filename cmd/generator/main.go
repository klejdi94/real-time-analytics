package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/user/analytic-dashboard/pkg/data"
)

var (
	endpoint  = flag.String("endpoint", "http://localhost:8080/api/ingest", "API endpoint URL")
	interval  = flag.Int("interval", 3, "Interval between data points in seconds")
	count     = flag.Int("count", 0, "Number of data points to generate (0 for infinite)")
	batchSize = flag.Int("batch", 1, "Number of data points to send in each batch")
)

func main() {
	flag.Parse()

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("Starting data generator:\n")
	fmt.Printf("- Endpoint: %s\n", *endpoint)
	fmt.Printf("- Interval: %d seconds\n", *interval)
	fmt.Printf("- Count: %s\n", func() string {
		if *count == 0 {
			return "infinite"
		}
		return fmt.Sprintf("%d", *count)
	}())
	fmt.Printf("- Batch Size: %d\n\n", *batchSize)

	// Generate and send data
	i := 0
	for {
		// Break if we've reached the count
		if *count > 0 && i >= *count {
			break
		}

		// Generate a batch of data
		batch := make([]data.Payload, 0, *batchSize)
		for j := 0; j < *batchSize; j++ {
			payload := generateRandomPayload()
			batch = append(batch, payload)
			i++
		}

		// Send the batch
		if err := sendBatch(batch); err != nil {
			log.Printf("Error sending batch: %v\n", err)
		} else {
			log.Printf("Sent batch of %d data points (total: %d)\n", *batchSize, i)
		}

		// Sleep for the interval
		time.Sleep(time.Duration(*interval) * time.Second)
	}

	fmt.Println("Data generation complete.")
}

// generateRandomPayload creates a random data payload
func generateRandomPayload() data.Payload {
	now := time.Now()
	
	// Randomly choose between sales and users events
	eventTypes := []string{"sales", "users"}
	eventType := eventTypes[rand.Intn(len(eventTypes))]
	
	// Randomly choose a region
	regions := []string{"North America", "Europe", "Asia", "South America", "Africa", "Oceania"}
	region := regions[rand.Intn(len(regions))]
	
	// Create values based on event type
	values := make(map[string]interface{})
	values["region"] = region
	
	if eventType == "sales" {
		// Generate a random sales amount between 50 and 500
		values["amount"] = 50 + rand.Intn(450)
		values["units"] = 1 + rand.Intn(10)
		values["channel"] = []string{"online", "in-store", "partner"}[rand.Intn(3)]
	} else {
		// Generate random user metrics
		values["active"] = 100 + rand.Intn(900)
		values["new"] = 10 + rand.Intn(90)
		values["returning"] = 50 + rand.Intn(200)
	}
	
	// Create the payload
	return data.Payload{
		Timestamp: now,
		Source:    "generator",
		Type:      eventType,
		Values:    values,
	}
}

// sendBatch sends a batch of data to the API
func sendBatch(batch []data.Payload) error {
	// For a single item, send a regular payload
	if len(batch) == 1 {
		jsonData, err := json.Marshal(batch[0])
		if err != nil {
			return fmt.Errorf("error marshaling JSON: %v", err)
		}
		
		resp, err := http.Post(*endpoint, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("error sending request: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		
		return nil
	}
	
	// For multiple items, send as an array
	jsonData, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}
	
	resp, err := http.Post(*endpoint+"/batch", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	return nil
} 