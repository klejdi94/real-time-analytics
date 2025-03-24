document.addEventListener('DOMContentLoaded', function() {
    // Tab switching
    const tabButtons = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabId = button.getAttribute('data-tab');
            
            // Update active button
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            
            // Update active content
            tabContents.forEach(content => content.classList.remove('active'));
            document.getElementById(tabId).classList.add('active');
        });
    });

    // Initialize charts
    const charts = {
        eventsByType: initPieChart('events-by-type-chart', 'Events by Type'),
        salesTime: initTimeSeriesChart('sales-time-chart', 'Sales Over Time'),
        usersTime: initTimeSeriesChart('users-time-chart', 'Active Users Over Time'),
        regionChart: initBarChart('region-chart', 'Region Distribution'),
        metricsChart: initLineChart('metrics-chart', 'Metrics Comparison')
    };

    // Connect to WebSocket
    const socket = connectWebSocket();
    const statusIndicator = document.getElementById('status-text');
    const statusText = document.getElementById('status-text');

    // Initialize WebSocket connection
    function connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/api/ws`;
        const socket = new WebSocket(wsUrl);

        socket.onopen = function() {
            console.log('WebSocket connection established');
            statusText.innerText = 'Connected';
            statusIndicator.parentElement.classList.add('connected');
            statusIndicator.parentElement.classList.remove('disconnected');
        };

        socket.onclose = function() {
            console.log('WebSocket connection closed');
            statusText.innerText = 'Disconnected';
            statusIndicator.parentElement.classList.remove('connected');
            statusIndicator.parentElement.classList.add('disconnected');
            
            // Try to reconnect after 5 seconds
            setTimeout(() => {
                statusText.innerText = 'Reconnecting...';
                statusIndicator.parentElement.classList.remove('connected');
                statusIndicator.parentElement.classList.remove('disconnected');
                connectWebSocket();
            }, 5000);
        };

        socket.onerror = function(error) {
            console.error('WebSocket error:', error);
            statusText.innerText = 'Connection Error';
            statusIndicator.parentElement.classList.remove('connected');
            statusIndicator.parentElement.classList.add('disconnected');
        };

        socket.onmessage = function(event) {
            try {
                const data = JSON.parse(event.data);
                updateDashboard(data);
            } catch (e) {
                console.error('Error parsing WebSocket message:', e);
            }
        };

        return socket;
    }

    // Update dashboard with new data
    function updateDashboard(data) {
        // Update total events
        document.getElementById('total-events').innerText = data.totalEvents || 0;
        
        // Update events by type chart
        if (data.eventsByType) {
            const labels = Object.keys(data.eventsByType);
            const values = Object.values(data.eventsByType);
            updatePieChart(charts.eventsByType, labels, values);
        }
        
        // Update recent values table
        if (data.recentValues) {
            updateRecentValuesTable(data.recentValues);
        }
        
        // Update time series charts
        if (data.timeSeriesData) {
            if (data.timeSeriesData.sales) {
                updateTimeSeriesChart(charts.salesTime, data.timeSeriesData.sales, 'amount');
            }
            
            if (data.timeSeriesData.users) {
                updateTimeSeriesChart(charts.usersTime, data.timeSeriesData.users, 'active');
            }
            
            // Update region distribution
            updateRegionChart(data.timeSeriesData);
            
            // Update metrics comparison
            updateMetricsComparison(data.timeSeriesData);
        }
        
        // Update real-time metrics
        if (data.recentValues) {
            updateRealTimeMetrics(data.recentValues);
        }
    }

    // Initialize a pie chart
    function initPieChart(canvasId, title) {
        const ctx = document.getElementById(canvasId).getContext('2d');
        return new Chart(ctx, {
            type: 'pie',
            data: {
                labels: [],
                datasets: [{
                    data: [],
                    backgroundColor: [
                        '#3498db',
                        '#2ecc71',
                        '#e74c3c',
                        '#f39c12',
                        '#9b59b6',
                        '#1abc9c'
                    ],
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom'
                    },
                    title: {
                        display: false,
                        text: title
                    }
                }
            }
        });
    }

    // Initialize a time series chart
    function initTimeSeriesChart(canvasId, title) {
        const ctx = document.getElementById(canvasId).getContext('2d');
        return new Chart(ctx, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                    label: title,
                    data: [],
                    borderColor: '#3498db',
                    backgroundColor: 'rgba(52, 152, 219, 0.1)',
                    borderWidth: 2,
                    fill: true,
                    tension: 0.4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        type: 'time',
                        time: {
                            unit: 'minute',
                            displayFormats: {
                                minute: 'HH:mm'
                            }
                        },
                        title: {
                            display: true,
                            text: 'Time'
                        }
                    },
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Value'
                        }
                    }
                },
                plugins: {
                    legend: {
                        display: false
                    },
                    title: {
                        display: true,
                        text: title
                    }
                }
            }
        });
    }

    // Initialize a bar chart
    function initBarChart(canvasId, title) {
        const ctx = document.getElementById(canvasId).getContext('2d');
        return new Chart(ctx, {
            type: 'bar',
            data: {
                labels: [],
                datasets: [{
                    label: title,
                    data: [],
                    backgroundColor: 'rgba(52, 152, 219, 0.7)',
                    borderColor: '#3498db',
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: true
                    }
                },
                plugins: {
                    legend: {
                        display: false
                    },
                    title: {
                        display: true,
                        text: title
                    }
                }
            }
        });
    }

    // Initialize a line chart
    function initLineChart(canvasId, title) {
        const ctx = document.getElementById(canvasId).getContext('2d');
        return new Chart(ctx, {
            type: 'line',
            data: {
                labels: [],
                datasets: [
                    {
                        label: 'Sales',
                        data: [],
                        borderColor: '#3498db',
                        backgroundColor: 'transparent',
                        borderWidth: 2,
                        tension: 0.4
                    },
                    {
                        label: 'Users',
                        data: [],
                        borderColor: '#2ecc71',
                        backgroundColor: 'transparent',
                        borderWidth: 2,
                        tension: 0.4
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: true
                    }
                },
                plugins: {
                    title: {
                        display: true,
                        text: title
                    }
                }
            }
        });
    }

    // Update pie chart with new data
    function updatePieChart(chart, labels, values) {
        chart.data.labels = labels;
        chart.data.datasets[0].data = values;
        chart.update();
    }

    // Update time series chart with new data
    function updateTimeSeriesChart(chart, timeSeriesData, valueKey) {
        if (!timeSeriesData || timeSeriesData.length === 0) return;
        
        const labels = [];
        const data = [];
        
        timeSeriesData.forEach(point => {
            // Parse the timestamp if it's a string
            const timestamp = typeof point.timestamp === 'string' 
                ? new Date(point.timestamp) 
                : point.timestamp;
            
            labels.push(timestamp);
            data.push(point.values[valueKey] || 0);
        });
        
        chart.data.labels = labels;
        chart.data.datasets[0].data = data.map((value, index) => ({
            x: labels[index],
            y: value
        }));
        
        chart.update();
    }

    // Update region chart
    function updateRegionChart(timeSeriesData) {
        const regionCounts = {};
        
        // Count occurrences of each region
        for (const type in timeSeriesData) {
            timeSeriesData[type].forEach(point => {
                if (point.values.region) {
                    const region = point.values.region;
                    regionCounts[region] = (regionCounts[region] || 0) + 1;
                }
            });
        }
        
        // Update chart
        const labels = Object.keys(regionCounts);
        const values = Object.values(regionCounts);
        
        charts.regionChart.data.labels = labels;
        charts.regionChart.data.datasets[0].data = values;
        charts.regionChart.update();
    }

    // Update metrics comparison chart
    function updateMetricsComparison(timeSeriesData) {
        const labels = [];
        const salesData = [];
        const usersData = [];
        
        // Get timestamps from sales data if available
        if (timeSeriesData.sales && timeSeriesData.sales.length > 0) {
            timeSeriesData.sales.forEach(point => {
                const timestamp = new Date(point.timestamp);
                const formattedTime = timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
                labels.push(formattedTime);
                salesData.push(point.values.amount || 0);
            });
        }
        
        // Get users data for the same timestamps
        if (timeSeriesData.users && timeSeriesData.users.length > 0) {
            // Align with the timestamps we already have
            for (let i = 0; i < labels.length; i++) {
                if (i < timeSeriesData.users.length) {
                    usersData.push(timeSeriesData.users[i].values.active || 0);
                } else {
                    usersData.push(0);
                }
            }
        }
        
        // Update chart
        charts.metricsChart.data.labels = labels;
        charts.metricsChart.data.datasets[0].data = salesData;
        charts.metricsChart.data.datasets[1].data = usersData;
        charts.metricsChart.update();
    }

    // Update recent values table
    function updateRecentValuesTable(recentValues) {
        const tableBody = document.querySelector('#recent-data-table tbody');
        tableBody.innerHTML = '';
        
        for (const key in recentValues) {
            const [type, metric] = key.split('.');
            const value = recentValues[key];
            
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${type}</td>
                <td>${metric}</td>
                <td>${value}</td>
            `;
            
            tableBody.appendChild(row);
        }
    }

    // Update real-time metrics
    function updateRealTimeMetrics(recentValues) {
        const container = document.getElementById('real-time-metrics');
        container.innerHTML = '';
        
        for (const key in recentValues) {
            const value = recentValues[key];
            
            const metricItem = document.createElement('div');
            metricItem.className = 'metric-item';
            metricItem.innerHTML = `
                <span class="metric-name">${key}</span>
                <span class="metric-value">${value}</span>
            `;
            
            container.appendChild(metricItem);
        }
    }

    // Fetch initial data from API
    fetch('/api/metrics')
        .then(response => response.json())
        .then(data => {
            updateDashboard(data);
        })
        .catch(error => {
            console.error('Error fetching initial data:', error);
        });
}); 