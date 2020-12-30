// https://github.com/m5stack/M5-ProductExampleCodes/blob/0599f9ee322385df47eb8f302d30b59bffe47b92/Unit/ENVII/Arduino/ENVII/SHT3X.h

#ifndef __SHT3X_H
#define __HT3X_H


#if ARDUINO >= 100
 #include "Arduino.h"
#else
 #include "WProgram.h"
#endif

#include "Wire.h"

class SHT3X{
public:
  SHT3X(uint8_t address=0x44);
  byte get(void);
  float cTemp=0;
  float fTemp=0;
  float humidity=0;

private:
  uint8_t _address;

};


#endif
