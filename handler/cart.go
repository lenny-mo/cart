package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/lenny-mo/order/proto/order"
	"github.com/lenny-mo/order/utils"

	"github.com/lenny-mo/cart/domain/models"
	"github.com/lenny-mo/cart/domain/services"
	"github.com/lenny-mo/cart/global"
	"github.com/lenny-mo/cart/proto/cart"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// 这个handler 要实现micro server rpc 的接口
//
//	type CartHandler interface {
//		// 加入购物车，status设置为0 表示等待支付状态
//		Add(context.Context, *AddCartRequest, *AddCartResponse) error
//		// 更新数量
//		Update(context.Context, *UpdateRequest, *UpdateResponse) error
//		// 从购物车中删除某个商品，把这个商品的status设置为2，表示discard
//		Delete(context.Context, *DeleteRequest, *DeleteResponse) error
//		FindAll(context.Context, *FindAllCartRequest, *FindAllCartResponse) error
//		// checkout cart 下单，需要创建一条order表的记录并且插入到order表中
//		// 把所有这个user的商品信息搜集起来，并且计算总金额， 在这个过程中需要把每个商品的status标记为1，也就是“被结算”状态
//		// 写一条记录，这个记录是要插入到order表中，并且等待支付
//		CheckOutCart(context.Context, *CheckOutCartRequest, *CheckOutCartResponse) error
//	}
type CartHandler struct {
	CartService services.CartService
}

func (c *CartHandler) Add(ctx context.Context, request *cart.AddCartRequest, response *cart.AddCartResponse) error {
	// 获取userid
	reqUserId, err := strconv.ParseInt(request.UserId, 10, 64)
	if err != nil {
		return err
	}

	// 获取skuid
	reqSKUId, err := strconv.ParseInt(request.Item.Skuid, 10, 64)
	if err != nil {
		return err
	}

	// 构造用于存储到存储层的数据
	var cart = &models.Cart{
		UserId:    reqUserId,
		SKUId:     reqSKUId,
		Count:     request.Item.Quantity,
		Timestamp: request.Item.Time.AsTime(), // 要导入google的包："google.golang.org/protobuf/types/known/timestamppb"
		Status:    int8(request.Item.Status),
	}

	// 调用service层的方法
	rowAffected, err := c.CartService.AddCart(cart)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if rowAffected == 0 {
		return errors.New("rowAffected == 0")
	}

	return nil
}

func (c *CartHandler) FindAll(ctx context.Context, req *cart.FindAllCartRequest, res *cart.FindAllCartResponse) error {
	reqUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	carts, err := c.CartService.FindAll(reqUserId)
	if err != nil {
		return err
	}

	for _, v := range carts {
		res.CartItems = append(res.CartItems, &cart.CartItem{
			Skuid:    strconv.FormatInt(v.SKUId, 10),
			Quantity: v.Count,
			Time:     timestamppb.New(v.Timestamp),
			Status:   cart.CartStatus(v.Status),
		})
	}
	res.Msg = "success"

	return nil
}

func (c *CartHandler) Update(ctx context.Context, req *cart.UpdateRequest, res *cart.UpdateResponse) error {
	// 根据rpc传递的参数item 构造用于数据库更新的model
	// 获取userid
	uid, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		fmt.Println(err)
		return err
	}
	skuid, err := strconv.ParseInt(req.Item.Skuid, 10, 64)
	if err != nil {
		fmt.Println(err)
		return err
	}
	data := &models.Cart{
		UserId:    uid,
		SKUId:     skuid,
		Count:     req.Item.Quantity,
		Timestamp: req.Item.Time.AsTime(),
		Status:    int8(req.Item.Status),
	}
	rowaffected, err := c.CartService.UpdateCart(data)
	if err != nil {
		fmt.Println(err)
		res.Msg = err.Error()
		return err
	}
	if rowaffected == 0 {
		res.Msg = "rowaffected == 0"
		return errors.New("rowaffected == 0")
	}
	res.Msg = "success"
	return nil
}

// 从购物车中删除某个商品，把这个商品的status设置为2，表示discard
func (c *CartHandler) Delete(ctx context.Context, req *cart.DeleteRequest, res *cart.DeleteResponse) error {
	// 找到这条记录
	data, err := c.CartService.FindCartByUserIDandSKUID(req.Userid, req.Skuid)
	if err != nil {
		res.Msg = "删除失败"
		return err
	}
	// 修改data的status字段为discard
	data.Status = 2
	// 调用update方法
	rowaffected, err := c.CartService.UpdateCart(data)
	if err != nil || rowaffected == 0 {
		fmt.Println("更新失败")
		return err
	}
	return nil
}

// CheckOutCart 清空购物车
func (c *CartHandler) CheckOutCart(ctx context.Context, req *cart.CheckOutCartRequest, res *cart.CheckOutCartResponse) error {
	reqUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 1. 先去cart 表根据userid, status= 0 的商品信息搜集item slice，在这个过程中需要把每个商品的status标记为1，也就是“被结算”状态
	items, err := c.CartService.FindAllByUserIdForCheckout(reqUserId)

	// 2. 把item序列化
	bytes, err := json.Marshal(items)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 3. 调用order的create rpc请求, 生成一条order记录 插入到order表中
	client := order.NewOrderService("go.micro.service.order", global.GetGlobalRPCService().Client())
	orderRes, err := client.InsertOrder(context.TODO(), &order.InserRequest{
		OrderData: &order.OrderInfo{
			OrderId:      utils.UUID(),
			OrderVersion: 1,
			UserId:       reqUserId,
			OrderData:    string(bytes),
			Status:       order.OrderStatus_UNPAID,
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	res.Msg = "order row affected: " + strconv.Itoa(int(orderRes.RowsAffected))
	return nil
}
