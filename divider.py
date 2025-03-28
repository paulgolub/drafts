import shutil
import re
import os

def split_files_by_number(directory):
    even_dir = os.path.join(directory, "even")
    odd_dir = os.path.join(directory, "odd")
    os.makedirs(even_dir, exist_ok=True)
    os.makedirs(odd_dir, exist_ok=True)
    
    pattern = re.compile(r'.*?(\d{3}).*')
    
    for filename in os.listdir(directory):
        filepath = os.path.join(directory, filename)
        
        if os.path.isfile(filepath):
            match = pattern.match(filename)
            if match:
                number = int(match.group(1))
                target_dir = even_dir if number % 2 == 0 else odd_dir
                shutil.move(filepath, os.path.join(target_dir, filename))
                print(f"Moved {filename} to {'even' if number % 2 == 0 else 'odd'} folder")

if __name__ == "__main__":
    directory = input("Enter directory path: ")
    split_files_by_number(directory)
