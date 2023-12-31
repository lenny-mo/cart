package dao

import (
	"strconv"

	"github.com/lenny-mo/cart/domain/models"

	"gorm.io/gorm"
)

// 接口驱动设计（Interface-Driven Design）
type CartDAOInterface interface {
	InitTable()
	CreateCart(*models.Cart) (int64, error)
	DeleteCart(int64) error
	UpdateCart(*models.Cart) (int64, error)
	FindCartByUserIDandSKUID(userid, skuid string) (*models.Cart, error)
	FindAll(int64) ([]*models.Cart, error)
	FindAllByUserIdForCheckout(int64) ([]*models.Cart, error)
}

type CartDAO struct {
	db *gorm.DB
}

func NewCartDAO(db *gorm.DB) CartDAOInterface {
	return &CartDAO{db: db}
}

func (c *CartDAO) InitTable() {
	if !c.db.Migrator().HasTable(&models.Cart{}) {
		c.db.Migrator().CreateTable(&models.Cart{})
	}
}

func (c *CartDAO) CreateCart(cart *models.Cart) (int64, error) {
	res := c.db.Create(cart)
	if res.Error != nil {
		return res.RowsAffected, res.Error
	}
	return res.RowsAffected, res.Error
}

// DeleteCart 删除购物车根据用户id
func (c *CartDAO) DeleteCart(userId int64) error {
	if err := c.db.Where("user_id = ?", userId).Delete(&models.Cart{}).Error; err != nil {
		return err
	}
	return nil
}

// 使用map更新，因为status默认零值，使用结构体会对零值不更新
func (c *CartDAO) UpdateCart(cart *models.Cart) (int64, error) {
	res := c.db.Model(&models.Cart{}).Where("user_id = ? AND sku_id = ?", cart.UserId, cart.SKUId).Updates(
		map[string]interface{}{"count": cart.Count, "status": cart.Status, "timestamp": cart.Timestamp})
	return res.RowsAffected, res.Error
}

func (c *CartDAO) FindCartByUserIDandSKUID(userid, skuid string) (*models.Cart, error) {
	cart := &models.Cart{}
	if err := c.db.Where(" user_id = ? AND sku_id = ?", userid, skuid).Find(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

func (c *CartDAO) FindAll(UserId int64) ([]*models.Cart, error) {
	var carts []*models.Cart
	res := c.db.Find(&carts, "user_id = ?", UserId)
	if res.Error != nil {
		return nil, res.Error
	}
	return carts, nil
}

func (c *CartDAO) FindAllByUserIdForCheckout(userid int64) ([]*models.Cart, error) {
	var carts []*models.Cart
	res := c.db.Where("user_id = ? and status = ?", strconv.FormatInt(userid, 10), "0").Find(&carts)
	if res.Error != nil {
		return nil, res.Error
	}
	return carts, nil
}
