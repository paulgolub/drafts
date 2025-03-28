#this code count media brightness of image
import cv2
import numpy as np

image = cv2.imread("image.jpg")

gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

mean_brightness = np.mean(gray)

print(f"brightness med: {mean_brightness}")
