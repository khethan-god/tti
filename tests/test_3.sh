#!/bin/bash

# performance and stress testing

# Multiple rapid executions concurrently
for i in {1..10}; do
  go run . "Test $i" &
done
wait

# Large batch with different parameters
for i in {1..50}; do
  width=$((400 + i * 10))
  height=$((300 + i * 5))
  size=$((20 + i))
  go run . -width=$width -height=$height -font-size=$size "Batch Test $i"
done

# Animation stress test
for bg in default perlin radial diagonal; do
  go run . -animate -bg=$bg "Animation Stress $bg"
done