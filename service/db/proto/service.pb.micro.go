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
	// user 相关服务
	InsertUser(ctx context.Context, in *InsertUserReq, opts ...client.CallOption) (*InsertUserResp, error)
	QueryUser(ctx context.Context, in *QueryUserReq, opts ...client.CallOption) (*QueryUserResp, error)
	// file 相关服务
	InsertFileMeta(ctx context.Context, in *InsertFileMetaReq, opts ...client.CallOption) (*InsertFileMetaResp, error)
	QueryFileMeta(ctx context.Context, in *QueryFileMetaReq, opts ...client.CallOption) (*QueryFileMetaResp, error)
	// user_file 相关服务
	ListUserFile(ctx context.Context, in *ListUserFileMetaReq, opts ...client.CallOption) (*ListUserFileMetaResp, error)
	InsertUserFile(ctx context.Context, in *InsertUserFileMetaReq, opts ...client.CallOption) (*InsertUserFileMetaResp, error)
	// session 相关服务
	InsertSession(ctx context.Context, in *InsertSessionReq, opts ...client.CallOption) (*InsertSessionResp, error)
	GetUserSession(ctx context.Context, in *GetUserSessionReq, opts ...client.CallOption) (*GetUserSessionResp, error)
	DeleteUserSession(ctx context.Context, in *DeleteUserSessionReq, opts ...client.CallOption) (*DeleteUserSessionResp, error)
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

