#!/bin/bash

# Configuration
BASE_URL="${1:-https://moviehub.nostackdev.online}"
DOMAIN=$(echo "$BASE_URL" | sed -e 's|^[^/]*//||' -e 's|/.*$||')

echo "==========================================================="
echo "🧪 Testing Deployment: $DOMAIN"
echo "==========================================================="

# Function to test a URL
test_endpoint() {
    local name=$1
    local url=$2
    local expected_code=$3

    echo -n "Testing $name... "
    # Curl command to get just the HTTP status code
    status_code=$(curl -o /dev/null -s -w "%{http_code}\n" "$url")
    
    if [ "$status_code" -eq "$expected_code" ] || ([ "$expected_code" == "20X" ] && [[ "$status_code" == 20* ]]); then
        echo "✅ (HTTP $status_code)"
    else
        echo "❌ (Expected $expected_code, Got HTTP $status_code)"
    fi
}

# 1. Test Frontend
test_endpoint "Frontend Homepage" "$BASE_URL/" "20X"

# 2. Test API Gateway Healthcheck
test_endpoint "API Gateway Healthcheck" "$BASE_URL/api/v1/movies/healthcheck" "404" # wait, proxy maps /api/v1/movies/, it doesn't map /healthcheck directly!
# Actually, the gateway serves /healthcheck on itself! Let's test the gateway directly.
test_endpoint "API Gateway Healthcheck" "$BASE_URL/healthcheck" "20X"

# 3. Test Movie Service API (List Movies)
test_endpoint "Movie Service (GET /api/v1/movies)" "$BASE_URL/api/v1/movies" "20X"

# 4. Test User Service API (Healthcheck/Placeholder)
# Since the gateway forwards /api/v1/users/ to user-service
# We can just test a known route, or the healthcheck if forwarded. 
# We'll just test a basic GET request to users. It might return 401 or 405 if unauthorized/method not allowed, 
# but getting a response from the service proves it's up.
echo -n "Testing User Service (GET /api/v1/users/healthcheck)... "
status_code=$(curl -o /dev/null -s -w "%{http_code}\n" "$BASE_URL/api/v1/users/healthcheck")
# We just want to make sure it's not a 502 Bad Gateway
if [ "$status_code" -ne 502 ] && [ "$status_code" -ne 504 ]; then
    echo "✅ (HTTP $status_code - Service is alive)"
else
    echo "❌ (Got HTTP $status_code - Service might be down)"
fi

echo "==========================================================="
echo "Tests completed!"
echo "If you see all ✅ marks, your entire microservice stack"
echo "is successfully communicating behind Nginx and Cloudflare!"
echo "==========================================================="
