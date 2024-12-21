# display this help message
@help:
    echo "Usage: just <recipe-name>"
    echo ""
    just --list --unsorted

# Build the program
[group('build')]
build:
    go build -o bin/netmon

# Clean up build artifacts
[group('build')]
clean:
    rm -f bin/netmon

# Install the service
[group('install')]
@install:
    mkdir -p log
    bin/netmon --service install

# Uninstall the service
[group('install')]
@uninstall:
    bin/netmon --service uninstall

# Start the program in the background
[group('run')]
@start *args:
    mkdir -p log
    nohup bin/netmon {{ args }} >> log/netmon.log && echo "netmon started" || echo "netmon failed to start. See netmon.log" &

# Stop the program
[group('run')]
@stop:
    pkill netmon && echo "netmon stopped" || echo "netmon is not running"

# View the status of the program
[group('run')]
@status:
    ps aux | grep netmon | grep -v grep || echo "netmon is not running"

# View the last 20 lines of the program logs
[group('logs')]
@logs:
    tail -n 20 log/netmon.log 2>&1 | grep -v "No such file or directory" || echo "No logs found"

# View the last 20 lines of the connection metrics log
[group('logs')]
@metrics:
    tail -n 20 data/metrics.json 2>&1 | jq . 2>/dev/null || echo "No metrics found"
