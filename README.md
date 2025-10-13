# HTTPRunner

**HTTPRunner** is a powerful tool for load testing HTTP requests, developed in Go. It allows developers and testers to effectively assess the performance of web applications, identify bottlenecks, and ensure stability under high traffic.

## Usage

- **Load Generation**: Create and send a multitude of HTTP requests to simulate real traffic.
- **Custom Scenarios**: Ability to create user-defined testing scenarios for various types of requests and parameters.
- **Performance Reports**: Generate detailed reports on response times, including average, minimum, and maximum response times.
- **Success Rate Calculation**: Calculate the percentage of successful responses.

### Command-Line Flags

- `-method`: HTTP method to use (e.g., GET, POST). Default is GET.
- `-url`: Target URL for the requests. This flag is required.
- `-count`: Number of requests to send. Default is 1.
- `-verbose`: Enable verbose output for detailed logging.
- `-concurrency`: Number of concurrent requests to send. Default is 10.
- `-headers`: Comma-separated list of headers in the format key:value.
- `-data`: JSON string of data to send in the request body. This flag is useful for POST requests where data needs to be sent to the server.
- `-config-file`: Path to the configuration file in YAML format. If this flag is provided, other flags will be ignored.

### Examples CLI commands
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

### Examples config file

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


## Installation
You can install http-runner using the following command:

```sh
go install github.com/idesyatov/http-runner@latest
```

Also download, modify and compile it yourself as you wish