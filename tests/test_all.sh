#!/bin/bash

echo "🚀 Starting comprehensive CLI testing..."

# Create test results directory
mkdir -p test-results
cd test-results

# Test counter
test_count=0
success_count=0
fail_count=0

# Function to run test and count results
run_test() {
    local cmd="$1"
    local description="$2"
    
    echo "🧪 Testing: $description"
    echo "   Command: $cmd"
    
    if eval "$cmd" 2>/dev/null; then
        echo "   ✅ SUCCESS"
        ((success_count++))
    else
        echo "   ❌ FAILED"
        ((fail_count++))
    fi
    
    ((test_count++))
    echo ""
}

# Run all test categories
echo "📝 Category 1: Basic Functionality"
run_test '../tti "Hello World"' "Basic text rendering"
run_test '../tti "A"' "Single character"
run_test '../tti "This is a very long text that might overflow"' "Long text"

echo "📝 Category 2: Font Styles"
for font in roboto_bold roboto_regular roboto_medium roboto_black; do
    run_test "../tti -font=$font \"Font Test: $font\"" "Font: $font"
done

echo "📝 Category 3: All Backgrounds"
for bg in default perlin perlin-s radial diagonal; do
    run_test "../tti -bg=$bg \"Background: $bg\"" "Background: $bg"
done

echo "📝 Category 4: Dimension Tests"
run_test '../tti -width=800 -height=600 "Custom Size"' "Custom dimensions"
run_test '../tti -width=1920 -height=1080 "HD Size"' "HD dimensions"

echo "📝 Category 5: Font Size Tests"
run_test '../tti -font-size=12 "Small Font"' "Small font size"
run_test '../tti -font-size=72 "Large Font"' "Large font size"

echo "📝 Category 6: Special Effects"
run_test '../tti -reveal-bg "Reveal Test"' "Reveal background"
run_test '../tti -animate "Animation Test"' "Animation"

echo "📝 Category 7: Complex Combinations"
run_test '../tti -width=1024 -height=768 -font-size=48 -font=roboto_italic -bg=perlin -animate "Full Combo"' "All parameters"

# Final results
echo "🏁 Test Summary:"
echo "   Total tests: $test_count"
echo "   Successful: $success_count"
echo "   Failed: $fail_count"
echo "   Success rate: $(( success_count * 100 / test_count ))%"

cd ..