#include <OneWire.h>
#include <DallasTemperature.h>
#include <WiFi.h>
#include <UniversalTelegramBot.h>

#define PRESSURE_SENSOR_PIN A0
#define WATER_LEVEL_PIN A1
#define MOTOR_TEMP_PIN 2
#define RELAY_PIN 3
#define FLOW_SENSOR_PIN 4
#define BUTTON_PIN 5 // Button for manual control

#define LED_BLUE_PIN 6   // Normal operation mode
#define LED_RED_PIN 7    // Low pressure error
#define LED_YELLOW_PIN 8 // Pump priming mode
#define LED_GREEN_PIN 9  // Overheating protection
#define LED_ORANGE_PIN 10 // Manual mode indicator
#define LED_WIFI_PIN 11   // Wi-Fi connection indicator
#define LED_ERROR_PIN 12  // Telegram message send error indicator

// Wi-Fi settings
const char* ssid = "YOUR_SSID";           // Enter your Wi-Fi SSID
const char* password = "YOUR_PASSWORD";   // Enter your Wi-Fi password

// Telegram Bot Token and Chat ID
#define BOT_TOKEN "YOUR_BOT_API_TOKEN"     // Your bot token
#define CHAT_ID "YOUR_CHAT_ID"             // Your chat ID

WiFiClientSecure client;
UniversalTelegramBot bot(BOT_TOKEN, client);

// Temperature sensor object
OneWire oneWire(MOTOR_TEMP_PIN);
DallasTemperature sensors(&oneWire);

// Calibration for pressure sensor
const float PRESSURE_SENSOR_VOLTAGE_TO_BAR = 1.0; // For example, 1V = 1 bar
const float MIN_PRESSURE_THRESHOLD = 0.2; // Minimum pressure to avoid dry running

// Parameters for checking
const int MAX_PUMP_TIME = 5000; // Max priming time in milliseconds (5 seconds)
const int MAX_PUMP_RETRIES = 3; // Number of priming retries

int pumpRetries = 0;
unsigned long pumpStartTime = 0;

void setup() {
  pinMode(RELAY_PIN, OUTPUT);
  pinMode(BUTTON_PIN, INPUT_PULLUP); // Button for manual operation
  pinMode(LED_BLUE_PIN, OUTPUT);
  pinMode(LED_RED_PIN, OUTPUT);
  pinMode(LED_YELLOW_PIN, OUTPUT);
  pinMode(LED_GREEN_PIN, OUTPUT);
  pinMode(LED_ORANGE_PIN, OUTPUT);
  pinMode(LED_WIFI_PIN, OUTPUT);    // Wi-Fi connection indicator
  pinMode(LED_ERROR_PIN, OUTPUT);   // Telegram message send error indicator

  Serial.begin(9600);
  sensors.begin();

  // Connect to Wi-Fi
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    digitalWrite(LED_WIFI_PIN, LOW);   // Turn off LED if not connected
    delay(1000);
    Serial.println("Connecting to WiFi...");
  }
  digitalWrite(LED_WIFI_PIN, HIGH);    // Turn on LED when connected
  Serial.println("Connected to WiFi");

  // Send a message to Telegram that the system has started
  sendTelegramMessage("System initialized and connected to WiFi.");
}

