package data

import (
	"sync"
	"time"
)

// Payload represents incoming data to be processed
type Payload struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Type      string            `json:"type"`
	Values    map[string]interface{} `json:"values"`
}

// DataPoint represents a processed data point
type DataPoint struct {
	Timestamp time.Time
	Source    string
	Type      string
	Values    map[string]interface{}
}

// Service handles data storage and retrieval
type Service struct {
	mu         sync.RWMutex
	dataPoints []DataPoint
	callbacks  []func(DataPoint)
}

// NewService creates a new data service
func NewService() *Service {
	return &Service{
		dataPoints: make([]DataPoint, 0),
		callbacks:  make([]func(DataPoint), 0),
	}
}

// Store adds a new data payload to the service
func (s *Service) Store(payload Payload) error {
	// Convert payload to data point
	dataPoint := DataPoint{
		Timestamp: payload.Timestamp,
		Source:    payload.Source,
		Type:      payload.Type,
		Values:    payload.Values,
	}

	// If no timestamp provided, use current time
	if dataPoint.Timestamp.IsZero() {
		dataPoint.Timestamp = time.Now()
	}

	// Store the data point
	s.mu.Lock()
	s.dataPoints = append(s.dataPoints, dataPoint)
	// Call registered callbacks
	callbacks := s.callbacks
	s.mu.Unlock()

	// Execute callbacks outside the lock
	for _, callback := range callbacks {
		callback(dataPoint)
	}

	return nil
}

// GetData returns all data points within a given time frame
func (s *Service) GetData(since time.Time) []DataPoint {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]DataPoint, 0)
	for _, dp := range s.dataPoints {
		if dp.Timestamp.After(since) {
			result = append(result, dp)
		}
	}
	return result
}

// GetDataByType returns all data points of a specific type within a given time frame
func (s *Service) GetDataByType(dataType string, since time.Time) []DataPoint {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]DataPoint, 0)
	for _, dp := range s.dataPoints {
		if dp.Type == dataType && dp.Timestamp.After(since) {
			result = append(result, dp)
		}
	}
	return result
}

// SubscribeToNewData registers a callback function to be called when new data arrives
func (s *Service) SubscribeToNewData(callback func(DataPoint)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callbacks = append(s.callbacks, callback)
}

// GenerateMockData creates sample data for testing
func (s *Service) GenerateMockData(count int) {
	now := time.Now()
	
	for i := 0; i < count; i++ {
		s.Store(Payload{
			Timestamp: now.Add(time.Duration(-i) * time.Minute),
			Source:    "mock",
			Type:      "sales",
			Values: map[string]interface{}{
				"amount": 100 + i,
				"region": "Europe",
			},
		})
		
		s.Store(Payload{
			Timestamp: now.Add(time.Duration(-i) * time.Minute),
			Source:    "mock",
			Type:      "users",
			Values: map[string]interface{}{
				"active": 50 + i,
				"region": "North America",
			},
		})
	}
} 