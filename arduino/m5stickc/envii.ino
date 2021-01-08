#include <WiFi.h>
#include <WiFiClient.h>
#include <WebServer.h>

#include <M5StickC.h>
#include <Wire.h>
#include <Adafruit_BMP280.h>
#include "Adafruit_Sensor.h"
#include "SHT3X.h"

#include <esp32-hal-gpio.h>

#define WIFI_SSID "somessid"
#define WIFI_PASSWORD "somepassword"
#define I2C_SDA 32
#define I2C_SCL 33
#define I2C_FREQ 100000
#define SERIAL_SPEED 115200
#define WEBSERVER_PORT 9000
#define TOP_BUTTON 37

Adafruit_BMP280 bme;
SHT3X sht30;

float c, h, p = 0.0;

WebServer server(WEBSERVER_PORT);

volatile int lcdTimeCounter;
hw_timer_t *lcdTimer = NULL;
portMUX_TYPE timerMux = portMUX_INITIALIZER_UNLOCKED;
bool isLcdOn = false;

// https://twitter.com/c_mos/status/1112633975552401408?lang=en
void IRAM_ATTR onPressed() {
  Serial.println("got interrupted");
  timerAlarmEnable(lcdTimer);
}


// https://55life555.blog.fc2.com/blog-entry-3194.html
void IRAM_ATTR onTimer() {
  portENTER_CRITICAL_ISR(&timerMux);
  lcdTimeCounter++;
  portENTER_CRITICAL_ISR(&timerMux);
  Serial.printf("now: %d\n", lcdTimeCounter);

}

void handleRoot() {
  server.send(200, "text/plain", "{\"c\":" + String(c) + ",\"h\":" + String(h) + ",\"p\":" + String(p) + "}");
}

void handleNotFound() {
  server.send(404, "text/plain", "");
}

void setup(void) {
  pinMode(TOP_BUTTON, INPUT_PULLUP);
  attachInterrupt(digitalPinToInterrupt(TOP_BUTTON), onPressed, FALLING);

  lcdTimer = timerBegin(0, 240, true);
  timerAttachInterrupt(lcdTimer, &onTimer, true);
  timerAlarmWrite(lcdTimer, 1000000, true);
  timerAlarmDisable(lcdTimer);

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

  M5.Axp.ScreenBreath(7);

  // SetLDO2()関数はボード自体がなぜか動作しなくなる
  // https://lang-ship.com/reference/unofficial/M5StickC/Class/AXP192/#screenbreath
}

void loop(void) {
  // panicを起こす
  //  if (lcdTimeCounter > 2) {
  //    portENTER_CRITICAL_ISR(&timerMux);
  //    lcdTimeCounter = 0;
  //    portENTER_CRITICAL_ISR(&timerMux);
  //    timerEnd(lcdTimer);
  //    isLcdOn = false;
  //  }

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