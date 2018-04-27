package main

import (
	"log"
	"net"

	google_protobuf "github.com/golang/protobuf/ptypes/any"
	"github.com/k0kubun/pp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials" // 引入grpc认证包
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata" // grpc metadata包
	pb "jvole.com/dsp/grpc/proto"
	"jvole.com/dsp/index"
	"jvole.com/dsp/util"
)

const (
	// Address gRPC服务地址
	Address = ":5005"
)

// 定义compaignService并实现约定的接口
type compaignService struct{}

// CompaignService ...
var CompaignService = compaignService{}

func (h compaignService) GetCompaign(ctx context.Context, in *pb.CompaignRequest) (*pb.CompaignReply, error) {
	resp := new(pb.CompaignReply)
	resp.Total = uint32(index.CPINDEX.GetCPLen())
	compaign := make([]*google_protobuf.Any, 0)
	for _, v := range in.Cid {
		cmp := new(google_protobuf.Any)
		cmp.Value, _ = util.EncodeStructToByte(index.CPINDEX.GetCompaign(v))
		compaign = append(compaign, cmp)
	}
	resp.Compaignlist = compaign
	return resp, nil
}

// auth 验证Token
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "无Token认证信息")
	}

	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "zaqwedcxs" {
		return grpc.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
	}

	return nil
}
func main() {
	log.Println("初始化索引...")
	index.CPINDEX.SetupIndex()
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("../keys/server.pem", "../keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	opts = append(opts, grpc.Creds(creds))
	// 注册interceptor
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	// 实例化grpc Server, 并开启TLS认证
	s := grpc.NewServer(opts...)

	// 注册CompaignServer
	pb.RegisterCompaignServer(s, CompaignService)

	// 开启trace
	// go startTrace()
	// grpclog.Println("Listen on " + Address + " with TLS")

	pp.Println("Listen on " + Address + " with TLS")
	s.Serve(listen)
}

// func startTrace() {
// 	grpc.EnableTracing = true
// 	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
// 		return true, true
// 	}
// 	go http.ListenAndServe(":50051", nil)
// 	grpclog.Println("Trace listen on 50051")
// }
