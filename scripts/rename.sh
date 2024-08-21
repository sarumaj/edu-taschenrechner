#!/bin/bash

set -euo pipefail

# Base directory containing the subfolders
base_dir="./fyne-cross/dist"

target_dir="./dist"
mkdir -p "$target_dir"

# Loop through each subfolder in the base directory
for subfolder in "$base_dir"/*/; do
    # Extract the platform or architecture name from the subfolder path
    platform=$(basename "$subfolder")

    # Loop through each file in the subfolder
    for file in "$subfolder"*; do
        if [ -f "$file" ]; then
            # Construct the new filename with the platform name
            if [[ "$file" == *.tar.gz ]]; then
                new_file="$target_dir/$(basename "$file" .tar.gz)_${platform}.tar.gz"
            elif [[ "$file" == *.exe.zip ]]; then
                new_file="$target_dir/$(basename "$file" .exe.zip)_${platform}.exe.zip"
            else
                extension="${file##*.}"
                if [ -z "$extension" ]; then
                    new_file="$target_dir/$(basename "$file")_${platform}"
                else
                    new_file="$target_dir/$(basename "$file" ."$extension")_${platform}.$extension"
                fi
            fi

            # Rename the file
            mv "$file" "$new_file"

            echo "Renamed $file to $new_file"
        fi
    done
done
