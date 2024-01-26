#!/bin/bash

# This script is to run all examples at once in order to check them
# Need to run this script from `./examples/` folder

set -e

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
            echo -e "\n----------------------------------------------------------------------\nRunning $file"
            go run "$file"
            echo -e "Finish...\n----------------------------------------------------------------------"
        done

        cd "$root_directory" || exit
    fi
done
