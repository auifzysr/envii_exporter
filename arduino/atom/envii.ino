// ref: https://github.com/m5stack/M5-ProductExampleCodes/blob/0599f9ee322385df47eb8f302d30b59bffe47b92/Unit/ENVII/Arduino/ENVII/ENVII.ino
// ref: https://gist.github.com/ksasao/cf601184691e59297fd464af910966b4

// need to install:
// https://github.com/m5stack/M5-ProductExampleCodes/blob/0599f9ee322385df47eb8f302d30b59bffe47b92/Unit/ENVII/Arduino/ENVII/SHT3X.h
// https://github.com/m5stack/M5-ProductExampleCodes/blob/0599f9ee322385df47eb8f302d30b59bffe47b92/Unit/ENVII/Arduino/ENVII/SHT3X.cpp
// https://github.com/m5stack/M5-ProductExampleCodes/blob/0599f9ee322385df47eb8f302d30b59bffe47b92/Unit/ENVII/Arduino/ENVII/Adafruit_Sensor.h

#include <WiFi.h>
#include <WiFiClient.h>
#include <WebServer.h>

#include <M5Atom.h>
#include <Wire.h>
#include <Adafruit_BMP280.h>
#include "Adafruit_Sensor.h"
#include "SHT3X.h"

const char* WIFI_SSID = "somessid";
const char* WIFI_PASSWORD = "somepassword";
const int I2C_SDA = 26;
const int I2C_SCL = 32;
const uint32_t I2C_FREQ = 100000;
const long SERIAL_SPEED = 115200;
const int WEBSERVER_PORT = 9000;

Adafruit_BMP280 bme;
SHT3X sht30;

float c, h, p = 0.0;

WebServer server(WEBSERVER_PORT);

uint8_t DisBuff[2 + 5 * 5 * 3];

void setBuff(uint8_t Rdata, uint8_t Gdata, uint8_t Bdata)
{
  DisBuff[0] = 0x05;
  DisBuff[1] = 0x05;
  for (int i = 0; i < 25; i++)
  {
    DisBuff[2 + i * 3 + 0] = Rdata;
    DisBuff[2 + i * 3 + 1] = Gdata;
    DisBuff[2 + i * 3 + 2] = Bdata;
  }
}

void handleRoot() {
  server.send(200, "text/plain", "{\"c\":" + String(c) + ",\"h\":" + String(h) + ",\"p\":" + String(p) + "}");
}

void handleNotFound() {
  server.send(404, "text/plain", "");
}

void setup(void) {
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
    setBuff(0x00, 0x40, 0x00);
  } else {
    setBuff(0x40, 0x00, 0x00);
  }

  M5.dis.displaybuff(DisBuff);
  server.handleClient();
}