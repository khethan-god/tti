#!/bin/bash

# Font + Background combinations (25 fonts × 4 backgrounds = 100 combinations)
for font in roboto_sc_ebold roboto_bold roboto_c_bold roboto_sc_italic roboto_ebold roboto_c_ebitalic roboto_sc_light roboto_elight roboto_c_elight roboto_sc_litalic roboto_italic roboto_c_elitalic roboto_sc_medium roboto_light roboto_c_regular roboto_sc_mitalic roboto_medium roboto_c_titalic roboto_sc_sbold roboto_regular roboto_sc_blitalic roboto_sc_sbitalic roboto_sbold roboto_sc_bitalic roboto_sc_thin ; do
  for bg in default perlin radial diagonal; do
    go run . -font=$font -bg=$bg "Font: $font, BG: $bg"
  done
done

# Dimension + Font Size combinations
go run . -width=800 -height=600 -font-size=64 "Large Canvas Big Font"
go run . -width=200 -height=150 -font-size=12 "Small Canvas Small Font"
go run . -width=1000 -height=100 -font-size=48 "Wide Canvas Medium Font"

# All parameters combined
go run . -width=1024 -height=768 -font-size=56 -font=bold -bg=perlin -output=test-output "Full Combo"
go run . -width=800 -height=600 -font-size=42 -font=times -bg=radial -reveal-bg -output=reveal-test "Reveal Combo"
go run . -width=600 -height=400 -font-size=36 -font=courier -bg=diagonal -animate -output=animation-test "Animation Combo"