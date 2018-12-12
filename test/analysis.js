function Decode(fPort, bytes) {
    var str = bytes;
    var start = 0;
    var strlen = str.length;
    if (str.substr(start, 4) != "ff02" || strlen < 18) {
        return;
    }
    start += 4;
    var len = parseInt(str.substr(start, 2), 16);
    start += 2;
    var ti = parseInt(str.substr(start, 8), 16);
    start += 8;
    var templen1 = parseInt(str.substr(start, 2), 16) * 2 + 18;
    start += 2;
    var templen2 = parseInt(str.substr(start, 2), 16) * 2 + templen1;
    start += 2;
    if (strlen < templen2) {
        return;
    }
    var tempA = new Array;
    tempA.push(parseInt(str.substr(start, 2), 16));
    start += 2;
    for (; start < templen1; i++) {
        if (str.substr(start, 1) != "a") {
            tempA.push(parseInt(str.substr(start, 2), 16));
            start += 2;
            continue;
        }
        start++;
        var lastT = tempA[tempA.length - 1];
        for (var i = 0; i < parseInt(str.substr(start, 1), 16); i++) {
            tempA.push(lastT);
        }
        start++;
    }
    var tempB = new Array;
    tempB.push(parseInt(str.substr(start, 2), 16));
    start += 2;
    for (; start < templen2; i++) {
        if (str.substr(start, 1) != "a") {
            tempB.push(parseInt(str.substr(start, 2), 16));
            start += 2;
            continue;
        }
        start++;
        var lastT = tempB[tempB.length - 1];
        for (var i = 0; i < parseInt(str.substr(start, 1), 16); i++) {
            tempB.push(lastT);
        }
        start++;
    }
    if (tempA.length != len || tempB.length != len) {
        return;
    }
    templen2 += parseInt(str.substr(start, 2), 16) * 2 + 2;
    if (strlen < templen2) {
        return
    }
    start += 2;
    var humi = new Array;
    humi.push(parseInt(str.substr(start, 2), 16));
    start += 2;
    for (; start < templen2; i++) {
        if (str.substr(start, 1) != "a") {
            humi.push(parseInt(str.substr(start, 2), 16));
            start += 2;
            continue;
        }
        start++;
        var lastT = humi[humi.length - 1];
        for (var i = 0; i < parseInt(str.substr(start, 1), 16); i++) {
            humi.push(lastT);
        }
        start++;
    }
    if (humi.length != len) {

    }
    templen2 += parseInt(str.substr(start, 2), 16) * 2 + 2;
    if (strlen < templen2) {
        return
    }
    start += 2;
    var elec = new Array;
    elec.push(parseInt(str.substr(start, 2), 16));
    start += 2;
    for (; start < templen2; i++) {
        if (str.substr(start, 1) != "a") {
            elec.push(parseInt(str.substr(start, 2), 16));
            start += 2;
            continue;
        }
        start++;
        var lastT = elec[elec.length - 1];
        for (var i = 0; i < parseInt(str.substr(start, 1), 16); i++) {
            elec.push(lastT);
        }
        start++;
    }
    var maxiiot = new Array;
    for (var i = 0; i < len; i++) {
        maxiiot.push({
            Temperature: tempA[i] + tempB[i] / 10,
            Humidity: humi[i],
            Electricity: elec[i],
            time: String(ti)
        })
        ti += 60;
    }
    return { maxiiot: maxiiot };
}
//
var b = "ff02055c1061f3020414a404a20506042a292aa2023fa4ff";
// ff02050101045c106757020200d5030127040118050114ff

var a = Decode("1", b);
console.log(a);