func (c *service) InsertUser(ctx context.Context, in *InsertUserReq, opts ...client.CallOption) (*InsertUserResp, error) {
	req := c.c.NewRequest(c.name, "Service.InsertUser", in)
	out := new(InsertUserResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) QueryUser(ctx context.Context, in *QueryUserReq, opts ...client.CallOption) (*QueryUserResp, error) {
	req := c.c.NewRequest(c.name, "Service.QueryUser", in)
	out := new(QueryUserResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) InsertFileMeta(ctx context.Context, in *InsertFileMetaReq, opts ...client.CallOption) (*InsertFileMetaResp, error) {
	req := c.c.NewRequest(c.name, "Service.InsertFileMeta", in)
	out := new(InsertFileMetaResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) QueryFileMeta(ctx context.Context, in *QueryFileMetaReq, opts ...client.CallOption) (*QueryFileMetaResp, error) {
	req := c.c.NewRequest(c.name, "Service.QueryFileMeta", in)
	out := new(QueryFileMetaResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) ListUserFile(ctx context.Context, in *ListUserFileMetaReq, opts ...client.CallOption) (*ListUserFileMetaResp, error) {
	req := c.c.NewRequest(c.name, "Service.ListUserFile", in)
	out := new(ListUserFileMetaResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) InsertUserFile(ctx context.Context, in *InsertUserFileMetaReq, opts ...client.CallOption) (*InsertUserFileMetaResp, error) {
	req := c.c.NewRequest(c.name, "Service.InsertUserFile", in)
	out := new(InsertUserFileMetaResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) InsertSession(ctx context.Context, in *InsertSessionReq, opts ...client.CallOption) (*InsertSessionResp, error) {
	req := c.c.NewRequest(c.name, "Service.InsertSession", in)
	out := new(InsertSessionResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) GetUserSession(ctx context.Context, in *GetUserSessionReq, opts ...client.CallOption) (*GetUserSessionResp, error) {
	req := c.c.NewRequest(c.name, "Service.GetUserSession", in)
	out := new(GetUserSessionResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) DeleteUserSession(ctx context.Context, in *DeleteUserSessionReq, opts ...client.CallOption) (*DeleteUserSessionResp, error) {
	req := c.c.NewRequest(c.name, "Service.DeleteUserSession", in)
	out := new(DeleteUserSessionResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Service service

type ServiceHandler interface {
	// user 相关服务
	InsertUser(context.Context, *InsertUserReq, *InsertUserResp) error
	QueryUser(context.Context, *QueryUserReq, *QueryUserResp) error
	// file 相关服务
	InsertFileMeta(context.Context, *InsertFileMetaReq, *InsertFileMetaResp) error
	QueryFileMeta(context.Context, *QueryFileMetaReq, *QueryFileMetaResp) error
	// user_file 相关服务
	ListUserFile(context.Context, *ListUserFileMetaReq, *ListUserFileMetaResp) error
	InsertUserFile(context.Context, *InsertUserFileMetaReq, *InsertUserFileMetaResp) error
	// session 相关服务
	InsertSession(context.Context, *InsertSessionReq, *InsertSessionResp) error
	GetUserSession(context.Context, *GetUserSessionReq, *GetUserSessionResp) error
	DeleteUserSession(context.Context, *DeleteUserSessionReq, *DeleteUserSessionResp) error
}

func RegisterServiceHandler(s server.Server, hdlr ServiceHandler, opts ...server.HandlerOption) error {
	type service interface {
		InsertUser(ctx context.Context, in *InsertUserReq, out *InsertUserResp) error
		QueryUser(ctx context.Context, in *QueryUserReq, out *QueryUserResp) error
		InsertFileMeta(ctx context.Context, in *InsertFileMetaReq, out *InsertFileMetaResp) error
		QueryFileMeta(ctx context.Context, in *QueryFileMetaReq, out *QueryFileMetaResp) error
		ListUserFile(ctx context.Context, in *ListUserFileMetaReq, out *ListUserFileMetaResp) error
		InsertUserFile(ctx context.Context, in *InsertUserFileMetaReq, out *InsertUserFileMetaResp) error
		InsertSession(ctx context.Context, in *InsertSessionReq, out *InsertSessionResp) error
		GetUserSession(ctx context.Context, in *GetUserSessionReq, out *GetUserSessionResp) error
		DeleteUserSession(ctx context.Context, in *DeleteUserSessionReq, out *DeleteUserSessionResp) error
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

func (h *serviceHandler) InsertUser(ctx context.Context, in *InsertUserReq, out *InsertUserResp) error {
	return h.ServiceHandler.InsertUser(ctx, in, out)
}

func (h *serviceHandler) QueryUser(ctx context.Context, in *QueryUserReq, out *QueryUserResp) error {
	return h.ServiceHandler.QueryUser(ctx, in, out)
}

func (h *serviceHandler) InsertFileMeta(ctx context.Context, in *InsertFileMetaReq, out *InsertFileMetaResp) error {
	return h.ServiceHandler.InsertFileMeta(ctx, in, out)
}

func (h *serviceHandler) QueryFileMeta(ctx context.Context, in *QueryFileMetaReq, out *QueryFileMetaResp) error {
	return h.ServiceHandler.QueryFileMeta(ctx, in, out)
}

func (h *serviceHandler) ListUserFile(ctx context.Context, in *ListUserFileMetaReq, out *ListUserFileMetaResp) error {
	return h.ServiceHandler.ListUserFile(ctx, in, out)
}

func (h *serviceHandler) InsertUserFile(ctx context.Context, in *InsertUserFileMetaReq, out *InsertUserFileMetaResp) error {
	return h.ServiceHandler.InsertUserFile(ctx, in, out)
}

func (h *serviceHandler) InsertSession(ctx context.Context, in *InsertSessionReq, out *InsertSessionResp) error {
	return h.ServiceHandler.InsertSession(ctx, in, out)
}

func (h *serviceHandler) GetUserSession(ctx context.Context, in *GetUserSessionReq, out *GetUserSessionResp) error {
	return h.ServiceHandler.GetUserSession(ctx, in, out)
}

func (h *serviceHandler) DeleteUserSession(ctx context.Context, in *DeleteUserSessionReq, out *DeleteUserSessionResp) error {
	return h.ServiceHandler.DeleteUserSession(ctx, in, out)
}