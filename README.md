# sonoff-lan-api
 Allows control of sonoff devices over LAN via Rest API

## Configuration
Setup devices within `configuration.yaml`, see example.  
Also an environment variable is required for each device `KEY_x` which is your device key.  
Each key relates to the devices in the order within configuration.yaml.  
```
KEY_1=xxxxx
KEY_2=xxxxx
```

## API Endpoints

`/deviceList`  
`/turnOn/<device-name>`  
`/turnOff/<device-name>`  

## SONOFF Curl commands
You can turn/on and off and that's it. Unable to get status with Basic R2 (also Basic R2 doesn't support DIY mode!).  

`sequence` - Not required  
`deviceid` - Not required  
`selfApikey` - Not required  
`data` - encrypted with device key (md5 hash'd) and iv | base64 encoded
`iv` - random 16 byte string | base64 encoded

```
curl -X POST -H "Connection: close" -H "Content-Type: application/json" http://192.168.1.156:8081/zeroconf/switch \
  -d '{
      "sequence": "0000000000000",
      "deviceid": "0000000000",
      "selfApikey": "123",
      "data": "xxxxxxx+mdhVu5HND4UcA7Gun5Lo5OH/t6si/PZ0=",
      "encrypt": true,
      "iv": "xxxxxxxNCysdfm8FafEQ=="
    }'
```

## References
https://github.com/antnks/sonoff-bash
