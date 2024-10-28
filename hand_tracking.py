# hand tracking, x,y,z coordinates (with determination of the depth of the position)
# pip install opencv-python mediapipe

import cv2
import mediapipe as mp

# Initialize MediaPipe Hands
mp_hands = mp.solutions.hands
hands = mp_hands.Hands(static_image_mode=False, max_num_hands=1, min_detection_confidence=0.7)

# Initialize OpenCV for video capture
cap = cv2.VideoCapture(0)

while cap.isOpened():
    success, image = cap.read()
    if not success:
        print("Failed to get frame.")
        break

    # Convert the image color from BGR to RGB
    image_rgb = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
    results = hands.process(image_rgb)

    # If hands are detected
    if results.multi_hand_landmarks:
        for hand_landmarks in results.multi_hand_landmarks:
            # Get the coordinates of the index finger tip (tip, index 8)
            index_finger_tip = hand_landmarks.landmark[8]
            h, w, c = image.shape
            cx, cy = int(index_finger_tip.x * w), int(index_finger_tip.y * h)
            z = index_finger_tip.z  # Z-coordinate (depth)

            # Draw a circle on the image for the index finger tip
            cv2.circle(image, (cx, cy), 5, (0, 255, 0), -1)

            # Output the coordinates x, y, z
            print(f'Coordinates of the index finger: x={cx}, y={cy}, z={z}')

            # Draw lines between key points
            mp.solutions.drawing_utils.draw_landmarks(image, hand_landmarks, mp_hands.HAND_CONNECTIONS)

    # Display the image
    cv2.imshow('Hand Tracking', image)

    if cv2.waitKey(5) & 0xFF == 27:  # Exit on pressing 'Esc'
        break

cap.release()
cv2.destroyAllWindows()
