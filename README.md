# ws-upload

Influx exporter for weather stations running EasyWeather software.

Influenced by [Domoticz-PWS-Plugin](https://github.com/Xorfor/Domoticz-PWS-Plugin).

## Prerequisites

### WSView Plus
You can use the WSView Plus app to set up your  weather station to upload data to ws-upload. Follow the instructions
below.

1. Install the WSView Plus app.
  - [Play Store](https://play.google.com/store/apps/details?id=com.ost.wsautool)
  - [App Store](https://apps.apple.com/nl/app/wsview-plus/id1581353359)
1. Open the app and connect to your weather station.
2. Select *Customized* in the menu.
3. Select *Enable*.
4. Enter the server IP/hostname where ws-upload is running.
5. Enter the path `/api/v1/observe?` (including the final question mark).
6. Enter a station ID. This is an arbitrary string that identifies your weather station. It is used to create the
   measurement name in InfluxDB. You can use the name of your weather station, for example.
7. Enter a station key. This is an arbitrary string that is used to authenticate your weather station. This should be
   as secure as a password. You can use a password generator to create a random string. You will need this key later to
   configure ws-upload.
8. Enter the port where ws-upload is running. The default port is 9108.
9. Set an upload interval
10. Click on *Save*.

<details>
<summary>Older instructions for WS View or WS Tool</summary>

See [this section on the Domoticz-PWS-Plugin page](https://github.com/Xorfor/Domoticz-PWS-Plugin#prerequisites).

</details>

### InfluxDB

InfluxDB is not required. ws-upload can also export data via MQTT. However, if you are using InfluxDB, you will need
to configure ws-upload correctly.

### MQTT

MQTT is not required. ws-upload can also export data to InfluxDB. However, if you are using MQTT, you will need
to configure ws-upload correctly.

## Installation

### Docker

The easiest way to run ws-upload is by using Docker. You can use the following command to run ws-upload:

```shell
docker run -d \
  --name ws-upload \
  -p 9108:9108 \
  -e STATION_PASSWORD=your_station_key \
  --restart unless-stopped \
  ghcr.io/koesie10/ws-upload:latest
```

Replace `your_station_key` with the station key you entered in the WSView Plus app.

### Manual

Binaries are available on the [releases page](https://github.com/koesie10/ws-upload/releases).
