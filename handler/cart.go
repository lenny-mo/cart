package handler

import (
	"github.com/lenny-mo/cart/domain/models"
	"github.com/lenny-mo/cart/domain/services"
	"github.com/lenny-mo/cart/proto/cart"
	"context"
	"strconv"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// 这个handler 要实现micro server rpc 的接口
// Server API for Cart service
//
//	type CartHandler interface {
//		Add(context.Context, *AddCartRequest, *AddCartResponse) error
//		FindAll(context.Context, *FindAllCartRequest, *FindAllCartResponse) error
//		CheckOutCart(context.Context, *CheckOutCartRequest, *CheckOutCartResponse) error
//	}
type CartHandler struct {
	CartService services.CartService
}

func (c *CartHandler) Add(ctx context.Context, request *cart.AddCartRequest, response *cart.AddCartResponse) error {
	reqUserId, err := strconv.ParseInt(request.UserId, 10, 64)
	if err != nil {
		return err
	}

	reqSKUId, err := strconv.ParseInt(request.Skuid, 10, 64)
	if err != nil {
		return err
	}

	var cart = &models.Cart{
		UserId:    reqUserId,
		SKUId:     reqSKUId,
		Count:     request.Quantity,
		Timestamp: request.Time.AsTime(), // 要导入google的包："google.golang.org/protobuf/types/known/timestamppb"
	}

	// 调用service层的方法
	rowAffected, err := c.CartService.AddCart(cart)
	if rowAffected == 0 || err != nil {
		return err
	}

	return nil
}

func (c *CartHandler) FindAll(ctx context.Context, request *cart.FindAllCartRequest, response *cart.FindAllCartResponse) error {
	reqUserId, err := strconv.ParseInt(request.UserId, 10, 64)
	if err != nil {
		return err
	}

	carts, err := c.CartService.FindAll(reqUserId)
	if err != nil {
		return err
	}

	for _, v := range carts {
		response.CartItems = append(response.CartItems, &cart.CartItem{
			Skuid:    strconv.FormatInt(v.SKUId, 10),
			Quantity: v.Count,
			Time:     timestamppb.New(v.Timestamp),
		})
	}

	return nil
}

// CheckOutCart 清空购物车
func (c *CartHandler) CheckOutCart(ctx context.Context, request *cart.CheckOutCartRequest, response *cart.CheckOutCartResponse) error {
	reqUserId, err := strconv.ParseInt(request.UserId, 10, 64)
	if err != nil {
		return err
	}

	err = c.CartService.DeleteCart(reqUserId)

	return err
}
