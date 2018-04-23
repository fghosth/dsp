package util

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jinzhu/now"
	"github.com/k0kubun/pp"
)

func TestDealWithClickUrl(t *testing.T) {
	clickurl := "http://iytg3a.nbtrk6.com/22b90d96-8adf-4978-9229-4cc643a0e12f"
	url := DealWithClickURL(clickurl)
	fmt.Println(url)
}
func aTestTmp(t *testing.T) {
	count := 100
	for i := 0; i < count; i++ {
		uuid := NewUUID(time.Now().UTC().String())
		pp.Println("==========", uuid)
	}

}

func aTestArrInsertAfter(t *testing.T) {
	source := []interface{}{"apple", "orange", "plum", "banana", "grape", "derek"}

	res := ArrInsertAfter(1, "abccccc", source)
	fmt.Println(res)
}

func aTestArrDeletePos(t *testing.T) {
	source := []interface{}{"apple", "orange", "plum", "banana", "grape", "derek"}

	res := ArrDeletePos(2, source)
	fmt.Println(res)
}

func aTestZoneOffset(t *testing.T) {
	res := ZoneOffset("-01:00")

	fmt.Println("time", res)

	zoneName := "dddd"
	rest := time.Now().UTC().In(time.FixedZone(zoneName, res))
	name, offset := rest.Zone()
	timestamp := rest.Unix() - int64(offset)
	tm := time.Unix(timestamp, 0).UTC()

	fmt.Println("zone====", name, offset, tm)
	fmt.Println("fixedzone", rest, time.Now().UTC())
	// nows := now.BeginningOfMinute().Format("2006-01-02 15:04:05")
	nows := time.Now().UTC().Format("2006-01-02 15:04:05")

	//========方法一
	loc, _ := time.LoadLocation("Asia/Seoul")
	loc3, _ := time.LoadLocation("UTC")
	dt3, _ := time.ParseInLocation("2006-01-02 15:04:05", nows, loc3)
	dt, _ := time.ParseInLocation("2006-01-02 15:04:05", nows, loc)
	fmt.Println("方法一", dt, dt3, dt.Unix(), dt3.Unix())
	loc2, _ := time.LoadLocation("Asia/Seoul")

	dt2, _ := time.ParseInLocation("2006-01-02 15:04:05", nows, loc2)

	pp.Println(now.BeginningOfMinute(), dt, dt2, (dt2.Unix()-dt.Unix())/60/60)
	fmt.Println(time.Now().Unix())
	d, _ := time.ParseDuration("24h")
	fmt.Println("==", dt.Add(d).Unix()-dt.Unix())
	//===========方法二
	m, _ := time.ParseDuration("-2m")
	fmt.Println(time.Now().UTC().Add(m), time.Now().UTC().Unix())

	ti, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 12:00:09 +0800")
	ti2, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 12:00:09 +0900")
	fmt.Println("方法二", ti, ti2, ti.Unix(), ti2.Unix())
}

func TestIsInsegment(t *testing.T) {
	//http://www.ab126.com/goju/1840.html 网段计算器
	ip := "192.212.0.1"
	ips := "192.212.15.254/20"
	t1 := time.Now() // get current time
	count := 10
	var res bool
	for i := 0; i < count; i++ {
		res = IsInsegment(ip, ips)
	}
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)
	fmt.Println("ipsegment:", res)
}

func TestIsIP(t *testing.T) {
	// ip := []string{"192.1.3.4", "192.13.3.4", "192.1333.3.4"}
	test := []struct {
		ip     string
		wanted bool
	}{
		{"192.10.3.4", true},
		{"192.12.3.4", true},
		{"192.1.34.4", true},
		{"192.1.3.45", true},
		{"192.1.3.466", false},
	}
	for _, v := range test {
		if got := IsIP(v.ip); got != v.wanted {
			t.Errorf("IsIP(%s),wanted %t,but %t", v.ip, v.wanted, got)
		}
	}
}

func TestIsIPsegment(t *testing.T) {
	// ip := []string{"192.1.3.4", "192.13.3.4", "192.1333.3.4"}
	test := []struct {
		ip     string
		wanted bool
	}{
		{"192.10.3.4/0", false},
		{"192.12.3.4/1", true},
		{"192.1.34.4/3", true},
		{"192.1.3.45/33", false},
		{"192.1.388.46/21", false},
	}
	for _, v := range test {
		if got := IsIPSegment(v.ip); got != v.wanted {
			t.Errorf("util.IsIP(%s),wanted %t,but %t", v.ip, v.wanted, got)
		}
	}
}

func TestIs2N(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := Is2N(tt.args.n); got != tt.want {
			t.Errorf("%q. Is2N() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetMd5String(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := GetMd5String(tt.args.s); got != tt.want {
			t.Errorf("%q. GetMd5String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestArrInsertAfter(t *testing.T) {
	type args struct {
		pos    int
		data   interface{}
		source []interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := ArrInsertAfter(tt.args.pos, tt.args.data, tt.args.source); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ArrInsertAfter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestArrDeletePos(t *testing.T) {
	type args struct {
		pos    int
		source []interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := ArrDeletePos(tt.args.pos, tt.args.source); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ArrDeletePos() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetGuid(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := GetGuid(); got != tt.want {
			t.Errorf("%q. GetGuid() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDoZlibCompress(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := DoZlibCompress(tt.args.src); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. DoZlibCompress() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDoZlibUnCompress(t *testing.T) {
	type args struct {
		compressSrc []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := DoZlibUnCompress(tt.args.compressSrc); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. DoZlibUnCompress() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestEncodeStructToByte(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := EncodeStructToByte(tt.args.data)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. EncodeStructToByte() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. EncodeStructToByte() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDecodeByteToStruct(t *testing.T) {
	type args struct {
		data []byte
		to   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := DecodeByteToStruct(tt.args.data, tt.args.to); (err != nil) != tt.wantErr {
			t.Errorf("%q. DecodeByteToStruct() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestReadFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := ReadFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ReadFile() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := PathExists(tt.args.path)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. PathExists() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. PathExists() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestZoneOffset(t *testing.T) {
	type args struct {
		offset string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := ZoneOffset(tt.args.offset); got != tt.want {
			t.Errorf("%q. ZoneOffset() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIsIPSegment(t *testing.T) {
	type args struct {
		ips string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := IsIPSegment(tt.args.ips); got != tt.want {
			t.Errorf("%q. IsIPSegment() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestInt64ToBytes(t *testing.T) {
	type args struct {
		i uint64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := Int64ToBytes(tt.args.i); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Int64ToBytes() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBytesToInt64(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := BytesToInt64(tt.args.buf); got != tt.want {
			t.Errorf("%q. BytesToInt64() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDealXMLURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := DealXMLURL(tt.args.url); got != tt.want {
			t.Errorf("%q. DealXMLURL() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestTracefile(t *testing.T) {
	type args struct {
		str_content string
		path        string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		Tracefile(tt.args.str_content, tt.args.path)
	}
}

func TestGetResponseZip(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		w := &bytes.Buffer{}
		if err := GetResponseZip(w, tt.args.data); (err != nil) != tt.wantErr {
			t.Errorf("%q. GetResponseZip() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if gotW := w.String(); gotW != tt.wantW {
			t.Errorf("%q. GetResponseZip() = %v, want %v", tt.name, gotW, tt.wantW)
		}
	}
}
