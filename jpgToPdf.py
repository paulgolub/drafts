#This script converts all JPG, JPEG, TIF, and TIFF images from a specified folder into a single PDF file. It scans the folder, loads the images, converts them to RGB format, and merges them into a PDF. If no images are found, it notifies the user. The script prompts for the input folder and output PDF name, using the current directory and "output.pdf" as defaults.

from PIL import Image
import os

def convert_jpg_to_pdf(input_folder, output_pdf):
    images = []
    
    # Iterate through files in the specified folder
    for file in sorted(os.listdir(input_folder)):
        # Check if the file is an image with the allowed extensions
        if file.lower().endswith(".jpg") or file.lower().endswith(".jpeg") or file.lower().endswith(".tif") or file.lower().endswith(".tiff"):
            img_path = os.path.join(input_folder, file)
            img = Image.open(img_path).convert("RGB")  # Convert image to RGB format
            images.append(img)
    
    if images:
        # Save the first image and append the rest to the same PDF file
        images[0].save(output_pdf, save_all=True, append_images=images[1:])
        print(f"PDF saved: {output_pdf}")
    else:
        print("No JPG files found for conversion.")

# Example usage
default_folder = os.getcwd()  # Directory where the script is executed
input_folder = input("Enter a path to JPG (current dir as default): ").strip() or default_folder
output_pdf = input("Enter name for output PDF (output.pdf as default): ").strip() or "output.pdf"

convert_jpg_to_pdf(input_folder, output_pdf)
