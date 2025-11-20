#!/bin/bash

API_URL="http://localhost:8080"

echo "🔍 Testing Deep Research Agent API"
echo "===================================="
echo ""

# Test 1: Health Check
echo "1️⃣ Testing health endpoint..."
curl -s $API_URL/health | jq '.'
echo ""

# Test 2: Register User
echo "2️⃣ Registering test user..."
REGISTER_RESPONSE=$(curl -s -X POST $API_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123",
    "name": "Test User"
  }')

TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token')
echo "Token: ${TOKEN:0:20}..."
echo ""

# Test 3: Start Research
echo "3️⃣ Starting research..."
RESEARCH_RESPONSE=$(curl -s -X POST $API_URL/api/v1/research/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "query": "Go programming language",
    "depth": "shallow",
    "max_sources": 3
  }')

JOB_ID=$(echo $RESEARCH_RESPONSE | jq -r '.job_id')
echo "Job ID: $JOB_ID"
echo ""

# Test 4: Check Status
echo "4️⃣ Checking research status..."
sleep 2
curl -s $API_URL/api/v1/research/$JOB_ID/status \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

echo "✅ API is working! You can connect any UI to it."
