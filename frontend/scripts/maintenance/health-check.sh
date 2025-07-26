#!/bin/bash

# Health Check Script for POS QR System Frontend
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVICE_URL=${SERVICE_URL:-"http://localhost:3000"}
TIMEOUT=${TIMEOUT:-10}

echo -e "${GREEN}🏥 Starting Health Check for POS QR Frontend${NC}"

# Function to check HTTP endpoint
check_endpoint() {
    local url=$1
    local expected_status=${2:-200}
    local description=$3
    
    echo -e "${YELLOW}🔍 Checking $description: $url${NC}"
    
    local status_code=$(curl -s -o /dev/null -w "%{http_code}" --max-time $TIMEOUT "$url" || echo "000")
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ $description: OK (Status: $status_code)${NC}"
        return 0
    else
        echo -e "${RED}❌ $description: FAILED (Status: $status_code)${NC}"
        return 1
    fi
}

# Function to check response time
check_response_time() {
    local url=$1
    local max_time=${2:-2000}
    local description=$3
    
    echo -e "${YELLOW}⏱️  Checking response time for $description${NC}"
    
    local response_time=$(curl -s -o /dev/null -w "%{time_total}" --max-time $TIMEOUT "$url" 2>/dev/null || echo "999")
    local response_time_ms=$(echo "$response_time * 1000" | bc -l | cut -d. -f1)
    
    if [ "$response_time_ms" -lt "$max_time" ]; then
        echo -e "${GREEN}✅ $description response time: ${response_time_ms}ms (< ${max_time}ms)${NC}"
        return 0
    else
        echo -e "${RED}❌ $description response time: ${response_time_ms}ms (>= ${max_time}ms)${NC}"
        return 1
    fi
}

# Health check results
FAILED_CHECKS=0

# Check main application
if ! check_endpoint "$SERVICE_URL" 200 "Main Application"; then
    ((FAILED_CHECKS++))
fi

# Check login pages
if ! check_endpoint "$SERVICE_URL/auth/admin-login" 200 "Admin Login Page"; then
    ((FAILED_CHECKS++))
fi

if ! check_endpoint "$SERVICE_URL/auth/store-login" 200 "Store Login Page"; then
    ((FAILED_CHECKS++))
fi

# Check API health endpoint (if available)
if ! check_endpoint "$SERVICE_URL/api/health" 200 "API Health Endpoint"; then
    echo -e "${YELLOW}⚠️  API health endpoint not available (this may be expected)${NC}"
fi

# Check response times
if ! check_response_time "$SERVICE_URL" 2000 "Main Application"; then
    ((FAILED_CHECKS++))
fi

# Check static assets
if ! check_endpoint "$SERVICE_URL/_next/static/css" 404 "Static Assets Directory" || 
   ! check_endpoint "$SERVICE_URL/favicon.ico" 200 "Favicon"; then
    echo -e "${YELLOW}⚠️  Some static assets may not be accessible${NC}"
fi

# Memory and CPU check (if running locally)
if command -v ps >/dev/null 2>&1; then
    echo -e "${YELLOW}💻 Checking system resources${NC}"
    
    # Check if Node.js process is running
    if pgrep -f "node.*next" > /dev/null; then
        local node_pid=$(pgrep -f "node.*next" | head -1)
        local memory_usage=$(ps -p $node_pid -o %mem --no-headers 2>/dev/null | tr -d ' ' || echo "N/A")
        local cpu_usage=$(ps -p $node_pid -o %cpu --no-headers 2>/dev/null | tr -d ' ' || echo "N/A")
        
        echo -e "${GREEN}📊 Node.js Process (PID: $node_pid)${NC}"
        echo -e "   Memory Usage: ${memory_usage}%"
        echo -e "   CPU Usage: ${cpu_usage}%"
    else
        echo -e "${YELLOW}⚠️  Node.js process not found (may be running in container)${NC}"
    fi
fi

# Summary
echo -e "\n${GREEN}📋 Health Check Summary${NC}"
echo -e "Service URL: $SERVICE_URL"
echo -e "Timestamp: $(date)"

if [ $FAILED_CHECKS -eq 0 ]; then
    echo -e "${GREEN}✅ All health checks passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ $FAILED_CHECKS health check(s) failed${NC}"
    exit 1
fi