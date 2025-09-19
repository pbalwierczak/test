# Simulator Setup

The simulator has been configured to run as a standalone container that connects to the main scootin-app container via HTTP requests.

## Architecture

- **Main App**: Contains the API server and database (runs via `docker-compose.yml`)
- **Simulator**: Contains only the simulator app (runs via `docker-compose.simulator.yml`)
- **Communication**: Simulator connects to main app via HTTP API calls using `host.docker.internal`

## Usage

### 1. Start the Main Application

First, start the main application with database:

```bash
make app
```

This will start:
- PostgreSQL database
- Scootin API server

### 2. Test Connectivity (Optional)

Test that the simulator can connect to the main app:

```bash
make simulator-test
```

### 3. Start the Simulator

Once the main app is running, start the simulator:

```bash
make simulator
```

The simulator will connect to the main app via the host network.

The simulator will:
- Build the simulator Docker image
- Connect to the running scootin-app container
- Start simulating scooter and user behavior

## Configuration

The simulator connects to the main app using the `SIMULATOR_SERVER_URL` environment variable:

- Default: `http://host.docker.internal:8080`
- Can be overridden by setting `SIMULATOR_SERVER_URL` in your environment

## Docker Compose Files

- `docker-compose.yml`: Main application with database
- `docker-compose.simulator.yml`: Simulator only (connects to external app)

## Troubleshooting

If the simulator cannot connect to the main app:

1. Ensure the main app is running: `docker ps | grep scootin-app`
2. Check the API is accessible: `curl http://localhost:8080/health`
3. Run the connectivity test: `make simulator-test`
4. Check simulator logs: `docker logs scootin-simulator`

### Common Issues

**"no such host" error**: This means the simulator can't reach the main app. Make sure:
- The main app is running (`make app`)
- The main app is accessible on port 8080
- The simulator is using the correct URL (`host.docker.internal:8080`)

**Connection refused**: If you get a connection error:
- Ensure the main app is running and listening on port 8080
- Check that the main app is accessible from the host: `curl http://localhost:8080/health`

## Benefits

- **Separation of Concerns**: Simulator and main app are independent
- **Resource Efficiency**: Simulator doesn't need its own database
- **Scalability**: Can run multiple simulators against one main app
- **Development**: Easier to develop and test simulator independently
