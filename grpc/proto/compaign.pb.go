// Code generated by protoc-gen-go. DO NOT EDIT.
// source: compaign.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	compaign.proto

It has these top-level messages:
	CompaignRequest
	CompaignReply
	CompaignIDsRequest
	CompaignIDsReply
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/any"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// CompaignRequest 请求结构
type CompaignRequest struct {
	Cid []uint32 `protobuf:"varint,1,rep,packed,name=cid" json:"cid,omitempty"`
}

func (m *CompaignRequest) Reset()                    { *m = CompaignRequest{} }
func (m *CompaignRequest) String() string            { return proto1.CompactTextString(m) }
func (*CompaignRequest) ProtoMessage()               {}
func (*CompaignRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CompaignRequest) GetCid() []uint32 {
	if m != nil {
		return m.Cid
	}
	return nil
}

// CompaignReply 响应结构
type CompaignReply struct {
	Total        uint32                 `protobuf:"varint,1,opt,name=total" json:"total,omitempty"`
	Compaignlist []*google_protobuf.Any `protobuf:"bytes,2,rep,name=compaignlist" json:"compaignlist,omitempty"`
}

func (m *CompaignReply) Reset()                    { *m = CompaignReply{} }
func (m *CompaignReply) String() string            { return proto1.CompactTextString(m) }
func (*CompaignReply) ProtoMessage()               {}
func (*CompaignReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CompaignReply) GetTotal() uint32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *CompaignReply) GetCompaignlist() []*google_protobuf.Any {
	if m != nil {
		return m.Compaignlist
	}
	return nil
}

// CompaignIDsReply 响应结构
type CompaignIDsRequest struct {
	Object *google_protobuf.Any `protobuf:"bytes,1,opt,name=object" json:"object,omitempty"`
}

func (m *CompaignIDsRequest) Reset()                    { *m = CompaignIDsRequest{} }
func (m *CompaignIDsRequest) String() string            { return proto1.CompactTextString(m) }
func (*CompaignIDsRequest) ProtoMessage()               {}
func (*CompaignIDsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CompaignIDsRequest) GetObject() *google_protobuf.Any {
	if m != nil {
		return m.Object
	}
	return nil
}

// CompaignIDsReply 响应结构
type CompaignIDsReply struct {
	Cids []uint32 `protobuf:"varint,1,rep,packed,name=cids" json:"cids,omitempty"`
}

func (m *CompaignIDsReply) Reset()                    { *m = CompaignIDsReply{} }
func (m *CompaignIDsReply) String() string            { return proto1.CompactTextString(m) }
func (*CompaignIDsReply) ProtoMessage()               {}
func (*CompaignIDsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CompaignIDsReply) GetCids() []uint32 {
	if m != nil {
		return m.Cids
	}
	return nil
}

