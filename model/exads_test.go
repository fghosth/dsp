package model

import (
	"reflect"
	"testing"
)

func TestExads_GetImages(t *testing.T) {
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
		{"exads native", *&fields{"url", 2, dataexads[0]}, []string{"100x100"}},
		{"exads2 banner", *&fields{"url2", 2, dataexads[1]}, []string{"300x250"}},
		// {"smaato3", *&fields{"url3", 1, data2[2]}, []string{"320x50"}},
	}
	for _, tt := range tests {
		st := NewExadsADX(tt.fields.Data, tt.fields.Code, "demo.com")
		// st.SetData()

		if got := st.GetImages(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Smaato.GetImages() = %v, want %v,%d", tt.name, got, tt.want, st.GetType())
		}
	}
}

var dataexads = [][]byte{[]byte(`
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
	[]byte(`
{
    "id": "b91cdbbe1d373eb898361216b017ad7a1d01ed62",
    "imp": [
        {
            "id": "6145298465",
            "instl": 0,
            "banner": {
                "w": 300,
                "h": 250
            },
            "bidfloor": 0.5,
            "bidfloorcur": "USD"
        }
    ],
    "site": {
        "id": "12345",
        "domain": "sitedomain.com",
        "cat": ["IAB25-3"],
        "page": "https://sitedomain.com/page",
        "ext": {
            "exchangecat": 508,
            "idzone": 445566
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
