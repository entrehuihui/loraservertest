// Code generated by "stringer -type=CID"; DO NOT EDIT.

package lorawan

import "strconv"

const (
	_CID_name_0 = "LinkCheckReqLinkADRReqDutyCycleReqRXParamSetupReqDevStatusReqNewChannelReqRXTimingSetupReqTXParamSetupReqDLChannelReq"
	_CID_name_1 = "DeviceTimeReq"
	_CID_name_2 = "PingSlotInfoReqPingSlotChannelReq"
	_CID_name_3 = "BeaconFreqReq"
)

var (
	_CID_index_0 = [...]uint8{0, 12, 22, 34, 49, 61, 74, 90, 105, 117}
	_CID_index_2 = [...]uint8{0, 15, 33}
)

func (i CID) String() string {
	switch {
	case 2 <= i && i <= 10:
		i -= 2
		return _CID_name_0[_CID_index_0[i]:_CID_index_0[i+1]]
	case i == 13:
		return _CID_name_1
	case 16 <= i && i <= 17:
		i -= 16
		return _CID_name_2[_CID_index_2[i]:_CID_index_2[i+1]]
	case i == 19:
		return _CID_name_3
	default:
		return "CID(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
