#!/bin/bash

# This script is to run all examples at once in order to check them
# Need to run this script from `./examples/` folder

set -e
skip_items=("/util/")

# Root directory
root_directory=$(pwd)

# Iterate over subdirectories
find "$root_directory" -type d -print0 | while IFS= read -r -d '' directory; do
    echo -e "\n\nChecking $directory..."

    # Check if the directory contains a .mod file
    if [ -f "$directory/go.mod" ]; then
        echo -e "Running go mod tidy..."
        cd "$directory"
        go mod tidy

        # Run all .go files in the directory and its subdirectories
        echo -e "Running go files:"
        find "$directory" -type f -name '*.go' | while read -r file; do
            # Check if the file or directory should be skipped
            skip=false
            echo -e "Checking $file"
            for skip_item in "${skip_items[@]}"; do
                if [[ "$file" == *"$skip_item"* ]]; then
                    echo "Skipping $file"
                    skip=true
                    break
                fi
            done

            # Run the file if it should not be skipped
            if [ "$skip" = false ]; then
                echo -e "\n----------------------------------------------------------------------\nRunning $file"
                go run "$file"
                echo -e "Finish...\n----------------------------------------------------------------------"
            fi
        done

        cd "$root_directory" || exit
    fi
done
