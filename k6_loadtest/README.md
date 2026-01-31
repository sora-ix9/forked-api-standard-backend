# k6 Load Testing

This directory contains load tests for the API Standard Backend using [k6](https://k6.io/).

## Prerequisites

You need to have k6 installed on your machine.

### MacOS
```bash
brew install k6
```

### Windows
```bash
winget install k6
```

### Linux (Debian/Ubuntu)
```bash
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491B6B8D6D9
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

## Running the Load Test

1. Ensure your backend server is running locally on port 1323.
   ```bash
   # From the project root
   go run cmd/main.go
   ```

2. Run the k6 script:
   ```bash
   k6 run k6_loadtest/k6_loadtest.js
   ```

## Test Configuration

The current test configuration (`k6_loadtest.js`) is set up to:
- Ramp up to **20 virtual users** over 30 seconds.
- Sustain 20 users for **1 minute**.
- Ramp down to 0 users over 10 seconds.
- Verify that the health check endpoint (`GET /`) returns a 200 status code and contains "Online".
- Enforce a threshold where 95% of request durations must be under **500ms**.

## Reports

After running the test, an HTML report will be generated at `k6_loadtest/summary.html`. Open this file in your browser to view detailed metrics and graphs.

## Customization

You can modify `k6_loadtest.js` to add more complex scenarios, such as creating users or logging in. Example code for login is included in the script as comments.
