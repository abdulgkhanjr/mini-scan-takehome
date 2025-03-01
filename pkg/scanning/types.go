package scanning

import (
	"encoding/json"
)

const (
	Version = iota
	V1
	V2
)

type Scan struct {
	Ip          string      `json:"ip"`
	Port        uint32      `json:"port"`
	Service     string      `json:"service"`
	Timestamp   int64       `json:"timestamp"`
	DataVersion int         `json:"data_version"`
	Data        interface{} `json:"data"`
}

type V1Data struct {
	ResponseBytesUtf8 []byte `json:"response_bytes_utf8"`
}

type V2Data struct {
	ResponseStr string `json:"response_str"`
}


// parse out any scan's service data
func (scan Scan) parseData() ([]byte) {
    dataMap := scan.Data.(map[string]interface{})
    dataBytes, err := json.Marshal(dataMap)
    if err != nil {
        panic(err)
    }
    return dataBytes
}

func (scan Scan) ParseV1Data() string {
	var respData V1Data
    dataBytes := scan.parseData()
    err := json.Unmarshal(dataBytes, &respData)
    if err != nil {
        panic(err)
    }
    return string(respData.ResponseBytesUtf8)
}

func (scan Scan) ParseV2Data() string {
	var respData V2Data
    dataBytes := scan.parseData()
    err := json.Unmarshal(dataBytes, &respData)
    if err != nil {
        panic(err)
    }
    return string(respData.ResponseStr)
}


