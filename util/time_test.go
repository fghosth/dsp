package util

import (
	"reflect"
	"testing"
	"time"
)

func TestIsBetweenTime(t *testing.T) {
	since1, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 12:00:00 +0800")
	till1, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 13:00:00 +0800")
	t1, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 12:10:00 +0800")
	since2, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 16:00:00 +0800")
	till2, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-09 12:00:00 +0800")
	t2, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 15:00:00 +0800")
	since3, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 16:00:00 +0800")
	till3, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 17:00:00 +0800")
	t3, _ := time.Parse("2006-01-02 15:04:05 -0700", "2017-08-08 15:00:01 +0700")
	type args struct {
		since time.Time
		till  time.Time
		t     time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantRes bool
	}{
		{"IsBetweenTime1", *&args{since1, till1, t1}, true},
		{"IsBetweenTime2", *&args{since2, till2, t2}, false},
		{"IsBetweenTime3", *&args{since3, till3, t3}, true}, //不同时区比较
	}
	for _, tt := range tests {
		if gotRes := IsBetweenTime(tt.args.since, tt.args.till, tt.args.t); gotRes != tt.wantRes {
			t.Errorf("%q. IsBetweenTime() = %v, want %v", tt.name, gotRes, tt.wantRes)
		}
	}
}

func TestCovnNOWUTC2Location(t *testing.T) {
	nows := time.Now().UTC().Format("2006-01-02 15:04")
	loc, _ := time.LoadLocation("America/Bahia")
	wantdt1, _ := time.ParseInLocation("2006-01-02 15:04", nows, loc)
	type args struct {
		location string
		zone     string
	}
	tests := []struct {
		name    string
		args    args
		wantDt  time.Time
		wantErr bool
	}{
		{"CovnNOWUTC2Location1", *&args{"America/Bahia", "-03:00"}, wantdt1, false},
	}
	for _, tt := range tests {
		gotDt, err := CovnNOWUTC2Location(tt.args.location)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CovnNOWUTC2Location() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotDt.Format("2006-01-02 15:04"), tt.wantDt.Format("2006-01-02 15:04")) {
			t.Errorf("%q. CovnNOWUTC2Location() = %v, want %v", tt.name, gotDt.Format("2006-01-02 15:04"), tt.wantDt.Format("2006-01-02 15:04"))
		}
	}
}

func TestCovnTime2Location(t *testing.T) {
	nows := time.Now().UTC().Format("2006-01-02 15:04:05")
	loc, _ := time.LoadLocation("America/Bahia")
	wantdt1, _ := time.ParseInLocation("2006-01-02 15:04:05", nows, loc)
	type args struct {
		t        string
		location string
		zone     string
	}
	tests := []struct {
		name    string
		args    args
		wantDt  time.Time
		wantErr bool
	}{
		{"CovnTime2Location1", *&args{nows, "America/Bahia", "-03:00"}, wantdt1, false},
	}
	for _, tt := range tests {
		gotDt, err := CovnTime2Location(tt.args.t, tt.args.location)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CovnTime2Location() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotDt, tt.wantDt) {
			t.Errorf("%q. CovnTime2Location() = %v, want %v", tt.name, gotDt, tt.wantDt)
		}
	}
}

func TestZone2Zone(t *testing.T) {
	t1, _ := time.Parse("2006-01-02 15:04:05 -0700", "2018-08-08 12:00:00 +0800")
	zone1 := "Asia/Bangkok" //+07:00
	wantdt1, _ := time.Parse("2006-01-02 15:04:05 -0700", "2018-08-08 4:00:00 +0700")
	type args struct {
		t    time.Time
		zone string
	}
	tests := []struct {
		name    string
		args    args
		wantDt  time.Time
		wantErr bool
	}{
		{"Zone2Zone1", *&args{t1, zone1}, wantdt1, false},
	}
	for _, tt := range tests {
		gotDt, err := Zone2Zone(tt.args.t, tt.args.zone)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Zone2Zone() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotDt.Unix(), tt.wantDt.Unix()) {
			t.Errorf("%q. Zone2Zone() = %v, want %v", tt.name, gotDt, tt.wantDt)
		}
	}
}

func TestBgeinAndEndDAYOfZone(t *testing.T) {
	today := time.Now().UTC().Format("2006-01-02")
	oneDay, _ := time.ParseDuration("24h") //一天后
	zone1 := "Asia/Colombo"                //+05:30
	wantBegin1, _ := time.Parse("2006-01-02 15:04:05 -0700", today+" 00:00:00 +0530")
	endstrDay := wantBegin1.Add(oneDay).Format("2006-01-02")
	wantEnd1, _ := time.Parse("2006-01-02 15:04:05 -0700", endstrDay+" 00:00:00 +0530")

	zone2 := "Asia/Bangkok" //+07:00
	wantBegin2, _ := time.Parse("2006-01-02 15:04:05 -0700", today+" 00:00:00 +0700")
	endstrDay = wantBegin2.Add(oneDay).Format("2006-01-02")
	wantEnd2, _ := time.Parse("2006-01-02 15:04:05 -0700", endstrDay+" 00:00:00 +0700")

	zone3 := "America/Chicago" //-06:00
	wantBegin3, _ := time.Parse("2006-01-02 15:04:05 -0700", today+" 00:00:00 -0500")
	endstrDay = wantBegin3.Add(oneDay).Format("2006-01-02")
	wantEnd3, _ := time.Parse("2006-01-02 15:04:05 -0700", endstrDay+" 00:00:00 -0500")
	type args struct {
		zone string
	}
	tests := []struct {
		name      string
		args      args
		wantBegin time.Time
		wantEnd   time.Time
		wantErr   bool
	}{
		{"BgeinAndEndDAYOfZone1", *&args{zone1}, wantBegin1, wantEnd1, false},
		{"BgeinAndEndDAYOfZone2", *&args{zone2}, wantBegin2, wantEnd2, false},
		{"BgeinAndEndDAYOfZone3", *&args{zone3}, wantBegin3, wantEnd3, false},
	}
	for _, tt := range tests {
		gotBegin, gotEnd, err := BgeinAndEndDAYOfZone(tt.args.zone)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. BgeinAndEndDAYOfZone() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotBegin.Unix(), tt.wantBegin.Unix()) {
			t.Errorf("%q. BgeinAndEndDAYOfZone() gotBegin = %v, want %v", tt.name, gotBegin, tt.wantBegin)
		}
		if !reflect.DeepEqual(gotEnd.Unix(), tt.wantEnd.Unix()) {
			t.Errorf("%q. BgeinAndEndDAYOfZone() gotEnd = %v, want %v", tt.name, gotEnd, tt.wantEnd)
		}
	}
}
