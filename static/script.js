let currentData = []; // Store the fetched data globally
let selectedNetwork = 'All'; // Initially show all networks

// Store chart instances globally
let latencyChart, packetLossChart, speedChart;

// Add an event listener to the dropdown
document.addEventListener('DOMContentLoaded', () => {
    fetchData();
    setInterval(fetchData, 30000); // Auto-refresh

    const networkSelect = document.getElementById('networkSelect');
    networkSelect.addEventListener('change', (event) => {
        selectedNetwork = event.target.value;
        updateCharts();
    });
});

async function fetchData() {
    const response = await fetch('/metrics');
    currentData = await response.json();

    populateNetworkDropdown();
    updateCharts();
}

function populateNetworkDropdown() {
    const networks = new Set(currentData.map(d => d.NetworkName));
    const dropdown = document.getElementById('networkSelect');
    dropdown.innerHTML = ''; // Clear existing options

    // Add an "All" option
    const allOption = document.createElement('option');
    allOption.value = 'All';
    allOption.text = 'All';
    dropdown.add(allOption);

    // Add options for each network
    networks.forEach(network => {
        const option = document.createElement('option');
        option.value = network;
        option.text = network;
        dropdown.add(option);
    });
}

function updateCharts() {
    latencyChart = updateChart("Latency", latencyChart, 'latencyChart', 'ms', [['PingMillis', 'rgb(23, 103, 224)'], ['JitterMillis', 'rgb(220, 146, 135)']]);
    packetLossChart = updateChart("Packet Loss", packetLossChart, 'packetLossChart', '%', [['PacketLossPercentage', 'rgb(252, 83, 72)']]);
    speedChart = updateChart("Speed", speedChart, 'speedChart', 'Mbps', [['DownloadMbps', 'rgb(114, 241, 235)'], ['UploadMbps', 'rgb(200, 138, 251)']]);
}

function updateChart(chartTitle, chart, chartId, unit, labelsAndColors) {
    const data = selectedNetwork === 'All'
        ? currentData
        : currentData.filter(d => d.NetworkName === selectedNetwork);

    const ctx = document.getElementById(chartId).getContext('2d');

    // Destroy existing chart if it exists
    if (chart) chart.destroy();

    // Format time labels with better precision
    const formatTime = (epochMillis) => {
        const date = new Date(epochMillis);
        return date.toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    // Create datasets for metrics
    const metricDatasets = labelsAndColors.map(([label, color]) => ({
        label: label,
        data: data.map(d => d.IsOnline ? d[label] : null),
        fill: true,
        borderColor: color,
        backgroundColor: color.replace('rgb', 'rgba').replace(')', ', 0.2)'),
        tension: 0.1,
        pointRadius: 3,
        pointHoverRadius: 5
    }));

    // Add offline indicator dataset
    const offlineDataset = {
        label: 'Offline',
        data: data.map(d => d.IsOnline ? null : unit === '%' ? 100 : unit === 'Mbps' ? 0 : unit === 'ms' ? 0 : 0),
        fill: true,
        backgroundColor: 'rgba(128, 128, 128, 0.2)',
        borderColor: 'rgba(128, 128, 128, 0.5)',
        pointRadius: 0,
        tension: 0
    };

    // Create new chart
    chart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: data.map(d => formatTime(d.EpochMillis)),
            datasets: [...metricDatasets, offlineDataset]
        },
        options: {
            responsive: true,
            interaction: {
                intersect: false,
                mode: 'index'
            },
            plugins: {
                title: {
                    display: true,
                    text: chartTitle,
                    font: {
                        size: 16
                    }
                },
                tooltip: {
                    callbacks: {
                        label: function (context) {
                            if (context.dataset.label === 'Offline' && context.raw !== null) {
                                return 'Network Offline';
                            }
                            return context.dataset.label + ': ' +
                                (context.raw !== null ? context.raw.toFixed(2) + ' ' + unit : 'Offline');
                        }
                    }
                },
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: 'Time'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: unit
                    },
                    beginAtZero: true,
                }
            },
        }
    });

    return chart;
}
