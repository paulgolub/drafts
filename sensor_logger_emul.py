import socket
import time

DEVICE_ID = "RTNode6"
MAX_REPEAT = 3

LOCAL_PORT = 1702
RECEIVER_PORT = 1703

DATA_LOGGER_IP = "192.168.0.160"
BUFFER_SIZE = 32

# Initialize UDP socket
udp_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
udp_socket.bind(("", LOCAL_PORT))

sensor_status = [True, True, True, True]
send_msg = True
repeat_msg = 0
last_message_ts = 0

sensors = [True, True, True, True]


def send_udp_message(buffer):
    udp_socket.sendto(buffer, (DATA_LOGGER_IP, RECEIVER_PORT))


while True:
    for i in range(4):
        sensor_val = sensors[i]

    elapsed = time.time() * 1000 - last_message_ts
    if elapsed > 1000:
        send_msg = True
        repeat_msg = 0

    if send_msg:
        message = f"#{DEVICE_ID}={int(not sensor_status[0])}{int(not sensor_status[1])}{int(not sensor_status[2])}{int(not sensor_status[3])}"
        send_udp_message(message.encode('utf-8'))
        last_message_ts = time.time() * 1000
        repeat_msg += 1
        if repeat_msg >= MAX_REPEAT:
            send_msg = False

    time.sleep(0.1)
