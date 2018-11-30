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
                // console.log("temp", num, start);
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
                // console.log(num, humi);
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
                // console.log(num, elec);
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
let u = U("ff0205c801043b9aca00026f0149014b0149014701450147014c014b014d0152015601570158015401580157f1015b0157015a015cf2015801570158015d0162016401670168016df1016b016e016f016c017001740177017c017f018301840188018c018b0190018f0191018c018af1018bf201860183017df1017b017a017d01800182017d017701740171f1016e016d016b016c016b016a01650164015f015c015d0160015d0160015df1015c015e0160015c01570155f10152014e0150014e014a014b014c014a0148014501440146f1014301460143014501410142f1013d01390137013501330130012c01270123011e01180117f10112010ef10111f1010f0109010300fe00ff010200ff00fc00fef100fb00fa00f900fcf100f700f200f500f100f400f000f2f300f400f5f200f6f100f700f600f400f000f200f300ef00f000ebf100ed00e700eaf100e700e100e200e500e700e100e200e400e500e6f100e200dd00de00d900d3f100d100d200d500d000d100d000d303a83c3d3ef13f40f24142f14344f145f4464748494af14bf14cf14d4ef14ff650f2515253f354f255565758595af1595856555452504f4d4b494846444342413f3e3d3b3a38363432312f2d2b2a292826252322201f1ef21f202223252628292a2b2c2e2f3032343637393a3b3c3d3e3f404243444648494b4d4f5152535557585af159f458f1575655f35453f3525150f34f4ef64d4c4bf14af249f1484746f44544f143f14241f14004bc34363537f1393a3c3b3d3e404241f140f142f344f145474645464546f14849f14bf24df14f515052f25153545352f153525152515354535455565557f15658f15756555759f1585a595b5af15c5b5a5c5b5c5ef15f5e5df15cf25e605ff161f1636463f162646566f26867f16869f16af1696a6bf16d6c6e6d6e6ff171706f71f1737476f17574f175f174f17675f176787675f1737472706ef16df16c6d6ef16cf26a6869f16a686967f168666564f165f366f164f363f26462605f050300ffa8ff00")
console.log(u.length)
// let arr = Decode("1", [102, 102, 48, 50, 48, 53, 49, 100, 48, 49, 48, 52, 53, 99, 48, 48, 97, 52, 98, 97, 48, 50, 48, 50, 48, 49, 52, 55, 48, 51, 48, 49, 51, 99, 48, 52, 48, 49, 51, 49, 48, 53, 48, 51, 48, 48, 102, 102, 97, 56, 102, 102, 48, 48])
// let arr = Decode("1", [255, 2, 5, 29, 1, 4, 92, 0, 161, 166, 2, 2, 1, 70, 3, 1, 60, 4, 1, 52, 5, 3, 0, 255, 168, 255, 0])
// console.log(arr, arr.maxiiot.length);
// let m = U("ff02055801045be3da8002050113ffffa8030330ffa804032fffa8050300ffa8ff00")
// console.log(m)
