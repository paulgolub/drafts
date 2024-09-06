import socket
import time
import random

DEVICE_ID = "Node6"
MAX_REPEAT = 3

LOCAL_PORT = 1702
RECEIVER_PORT = 8085

DATA_LOGGER_IP = "192.168.0.2"
BUFFER_SIZE = 32

# Initialize UDP socket
udp_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
udp_socket.bind(("", LOCAL_PORT))

sensor_status = [False, False, False, False]
send_msg = True
repeat_msg = 0
last_message_ts = 0
status_change_ts = time.time() * 1000

def send_udp_message(buffer):
    udp_socket.sendto(buffer, (DATA_LOGGER_IP, RECEIVER_PORT))

def update_sensor_status():
    # Ensure only one of the first three sensors is True
    true_indices = [i for i in range(4) if sensor_status[i]]
    if true_indices:
        for i in true_indices:
            sensor_status[i] = False
    
    # Randomly activate one of the first three sensors
    active_sensor = random.randint(0, 3)
    sensor_status[active_sensor] = True

while True:
    current_time = time.time() * 1000
    elapsed = current_time - last_message_ts
    status_elapsed = current_time - status_change_ts

    if status_elapsed > 3000:  # Change sensor status every 3 seconds
        update_sensor_status()
        status_change_ts = current_time

    if elapsed > 1000:
        send_msg = True
        repeat_msg = 0

    if send_msg:
        message = f"#{DEVICE_ID}={int(not sensor_status[0])}{int(not sensor_status[1])}{int(not sensor_status[2])}{int(not sensor_status[3])}"
        print(message)
        send_udp_message(message.encode('utf-8'))
        last_message_ts = current_time
        repeat_msg += 1
        if repeat_msg >= MAX_REPEAT:
            send_msg = False

    time.sleep(0.1)
