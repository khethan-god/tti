#!/bin/bash

# complete feature matix testing

# Generate a comprehensive test matrix
declare -a fonts=("roboto_bold" "roboto_medium" "roboto_bold" "roboto_italic")
declare -a backgrounds=("default" "perlin" "radial" "diagonal")
declare -a sizes=("24" "48" "72")
declare -a dimensions=("600x300" "800x400" "1024x768")

for font in "${fonts[@]}"; do
  for bg in "${backgrounds[@]}"; do
    for size in "${sizes[@]}"; do
      for dim in "${dimensions[@]}"; do
        IFS='x' read -ra DIM <<< "$dim"
        width=${DIM[0]}
        height=${DIM[1]}
        
        # Static image
        go run . -font=$font -bg=$bg -font-size=$size -width=$width -height=$height "Matrix: $font-$bg-$size-$dim"
        
        # Animated version
        go run . -font=$font -bg=$bg -font-size=$size -width=$width -height=$height -animate "Anim: $font-$bg-$size-$dim"
        
        # Reveal background version
        go run . -font=$font -bg=$bg -font-size=$size -width=$width -height=$height -reveal-bg "Reveal: $font-$bg-$size-$dim"
      done
    done
  done
done