# display this help message
[group('info')]
@help:
    echo "Usage: just <recipe-name>"
    echo ""
    just --list --unsorted

# build the binary
[group('build')]
build:
    go mod tidy
    go mod vendor
    # trim debug information from the binary
    go build -ldflags="-s -w" -o bin/netmon

# clean the binary
[group('build')]
clean:
    rm -f bin/netmon

# install and start the service
[group('service')]
@install:
    pgrep -f bin/netmon > /dev/null && { echo "netmon is already running. try stop or uninstall"; exit 1; } || true
    sudo bin/netmon --service install
    sudo bin/netmon --service start

# stop and uninstall the service
[group('service')]
@uninstall:
    sudo bin/netmon --service stop
    sudo bin/netmon --service uninstall

# start the program in the foreground
[group('standalone')]
@run *args:
    pgrep -f bin/netmon > /dev/null && { echo "netmon is already running. try stop or uninstall"; exit 1; } || true
    bin/netmon {{ args }} -interval 30s

# start the program in the background
[group('standalone')]
@start *args:
    pgrep -f bin/netmon > /dev/null && { echo "netmon is already running. try stop or uninstall"; exit 1; } || true
    mkdir -p log
    nohup bin/netmon {{ args }} >> log/netmon.log 2>&1 & disown
    sleep 1
    pgrep -f bin/netmon > /dev/null && echo "netmon started" || echo "netmon failed to start. See netmon.log"

# stop the program
[group('standalone')]
@stop:
    pgrep -f bin/netmon > /dev/null && { pkill -f bin/netmon && echo "netmon stopped" || echo "netmon failed to stop. try uninstall"; } || echo "netmon is not running"

# check the running status of the program (both service and standalone)
[group('info')]
@status:
    pgrep -f bin/netmon > /dev/null && echo "netmon is running" || echo "netmon is not running"
