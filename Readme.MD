# LetServerCook

LetServerCook is a high-performance HTTP benchmarking tool written in Go, designed to test REST APIs under load with configurable concurrency and request counts. It allows developers to measure server performance by tracking request success rates, average and maximum latency, and requests per second (RPS) in real time.
The tool leverages Go’s goroutines and channels to efficiently manage worker pools, ensuring accurate simulation of concurrent requests. It includes robust timeout handling, thread-safe result aggregation, and flexible configuration of HTTP methods and request parameters.
Unit tests using testify validate benchmark logic and timeout handling, ensuring reliability across different scenarios. LetServerCook is a practical project for anyone looking to showcase skills in Go concurrency, HTTP programming, performance testing, and clean, testable code.


---

## Features

* **Concurrent Workers:** Spawn multiple workers to send HTTP requests concurrently.
* **Timeout Support:** Set request timeout (TTL) to simulate slow responses.
* **Detailed Metrics:** Captures success/failure counts, average latency, maximum latency, and RPS.
* **Cross-Platform:** Buildable for Windows, Linux, and MacOS.
* **Minimal Dependencies:** Pure Go, no heavy libraries.

---

## Installation

Clone the repository and build:

```bash
git clone https://github.com/Rishi-Mishra0704/LetServerCook.git
cd LetServerCook
make build-windows   # For Windows
make build           # for rest of os
```

Binaries will be available in the `./bin/` folder.

---

## Usage

```bash
./bin/lsc -u <URL> -n <TOTAL_REQUESTS> -w <WORKERS> -t <TIMEOUT> -m <METHOD>
```

### Example

```bash
./bin/lsc -u https://jsonplaceholder.typicode.com/todos -n 50 -w 5 -t 1.5s -m GET
```

**Output:**

```
Starting LetServerCook benchmark...
URL: https://jsonplaceholder.typicode.com/todos | Method: GET | Requests: 50 | Workers: 5 | Timeout: 1.5s
Spawning workers...
Enqueuing tasks...
Running benchmark...
Worker 0 processed 10 requests...
Worker 1 processed 10 requests...
Worker 2 processed 10 requests...
Worker 3 processed 10 requests...
Worker 4 processed 10 requests...
Calculating results...
Benchmark completed!
Server Cooked
🔥 LetServerCook Results 🔥
Requests:     50
Concurrency:  5
Success:      50 | Failures: 0
Success Rate: 100 percent
Avg Latency:  141ms
Max Latency:  412ms
RPS:          35.32
```

---

## Flags

| Flag | Description                        |
| ---- | ---------------------------------- |
| `-h` | Help                               |
| `-u` | Target URL to benchmark            |
| `-n` | Total number of requests           |
| `-w` | Number of concurrent workers       |
| `-t` | Timeout per request (e.g., `1.5s`) |
| `-m` | HTTP method (GET, POST, etc.)      |

---

## Testing

Unit tests are provided in the `cmd` package. Run:

```bash
go test ./cmd -v
```

Tests cover:

* Successful requests
* Timeout handling

---

## Contribution

Feel free to submit issues or pull requests. Follow these guidelines:

* Keep the benchmark code lean and idiomatic Go.
* Add tests for new features.
* Avoid breaking the current request/worker pattern.

---

## License

[MIT LICENCE]("./LICENSE")
