// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: service.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Service service

func NewServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Service service

type Service interface {
	// JWT 相关服务
	DownloadJWTEncode(ctx context.Context, in *DownloadFileMeta, opts ...client.CallOption) (*EncodeResp, error)
	DownloadJWTDecode(ctx context.Context, in *DecodeReq, opts ...client.CallOption) (*DownloadJWTDecodeResp, error)
	UploadJWTEncode(ctx context.Context, in *UploadFileMeta, opts ...client.CallOption) (*EncodeResp, error)
	UploadJWTDecode(ctx context.Context, in *DecodeReq, opts ...client.CallOption) (*UploadJWTDecodeResp, error)
}

type service struct {
	c    client.Client
	name string
}

func NewService(name string, c client.Client) Service {
	return &service{
		c:    c,
		name: name,
	}
}

func (c *service) DownloadJWTEncode(ctx context.Context, in *DownloadFileMeta, opts ...client.CallOption) (*EncodeResp, error) {
	req := c.c.NewRequest(c.name, "Service.DownloadJWTEncode", in)
	out := new(EncodeResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) DownloadJWTDecode(ctx context.Context, in *DecodeReq, opts ...client.CallOption) (*DownloadJWTDecodeResp, error) {
	req := c.c.NewRequest(c.name, "Service.DownloadJWTDecode", in)
	out := new(DownloadJWTDecodeResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) UploadJWTEncode(ctx context.Context, in *UploadFileMeta, opts ...client.CallOption) (*EncodeResp, error) {
	req := c.c.NewRequest(c.name, "Service.UploadJWTEncode", in)
	out := new(EncodeResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) UploadJWTDecode(ctx context.Context, in *DecodeReq, opts ...client.CallOption) (*UploadJWTDecodeResp, error) {
	req := c.c.NewRequest(c.name, "Service.UploadJWTDecode", in)
	out := new(UploadJWTDecodeResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Service service

type ServiceHandler interface {
	// JWT 相关服务
	DownloadJWTEncode(context.Context, *DownloadFileMeta, *EncodeResp) error
	DownloadJWTDecode(context.Context, *DecodeReq, *DownloadJWTDecodeResp) error
	UploadJWTEncode(context.Context, *UploadFileMeta, *EncodeResp) error
	UploadJWTDecode(context.Context, *DecodeReq, *UploadJWTDecodeResp) error
}

func RegisterServiceHandler(s server.Server, hdlr ServiceHandler, opts ...server.HandlerOption) error {
	type service interface {
		DownloadJWTEncode(ctx context.Context, in *DownloadFileMeta, out *EncodeResp) error
		DownloadJWTDecode(ctx context.Context, in *DecodeReq, out *DownloadJWTDecodeResp) error
		UploadJWTEncode(ctx context.Context, in *UploadFileMeta, out *EncodeResp) error
		UploadJWTDecode(ctx context.Context, in *DecodeReq, out *UploadJWTDecodeResp) error
	}
	type Service struct {
		service
	}
	h := &serviceHandler{hdlr}
	return s.Handle(s.NewHandler(&Service{h}, opts...))
}

type serviceHandler struct {
	ServiceHandler
}

func (h *serviceHandler) DownloadJWTEncode(ctx context.Context, in *DownloadFileMeta, out *EncodeResp) error {
	return h.ServiceHandler.DownloadJWTEncode(ctx, in, out)
}

func (h *serviceHandler) DownloadJWTDecode(ctx context.Context, in *DecodeReq, out *DownloadJWTDecodeResp) error {
	return h.ServiceHandler.DownloadJWTDecode(ctx, in, out)
}

func (h *serviceHandler) UploadJWTEncode(ctx context.Context, in *UploadFileMeta, out *EncodeResp) error {
	return h.ServiceHandler.UploadJWTEncode(ctx, in, out)
}

func (h *serviceHandler) UploadJWTDecode(ctx context.Context, in *DecodeReq, out *UploadJWTDecodeResp) error {
	return h.ServiceHandler.UploadJWTDecode(ctx, in, out)
}