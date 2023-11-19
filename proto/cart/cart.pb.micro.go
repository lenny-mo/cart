// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/cart/cart.proto

package cart

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

// Api Endpoints for Cart service

func NewCartEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Cart service

type CartService interface {
	Add(ctx context.Context, in *AddCartRequest, opts ...client.CallOption) (*AddCartResponse, error)
	FindAll(ctx context.Context, in *FindAllCartRequest, opts ...client.CallOption) (*FindAllCartResponse, error)
	CheckOutCart(ctx context.Context, in *CheckOutCartRequest, opts ...client.CallOption) (*CheckOutCartResponse, error)
}

type cartService struct {
	c    client.Client
	name string
}

func NewCartService(name string, c client.Client) CartService {
	return &cartService{
		c:    c,
		name: name,
	}
}

func (c *cartService) Add(ctx context.Context, in *AddCartRequest, opts ...client.CallOption) (*AddCartResponse, error) {
	req := c.c.NewRequest(c.name, "Cart.Add", in)
	out := new(AddCartResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartService) FindAll(ctx context.Context, in *FindAllCartRequest, opts ...client.CallOption) (*FindAllCartResponse, error) {
	req := c.c.NewRequest(c.name, "Cart.FindAll", in)
	out := new(FindAllCartResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartService) CheckOutCart(ctx context.Context, in *CheckOutCartRequest, opts ...client.CallOption) (*CheckOutCartResponse, error) {
	req := c.c.NewRequest(c.name, "Cart.CheckOutCart", in)
	out := new(CheckOutCartResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Cart service

type CartHandler interface {
	Add(context.Context, *AddCartRequest, *AddCartResponse) error
	FindAll(context.Context, *FindAllCartRequest, *FindAllCartResponse) error
	CheckOutCart(context.Context, *CheckOutCartRequest, *CheckOutCartResponse) error
}

func RegisterCartHandler(s server.Server, hdlr CartHandler, opts ...server.HandlerOption) error {
	type cart interface {
		Add(ctx context.Context, in *AddCartRequest, out *AddCartResponse) error
		FindAll(ctx context.Context, in *FindAllCartRequest, out *FindAllCartResponse) error
		CheckOutCart(ctx context.Context, in *CheckOutCartRequest, out *CheckOutCartResponse) error
	}
	type Cart struct {
		cart
	}
	h := &cartHandler{hdlr}
	return s.Handle(s.NewHandler(&Cart{h}, opts...))
}

type cartHandler struct {
	CartHandler
}

func (h *cartHandler) Add(ctx context.Context, in *AddCartRequest, out *AddCartResponse) error {
	return h.CartHandler.Add(ctx, in, out)
}

func (h *cartHandler) FindAll(ctx context.Context, in *FindAllCartRequest, out *FindAllCartResponse) error {
	return h.CartHandler.FindAll(ctx, in, out)
}

func (h *cartHandler) CheckOutCart(ctx context.Context, in *CheckOutCartRequest, out *CheckOutCartResponse) error {
	return h.CartHandler.CheckOutCart(ctx, in, out)
}