void loop() {
  // Read pressure from the sensor
  int pressureRaw = analogRead(PRESSURE_SENSOR_PIN);
  float pressureVoltage = pressureRaw * (5.0 / 1023.0);
  float pressure = pressureVoltage * PRESSURE_SENSOR_VOLTAGE_TO_BAR;

  // Read water level
  int waterLevel = analogRead(WATER_LEVEL_PIN);

  // Read motor temperature
  sensors.requestTemperatures();
  float motorTemp = sensors.getTempCByIndex(0);

  // Read button state
  bool buttonPressed = digitalRead(BUTTON_PIN) == LOW;

  // Control LEDs based on the current state
  if (buttonPressed) {
    // If the button is pressed, turn on the orange LED for manual mode
    digitalWrite(LED_ORANGE_PIN, HIGH);
    digitalWrite(LED_BLUE_PIN, LOW); 
    digitalWrite(LED_RED_PIN, LOW); 
    digitalWrite(LED_YELLOW_PIN, LOW); 
    digitalWrite(LED_GREEN_PIN, LOW); 

    // Turn on the pump in manual mode
    digitalWrite(RELAY_PIN, HIGH);  
    Serial.println("Button pressed. Pump running manually.");
    sendTelegramMessage("Manual mode activated. Pump running.");
  } 
  else {
    if (pressure < MIN_PRESSURE_THRESHOLD) {
      // If pressure is low, show waiting mode (add water to the pump)
      digitalWrite(LED_ORANGE_PIN, HIGH);  // Turn on the orange LED (waiting mode)
      digitalWrite(LED_BLUE_PIN, LOW);
      digitalWrite(LED_RED_PIN, LOW);
      digitalWrite(LED_YELLOW_PIN, LOW);
      digitalWrite(LED_GREEN_PIN, LOW);

      Serial.println("Low pressure detected. Please add water to the pump.");
      sendTelegramMessage("Low pressure detected. Please add water to the pump.");

      // Wait for pressure to increase after water is added
      digitalWrite(RELAY_PIN, LOW);  // Pump is off in waiting mode

      // Once the water is added and pressure increases, system will switch to normal mode
      if (pressure >= MIN_PRESSURE_THRESHOLD) {
        digitalWrite(LED_ORANGE_PIN, LOW);  // Turn off the orange LED (waiting mode)
        Serial.println("System pressurized. You can start the pump.");
        sendTelegramMessage("System pressurized. You can start the pump.");
      }
    } 
    else if (pressure < MIN_PRESSURE_THRESHOLD && pumpRetries < MAX_PUMP_RETRIES) {
      // If pressure is still low, start priming
      digitalWrite(LED_YELLOW_PIN, HIGH);  // Turn on yellow LED (priming mode)
      digitalWrite(LED_BLUE_PIN, LOW);
      digitalWrite(LED_RED_PIN, LOW);
      digitalWrite(LED_GREEN_PIN, LOW);

      if (millis() - pumpStartTime >= MAX_PUMP_TIME) {
        // If priming time is exceeded, check pressure
        pumpRetries++;
        pumpStartTime = millis();
        digitalWrite(RELAY_PIN, LOW); // Turn off the pump for a pause

        // Check if pressure is still low
        if (pressure < MIN_PRESSURE_THRESHOLD) {
          Serial.println("Low pressure. Retrying pump.");
          sendTelegramMessage("Low pressure detected. Retrying pump.");
          digitalWrite(RELAY_PIN, HIGH); // Turn on the pump again for priming
        } else {
          Serial.println("System pressurized. Pumping complete.");
          sendTelegramMessage("System pressurized. Pumping complete.");
        }
      }
    } 
    else {
      // After successful priming, the system should operate normally
      if (pressure >= MIN_PRESSURE_THRESHOLD) {
        digitalWrite(LED_BLUE_PIN, HIGH);  // Turn on blue LED (normal mode)
        digitalWrite(LED_YELLOW_PIN, LOW);
        digitalWrite(LED_RED_PIN, LOW);
        digitalWrite(LED_GREEN_PIN, LOW);

        // Pressure is normal, turn on the pump
        digitalWrite(RELAY_PIN, HIGH);
        Serial.println("Pump on. System pressurized.");
        sendTelegramMessage("Pump on. System pressurized.");
      }
      else {
        // If pressure is low, turn off the pump to protect from dry running
        digitalWrite(LED_RED_PIN, HIGH);  // Turn on red LED (error)
        digitalWrite(LED_BLUE_PIN, LOW);
        digitalWrite(LED_YELLOW_PIN, LOW);
        digitalWrite(LED_GREEN_PIN, LOW);

        digitalWrite(RELAY_PIN, LOW);
        Serial.println("Pump off due to low pressure.");
        sendTelegramMessage("Pump off due to low pressure.");
      }
    }
  }

  // Overheating protection logic
  if (motorTemp > 50) {
    digitalWrite(LED_GREEN_PIN, HIGH);  // Turn on green LED (overheating)
    digitalWrite(LED_BLUE_PIN, LOW);
    digitalWrite(LED_RED_PIN, LOW);
    digitalWrite(LED_YELLOW_PIN, LOW);
    digitalWrite(LED_ORANGE_PIN, LOW);

    digitalWrite(RELAY_PIN, LOW); // Turn off the pump in case of overheating
    Serial.println("Motor overheating. Pump off.");
    sendTelegramMessage("Motor overheating. Pump off.");
  }

  // Print values to serial monitor for debugging
  Serial.print("Pressure: ");
  Serial.print(pressure);
  Serial.print(" bar, Water level: ");
  Serial.print(waterLevel);
  Serial.print(", Motor temp: ");
  Serial.println(motorTemp);

  delay(1000); // Delay between measurements
}

// Function to send messages to Telegram
void sendTelegramMessage(String message) {
  if (WiFi.status() == WL_CONNECTED) {
    bool sent = bot.sendMessage(CHAT_ID, message, "");
    if (!sent) {
      // If message failed to send, turn on red LED
      digitalWrite(LED_ERROR_PIN, HIGH);
      Serial.println("Error sending Telegram message.");
    } else {
      // If message sent successfully, turn off red LED
      digitalWrite(LED_ERROR_PIN, LOW);
    }
  } else {
    Serial.println("WiFi not connected, message not sent.");
    digitalWrite(LED_ERROR_PIN, HIGH);  // Turn on red LED if no Wi-Fi connection
  }
}
