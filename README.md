# netmon

![Screenshot](img/ss.png)

**netmon** is a simple network connection latency & speed monitor.

- **Single binary**, no dependencies
- **Cross-platform** (Linux, macOS, Windows)
- **Minimalist**, no configuration, no setup, no installation
- **Periodically** checks network connection latency & speed
- Stores results in an **SQLite** database
- Provides a simple **web interface** to view results as visual charts
- Supports both running in **standalone** mode and installation as a **system service**

## Installation

- **Option 1:** Download and extract the latest
  release: [Releases](https://github.com/mmdemirbas/netmon/releases)

- **Option 2:** Build from source:

   ```bash
   git clone https://github.com/mmdemirbas/netmon.git
   cd netmon
   just build # OR: go build, if you like
   ```
  To see all available commands, run `just`.

## Usage — Standalone

1. Run the binary.
   ```bash
   ./netmon
   ```

2. Open the web interface at [http://localhost:9898](http://localhost:9898).
3. To stop, press `Ctrl+C`.

- If you cloned the source code, you can also use the `just` commands:
    ```bash
    just run   # to run in foreground
    just start # to run in background
    just stop  # to stop
    ```

## Usage — System Service

1. Install and start the service.
   ```bash
   just install
   ```
   or
   ```bash
   sudo ./netmon -service install
   sudo ./netmon -service start
   ```
2. Open the web interface at [http://localhost:9898](http://localhost:9898).
3. Stop and uninstall the service.
   ```bash
    just uninstall
    ```
   or
    ```bash
   sudo ./netmon -service stop
   sudo ./netmon -service uninstall
   ```

## Configuration

| Flag              | Description                                                 | Default          |
|-------------------|-------------------------------------------------------------|------------------|
| `-db-file <path>` | Path to the SQLite database file.                           | `data/netmon.db` |
| `-interval <s>`   | Interval between checks.                                    | `5m`             |
| `-port <n>`       | Port for the web interface.                                 | `9898`           |
| `-service <s>`    | Control servicel (`install`, `start`, `stop`, `uninstall`). |                  |

## Contributing

Contributions are welcome! For major changes, please open an issue first to discuss what you would
like to change.

### Pre-requisites

- [Go](https://go.dev/)
- [Just](https://just.systems)

### Development

Clone the repository.

```bash
git clone https://github.com/mmdemirbas/netmon.git
cd netmon
```

Run `just` to see available commands.

```bash
Usage: just <recipe-name>

Available recipes:
    [info]
    help        # display this help message
    status      # check the running status of the program (both service and standalone)

    [build]
    build       # build the binary
    clean       # clean the binary

    [service]
    install     # install and start the service
    uninstall   # stop and uninstall the service

    [standalone]
    run *args   # start the program in the foreground
    start *args # start the program in the background
    stop        # stop the program
```

You can combine multiple commands to streamline the development cycle like this:

```bash
just stop clean build run
```

## Author

Muhammed Demirbaş - [GitHub](https://github.com/mmdemirbas)

Thanks to the following open-source projects:

- [Go](https://go.dev/)
- [speedtest-go](https://showwin/speedtest-go)
- [Chart.js](https://www.chartjs.org/)
- [go-sqlite3](https://github.com/mattn/go-sqlite3)
- [service](https://github.com/kardianos/service)
- [justfile](https://just.systems)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
