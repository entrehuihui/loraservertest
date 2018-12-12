//WGS-84 to GCJ-02
function gcj_encrypt(wgsLat, wgsLon) {
    if (wgsLon < 72.004 || wgsLon > 137.8347 || wgsLat < 0.8293 || wgsLat > 55.8271)
        return { 'lat': wgsLat, 'lon': wgsLon };
    var PI = 3.14159265358979324;
    // var lat = wgsLat;
    // var lon = wgsLon;
    var a = 6378245.0; //  a: 卫星椭球坐标投影到平面地图坐标系的投影因子。
    var ee = 0.00669342162296594323; //  ee: 椭球的偏心率。

    var x = wgsLon - 105.0;
    var y = wgsLat - 35.0;
    var ret = -100.0 + 2.0 * x + 3.0 * y + 0.2 * y * y + 0.1 * x * y + 0.2 * Math.sqrt(Math.abs(x));
    ret += (20.0 * Math.sin(6.0 * x * PI) + 20.0 * Math.sin(2.0 * x * PI)) * 2.0 / 3.0;
    ret += (20.0 * Math.sin(y * PI) + 40.0 * Math.sin(y / 3.0 * PI)) * 2.0 / 3.0;
    ret += (160.0 * Math.sin(y / 12.0 * PI) + 320 * Math.sin(y * PI / 30.0)) * 2.0 / 3.0;
    var dLat = ret;

    x = wgsLon - 105.0;
    y = wgsLat - 35.0;
    ret = 300.0 + x + 2.0 * y + 0.1 * x * x + 0.1 * x * y + 0.1 * Math.sqrt(Math.abs(x));
    ret += (20.0 * Math.sin(6.0 * x * PI) + 20.0 * Math.sin(2.0 * x * PI)) * 2.0 / 3.0;
    ret += (20.0 * Math.sin(x * PI) + 40.0 * Math.sin(x / 3.0 * PI)) * 2.0 / 3.0;
    ret += (150.0 * Math.sin(x / 12.0 * PI) + 300.0 * Math.sin(x / 30.0 * PI)) * 2.0 / 3.0;
    var dLon = ret;

    var radLat = wgsLat / 180.0 * PI;
    var magic = Math.sin(radLat);
    magic = 1 - ee * magic * magic;
    var sqrtMagic = Math.sqrt(magic);
    dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * PI);
    dLon = (dLon * 180.0) / (a / sqrtMagic * Math.cos(radLat) * PI);
    return { 'lat': wgsLat + dLat, 'lon': wgsLon + dLon };
}
var a = gcj_encrypt(22.585152, 113.913083);
console.log(a);
// { lat: 22.582105469397508, lon: 113.91794594472783 }
// var x_pi = 3.14159265358979324 * 3000.0 / 180.0;
// function delta(lat, lon) {
//     var a = 6378245.0; //  a: 卫星椭球坐标投影到平面地图坐标系的投影因子。
//     var ee = 0.00669342162296594323; //  ee: 椭球的偏心率。
//     var dLat = transformLat(lon - 105.0, lat - 35.0);
//     var dLon = transformLon(lon - 105.0, lat - 35.0);
//     var radLat = lat / 180.0 * PI;
//     var magic = Math.sin(radLat);
//     magic = 1 - ee * magic * magic;
//     var sqrtMagic = Math.sqrt(magic);
//     dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * PI);
//     dLon = (dLon * 180.0) / (a / sqrtMagic * Math.cos(radLat) * PI);
//     return { 'lat': dLat, 'lon': dLon };
// }
// function outOfChina(lat, lon) {
//     if (lon < 72.004 || lon > 137.8347)
//         return true;
//     if (lat < 0.8293 || lat > 55.8271)
//         return true;
//     return false;
// }
// function transformLat(x, y) {
//     var ret = -100.0 + 2.0 * x + 3.0 * y + 0.2 * y * y + 0.1 * x * y + 0.2 * Math.sqrt(Math.abs(x));
//     ret += (20.0 * Math.sin(6.0 * x * PI) + 20.0 * Math.sin(2.0 * x * PI)) * 2.0 / 3.0;
//     ret += (20.0 * Math.sin(y * PI) + 40.0 * Math.sin(y / 3.0 * PI)) * 2.0 / 3.0;
//     ret += (160.0 * Math.sin(y / 12.0 * PI) + 320 * Math.sin(y * PI / 30.0)) * 2.0 / 3.0;
//     return ret;
// }
// function transformLon(x, y) {
//     var ret = 300.0 + x + 2.0 * y + 0.1 * x * x + 0.1 * x * y + 0.1 * Math.sqrt(Math.abs(x));
//     ret += (20.0 * Math.sin(6.0 * x * PI) + 20.0 * Math.sin(2.0 * x * PI)) * 2.0 / 3.0;
//     ret += (20.0 * Math.sin(x * PI) + 40.0 * Math.sin(x / 3.0 * PI)) * 2.0 / 3.0;
//     ret += (150.0 * Math.sin(x / 12.0 * PI) + 300.0 * Math.sin(x / 30.0 * PI)) * 2.0 / 3.0;
//     return ret;
// }
// //GCJ-02 to WGS-84
// function gcj_decrypt(gcjLat, gcjLon) {
//     if (this.outOfChina(gcjLat, gcjLon))
//         return { 'lat': gcjLat, 'lon': gcjLon };

