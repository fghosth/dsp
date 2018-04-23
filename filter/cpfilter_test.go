package filter

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"jvole.com/dsp/config"
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"
)

func aTestCPFiltercomp(t *testing.T) {

	data := []byte(`{"cmd":"select * from BaseDspCampaign where id=87  limit 1"}`)
	res, err := util.DoBytesPost(config.CompaignURL, data)
	if err != nil {
		log.Printf("request err:%s", err)
	}
	mcom := &model.Compaign{}
	compaign := mcom.GetData(res)

	smaato := model.NewSmaatoADX(nil, 1, "domain.com")
	smaato.SetData(data2[0])
	smaato.SetOfferInfo()
	// filter := new(DSPFilter)
	t1 := time.Now() // get current time
	rest := CPFiltercomp(compaign[0], smaato)
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)
	fmt.Println(compaign[0].FilterSet, rest)
	count := 5000
	var wg sync.WaitGroup
	// runtime.GOMAXPROCS(4)
	t2 := time.Now() // get current time
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			CPFiltercomp(compaign[0], smaato)
			wg.Add(-1)
		}()
	}
	wg.Wait()
	elapsed2 := time.Since(t2)
	fmt.Println("App elapsed: ", count, elapsed2)
}
func TestMaxBidFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &MaxBidFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. MaxBidFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDailyBudgetFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &DailyBudgetFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. DailyBudgetFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSpendStrategyFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &SpendStrategyFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. SpendStrategyFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestTotalBudgetFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &TotalBudgetFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. TotalBudgetFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFreqCapFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &FreqCapFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. FreqCapFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBidderTimeFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &BidderTimeFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. BidderTimeFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDayPartingFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &DayPartingFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. DayPartingFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCountriesFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &CountriesFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. CountriesFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestStatesFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &StatesFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. StatesFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCitiesFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &CitiesFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. CitiesFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestADXFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &ADXFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. ADXFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestControlListFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &ControlListFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. ControlListFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCategoriesFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &CategoriesFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. CategoriesFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCarriersFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &CarriersFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. CarriersFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestOSFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &OSFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. OSFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAudiencesFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &AudiencesFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. AudiencesFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDailyPPBFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mbf := &DailyPPBFilter{
			Code: tt.fields.Code,
		}
		if got := mbf.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. DailyPPBFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAndCPFilter_Filter(t *testing.T) {
	type fields struct {
		filter []CPFilter
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		af := &AndCPFilter{
			filter: tt.fields.filter,
		}
		if got := af.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. AndCPFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAndCPFilter_AndFilter(t *testing.T) {
	type fields struct {
		filter []CPFilter
	}
	type args struct {
		filter []CPFilter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		af := &AndCPFilter{
			filter: tt.fields.filter,
		}
		af.AndFilter(tt.args.filter...)
	}
}

func TestOrCPFilter_Filter(t *testing.T) {
	type fields struct {
		filter []CPFilter
	}
	type args struct {
		compaign model.Compaign
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		af := &OrCPFilter{
			filter: tt.fields.filter,
		}
		if got := af.Filter(tt.args.compaign, tt.args.offer); got != tt.want {
			t.Errorf("%q. OrCPFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestOrCPFilter_OrFilter(t *testing.T) {
	type fields struct {
		filter []CPFilter
	}
	type args struct {
		filter []CPFilter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		af := &OrCPFilter{
			filter: tt.fields.filter,
		}
		af.OrFilter(tt.args.filter...)
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
}
