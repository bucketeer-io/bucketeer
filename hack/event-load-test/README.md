# Event Load Testing Tool

This tool is designed to test the Bucketeer event processing system by sending large volumes of evaluation and goal events through the API gateway. It simulates real-world traffic patterns and measures system performance under load.

## Features

- Generates configurable volumes of evaluation and goal events
- Sends events through the standard Bucketeer API gateway interface
- Uses concurrent workers for high throughput testing
- Provides detailed metrics on performance and success rates
- Supports customizable event attributes (features, users, goals, tags)

## Prerequisites

- Go 1.16+
- API key for a Bucketeer environment

## Usage

```bash
go run main.go -api-key-path=/path/to/api-key [options]
```

### Required Flags

- `--api-key-path`: Path to file containing the API key

### Connection Parameters

- `--gateway-addr`: Gateway host:port (default: "localhost:8443")
- `--insecure-tls`: Skip TLS certificate verification (default: true)

### Load Test Parameters

- `--num-eval-events`: Number of evaluation events to generate (default: 10000)
- `--num-goal-events`: Number of goal events to generate (default: 5000)
- `--concurrent-workers`: Number of concurrent workers for sending events (default: 10)
- `--environment-id`: Environment ID to use for events
- `--feature-count`: Number of unique features to use (default: 10)
- `--user-count`: Number of unique users to use (default: 100)
- `--tag-count`: Number of unique tags to use (default: 5)
- `--goal-count`: Number of unique goals to use (default: 5)

### Data Generation Parameters

- `--user-data-size`: Size of user data map (default: 10)
- `--json-field-size`: Size of random string fields (default: 20)

## Examples

### Basic Usage

```bash
go run main.go --api-key-path=/path/to/api.key --gateway-addr=api.example.com:443
```

### High Volume Testing

```bash
go run main.go --api-key-path=/path/to/api.key --num-eval-events=100000 --num-goal-events=50000 --concurrent-workers=20
```

### Testing with Specific Environment

```bash
go run main.go --api-key-path=/path/to/api.key --environment-id=your-environment-id
```

## Interpreting Results

The tool provides detailed logging and performance metrics including:

- Number of events successfully sent
- Number of events that failed
- Total duration for each event type
- Events per second throughput rate

These metrics can help identify performance bottlenecks and validate system capacity. 