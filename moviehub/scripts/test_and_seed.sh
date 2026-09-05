#!/usr/bin/env bash
# =============================================================================
# MovieHub Microservices Smoke Test & Initial Seeder Wrapper
# =============================================================================
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_URL="${1:-http://localhost:8080}"

echo "=================================================="
echo " Running MovieHub Test & Seeder against: $TARGET_URL"
echo "=================================================="

if command -v python3 &>/dev/null; then
    python3 "$SCRIPT_DIR/test_and_seed.py" --base-url "$TARGET_URL"
else
    echo "python3 is not installed. Please install python3 to run the test and seed script."
    exit 1
fi
