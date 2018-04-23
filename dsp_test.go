package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"jvole.com/dsp/index"
	"jvole.com/dsp/util"
)

func aTestDspWin(t *testing.T) {
	url := "http://localhost:9900/v1/ADXNotify?oid=oid-adfawer&price=1.23&cid=32&uid=14&postion=postion-adsf&t=1241341234&user=user-adfwerqwer&device=device-asdfaerqwer"

	res, err := util.DoHTTPGet(url)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("DSPwin: ", string(res))
}

func aTestDspB(t *testing.T) {
	url := "http://localhost:9900/v1/Bidder"
	var wg sync.WaitGroup
	count := 100
	t2 := time.Now() // get current time
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			util.DoBytesPost(url, adxcontext)
			wg.Done()
		}()
	}

	wg.Wait()
	elapsed2 := time.Since(t2)
	fmt.Println("DSP: ", elapsed2)
}

func aTestDsp(t *testing.T) {
	url := "http://localhost:9900/v1/Bidder"

	t2 := time.Now() // get current time

	resu, err := util.DoBytesPost(url, adxcontext)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("====", string(resu))

	elapsed2 := time.Since(t2)
	fmt.Println("DSP: ", elapsed2)
}

func aTestSyncIndex(t *testing.T) {
	url := "http://localhost:9900/v1/SyncIndex"

	t2 := time.Now() // get current time
	resu, err := util.DoBytesPost(url, []byte(`{"server":0}`))
	if err != nil {
		log.Println(err)
	}
	data := util.DoZlibUnCompress(resu)
	index.CPINDEX.SetIndexFromByte(data)
	// fmt.Println("====", string(resu))
	fmt.Println("====", index.CPINDEX.GetCPLen())

	elapsed2 := time.Since(t2)
	fmt.Println("DSP: ", elapsed2)
}

var adxcontext = []byte(`
{
  "id": "DxU0032U8a",
  "at": 2,
  "allimps": 0,
  "imp": [
    {
      "id": "1",
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
      "country": "CHN",
      "region": "04",
      "zip": "20099",
      "metro": "0",
      "city": "Zaoyang",
      "type": 2
    },
    "make": "Apple",
    "model": "iPhone",
    "os": "IOS",
    "osv": "4.3.2",
    "ua": "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5",
    "ip": "95.90.255.47",
    "js": 0,
    "connectiontype": 0,
    "ifa": "e4273e31-97a9-4b29-93a8-8a99f0cea068",
    "devicetype": 1
  },
  "app": {
    "id": "101000415",
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
    "udi": {},
    "operaminibrowser": 0,
    "carriername": "China Mobile"
  },
  "regs": {
    "coppa": 0
  }
}
  `)
var adxdata = [][]byte{[]byte(`
{
    "allimps": 0,
    "app": {
        "bundle": "com.grindrapp.android",
        "cat": [
            "IAB14"
        ],
        "id": "120000491",
        "name": "Grindr_Android_Grindr_A_320_T3_SM_ALL_50_Android_XXLARGE_320x50_IAB14",
        "publisher": {
            "id": "1100000436",
            "name": "Grindr"
        },
        "storeurl": "https://play.google.com/store/apps/details?id=com.grindrapp.android&hl=en"
    },
    "at": 2,
    "bcat": [
        "IAB12",
        "IAB23",
        "IAB7-28",
        "IAB26",
        "IAB25",
        "IAB24",
        "IAB17-18",
        "IAB14-8",
        "IAB7-42",
        "IAB11-5",
        "IAB14-2",
        "IAB14-1",
        "IAB11-4",
        "IAB14-7",
        "IAB14-6",
        "IAB14-5",
        "IAB14-4",
        "IAB9-9"
    ],
    "device": {
        "carrier": "unknown - probably WLAN",
        "connectiontype": 6,
        "devicetype": 1,
        "geo": {
            "city": "Quanzhou",
            "country": "TWN",
            "ipservice": 3,
            "lat": 23.172384,
            "lon": 120.248032,
            "metro": "0",
            "region": "07",
            "type": 1,
            "zip": "721"
        },
        "ifa": "5087e76e-8173-4047-ad53-c52e37e8ab14",
        "ip": "115.82.241.172",
        "js": 1,
        "lmt": 0,
        "make": "Samsung",
        "model": "SM-J510UN",
        "os": "Android",
        "osv": "6.0",
        "ua": "Mozilla/5.0 (Linux; Android 6.0.1; SM-J510UN Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/58.0.3029.83 Mobile Safari/537.36"
    },
    "ext": {
        "carriername": "unknown - probably WLAN",
        "operaminibrowser": 0,
        "udi": {
            "googleadid": "5087e76e-8173-4047-ad53-c52e37e8ab14",
            "googlednt": 0,
            "idfatracking": 1
        }
    },
    "id": "77a725c6-ae98-46d4-812a-98e02e60ef2a",
    "imp": [
        {
            "banner": {
                "api": [
                    3,
                    5
                ],
                "battr": [
                    1,
                    2,
                    3,
                    5,
                    6,
                    8,
                    9,
                    10,
                    14
                ],
                "btype": [
                    1
                ],
                "format": [
                    {
                        "h": 50,
                        "w": 320
                    }
                ],
                "h": 50,
                "mimes": [
                    "text/javascript",
                    "application/javascript",
                    "image/jpeg",
                    "image/png",
                    "image/gif"
                ],
                "pos": 1,
                "w": 320
            },
            "displaymanager": "SOMA",
            "displaymanagerver": "sdkandroid_7-2-0",
            "ext": {
                "strictbannersize": 0
            },
            "id": "1",
            "instl": 0,
            "secure": 0,
            "tagid": "130045925"
        }
    ],
    "regs": {
        "coppa": 0
    },
    "user": {
        "gender": "M",
        "yob": 2000
    }
}
	`)}
