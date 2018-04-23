package model

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func aTestGetBidderRequest(tt *testing.T) {
	smaato := NewSmaatoADX(nil, 1, "demo.com")
	smaato.SetData(data2[0])
	smaato.SetOfferInfo()
	cp := *&Compaign{}
	cp = cp.GetComByID(81)
	t1 := time.Now()
	s := smaato.GetBidderResponse(cp)
	es := time.Since(t1)
	fmt.Println("---------", es, s)

}

func TestSmaato_GetUA(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5"},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetUA(); got != tt.want {
			t.Errorf("%q. Smaato.GetPostion() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetOffer(t *testing.T) {
	smaato := NewSmaatoADX(nil, 1, "demo.com")
	smaato.SetData(data2[0])
	smaato.SetOfferInfo()
	// pp.Println(smaato.GetOfferInfo())
}

func TestSmaato_GetImages(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	// smaato := NewSmaatoADX(nil, 1)
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, []string{"320x50"}},
		{"smaato2", *&fields{"url2", 1, data2[1]}, []string{"320x50"}},
		{"smaato3", *&fields{"url3", 1, data2[3]}, []string{"100x100"}},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")

		if got := st.GetImages(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Smaato.GetImages() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetOfferID(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	// smaato := NewSmaatoADX(nil, 1)
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, "DxU0032U8a"},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "DxU0032U8a"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "DxU0032U8a1"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetOfferID(); got != tt.want {
			t.Errorf("%q. Smaato.GetOfferID() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetType(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, 2},
		{"smaato2", *&fields{"url2", 1, data2[1]}, 3},
		{"smaato3", *&fields{"url3", 1, data2[2]}, 1},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetType(); got != tt.want {
			t.Errorf("%q. Smaato.GetType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetBidFloor(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, 1300000},
		{"smaato2", *&fields{"url2", 1, data2[1]}, 0},
		{"smaato3", *&fields{"url3", 1, data2[2]}, 300000},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetBidFloor(); got != tt.want {
			t.Errorf("%q. Smaato.GetBidFloor() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetPostion(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, "101000dd415"},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "101008--4563"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "1010084563"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetPostion(); got != tt.want {
			t.Errorf("%q. Smaato.GetPostion() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetDeviceID(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, ""},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "e4273e31-97a9-4b29-93a8-8a99f0cea068"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "e4273e31-97a9-4b29-93a8-8a99f0cea068"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetDeviceID(); got != tt.want {
			t.Errorf("%q. Smaato.GetDeviceID() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetUserID(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, ""},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "adsfaerwrq"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, ""},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetUserID(); got != tt.want {
			t.Errorf("%q. Smaato.GetUserID() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetCountry(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, "AF"},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "AE"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "DEU"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetCountry(); got != tt.want {
			t.Errorf("%q. Smaato.GetCountry() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetCity(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, "Hamburg"},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "Hamburg2"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "Hamburg3"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetCity(); got != tt.want {
			t.Errorf("%q. Smaato.GetCity() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetRegion(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, "01"},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "05"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "04"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetRegion(); got != tt.want {
			t.Errorf("%q. Smaato.GetRegion() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetADX(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, 1},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetADX(); got != tt.want {
			t.Errorf("%q. Smaato.GetADX() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetSourceType(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, 2},
		{"smaato2", *&fields{"url2", 1, data2[1]}, 2},
		{"smaato3", *&fields{"url3", 1, data2[2]}, 1},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetSourceType(); got != tt.want {
			t.Errorf("%q. Smaato.GetSourceType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetContentCat(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, []string{"IAB13", "BAC"}},
		{"smaato2", *&fields{"url2", 1, data2[1]}, []string{"DSW"}},
		{"smaato3", *&fields{"url3", 1, data2[2]}, []string{"IAB1"}},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetContentCat(); fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
			t.Errorf("%q. Smaato.GetContentCat() = %v, want %v", tt.name, got, tt.want)
		}

	}
}

func TestSmaato_GetConnType(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, 0},
		{"smaato2", *&fields{"url2", 1, data2[1]}, 1},
		{"smaato3", *&fields{"url3", 1, data2[2]}, 5},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetConnType(); got != tt.want {
			t.Errorf("%q. Smaato.GetConnType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetCarriers(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, "unknown - probably WLAN"},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "unknown - probably WLAN2"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "unknown - probably WLAN4"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetCarriers(); got != tt.want {
			t.Errorf("%q. Smaato.GetCarriers() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetDeviceType(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, 1},
		{"smaato2", *&fields{"url2", 1, data2[1]}, 1},
		{"smaato3", *&fields{"url3", 1, data2[2]}, 1},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetDeviceType(); got != tt.want {
			t.Errorf("%q. Smaato.GetDeviceType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetOS(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, "ios 4.3.2"},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "ios2 4.3.2"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "ios3 4.3.2"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetOS(); got != tt.want {
			t.Errorf("%q. Smaato.GetOS() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetIdfaGaid(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, ""},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "e4273e31-97a9-4b29-93a8-8a99f0cea068"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "e4273e31-97a9-4b29-93a8-8a99f0cea068"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetIdfaGaid(); got != tt.want {
			t.Errorf("%q. Smaato.GetIdfaGaid() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSmaato_GetIP(t *testing.T) {
	type fields struct {
		CallbackURL string
		Code        uint32
		Data        []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"smaato", *&fields{"url", 1, data2[0]}, "192.90.255.47"},
		{"smaato2", *&fields{"url2", 1, data2[1]}, "95.90.255.47"},
		{"smaato3", *&fields{"url3", 1, data2[2]}, "203.12.34.147"},
	}
	for _, tt := range tests {
		st := NewSmaatoADX(tt.fields.Data, tt.fields.Code, "demo.com")
		if got := st.GetIP(); got != tt.want {
			t.Errorf("%q. Smaato.GetIP() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

var data2 = [][]byte{[]byte(`{
    "id": "DxU0032U8a",
    "at": 2,
    "allimps": 0,
    "imp": [
        {
            "id": "1",
						"bidfloor":1.3,
            "banner": {
                "w": 320,
                "h": 50,
                "format": [
                    {
                        "w": 320,
                        "h": 50
                    }
                ],
                "btype": [
                    1,
                    3
                ],
                "battr": [
                    1,
                    3,
                    5,
                    6,
                    8,
                    9,
                    10,
                    11
                ],
                "pos": 3,
                "mimes": [
                    "image/jpeg",
                    "image/png",
                    "image/gif"
                ],
                "api": []
            },
            "ext": {
                "strictbannersize": 0
            },
            "instl": 0,
            "displaymanager": "SOMA",
            "tagid": "101000415",
            "secure": 0
        }
    ],
    "device": {
        "geo": {
            "lat": 53.550003,
            "lon": 10,
            "ipservice": 3,
            "country": "AF",
            "region": "01",
            "zip": "20099",
            "metro": "0",
            "city": "Hamburg",
            "type": 2
        },
        "make": "Apple",
        "model": "iPhone",
        "os": "iOS",
        "osv": "4.3.2",
        "ua": "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5",
        "ip": "192.90.255.47",
        "js": 0,
        "connectiontype": 0,
        "devicetype": 1
    },
    "app": {
        "id": "101000dd415",
        "name": "OpenRTB_2_4_UATest_openRtb_2_4_iOS_XXLARGE_320x50_IAB1",
        "domain": "example.com",
        "cat": [
            "IAB13",
						"BAC"
        ],
        "storeurl": "http://example.com",
        "keywords": "",
        "publisher": {
            "id": "1001028764",
            "name": "OpenRTB_2_4_UATest_openRtb_2_4_iOS_XXLARGE_320x50_IAB1"
        }
    },
    "user": {
        "keywords": ""
    },
    "bcat": [
        "IAB17-18",
        "IAB7-42",
        "IAB23",
        "IAB7-28",
        "IAB26",
        "IAB25",
        "IAB9-9",
        "IAB24"
    ],
    "badv": [],
    "ext": {
        "udi": {
        },
        "operaminibrowser": 0,
        "carriername": "unknown - probably WLAN"
    },
    "regs": {
        "coppa": 0
    }
}`),
	[]byte(`{
    "id": "DxU0032U8a",
    "at": 2,
    "allimps": 0,
    "imp": [
        {
            "id": "1",
            "popup": {
                "w": 320,
                "h": 50,
                "format": [
                    {
                        "w": 320,
                        "h": 50
                    }
                ],
                "btype": [
                    1,
                    3
                ],
                "battr": [
                    1,
                    3,
                    5,
                    6,
                    8,
                    9,
                    10,
                    11
                ],
                "pos": 3,
                "mimes": [
                    "image/jpeg",
                    "image/png",
                    "image/gif"
                ],
                "api": []
            },
            "ext": {
                "strictbannersize": 0
            },
            "instl": 0,
            "displaymanager": "SOMA",
            "tagid": "101000415",
            "secure": 0
        }
    ],
    "device": {
        "geo": {
            "lat": 53.550003,
            "lon": 10,
            "ipservice": 3,
            "country": "AE",
            "region": "05",
            "zip": "20099",
            "metro": "0",
            "city": "Hamburg2",
            "type": 2
        },
        "make": "Apple",
        "model": "iPhone",
        "os": "iOS2",
        "osv": "4.3.2",
        "ua": "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5",
        "ip": "95.90.255.47",
        "js": 0,
        "connectiontype": 1,
        "ifa": "e4273e31-97a9-4b29-93a8-8a99f0cea068",
        "devicetype": 1
    },
    "app": {
        "id": "101008--4563",
        "name": "OpenRTB_2_4_UATest_openRtb_2_4_iOS_XXLARGE_320x50_IAB1",
        "domain": "example.com",
        "cat": [
            "DSW"
        ],
        "storeurl": "http://example.com",
        "keywords": "",
        "publisher": {
            "id": "1001028764",
            "name": "OpenRTB_2_4_UATest_openRtb_2_4_iOS_XXLARGE_320x50_IAB1"
        }
    },
    "user": {
			"id":"adsfaerwrq",
        "keywords": ""
    },
    "bcat": [
        "IAB17-18",
        "IAB7-42",
        "IAB23",
        "IAB7-28",
        "IAB26",
        "IAB25",
        "IAB9-9",
        "IAB24"
    ],
    "badv": [],
    "ext": {
        "udi": {
        },
        "operaminibrowser": 0,
        "carriername": "unknown - probably WLAN2"
    },
    "regs": {
        "coppa": 0
    }
}`),
	[]byte(`{
    "id": "DxU0032U8a1",
    "at": 2,
    "allimps": 0,
    "imp": [
        {
            "id": "1",
						"bidfloor":0.3,
            "native": {
                "w": 320,
                "h": 50,
                "format": [
                    {
                        "w": 320,
                        "h": 50
                    }
                ],
                "request": [
                    1,
                    3
                ],
                "battr": [
                    1,
                    3,
                    5,
                    6,
                    8,
                    9,
                    10,
                    11
                ],
                "pos": 3,
                "mimes": [
                    "image/jpeg",
                    "image/png",
                    "image/gif"
                ],
                "api": []
            },
            "ext": {
                "strictbannersize": 0
            },
            "instl": 0,
            "displaymanager": "SOMA",
            "tagid": "101000415",
            "secure": 0
        }
    ],
    "device": {
        "geo": {
            "lat": 53.550003,
            "lon": 10,
            "ipservice": 3,
            "country": "DEU",
            "region": "04",
            "zip": "20099",
            "metro": "0",
            "city": "Hamburg3",
            "type": 2
        },
        "make": "Apple",
        "model": "iPhone",
        "os": "iOS3",
        "osv": "4.3.2",
        "ua": "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5",
        "ip": "203.12.34.147",
        "js": 0,
        "connectiontype": 5,
        "ifa": "e4273e31-97a9-4b29-93a8-8a99f0cea068",
        "devicetype": 1
    },
    "site": {
        "id": "1010084563",
        "name": "OpenRTB_2_4_UATest_openRtb_2_4_iOS_XXLARGE_320x50_IAB1",
        "domain": "example.com",
        "cat": [
            "IAB1"
        ],
        "storeurl": "http://example.com",
        "keywords": "",
        "publisher": {
            "id": "1001028764",
            "name": "OpenRTB_2_4_UATest_openRtb_2_4_iOS_XXLARGE_320x50_IAB1"
        }
    },
    "user": {
        "keywords": ""
    },
    "bcat": [
        "IAB17-18",
        "IAB7-42",
        "IAB23",
        "IAB7-28",
        "IAB26",
        "IAB25",
        "IAB9-9",
        "IAB24"
    ],
    "badv": [],
    "ext": {
        "udi": {
        },
        "operaminibrowser": 0,
        "carriername": "unknown - probably WLAN4"
    },
    "regs": {
        "coppa": 0
    }
}`),
	[]byte(`
{
  "id": "2779d50acaaa9de931d87e19b1af629b3baf6f3f",
  "imp": [
    {
      "id": "90180978",
      "native": {
        "request": "{\"native\":{\"ver\":\"1.2\",\"context\":1,\"plcmttype\":4,\"plcmtcnt\":5,\"assets\":[{\"id\":1,\"required\":1,\"img\":{\"type\":3,\"wmin\":100,\"hmin\":100}},{\"id\":1,\"title\":{\"len\":60}},{\"id\":1,\"data\":{\"type\":2,\"len\":110}}],\"privacy\":0}}",
        "ver": 1.2
      }
    }
  ],
  "site": {
    "id": "12345",
    "domain": "sitedomain.com",
    "cat": [
      "IAB25-3"
    ],
    "page": "https://sitedomain.com/page",
    "ext": {
      "exchangecat": 508,
      "idzone": 112233
    }
  },
  "device": {
    "ua": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.63 Safari/537.36",
    "ip": "131.34.123.159",
    "geo": {
      "country": "IRL"
    },
    "language": "en",
    "os": "Linux & UNIX",
    "js": 0,
    "ext": {
      "remote_addr": "131.34.123.159",
      "x_forwarded_for": "",
      "accept_language": "en-GB;q=0.8,pt-PT;q=0.6,en;q=0.4,en-US;q=0.2,de;q=0.2,es;q=0.2,fr;q=0.2"
    }
  },
  "user": {
    "id": "57592f333f8983.043587162282415065"
  }
}
  `),
}
