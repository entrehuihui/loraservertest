//[255, 2, 5, 29, 1, 4, 92, 0, 161, 166, 2, 2, 1, 70, 3, 1, 60, 4, 1, 52, 5, 3, 0, 255, 168, 255, 0]
function Decode(fPort, bytes) {
    str = "A";
    code = str.charCodeAt();
    str2 = String.fromCharCode(bytes[0]);
    return str2
}
// ff02051d01045c00ea070237014a014d0150014e014a01460148014901450146014b0150014f014b014f0151014f014c014bf101490147014b014c014f014e0150014e03183df13e3f4041f1424344f245f146f14748f249f24af24b4c041a32333537393a3bf13c3df13f3e40f142f344454446f14544f143050300ffa8ff00
function U(str1) {
    const str = str1
    let len = str.length
    //分包数
    let dataNum = parseInt(str.substr(4, 2), 16);
    //数据长度
    let dataLen = parseInt(str.substr(6, 2), 16);
    if (dataLen == 0) {
        return
    }
    let ti = 0;
    let num = parseInt(str.substr(8, 2), 16);
    let start = 10;
    let temp = new Array;
    let humi = new Array;
    let elec = new Array;
    for (; ;) {
        switch (num) {
            case 1:
                let l1 = parseInt(str.substr(start, 2), 16);
                if (l1 + start > len) {
                    return;
                }
                start += 2;
                ti = parseInt(str.substr(start, 8), 16);
                start += 8;
                num = parseInt(str.substr(start, 2), 16);
                start += 2;
                break;
            case 2:
                let l2 = parseInt(str.substr(start, 2), 16) * 2;
                if (l2 + start > len) {
                    return;
                }
                start += 2;
                let ii2 = start + l2;
                for (; start < ii2;) {
                    if (parseInt(str.substr(start, 1), 16) >= 10) {
                        start++;
                        if (temp.length == 0) {
                            return;
                        }
                        let temp1 = temp[temp.length - 1];
                        for (let i = 0; i < parseInt(str.substr(start, 1), 16); i++) {
                            temp.push(temp1);
                        }
                        start++;
                    } else {
                        temp.push(parseInt(str.substr(start, 4), 16) / 10);
                        start += 4;
                    }
                }
                num = parseInt(str.substr(start, 2), 16);
                start += 2;
                console.log("temp", temp.length, num, start);
                break;
            case 3:
                let l3 = parseInt(str.substr(start, 2), 16) * 2;
                if (l3 + start > len) {
                    return;
                }
                start += 2;
                let ii3 = start + l3;
                for (; start < ii3;) {
                    if (parseInt(str.substr(start, 1), 16) >= 10) {
                        start++;
                        if (humi.length == 0) {
                            return;
                        }
                        let humi1 = humi[humi.length - 1];
                        for (let i = 0; i < parseInt(str.substr(start, 1), 16); i++) {
                            humi.push(humi1);
                        }
                        start++;
                    } else {
                        humi.push(parseInt(str.substr(start, 2), 16));
                        start += 2;
                    }
                }
                num = parseInt(str.substr(start, 2), 16);
                start += 2
                console.log("humi", humi.length, num, start);
                break;
            case 4:
                let l4 = parseInt(str.substr(start, 2), 16) * 2;
                if (l4 + start > len) {
                    return;
                }
                start += 2;
                let ii4 = start + l4;
                for (; start < ii4;) {
                    if (parseInt(str.substr(start, 1), 16) >= 10) {
                        start++;
                        if (elec.length == 0) {
                            return;
                        }
                        let elec1 = elec[elec.length - 1];
                        for (let i = 0; i < parseInt(str.substr(start, 1), 16); i++) {
                            elec.push(elec1);
                        }
                        start++;
                    } else {
                        elec.push(parseInt(str.substr(start, 2), 16));
                        start += 2;
                    }
                }
                num = parseInt(str.substr(start, 2), 16);
                start += 2
                console.log("elec", elec.length, num, start);
                break;
            case 5:
                num = 6;
                break;
            default:
                break;
        }
        if (num > 5) {
            break;
        }
    }
    let alen = temp.length > humi.length ? humi.length : temp.length;
    alen = alen > elec.length ? elec.length : alen;
    let maxiiot = new Array;
    for (let i = 0; i < alen; i++) {
        maxiiot.push({
            Temperature: temp[i],
            Humidity: humi[i],
            Electricity: elec[i],
            time: String(ti)
        })
        ti += 1;
    }
    return maxiiot
}
// ff02051d01045c00a4ba0202014703013c040131050300ffa8ff00
let u = U("ff0205fa01043b9aca000213014afffffffffffffffffffffffffffffffff903123cfffffffffffffffffffffffffffffffff9041232fffffffffffffffffffffffffffffffff9050300ffa8ff00")
console.log(u.length)
// let arr = Decode("1", [102, 102, 48, 50, 48, 53, 49, 100, 48, 49, 48, 52, 53, 99, 48, 48, 97, 52, 98, 97, 48, 50, 48, 50, 48, 49, 52, 55, 48, 51, 48, 49, 51, 99, 48, 52, 48, 49, 51, 49, 48, 53, 48, 51, 48, 48, 102, 102, 97, 56, 102, 102, 48, 48])
// let arr = Decode("1", [255, 2, 5, 29, 1, 4, 92, 0, 161, 166, 2, 2, 1, 70, 3, 1, 60, 4, 1, 52, 5, 3, 0, 255, 168, 255, 0])
// console.log(arr, arr.maxiiot.length);
// let m = U("ff02055801045be3da8002050113ffffa8030330ffa804032fffa8050300ffa8ff00")
// console.log(m)
