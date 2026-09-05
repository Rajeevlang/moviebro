BASE_URL="https://moviehub.nostack.online"
    
    echo "==========================================================="
    echo "🧪 Testing Deployment: $BASE_URL"
    echo "==========================================================="
    
    test_endpoint() {
        local name=$1
        local url=$2
        
        echo -n "Testing $name... "
        # Curl command to get just the HTTP status code
        status_code=$(curl -o /dev/null -s -w "%{http_code}\n" "$url")
        
        if [ "$status_code" -eq 200 ] || [ "$status_code" -eq 201 ] || [ "$status_code" -eq 401 ]; then
            echo "✅ (HTTP $status_code)"
        elif [ "$status_code" -eq 502 ] || [ "$status_code" -eq 504 ]; then
            echo "❌ (HTTP $status_code - Bad Gateway/Timeout. Service is down!)"
        else
            echo "⚠️ (HTTP $status_code - Unexpected response, but service is alive)"
        fi
    }

    # 1. Test Frontend React App
    test_endpoint "Frontend Homepage" "$BASE_URL/"

    # 2. Test API Gateway Healthcheck
    test_endpoint "API Gateway Healthcheck" "$BASE_URL/healthcheck"

    # 3. Test Movie Service (Routed through Gateway)
    test_endpoint "Movie Service (GET /api/v1/movies)" "$BASE_URL/api/v1/movies"

    # 4. Test User Service (Routed through Gateway)
    test_endpoint "User Service (GET /api/v1/users/healthcheck)" "$BASE_URL/api/v1/users/healthcheck"

    echo "==========================================================="
