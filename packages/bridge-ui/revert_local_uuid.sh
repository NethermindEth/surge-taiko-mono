#!/bin/bash

# Find all files in the bridge-ui directory that contain Math.random().toString(36).substring(2, 15)
find . -type f \( -name "*.ts" -o -name "*.svelte" \) | while read -r file; do
    # Check if the file contains Math.random().toString(36).substring(2, 15)
    if grep -q "Math\.random()\.toString(36)\.substring(2, 15)" "$file"; then
        echo "Processing $file"
        # Replace Math.random().toString(36).substring(2, 15) with crypto.randomUUID()
        sed -i '' 's/Math\.random()\.toString(36)\.substring(2, 15)/crypto.randomUUID()/g' "$file"
    fi
done