

function Decode(fPort, bytes) {
    var str = "";
    for (var i = 0; i < bytes.length; i++) {
        var s = bytes[i].toString(16);
        if (s.length < 2) {
            str += "0" + s;
        } else {
            str += s;
        }
    }
    var start = 0;
    if (str.substr(start, 2) != "aa" || str.length < 40) {
        return
    }
    start += 2;
    var ti = parseInt(str.substr(start, 8), 16);
    start += 8;
    var latitude = parseInt(str.substr(start, 8), 16) / 1000000;
    start += 8;
    var longitude = parseInt(str.substr(start, 8), 16) / 1000000;
    start += 8;
    var altitude = parseInt(str.substr(start, 4), 16);
    start += 4;
    var speed = parseInt(str.substr(start, 2), 16);
    start += 2;
    var azimuth = parseInt(str.substr(start, 4), 16);
    start += 4;
    var RFU = parseInt(str.substr(start, 4), 16);

    if (longitude < 72.004 || longitude > 137.8347 || latitude < 0.8293 || latitude > 55.8271)
        return {
            GPS: [
                {
                    time: ti,
                    latitude: latitude,
                    longitude: longitude,
                    altitude: altitude,
                    speed: speed,
                    azimuth: azimuth,
                    RFU: RFU
                }
            ]
        }
    var PI = 3.14159265358979324;
    var a = 6378245.0; //  a: 卫星椭球坐标投影到平面地图坐标系的投影因子。
    var ee = 0.00669342162296594323; //  ee: 椭球的偏心率。

    var x = longitude - 105.0;
    var y = latitude - 35.0;
    var ret = -100.0 + 2.0 * x + 3.0 * y + 0.2 * y * y + 0.1 * x * y + 0.2 * Math.sqrt(Math.abs(x));
    ret += (20.0 * Math.sin(6.0 * x * PI) + 20.0 * Math.sin(2.0 * x * PI)) * 2.0 / 3.0;
    ret += (20.0 * Math.sin(y * PI) + 40.0 * Math.sin(y / 3.0 * PI)) * 2.0 / 3.0;
    ret += (160.0 * Math.sin(y / 12.0 * PI) + 320 * Math.sin(y * PI / 30.0)) * 2.0 / 3.0;
    var dLat = ret;

    x = longitude - 105.0;
    y = latitude - 35.0;
    ret = 300.0 + x + 2.0 * y + 0.1 * x * x + 0.1 * x * y + 0.1 * Math.sqrt(Math.abs(x));
    ret += (20.0 * Math.sin(6.0 * x * PI) + 20.0 * Math.sin(2.0 * x * PI)) * 2.0 / 3.0;
    ret += (20.0 * Math.sin(x * PI) + 40.0 * Math.sin(x / 3.0 * PI)) * 2.0 / 3.0;
    ret += (150.0 * Math.sin(x / 12.0 * PI) + 300.0 * Math.sin(x / 30.0 * PI)) * 2.0 / 3.0;
    var dLon = ret;

    var radLat = latitude / 180.0 * PI;
    var magic = Math.sin(radLat);
    magic = 1 - ee * magic * magic;
    var sqrtMagic = Math.sqrt(magic);
    dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * PI);
    dLon = (dLon * 180.0) / (a / sqrtMagic * Math.cos(radLat) * PI);

    return {
        GPS: [
            {
                time: ti,
                latitude: latitude + dLat,
                longitude: longitude + dLon,
                altitude: altitude,
                speed: speed,
                azimuth: azimuth,
                RFU: RFU
            }
        ]
    }
}
var a = [170, 92, 7, 101, 123, 1, 88, 159, 132, 6, 202, 45, 17, 0, 0, 20, 0, 0, 2, 7, 0, 33, 170, 170, 170, 170, 40, 219];
console.log(Decode("1", a));