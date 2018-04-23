package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"time"

	"github.com/k0kubun/pp"
	"jvole.com/dsp/util"
)

var success, failed, total uint64

func requestDspBidder(url string, count int, t int, data []byte, medth string) {
	var wg sync.WaitGroup
	var num int
	t2 := time.Now() // get current time
	sleep := 1000 * 1000 / count
	for range time.Tick(1 * time.Second) {
		if num >= t {
			break
		}
		for i := 0; i < count; i++ {
			time.Sleep(time.Microsecond * time.Duration(sleep))
			wg.Add(1)
			go func() {
				var err error
				switch medth {
				case "post":
					_, err = util.DoBytesPost(url, data)
				case "get":
					_, err = util.DoHTTPGet(url)
				}

				// pp.Println(string(res))
				if err != nil {
					atomic.AddUint64(&failed, 1)
					// log.Println(err)
				} else {
					atomic.AddUint64(&success, 1)
				}
				wg.Done()
			}()
		}
		num++
	}

	wg.Wait()
	elapsed2 := time.Since(t2)
	fmt.Println("DSP: ", elapsed2)
}

func main() {
	medth := flag.String("m", "post", "提交方法")
	num := flag.Int("n", 100, "1秒内并发数量")
	t := flag.Int("t", 3, "持续次数")
	url := flag.String("url", "http://adx-exads.newbidder.com/v1/Bidder", "测试的接口")
	file := flag.String("f", "body.json", "post传输时body中json数据")
	flag.Parse()
	// url := "http://adx-exads.newbidder.com/v1/Bidder"
	// url := "http://localhost:9900/v1/Bidder"
	// url := "http://www.sunya.com.cn"
	// if err != nil {
	// 	log.Fatalln("参数不正确", err)
	// }
	data := util.ReadFile(*file)
	total = uint64(*num * *t)
	// go func() {
	// 	http.ListenAndServe("localhost:8080", nil)
	// }()
	// fm, err := os.OpenFile("mem.out", os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.WriteHeapProfile(fm)
	// fm.Close()
	requestDspBidder(*url, *num, *t, data, *medth)
	pp.Println("total:", int(atomic.LoadUint64(&total)))
	pp.Println("success:", int(atomic.LoadUint64(&success)))
	pp.Println("failed:", int(atomic.LoadUint64(&failed)))
}

var adxcontext = []byte(`
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
  `)
