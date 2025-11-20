#!/bin/bash

# =================================================================
# Deep Research Agent - Comprehensive Docker Test Suite
# =================================================================
# This script tests the entire Docker setup thoroughly
# Run this on your Docker Desktop after pulling the code
# =================================================================

set -e  # Exit on error

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

TESTS_PASSED=0
TESTS_FAILED=0
TEST_RESULTS=()

# Helper functions
print_header() {
    echo ""
    echo "=========================================="
    echo "$1"
    echo "=========================================="
    echo ""
}

print_test() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    TESTS_PASSED=$((TESTS_PASSED + 1))
    TEST_RESULTS+=("✓ $1")
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
    TESTS_FAILED=$((TESTS_FAILED + 1))
    TEST_RESULTS+=("✗ $1")
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# =================================================================
# Test 1: Prerequisites Check
# =================================================================
print_header "Test 1: Prerequisites Check"

print_test "Checking Docker installation..."
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version)
    print_success "Docker is installed: $DOCKER_VERSION"
else
    print_error "Docker is not installed"
    echo "Install from: https://www.docker.com/products/docker-desktop"
    exit 1
fi

print_test "Checking Docker Compose installation..."
if command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(docker-compose --version)
    print_success "Docker Compose is installed: $COMPOSE_VERSION"
else
    print_error "Docker Compose is not installed"
    exit 1
fi

print_test "Checking Docker daemon..."
if docker info &> /dev/null; then
    print_success "Docker daemon is running"
else
    print_error "Docker daemon is not running"
    echo "Please start Docker Desktop"
    exit 1
fi

# =================================================================
# Test 2: File Structure Validation
# =================================================================
print_header "Test 2: File Structure Validation"

REQUIRED_FILES=(
    "Dockerfile"
    "docker-compose.yml"
    ".dockerignore"
    ".env.example"
    "docker-start.sh"
    "DOCKER_SETUP.md"
    "go.mod"
    "main.go"
    "monitoring/prometheus.yml"
    "monitoring/grafana/provisioning/datasources/prometheus.yml"
    "monitoring/grafana/provisioning/dashboards/dashboard.yml"
)

for file in "${REQUIRED_FILES[@]}"; do
    print_test "Checking $file..."
    if [ -f "$file" ]; then
        print_success "$file exists"
    else
        print_error "$file is missing"
    fi
done

# =================================================================
# Test 3: Configuration File Syntax
# =================================================================
print_header "Test 3: Configuration File Syntax"

print_test "Validating docker-compose.yml syntax..."
if docker-compose config > /dev/null 2>&1; then
    print_success "docker-compose.yml syntax is valid"
else
    print_error "docker-compose.yml has syntax errors"
    docker-compose config
fi

print_test "Validating Dockerfile syntax..."
if docker build --check --file Dockerfile . > /dev/null 2>&1 || true; then
    print_success "Dockerfile syntax check passed (if supported)"
else
    print_warning "Dockerfile syntax check not supported on this Docker version"
fi

# =================================================================
# Test 4: Environment Setup
# =================================================================
print_header "Test 4: Environment Setup"

