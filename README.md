## Prerequisites
Your PWS needs to be connected to your network. This can be done by using [WS View](#ws-view-ws-tool) or the [Garni Technology](#garni) app.

### WS View (WS Tool)
If supported by your PWS, connect your PWS with `WS View` (and also the 'older' `WS Tool`) to your router by wifi, so that your PWS can upload weather data to Domoticz.

1. Install `WS View` on your mobile device
    * [Google Play Store](https://play.google.com/store/apps/details?id=com.ost.wsview)
    * [Apple App Store](https://apps.apple.com/us/app/ws-view/id1362944193)
1. , or `WS Tool`
    * [Google Play Store](https://play.google.com/store/apps/details?id=com.dtston.wstool)
    * [Apple App Store](https://apps.apple.com/nl/app/ws-tool/id1125344077)
1. Follow the instructions to connect your PWS to your router
1. Goto to Device List in Menu and choose your PWS
1. Click on Next untill you are on on the `Customized` page
1. Choose `Enable`
1. For `Protocol Type Same As` choose `Wunderground` (preferred)
1. For `Server IP / Hostname` enter your Domoticz Server ip address, eg. 192.168.0.10
1. If you choose for `Wunderground` protocol:
    * Fill in `Station ID` with a value
    * Fill in `Station Key` with a value
1. `Port` enter a free port number, eg. `5000`
1. `Upload Interval`, leave it `60` seconds
1. Click on `Save`

![Screenshot](/images/screendump2.png) ![Screenshot](/images/screendump3.png)
