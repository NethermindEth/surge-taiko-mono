#!/bin/bash

# Find all files in the bridge-ui directory that contain crypto.randomUUID()
find . -type f \( -name "*.ts" -o -name "*.svelte" \) | while read -r file; do
    # Check if the file contains crypto.randomUUID()
    if grep -q "crypto\.randomUUID()" "$file"; then
        echo "Processing $file"
        # Replace crypto.randomUUID() with Math.random().toString(36).substring(2, 15)
        sed -i '' 's/crypto\.randomUUID()/Math.random().toString(36).substring(2, 15)/g' "$file"
    fi
done