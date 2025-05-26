#!/bin/bash

# edge cases and error testing

# Negative values (should fail gracefully)
go run . -width=-100 "Negative Width"
go run . -height=-100 "Negative Height"
go run . -font-size=-10 "Negative Font Size"

# Zero values (should fail gracefully)
go run . -width=0 "Zero Width"
go run . -height=0 "Zero Height"
go run . -font-size=0 "Zero Font Size"

# Extremely large values (memory stress test)
go run . -width=10000 -height=10000 "Huge Image"
go run . -font-size=500 "Huge Font"

# Text that might cause issues
go run . "Text with\nnewlines"
go run . "Text	with	tabs"
go run . "Text with / slashes \\ and \\ backslashes"
go run . "Filenames: <>:|?*"