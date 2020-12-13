// ref: https://github.com/m5stack/M5-ProductExampleCodes/blob/0599f9ee322385df47eb8f302d30b59bffe47b92/Unit/ENVII/Arduino/ENVII/ENVII.ino
// ref: https://gist.github.com/ksasao/cf601184691e59297fd464af910966b4

#include <WiFi.h>
#include <WiFiClient.h>
#include <WebServer.h>

#include <M5Atom.h>
#include <Wire.h>
#include <Adafruit_BMP280.h>
#include "Adafruit_Sensor.h"
#include "SHT3X.h"

const char* ssid = "wifissid";
const char* password = "wifipassword";

WebServer server(9000);

Adafruit_BMP280 bme;
SHT3X sht30;

float tmp = 0.0;
float hum = 0.0;
float pressure = 0.0;

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
  if (sht30.get() == 0) {
    setBuff(0x00, 0x40, 0x00);
    M5.dis.displaybuff(DisBuff);
    M5.update();
    server.send(200, "text/plain",
                "{\"c\":" + String(sht30.cTemp)
                + ",\"h\":" + String(sht30.humidity)
                + ",\"p\":" + String(bme.readPressure())
                + "}");
    return;
  }
  setBuff(0x40, 0x00, 0x00);
  M5.dis.displaybuff(DisBuff);
  M5.update();
  server.send(500, "text/plain", "");
}

void handleNotFound() {
  String message = "File Not Found\n\n";
  message += "URI: ";
  message += server.uri();
  message += "\nMethod: ";
  message += (server.method() == HTTP_GET) ? "GET" : "POST";
  message += "\nArguments: ";
  message += server.args();
  message += "\n";
  for (uint8_t i = 0; i < server.args(); i++) {
    message += " " + server.argName(i) + ": " + server.arg(i) + "\n";
  }
  server.send(404, "text/plain", message);
}

void setup(void) {
  Serial.begin(115200);
  WiFi.begin(ssid, password);
  Serial.println("");

  // Wait for connection
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.print("Connected to ");
  Serial.println(ssid);
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

  server.on("/", handleRoot);

  server.on("/inline", []() {
    server.send(200, "text/plain", "this works as well");
  });

  server.onNotFound(handleNotFound);

  server.begin();
  Serial.println("HTTP server started");

  M5.begin(true, true, true);
  Wire.begin(26, 32);
  Serial.println(F("ENV Unit(SHT30 and BMP280) test..."));

  while (!bme.begin(0x76)) {
    Serial.println("Could not find a valid BMP280 sensor, check wiring!");
  }

}

void loop(void) {
  server.handleClient();
}