func init() {
	proto1.RegisterType((*CompaignRequest)(nil), "proto.CompaignRequest")
	proto1.RegisterType((*CompaignReply)(nil), "proto.CompaignReply")
	proto1.RegisterType((*CompaignIDsRequest)(nil), "proto.CompaignIDsRequest")
	proto1.RegisterType((*CompaignIDsReply)(nil), "proto.CompaignIDsReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Compaign service

type CompaignClient interface {
	// 定义GetCompaign方法
	GetCompaign(ctx context.Context, in *CompaignRequest, opts ...grpc.CallOption) (*CompaignReply, error)
	GetCompaignIDs(ctx context.Context, in *CompaignIDsRequest, opts ...grpc.CallOption) (*CompaignIDsReply, error)
}

type compaignClient struct {
	cc *grpc.ClientConn
}

func NewCompaignClient(cc *grpc.ClientConn) CompaignClient {
	return &compaignClient{cc}
}

func (c *compaignClient) GetCompaign(ctx context.Context, in *CompaignRequest, opts ...grpc.CallOption) (*CompaignReply, error) {
	out := new(CompaignReply)
	err := grpc.Invoke(ctx, "/proto.Compaign/GetCompaign", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compaignClient) GetCompaignIDs(ctx context.Context, in *CompaignIDsRequest, opts ...grpc.CallOption) (*CompaignIDsReply, error) {
	out := new(CompaignIDsReply)
	err := grpc.Invoke(ctx, "/proto.Compaign/GetCompaignIDs", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Compaign service

type CompaignServer interface {
	// 定义GetCompaign方法
	GetCompaign(context.Context, *CompaignRequest) (*CompaignReply, error)
	GetCompaignIDs(context.Context, *CompaignIDsRequest) (*CompaignIDsReply, error)
}

func RegisterCompaignServer(s *grpc.Server, srv CompaignServer) {
	s.RegisterService(&_Compaign_serviceDesc, srv)
}

func _Compaign_GetCompaign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompaignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaignServer).GetCompaign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Compaign/GetCompaign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaignServer).GetCompaign(ctx, req.(*CompaignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Compaign_GetCompaignIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompaignIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaignServer).GetCompaignIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Compaign/GetCompaignIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaignServer).GetCompaignIDs(ctx, req.(*CompaignIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Compaign_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Compaign",
	HandlerType: (*CompaignServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCompaign",
			Handler:    _Compaign_GetCompaign_Handler,
		},
		{
			MethodName: "GetCompaignIDs",
			Handler:    _Compaign_GetCompaignIDs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "compaign.proto",
}

func init() { proto1.RegisterFile("compaign.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0xce, 0xcf, 0x2d,
	0x48, 0xcc, 0x4c, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0x92,
	0xe9, 0xf9, 0xf9, 0xe9, 0x39, 0xa9, 0xfa, 0x60, 0x5e, 0x52, 0x69, 0x9a, 0x7e, 0x62, 0x5e, 0x25,
	0x44, 0x85, 0x92, 0x32, 0x17, 0xbf, 0x33, 0x54, 0x4f, 0x50, 0x6a, 0x61, 0x69, 0x6a, 0x71, 0x89,
	0x90, 0x00, 0x17, 0x73, 0x72, 0x66, 0x8a, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0x6f, 0x10, 0x88, 0xa9,
	0x14, 0xcf, 0xc5, 0x8b, 0x50, 0x54, 0x90, 0x53, 0x29, 0x24, 0xc2, 0xc5, 0x5a, 0x92, 0x5f, 0x92,
	0x98, 0x23, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1b, 0x04, 0xe1, 0x08, 0x59, 0x70, 0xf1, 0xc0, 0xec,
	0xcf, 0xc9, 0x2c, 0x2e, 0x91, 0x60, 0x52, 0x60, 0xd6, 0xe0, 0x36, 0x12, 0xd1, 0x83, 0xd8, 0xae,
	0x07, 0xb3, 0x5d, 0xcf, 0x31, 0xaf, 0x32, 0x08, 0x45, 0xa5, 0x92, 0x13, 0x97, 0x10, 0xcc, 0x02,
	0x4f, 0x97, 0x62, 0x98, 0x43, 0x74, 0xb8, 0xd8, 0xf2, 0x93, 0xb2, 0x52, 0x93, 0x4b, 0xc0, 0xd6,
	0xe0, 0x32, 0x09, 0xaa, 0x46, 0x49, 0x8d, 0x4b, 0x00, 0xc5, 0x0c, 0x90, 0x3b, 0x85, 0xb8, 0x58,
	0x92, 0x33, 0x53, 0x8a, 0xa1, 0x7e, 0x01, 0xb3, 0x8d, 0x26, 0x32, 0x72, 0x71, 0xc0, 0x14, 0x0a,
	0xd9, 0x72, 0x71, 0xbb, 0xa7, 0x96, 0xc0, 0xb9, 0x62, 0x10, 0xa3, 0xf5, 0xd0, 0x82, 0x44, 0x4a,
	0x04, 0x43, 0xbc, 0x20, 0xa7, 0x52, 0x89, 0x41, 0xc8, 0x8d, 0x8b, 0x0f, 0x49, 0xbb, 0xa7, 0x4b,
	0xb1, 0x90, 0x24, 0x9a, 0x4a, 0x84, 0x77, 0xa4, 0xc4, 0xb1, 0x49, 0x81, 0xcd, 0x49, 0x62, 0x03,
	0xcb, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4e, 0x32, 0xde, 0x52, 0xc0, 0x01, 0x00, 0x00,
}