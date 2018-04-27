package main

import (
	// "context"

	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	pb "jvole.com/dsp/grpc/proto" // 引入proto包
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"

	"golang.org/x/net/context"
)

var (
	password = "zaqwedcxs"
	islogin  = false
	// Address gRPC服务地址
	Address = "127.0.0.1:5005"
	// OpenTLS 是否开启TLS认证
	OpenTLS        = true
	CID     uint32 = 0
	stop           = false //停止信号
	pemfile        = "../grpc/keys/server.pem"
	Conn    *grpc.ClientConn
	fixcash float64 = 1000000
)

func usage(w io.Writer) {
	io.WriteString(w, "监控dsp内存数据【set compaign】 设置当前监控的compaignid:\n【show】显示正在监控的compaignid\n【watch】每隔一秒显示需要监控的内容\n【get】获取需要监控的内容\np.s.输入[$]切换是否每秒刷新")
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

var completer = readline.NewPrefixCompleter(

	readline.PcItem("login"),
	readline.PcItem("set",
		readline.PcItem("compaign"),
	),
	readline.PcItem("show"),
	readline.PcItem("watch",
		readline.PcItem("dailyBudgetRecore"),
		readline.PcItem("dailyPPBRecord"),
		readline.PcItem("totalBudgetRecord"),
		readline.PcItem("freqRecord"),
		readline.PcItem("status"),
	),
	readline.PcItem("get",
		readline.PcItem("allCid"),
		readline.PcItem("number"),
		readline.PcItem("compaign"),
		readline.PcItem("dailyBudgetRecore"),
		readline.PcItem("dailyPPBRecord"),
		readline.PcItem("totalBudgetRecord"),
		readline.PcItem("freqRecord"),
		readline.PcItem("status"),
	),
	readline.PcItem("exit"),
	readline.PcItem("help"),
)

type stopListener struct{}

func (sl *stopListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	if "$" == string(line) {
		if stop {
			stop = false
		} else {
			stop = true
		}
		return nil, 0, false
	}
	return line, pos, false
}

func main() {
	defer close()
	server := flag.String("s", "", "server")
	keyfile := flag.String("f", "server.pem", "pem file")
	flag.Parse()
	Address = *server
	pemfile = *keyfile
	conn()

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[39m»\033[0m ",
		HistoryFile:     "readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		Listener:        new(stopListener),
		// HistorySearchFold:   true,
		// FuncFilterInputRune: filterInput,
	})

	if err != nil {
		panic(err)
	}
	defer l.Close()
login:
	pswd, err := l.ReadPassword("please enter your password: ")
	if err != nil {
		os.Exit(0)
	}
	if string(pswd) != password {
		println("you passowrd is wrong:")
		goto login
	}
	log.SetOutput(l.Stderr())

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		switch {

		//设置全局CID
		case strings.HasPrefix(line, "set compaign "):
			n, _ := strconv.Atoi(line[13:])
			CID = uint32(n)
			//显示全局CID
		case line == "show":
			pp.Println("currentCID:", int(CID))
			//登录
		case line == "login":
			pswd, err := l.ReadPassword("please enter your password: ")
			if err != nil {
				break
			}
			if strconv.Quote(string(pswd)) != password {
				println("you passowrd is wrong:")
			}
			//监控相关内容
		case strings.HasPrefix(line, "watch "):
			switch line[6:] {
			case "status":
				go func() {
					for range time.Tick(1 * time.Second) {
						pp.Println(getstatus())
						if stop {
							return
						}
					}
				}()
			case "dailyBudgetRecore":
				go func() {
					for range time.Tick(1 * time.Second) {

						cmp := make([]*model.Compaign, 0)
						if CID != 0 {
							_, cmp = query([]uint32{CID})
						}
						if len(cmp) > 0 {
							fmt.Printf("cost:%f, SinceTime:%s, TillTime:%s\n", (float64(cmp[0].DailyBudgetRecores.Cost) / fixcash), cmp[0].DailyBudgetRecores.SinceTime, cmp[0].DailyBudgetRecores.TillTime)
						}
						if stop {
							return
						}
					}
				}()
			case "dailyPPBRecord":
				go func() {
					for range time.Tick(1 * time.Second) {

						cmp := make([]*model.Compaign, 0)
						if CID != 0 {
							_, cmp = query([]uint32{CID})
						}
						if len(cmp) > 0 {
							pp.Println(cmp[0].DailyPPBRecords)
						}
						if stop {
							return
						}
					}
				}()
			case "totalBudgetRecord":
				go func() {
					for range time.Tick(1 * time.Second) {

						cmp := make([]*model.Compaign, 0)
						if CID != 0 {
							_, cmp = query([]uint32{CID})
						}
						if len(cmp) > 0 {
							fmt.Printf("cost:%f\n", (float64(cmp[0].TotalBudgetRecords.Cost) / fixcash))
						}
						if stop {
							return
						}
					}
				}()
			case "freqRecord":
				go func() {
					for range time.Tick(1 * time.Second) {

						cmp := make([]*model.Compaign, 0)
						if CID != 0 {
							_, cmp = query([]uint32{CID})
						}
						if len(cmp) > 0 {
							pp.Println(cmp[0].FreqRecords)
						}
						if stop {
							return
						}
					}
				}()
			}
			//获取相关内容
		case strings.HasPrefix(line, "get "):
			sub := line[4:]
			cmp := make([]*model.Compaign, 0)
			if CID != 0 {
				_, cmp = query([]uint32{CID})
			}
			switch sub {
			case "status":
				pp.Println(getstatus())
			case "allCid":
				fmt.Println(getCompaignIDs())
			case "dailyBudgetRecore":
				if len(cmp) > 0 {
					fmt.Printf("cost:%f, SinceTime:%s, TillTime:%s", (float64(cmp[0].DailyBudgetRecores.Cost) / fixcash), cmp[0].DailyBudgetRecores.SinceTime, cmp[0].DailyBudgetRecores.TillTime)
				}
			case "dailyPPBRecord":
				if len(cmp) > 0 {
					pp.Println(cmp[0].DailyPPBRecords)
				}
			case "totalBudgetRecord":
				if len(cmp) > 0 {
					fmt.Printf("cost:%f", (float64(cmp[0].TotalBudgetRecords.Cost) / fixcash))
				}
			case "freqRecord":
				if len(cmp) > 0 {
					pp.Println(cmp[0].FreqRecords)
				}
			}
			if strings.HasPrefix(sub, "number") {
				total, _ := query([]uint32{0})
				pp.Println(total)
			} else if strings.HasPrefix(sub, "compaign ") {
				cidstr := strings.Split(sub[9:], ",")
				cid := make([]uint32, len(cidstr))
				for k, v := range cidstr { //拆解已,分隔的cid
					n, _ := strconv.Atoi(v)
					cid[k] = uint32(n) //字符串转uint32
				}
				_, cmp = query(cid)
				for _, v := range cmp { //遍历compaign列表
					if v.ID != 0 { //有则显示
						pp.Println(cmp)
					}
				}
			}
		case line == "help":
			usage(l.Stderr())
		case line == "exit":
			os.Exit(0)
		}
	}
}