//     var d =delta(gcjLat, gcjLon);
//     return { 'lat': gcjLat - d.lat, 'lon': gcjLon - d.lon };
// }
// //GCJ-02 to WGS-84 exactly
// function gcj_decrypt_exact(gcjLat, gcjLon) {
//     var initDelta = 0.01;
//     var threshold = 0.000000001;
//     var dLat = initDelta, dLon = initDelta;
//     var mLat = gcjLat - dLat, mLon = gcjLon - dLon;
//     var pLat = gcjLat + dLat, pLon = gcjLon + dLon;
//     var wgsLat, wgsLon, i = 0;
//     while (1) {
//         wgsLat = (mLat + pLat) / 2;
//         wgsLon = (mLon + pLon) / 2;
//         var tmp =gcj_encrypt(wgsLat, wgsLon)
//         dLat = tmp.lat - gcjLat;
//         dLon = tmp.lon - gcjLon;
//         if ((Math.abs(dLat) < threshold) && (Math.abs(dLon) < threshold))
//             break;

//         if (dLat > 0) pLat = wgsLat; else mLat = wgsLat;
//         if (dLon > 0) pLon = wgsLon; else mLon = wgsLon;

//         if (++i > 10000) break;
//     }
//     //console.log(i);
//     return { 'lat': wgsLat, 'lon': wgsLon };
// }
// //GCJ-02 to BD-09
// function bd_encrypt(gcjLat, gcjLon) {
//     var x = gcjLon, y = gcjLat;
//     var z = Math.sqrt(x * x + y * y) + 0.00002 * Math.sin(y *x_pi);
//     var theta = Math.atan2(y, x) + 0.000003 * Math.cos(x *x_pi);
//     bdLon = z * Math.cos(theta) + 0.0065;
//     bdLat = z * Math.sin(theta) + 0.006;
//     return { 'lat': bdLat, 'lon': bdLon };
// }
// //BD-09 to GCJ-02
// function bd_decrypt(bdLat, bdLon) {
//     var x = bdLon - 0.0065, y = bdLat - 0.006;
//     var z = Math.sqrt(x * x + y * y) - 0.00002 * Math.sin(y *x_pi);
//     var theta = Math.atan2(y, x) - 0.000003 * Math.cos(x *x_pi);
//     var gcjLon = z * Math.cos(theta);
//     var gcjLat = z * Math.sin(theta);
//     return { 'lat': gcjLat, 'lon': gcjLon };
// }
// //WGS-84 to Web mercator
// //mercatorLat -> y mercatorLon -> x
// function mercator_encrypt(wgsLat, wgsLon) {
//     var x = wgsLon * 20037508.34 / 180.;
//     var y = Math.log(Math.tan((90. + wgsLat) *PI / 360.)) / (this.PI / 180.);
//     y = y * 20037508.34 / 180.;
//     return { 'lat': y, 'lon': x };
// }
// // Web mercator to WGS-84
// // mercatorLat -> y mercatorLon -> x
// function mercator_decrypt(mercatorLat, mercatorLon) {
//     var x = mercatorLon / 20037508.34 * 180.;
//     var y = mercatorLat / 20037508.34 * 180.;
//     y = 180 /PI * (2 * Math.atan(Math.exp(y *PI / 180.)) -PI / 2);
//     return { 'lat': y, 'lon': x };
// }
// // two point's distance
// function distance(latA, lonA, latB, lonB) {
//     var earthR = 6371000.;
//     var x = Math.cos(latA *PI / 180.) * Math.cos(latB *PI / 180.) * Math.cos((lonA - lonB) *PI / 180);
//     var y = Math.sin(latA *PI / 180.) * Math.sin(latB *PI / 180.);
//     var s = x + y;
//     if (s > 1) s = 1;
//     if (s < -1) s = -1;
//     var alpha = Math.acos(s);
//     var distance = alpha * earthR;
//     return distance;
// }