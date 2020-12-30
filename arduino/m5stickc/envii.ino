#include <WiFi.h>
#include <WiFiClient.h>
#include <WebServer.h>

#include <M5StickC.h>
#include <Wire.h>
#include <Adafruit_BMP280.h>
#include "Adafruit_Sensor.h"
#include "SHT3X.h"

const char* WIFI_SSID = "somessid";
const char* WIFI_PASSWORD = "somepassword";
const int I2C_SDA = 32;
const int I2C_SCL = 33;
const uint32_t I2C_FREQ = 100000;
const long SERIAL_SPEED = 115200;
const int WEBSERVER_PORT = 9000;

Adafruit_BMP280 bme;
SHT3X sht30;

float c, h, p = 0.0;

WebServer server(WEBSERVER_PORT);

void handleRoot() {
  server.send(200, "text/plain", "{\"c\":" + String(c) + ",\"h\":" + String(h) + ",\"p\":" + String(p) + "}");
}

void handleNotFound() {
  server.send(404, "text/plain", "");
}

void setup(void) {
  // starts LCD
  M5.Lcd.setRotation(3);
  M5.Lcd.setTextColor(WHITE, BLACK);

  // starts Serial Port
  Serial.begin(SERIAL_SPEED);
  M5.begin(true, true, true);
  Wire.begin(I2C_SDA, I2C_SCL, I2C_FREQ);
  while (!bme.begin(0x76)) {
    Serial.println("Could not find a valid BMP280 sensor, check wiring!");
  }

  // starts WiFi
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  Serial.println("");
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.printf("Connected to SSID:%s\n", WIFI_SSID);

  // starts HTTP server
  server.on("/", handleRoot);
  server.onNotFound(handleNotFound);
  server.begin();
  Serial.println("HTTP server started");
}

void loop(void) {
  // gets senser data
  if (sht30.get() == 0) {
    c = sht30.cTemp;
    h = sht30.humidity;
    p = bme.readPressure();
  }

  server.handleClient();
  M5.Lcd.setCursor(0, 15, 2);
  M5.Lcd.printf("Temperature: %.2f C\nHumidity: %.2f %%\nPressure: %.2f Pa", c, h, p);
}