// customCredential 自定义认证
type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "101010",
		"appkey": "zaqwedcxs",
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	if OpenTLS {
		return true
	}
	return false
}

//获得连接
func conn() {
	var opts []grpc.DialOption

	if OpenTLS {
		// TLS连接
		creds, err := credentials.NewClientTLSFromFile(pemfile, "DSP")
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// 指定自定义认证
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
	conn, err := grpc.Dial(Address, opts...)
	Conn = conn
	if err != nil {
		grpclog.Fatalln(err)
	}

}

func close() {
	Conn.Close()
}

//获得状态
func getstatus() (status string) {
	_, cmp := query([]uint32{CID})
	if len(cmp) > 0 {
		switch cmp[0].Status {
		case 0:
			status = "pending"
		case 1:
			status = "runing"
		}
	}
	return
}

func query(cid []uint32) (total int, compaign []*model.Compaign) {
	// 初始化客户端
	c := pb.NewCompaignClient(Conn)
	// 调用方法
	reqBody := new(pb.CompaignRequest)
	reqBody.Cid = cid
	r, err := c.GetCompaign(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}
	total = int(r.Total)

	for _, v := range r.Compaignlist {
		cmp := new(model.Compaign)
		util.DecodeByteToStruct(v.Value, cmp)
		compaign = append(compaign, cmp)
	}
	return
}

func getCompaignIDs() []uint32 {
	// 初始化客户端
	c := pb.NewCompaignClient(Conn)
	// 调用方法
	reqBody := new(pb.CompaignIDsRequest)
	r, err := c.GetCompaignIDs(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}
	return r.Cids
}
