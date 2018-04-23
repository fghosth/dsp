package util

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"
)

func aTestDoBytesPost(t *testing.T) {
	url := "http://beta.newbidder.com/dsp-campaign/getMany"
	data := []byte(`{"cmd":"select * from BaseDspCampaign where id=2  limit 1"}`)
	res, err := DoBytesPost(url, data)
	if err != nil {
		log.Printf("request err:%s", err)
	}
	fmt.Println(string(res))
}

func aTestInfluxdb(t *testing.T) {
	url := "http://influx.newbidder.com/influx/v2/query"
	data := []byte(`{"cmd":"select Visits  from adstatis_table_84 limit 1","db":"tracking","precision":"s","uid":[84]}`)
	var wg sync.WaitGroup

	count := 1
	t2 := time.Now() // get current time
	for i := 0; i < count; i++ {
		wg.Add(1)
		time.Sleep(1)
		go func() {
			fmt.Println("e")
			res, err := DoBytesPost(url, data)
			if err != nil {
				log.Printf("request err:%s", err)
			}
			fmt.Println(string(res))
			wg.Done()
		}()
	}
	wg.Wait()
	elapsed2 := time.Since(t2)
	fmt.Println("App elapsed2: ", count, elapsed2)
	// fmt.Println(string(res))
}

func TestDoBytesPost(t *testing.T) {
	type args struct {
		url  string
		data []byte
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
		got, err := DoBytesPost(tt.args.url, tt.args.data)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. DoBytesPost() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. DoBytesPost() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestDoHTTPGet(t *testing.T) {
	type args struct {
		url string
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
		got, err := DoHTTPGet(tt.args.url)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. DoHTTPGet() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. DoHTTPGet() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
