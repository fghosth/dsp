package main

import (
	"github.com/k0kubun/pp"
	pb "jvole.com/dsp/grpc/proto" // 引入proto包
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials" // 引入grpc认证包
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:5005"
	// OpenTLS 是否开启TLS认证
	OpenTLS = true
)

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

func main() {
	var opts []grpc.DialOption

	if OpenTLS {
		// TLS连接
		creds, err := credentials.NewClientTLSFromFile("/data/dsp/server.pem", "DSP")
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

	if err != nil {
		grpclog.Fatalln(err)
	}

	defer conn.Close()

	// 初始化客户端
	c := pb.NewCompaignClient(conn)

	// 调用方法
	reqBody := new(pb.CompaignRequest)
	reqBody.Cid = []uint32{12627, 12654}
	r, err := c.GetCompaign(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}
	for _, v := range r.Compaignlist {
		cmp := new(model.Compaign)
		util.DecodeByteToStruct(v.Value, cmp)
		pp.Println(cmp.EndDate)
	}

	pp.Println("total:", int(r.Total))
}