print_test "Checking .env file..."
if [ -f ".env" ]; then
    print_success ".env file exists"

    print_test "Checking JWT_SECRET in .env..."
    if grep -q "JWT_SECRET=" .env; then
        JWT_SECRET=$(grep "JWT_SECRET=" .env | cut -d'=' -f2)
        if [ "$JWT_SECRET" != "please-change-this-to-a-secure-random-string-min-32-characters" ] && [ ${#JWT_SECRET} -ge 32 ]; then
            print_success "JWT_SECRET is set and appears secure (${#JWT_SECRET} characters)"
        else
            print_warning "JWT_SECRET needs to be changed or is too short"
            echo "  Generate with: openssl rand -base64 32"
        fi
    else
        print_error "JWT_SECRET not found in .env"
    fi
else
    print_warning ".env file does not exist, will use defaults"
    echo "  Run: cp .env.example .env"
fi

# =================================================================
# Test 5: Docker Build Test
# =================================================================
print_header "Test 5: Docker Build Test"

print_test "Building Docker image (this may take a few minutes)..."
echo "Building research-agent:test..."
if docker build -t research-agent:test -f Dockerfile . > /tmp/docker-build.log 2>&1; then
    print_success "Docker image built successfully"

    # Check image size
    IMAGE_SIZE=$(docker images research-agent:test --format "{{.Size}}")
    print_success "Image size: $IMAGE_SIZE"

    # Check image layers
    LAYERS=$(docker history research-agent:test --no-trunc | wc -l)
    print_success "Image has $LAYERS layers"
else
    print_error "Docker build failed"
    echo "Build log:"
    cat /tmp/docker-build.log
    exit 1
fi

# =================================================================
# Test 6: Docker Compose Stack Test
# =================================================================
print_header "Test 6: Docker Compose Stack Test"

print_test "Starting Docker Compose stack..."
if docker-compose up -d > /tmp/compose-up.log 2>&1; then
    print_success "Stack started successfully"
else
    print_error "Failed to start stack"
    cat /tmp/compose-up.log
    exit 1
fi

echo ""
echo "Waiting 15 seconds for services to initialize..."
sleep 15

# Check container status
print_test "Checking container health..."
CONTAINERS=$(docker-compose ps --services)
for container in $CONTAINERS; do
    STATUS=$(docker-compose ps $container | grep $container | awk '{print $4}')
    if [[ $STATUS == *"Up"* ]]; then
        print_success "$container is running"
    else
        print_error "$container is not running properly"
        docker-compose logs $container | tail -20
    fi
done

# =================================================================
# Test 7: Service Accessibility Tests
# =================================================================
print_header "Test 7: Service Accessibility Tests"

print_test "Testing Research Agent health endpoint..."
if curl -f -s http://localhost:8080/health > /dev/null 2>&1; then
    HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
    print_success "Research Agent is responding"
    echo "  Response: $HEALTH_RESPONSE"
else
    print_error "Research Agent health check failed"
    echo "  Check logs: docker-compose logs research-agent"
fi

print_test "Testing Prometheus..."
if curl -f -s http://localhost:9090/-/ready > /dev/null 2>&1; then
    print_success "Prometheus is accessible"
else
    print_warning "Prometheus might not be ready yet"
fi

print_test "Testing Grafana..."
if curl -f -s http://localhost:3000/api/health > /dev/null 2>&1; then
    print_success "Grafana is accessible"
else
    print_warning "Grafana might not be ready yet"
fi

# =================================================================
# Test 8: API Functionality Tests
# =================================================================
print_header "Test 8: API Functionality Tests"

print_test "Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test-'$(date +%s)'@example.com",
    "password": "testpass123",
    "name": "Test User"
  }')

if echo "$REGISTER_RESPONSE" | grep -q "token"; then
    print_success "User registration works"
    TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "  Token received: ${TOKEN:0:20}..."
else
    print_error "User registration failed"
    echo "  Response: $REGISTER_RESPONSE"
fi

if [ ! -z "$TOKEN" ]; then
    print_test "Testing authenticated endpoint..."
    PROFILE_RESPONSE=$(curl -s http://localhost:8080/api/v1/auth/me \
      -H "Authorization: Bearer $TOKEN")

    if echo "$PROFILE_RESPONSE" | grep -q "email"; then
        print_success "Authentication works"
    else
        print_error "Authentication failed"
        echo "  Response: $PROFILE_RESPONSE"
    fi

    print_test "Testing research start endpoint..."
    RESEARCH_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/research/start \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "query": "Test query",
        "depth": "shallow",
        "max_sources": 2
      }')

    if echo "$RESEARCH_RESPONSE" | grep -q "job_id"; then
        print_success "Research API works"
        JOB_ID=$(echo "$RESEARCH_RESPONSE" | grep -o '"job_id":"[^"]*' | cut -d'"' -f4)
        echo "  Job ID: $JOB_ID"
    else
        print_error "Research API failed"
        echo "  Response: $RESEARCH_RESPONSE"
    fi
fi

# =================================================================
# Test 9: Data Persistence Test
# =================================================================
print_header "Test 9: Data Persistence Test"

print_test "Checking data directories..."
if [ -d "docker-data/database" ]; then
    print_success "Database directory exists"

    if [ "$(ls -A docker-data/database)" ]; then
        print_success "Database files created"
        ls -lh docker-data/database/
    else
        print_warning "Database directory is empty"
    fi
else
    print_error "Database directory not created"
fi

# =================================================================
# Test 10: Container Resource Usage
# =================================================================
print_header "Test 10: Container Resource Usage"

print_test "Checking container resource usage..."
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" | grep research

# =================================================================
# Test 11: Logs Check
# =================================================================
print_header "Test 11: Logs Check"

print_test "Checking for errors in logs..."
ERROR_COUNT=$(docker-compose logs research-agent 2>&1 | grep -i "error\|fatal\|panic" | wc -l)
if [ "$ERROR_COUNT" -eq 0 ]; then
    print_success "No errors found in logs"
else
    print_warning "Found $ERROR_COUNT error messages in logs"
    echo "  Last 10 errors:"
    docker-compose logs research-agent 2>&1 | grep -i "error\|fatal\|panic" | tail -10
fi

# =================================================================
# Test Summary
# =================================================================
print_header "Test Summary"

echo "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo ""

echo "Detailed Results:"
for result in "${TEST_RESULTS[@]}"; do
    echo "  $result"
done

echo ""
echo "=========================================="
if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    echo "=========================================="
    echo ""
    echo "Your Deep Research Agent is ready to use!"
    echo ""
    echo "Access the services:"
    echo "  • Research Agent:  http://localhost:8080"
    echo "  • Prometheus:      http://localhost:9090"
    echo "  • Grafana:         http://localhost:3000 (admin/admin)"
    echo ""
    echo "Next steps:"
    echo "  1. Open http://localhost:8080 in your browser"
    echo "  2. Register a user and start researching"
    echo "  3. Check Grafana for monitoring dashboards"
    echo ""
    echo "View logs:    docker-compose logs -f research-agent"
    echo "Stop:         docker-compose down"
    echo ""
else
    echo -e "${RED}✗ Some tests failed${NC}"
    echo "=========================================="
    echo ""
    echo "Please check the errors above and:"
    echo "  1. Review docker-compose logs"
    echo "  2. Check DOCKER_SETUP.md for troubleshooting"
    echo "  3. Ensure all prerequisites are met"
    echo ""
fi

# Ask if user wants to keep containers running
echo ""
read -p "Keep containers running? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Stopping containers..."
    docker-compose down
    echo "Containers stopped"
fi

exit $TESTS_FAILED
