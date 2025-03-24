package processing

import (
	"sync"
	"time"

	"github.com/klejdi94/real-time-analytics/pkg/data"
)

// Metrics represents computed business metrics
type Metrics struct {
	TotalEvents    int                          `json:"totalEvents"`
	EventsByType   map[string]int               `json:"eventsByType"`
	RecentValues   map[string]interface{}       `json:"recentValues"`
	TimeSeriesData map[string][]TimeSeriesPoint `json:"timeSeriesData"`
}

// TimeSeriesPoint represents a point in a time series
type TimeSeriesPoint struct {
	Timestamp time.Time              `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
}

// Processor is responsible for processing incoming data
type Processor struct {
	dataService *data.Service
	metrics     Metrics
	mu          sync.RWMutex
	stop        chan struct{}
	listeners   []func(Metrics)
}

// NewProcessor creates a new data processor
func NewProcessor(dataService *data.Service) *Processor {
	p := &Processor{
		dataService: dataService,
		metrics: Metrics{
			EventsByType:   make(map[string]int),
			RecentValues:   make(map[string]interface{}),
			TimeSeriesData: make(map[string][]TimeSeriesPoint),
		},
		stop:      make(chan struct{}),
		listeners: make([]func(Metrics), 0),
	}

	// Subscribe to new data
	dataService.SubscribeToNewData(p.processDataPoint)

	return p
}

// Start begins background processing
func (p *Processor) Start() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.updateMetrics()
		case <-p.stop:
			return
		}
	}
}

// Stop halts the processor
func (p *Processor) Stop() {
	close(p.stop)
}

// GetMetrics returns the current metrics
func (p *Processor) GetMetrics() (Metrics, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Return a copy to prevent race conditions
	metricsCopy := Metrics{
		TotalEvents:    p.metrics.TotalEvents,
		EventsByType:   make(map[string]int),
		RecentValues:   make(map[string]interface{}),
		TimeSeriesData: make(map[string][]TimeSeriesPoint),
	}

	for k, v := range p.metrics.EventsByType {
		metricsCopy.EventsByType[k] = v
	}

	for k, v := range p.metrics.RecentValues {
		metricsCopy.RecentValues[k] = v
	}

	for k, v := range p.metrics.TimeSeriesData {
		metricsCopy.TimeSeriesData[k] = make([]TimeSeriesPoint, len(v))
		copy(metricsCopy.TimeSeriesData[k], v)
	}

	return metricsCopy, nil
}

// AddListener registers a callback function to be called when metrics are updated
func (p *Processor) AddListener(listener func(Metrics)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.listeners = append(p.listeners, listener)
}

// processDataPoint processes a single data point
func (p *Processor) processDataPoint(point data.DataPoint) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Update total events
	p.metrics.TotalEvents++

	// Update events by type
	p.metrics.EventsByType[point.Type]++

	// Update recent values
	for k, v := range point.Values {
		key := point.Type + "." + k
		p.metrics.RecentValues[key] = v
	}

	// Update time series data
	timeSeriesKey := point.Type
	timeSeriesPoint := TimeSeriesPoint{
		Timestamp: point.Timestamp,
		Values:    point.Values,
	}

	if _, exists := p.metrics.TimeSeriesData[timeSeriesKey]; !exists {
		p.metrics.TimeSeriesData[timeSeriesKey] = make([]TimeSeriesPoint, 0)
	}

	// Add new point and maintain a maximum of 100 points
	p.metrics.TimeSeriesData[timeSeriesKey] = append(p.metrics.TimeSeriesData[timeSeriesKey], timeSeriesPoint)
	if len(p.metrics.TimeSeriesData[timeSeriesKey]) > 100 {
		p.metrics.TimeSeriesData[timeSeriesKey] = p.metrics.TimeSeriesData[timeSeriesKey][1:]
	}

	// Notify listeners
	metrics := p.metrics
	go func() {
		for _, listener := range p.listeners {
			listener(metrics)
		}
	}()
}

// updateMetrics periodically refreshes metrics from the data service
func (p *Processor) updateMetrics() {
	// Get data from the last hour
	since := time.Now().Add(-1 * time.Hour)
	recentData := p.dataService.GetData(since)

	p.mu.Lock()
	defer p.mu.Unlock()

	// Reset metrics
	p.metrics.TotalEvents = len(recentData)
	p.metrics.EventsByType = make(map[string]int)

	// Recalculate events by type
	for _, point := range recentData {
		p.metrics.EventsByType[point.Type]++
	}

	// Notify listeners
	metrics := p.metrics
	go func() {
		for _, listener := range p.listeners {
			listener(metrics)
		}
	}()
}
