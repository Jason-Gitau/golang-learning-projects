#!/bin/bash

# Agent Orchestrator API Test Script
# This script tests all the tools and endpoints

API_URL="http://localhost:8080"

echo "=================================="
echo "Agent Orchestrator API Test Suite"
echo "=================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to test endpoint
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4

    echo -e "${YELLOW}Testing: $name${NC}"

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$API_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ Success${NC}"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo -e "${RED}✗ Failed (HTTP $http_code)${NC}"
        echo "$body"
    fi
    echo ""
}

# Check if server is running
echo "Checking if server is running..."
if ! curl -s "$API_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}Error: Server is not running on $API_URL${NC}"
    echo "Please start the server first: go run main.go"
    exit 1
fi
echo -e "${GREEN}✓ Server is running${NC}"
echo ""

# Test health endpoint
test_endpoint "Health Check" "GET" "/health"

# Test root endpoint
test_endpoint "Root Endpoint" "GET" "/"

# Test tools list
test_endpoint "List Tools" "GET" "/api/v1/tools"

# Test agents list
test_endpoint "List Agents" "GET" "/api/v1/agents"

# Test statistics
test_endpoint "System Statistics" "GET" "/api/v1/stats"

echo "=================================="
echo "Testing Calculator Tool"
echo "=================================="
echo ""

test_endpoint "Calculator - Add" "POST" "/api/v1/request" \
'{
  "tool_name": "calculator",
  "params": {
    "operation": "add",
    "a": 10,
    "b": 5
  }
}'

test_endpoint "Calculator - Multiply" "POST" "/api/v1/request" \
'{
  "tool_name": "calculator",
  "params": {
    "operation": "multiply",
    "a": 7,
    "b": 6
  }
}'

test_endpoint "Calculator - Divide" "POST" "/api/v1/request" \
'{
  "tool_name": "calculator",
  "params": {
    "operation": "divide",
    "a": 100,
    "b": 5
  }
}'

echo "=================================="
echo "Testing Time Tool"
echo "=================================="
echo ""

test_endpoint "Time - Current" "POST" "/api/v1/request" \
'{
  "tool_name": "time",
  "params": {
    "action": "current"
  }
}'

test_endpoint "Time - Timezone" "POST" "/api/v1/request" \
'{
  "tool_name": "time",
  "params": {
    "action": "timezone",
    "timezone": "America/New_York"
  }
}'

test_endpoint "Time - Add Duration" "POST" "/api/v1/request" \
'{
  "tool_name": "time",
  "params": {
    "action": "add_duration",
    "duration": "2h30m"
  }
}'

echo "=================================="
echo "Testing Random Tool"
echo "=================================="
echo ""

test_endpoint "Random - Integer" "POST" "/api/v1/request" \
'{
  "tool_name": "random",
  "params": {
    "type": "int",
    "min": 1,
    "max": 100
  }
}'

test_endpoint "Random - UUID" "POST" "/api/v1/request" \
'{
  "tool_name": "random",
  "params": {
    "type": "uuid"
  }
}'

test_endpoint "Random - Dice" "POST" "/api/v1/request" \
'{
  "tool_name": "random",
  "params": {
    "type": "dice",
    "sides": 6,
    "count": 3
  }
}'

test_endpoint "Random - Choice" "POST" "/api/v1/request" \
'{
  "tool_name": "random",
  "params": {
    "type": "choice",
    "choices": ["red", "green", "blue", "yellow"]
  }
}'

echo "=================================="
echo "Testing Text Tool"
echo "=================================="
echo ""

test_endpoint "Text - Uppercase" "POST" "/api/v1/request" \
'{
  "tool_name": "text",
  "params": {
    "operation": "uppercase",
    "text": "hello world"
  }
}'

test_endpoint "Text - Reverse" "POST" "/api/v1/request" \
'{
  "tool_name": "text",
  "params": {
    "operation": "reverse",
    "text": "golang"
  }
}'

test_endpoint "Text - Word Count" "POST" "/api/v1/request" \
'{
  "tool_name": "text",
  "params": {
    "operation": "word_count",
    "text": "The quick brown fox jumps over the lazy dog"
  }
}'

echo "=================================="
echo "Testing Weather Tool"
echo "=================================="
echo ""

test_endpoint "Weather - London" "POST" "/api/v1/request" \
'{
  "tool_name": "weather",
  "params": {
    "location": "London",
    "units": "metric"
  }
}'

echo "=================================="
echo "Concurrency Test"
echo "=================================="
echo ""

echo "Sending 10 concurrent requests..."
for i in {1..10}; do
    curl -s -X POST "$API_URL/api/v1/request" \
        -H "Content-Type: application/json" \
        -d "{\"tool_name\":\"calculator\",\"params\":{\"operation\":\"add\",\"a\":$i,\"b\":$i}}" > /dev/null &
done
wait
echo -e "${GREEN}✓ All concurrent requests completed${NC}"
echo ""

# Check final statistics
test_endpoint "Final Statistics" "GET" "/api/v1/stats"

echo "=================================="
echo "Test Suite Completed"
echo "=================================="
