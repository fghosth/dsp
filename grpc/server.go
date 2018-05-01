package grpc

import (
	"net"

	google_protobuf "github.com/golang/protobuf/ptypes/any"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials" // 引入grpc认证包
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata" // grpc metadata包
	"jvole.com/dsp/config"
	pb "jvole.com/dsp/grpc/proto"
	"jvole.com/dsp/index"
	"jvole.com/dsp/util"
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

func (h compaignService) GetCompaignIDs(ctx context.Context, in *pb.CompaignIDsRequest) (*pb.CompaignIDsReply, error) {
	resp := new(pb.CompaignIDsReply)
	var err error
	cmps := index.CPINDEX.GetAllCompaigns()
	for _, v := range cmps {
		resp.Cids = append(resp.Cids, v.ID)
	}
	return resp, err
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

//Run 运行
func Run() {

	listen, err := net.Listen("tcp", config.GrpcAddress)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	// TLS认证
	creds, err := credentials.NewServerTLSFromFile(config.GRPCPEM, config.GRPCKEY)
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

	// pp.Println("Listen on " + Address + " with TLS")
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
