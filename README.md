# HTTPRunner

[![CI](https://github.com/idesyatov/http-runner/actions/workflows/ci.yml/badge.svg)](https://github.com/idesyatov/http-runner/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/idesyatov/http-runner)](https://github.com/idesyatov/http-runner/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/idesyatov/http-runner)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/idesyatov/http-runner)](https://goreportcard.com/report/github.com/idesyatov/http-runner)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENCE)

**HTTPRunner** is a powerful tool for load testing HTTP requests, developed in Go. It allows developers and testers to effectively assess the performance of web applications, identify bottlenecks, and ensure stability under high traffic.

<p align="center">
  <img src="assets/demo.gif" alt="HTTPRunner demo" width="700">
</p>

## Features

- **Load Generation** — create and send a multitude of HTTP requests to simulate real traffic.
- **Custom Scenarios** — define testing scenarios for various types of requests and parameters via YAML.
- **Performance Reports** — detailed reports on response times (average, minimum, maximum).
- **Success Rate Calculation** — the percentage of successful responses, per status code.

## Quick Start

```bash
# Install (Linux / macOS)
curl -fsSL https://raw.githubusercontent.com/idesyatov/http-runner/master/install.sh | sh

# Single GET request
http-runner -url "https://example.com"

# 100 requests with 50 concurrent workers
http-runner -url "https://example.com" -count 100 -concurrency 50

# POST with a JSON body
http-runner -method POST -url "https://example.com/api" -count 10 -data '{"key":"value"}'
```

> Run `http-runner -help` for the full list of flags.

<details>
<summary><strong>Installation</strong> (script · go install · prebuilt binaries)</summary>

### Install script (Linux / macOS)

Download the latest release binary and install it to `/usr/local/bin`:

```sh
curl -fsSL https://raw.githubusercontent.com/idesyatov/http-runner/master/install.sh | sh
```

The script detects your OS and architecture, fetches the latest release, verifies
the SHA-256 checksum, and installs the binary. You can override the defaults:

```sh
# Pin a specific version and/or change the install directory:
curl -fsSL https://raw.githubusercontent.com/idesyatov/http-runner/master/install.sh | VERSION=v1.4.0 BINDIR=$HOME/.local/bin sh
```

Prefer to review before running? Download `install.sh`, read it, then run `sh install.sh`.

### Go install

```sh
go install github.com/idesyatov/http-runner@latest
```

### Prebuilt binaries

Download an archive for your platform from the
[Releases](https://github.com/idesyatov/http-runner/releases) page, or build it
yourself with `make build`.

</details>

<details>
<summary><strong>Command-line flags</strong></summary>

- `-method`: HTTP method to use (e.g., GET, POST). Default is GET.
- `-url`: Target URL for the requests. This flag is required.
- `-count`: Number of requests to send. Default is 1.
- `-verbose`: Enable verbose output for detailed logging.
- `-concurrency`: Number of concurrent requests to send. Default is 10.
- `-headers`: Comma-separated list of headers in the format key:value.
- `-data`: JSON string of data to send in the request body. This flag is useful for POST requests where data needs to be sent to the server.
- `-config-file`: Path to the configuration file in YAML format. If this flag is provided, other flags will be ignored.
- `-version`: Show the application version and exit.

</details>

<details>
<summary><strong>Examples</strong> (CLI)</summary>

```bash
# To get information on all flags:
http-runner -help

# To send a single GET request to a specified URL:
http-runner -url "https://example.com"

# To send 100 GET requests to a specified URL:
http-runner -url "https://example.com" -count 100

# To send 10 POST requests to a specified URL with JSON data:
http-runner -method "POST" \
    -url "https://example.com/api" \
    -count 10 \
    -data '{"key1":"value1", "key2":"value2"}'

# To enable verbose output while sending requests:
http-runner -url "https://example.com" \
    -count 10 \
    -verbose

# To send 50 concurrent requests to a specified URL:
http-runner -url "https://example.com" \
    -count 50 \
    -concurrency 50

# To send a single GET request to a specified URL with custom headers:
http-runner -url "https://example.com/api" \
    -headers "Authorization: Bearer your_token, Content-Type: application/json"

# To load configuration from a YAML file:
http-runner -config-file "config.yaml"
```

</details>

<details>
<summary><strong>Configuration file</strong> (YAML, all parameters)</summary>

Pass `-config-file path.yml`; when it is set, all other flags are ignored.

```yml
# Configuration file for http-runner, demonstrating all possible parameters

endpoints:
  - url: "https://example.com/api"      # (Required) Target URL for requests.
    method: "POST"                      # (Optional, default: GET) HTTP method for the request.
    headers:                            # (Optional) Headers for the request in key:value format.
      Authorization: "Bearer your_token"
      Content-Type: "application/json"
    data:                               # (Optional) JSON string of data to send in the request body.
      key1: "value1"
      key2: "value2"
    count: 5                            # (Optional, default: 1) Number of requests to send.
    concurrency: 3                      # (Optional, default: 10) Number of concurrent requests.
    verbose: true                       # (Optional) Enables detailed output for logging.

  - url: "https://example.org"          # (Optional) Second example with a different URL.
    method: "GET"                       # (Optional) Default GET method.
    headers:                            # (Optional) Headers for the request.
      Accept: "application/json"
    count: 10                           # (Optional, default: 1) Number of requests.
    concurrency: 5                      # (Optional, default: 10) Number of concurrent requests.
    verbose: false                      # (Optional) Disables detailed output.

  - url: "https://api.example.com/data" # (Optional) Third example URL.
    method: "PUT"                       # (Optional) PUT method.
    headers:                            # (Optional) Headers for the request.
      Content-Type: "application/x-www-form-urlencoded"
    data:                               # (Optional) Data for updating.
      name: "Item"
      value: "UpdatedValue"
    count: 1                            # (Optional, default: 1) Only one request.
    concurrency: 1                      # (Optional, default: 10) One request at a time.
    verbose: true                       # (Optional) Enables detailed output.
```

> Note: `data` currently supports flat `key: value` pairs only (no nested JSON objects or arrays).

</details>

## License

[MIT](LICENCE)
