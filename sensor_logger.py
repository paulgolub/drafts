# This code listen port 1703 and get data from ESP-32 in "#RTNode1=1111" format

import socket
import datetime

# Server parameters
UDP_IP = "0.0.0.0"  # all interfaces
UDP_PORT = 1703
BUFFER_SIZE = 32  # buffer size

# Log file name
LOG_FILE = "sensor_data.log"

# Write on file
def log_data(data):
    with open(LOG_FILE, "a") as file:
        timestamp = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        file.write(f"{timestamp} - {data}\n")

# Create UDP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind((UDP_IP, UDP_PORT))

print(f"Listening on {UDP_IP}:{UDP_PORT}...")

def parse_string(input_str):
    # Delete '#' and split by '='
    node_part, values_part = input_str[1:].split('=')
    
    # Convert value into list
    values_list = list(values_part)
    
    # Combine the part with the node name and the list of values
    result = [node_part] + values_list
    
    return result

try:
    while True:
        # Get the data from client (ESP32)
        data, addr = sock.recvfrom(BUFFER_SIZE)
        data_str = data.decode('utf-8').strip()
        
        # logging the data
        log_data(data_str)

        parsed_list = parse_string(data_str)

        print(f"Received message from {addr}: {data_str}")
        print(f"Parsed message from {addr}: {parsed_list}")        

except KeyboardInterrupt:
    print("Server stopped.")

finally:
    # Close socket in the end
    sock.close()
