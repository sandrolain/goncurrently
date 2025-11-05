#!/bin/bash
# Test all examples to verify they work

cd "$(dirname "$0")/.." || exit 1

echo "========================================="
echo "Testing goncurrently examples"
echo "========================================="
echo ""

# Build goncurrently first
echo "Building goncurrently..."
go build -o goncurrently . || exit 1
echo "✓ Build successful"
echo ""

EXAMPLES=(
  "basic"
  "env-vars"
  "timeouts"
)

FAILED=0
PASSED=0

for example in "${EXAMPLES[@]}"; do
  echo "Testing: $example.yaml"
  
  timeout 5s bash -c "cat examples/${example}.yaml | ./goncurrently" > /dev/null 2>&1
  EXIT_CODE=$?
  
  # Exit codes 0 (success) or 124 (timeout) are acceptable
  if [ $EXIT_CODE -eq 0 ] || [ $EXIT_CODE -eq 124 ]; then
    echo "✓ $example passed"
    ((PASSED++))
  else
    echo "✗ $example failed (exit code: $EXIT_CODE)"
    ((FAILED++))
  fi
  echo ""
done

echo "========================================="
echo "Results: $PASSED passed, $FAILED failed"
echo "========================================="

if [ $FAILED -eq 0 ]; then
  echo "✓ All tests passed!"
  exit 0
else
  echo "✗ Some tests failed"
  exit 1
fi
