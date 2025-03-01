package scanning

import (
    "encoding/json"
    "testing"
    "github.com/stretchr/testify/assert"
)

// Mock data for V1 and V2 responses
var v1Response = `{"response_bytes_utf8": [65, 66, 67]}`
var v2Response = `{"response_str": "Success"}`

func TestParseData(t *testing.T) {
    scan := Scan{
        Ip:          "192.168.1.1",
        Port:        8080,
        Service:     "http",
        Timestamp:   1617181723,
        DataVersion: V1,
        Data: map[string]interface{}{
            "response_bytes_utf8": []byte{97, 98, 100, 117, 108},
        },
    }

    dataBytes := scan.parseData()
    var response V1Data
    err := json.Unmarshal(dataBytes, &response)
    assert.NoError(t, err, "unmarshaling data should not return an error")
    assert.Equal(t, "abdul", string(response.ResponseBytesUtf8), "response should match")
}

func TestParseV1Data(t *testing.T) {
    scan := Scan{
        Ip:          "192.168.1.1",
        Port:        8080,
        Service:     "http",
        Timestamp:   1617181723,
        DataVersion: V1,
        Data: map[string]interface{}{
            "response_bytes_utf8": []byte{65, 66, 67},
        },
    }

    result := scan.ParseV1Data()
    expected := "ABC" // Expected output for the byte slice [65, 66, 67]
    assert.Equal(t, expected, result, "Parsed V1 data should match expected result")
}

func TestParseV2Data(t *testing.T) {
    scan := Scan{
        Ip:          "192.168.1.1",
        Port:        8080,
        Service:     "http",
        Timestamp:   1617181723,
        DataVersion: V2,
        Data: map[string]interface{}{
            "response_str": "Success",
        },
    }

    result := scan.ParseV2Data()
    expected := "Success" // Expected output for response_str
    assert.Equal(t, expected, result, "Parsed V2 data should match expected result")
}

func TestParseV1Data_InvalidJSON(t *testing.T) {
    scan := Scan{
        Ip:          "192.168.1.1",
        Port:        8080,
        Service:     "http",
        Timestamp:   1617181723,
        DataVersion: V1,
        Data: map[string]interface{}{
            "response_bytes_utf8": "invalid_bytes",
        },
    }

    // Expecting a panic due to unmarshaling error (invalid type)
    assert.Panics(t, func() {
        scan.ParseV1Data()
    }, "Parsing invalid V1 data should panic")
}

func TestParseV2Data_InvalidJSON(t *testing.T) {
    scan := Scan{
        Ip:          "192.168.1.1",
        Port:        8080,
        Service:     "http",
        Timestamp:   1617181723,
        DataVersion: V2,
        Data: map[string]interface{}{
            "response_str": 123, // Invalid type, should be string
        },
    }

    // Expecting a panic due to unmarshaling error (invalid type)
    assert.Panics(t, func() {
        scan.ParseV2Data()
    }, "Parsing invalid V2 data should panic")
}


