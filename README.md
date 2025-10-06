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
- `-config-file`: Path to the configuration file in YAML format. If this flag is provided, other flags will be ignored.

### Examples CLI commands
```bash
# To get information on all flags:
http-runner -help

# To send a single GET request to a specified URL:
http-runner -url "https://example.com"

# To send 100 GET requests to a specified URL:
http-runner -url "https://example.com" -count 100

# To send 10 POST requests to a specified URL:
http-runner -method "POST" -url "https://example.com/api" -count 10

# To enable verbose output while sending requests:
http-runner -url "https://example.com" -count 10 -verbose

# To send 50 concurrent requests to a specified URL:
http-runner -url "https://example.com" -count 50 -concurrency 50

# To send a single GET request to a specified URL with custom headers:
http-runner -url "https://example.com/api" \
    -headers "Authorization: Bearer your_token, Content-Type: application/json"

# To load configuration from a YAML file:
http-runner -config-file "config.yaml"
```

### Examples config file

```yml
endpoints:
  - url: 'https://example.com/api'
    verbose: false
    method: 'GET'
    headers: 
      Authorization: 'Bearer your_token'
      Content-Type: 'application/json'
    count: 10
    concurrency: 20
  - url: 'https://example.com/v2/api'
    verbose: false
    method: 'POST'
    headers: 
      Authorization: 'Bearer your_token'
      Content-Type: 'application/json'
    data:
      key: 'value'
    count: 10
    concurrency: 20
```


## Installation
You can install http-runner using the following command:

```sh
go install github.com/idesyatov/http-runner@latest
```

Also download, modify and compile it yourself as you wish