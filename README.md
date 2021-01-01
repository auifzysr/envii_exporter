# envii_exporter
pull metrics from M5Stack ATOM Matrix with ENV II Unit and expose them to Prometheus

## setup
1. install the codes required to run with `envii.ino`
1. upload the codes to Atom or M5StickC
1. deploy an exporter which connects to the Atom or the M5StickC for each
1. integrate the exporter with a Prometheus instance