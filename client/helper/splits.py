#!/usr/bin/python3
import os
import zipfile
import sys

def create_n_splits(main_folder, n_splits):
    os.makedirs("splits", exist_ok=True)
    # Iterate through all subfolders in the main folder
    for root, dirs, files in os.walk(main_folder):
        # Exclude hidden folders (those starting with '.')
        dirs[:] = [d for d in dirs if not d.startswith('.')]

        num_files = len(files)
        split_amount = num_files // n_splits
        

        for i in range(n_splits):
            zip_file_path = f"./splits/split_{i+1}.zip"
            start_index = split_amount * i
            if i == 0:
                file_split = files[start_index:split_amount]
            elif i + 1 == n_splits:
                file_split = files[start_index:]
            else:
                file_split = files[start_index:start_index * 2]

            with zipfile.ZipFile(zip_file_path, 'a', zipfile.ZIP_DEFLATED) as zipf:
                for file in file_split:
                    file_path = os.path.join(root, file)
                    print(root + "/" + file)
                    zipf.write(file_path, arcname=root + "/" + file)
            
# Example usage:
n_peers = int(sys.argv[1])
training_data_folder = sys.argv[2]

create_n_splits(training_data_folder, n_peers)
