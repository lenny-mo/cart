package dao

import (
	"cart/domain/models"

	"gorm.io/gorm"
)

// 接口驱动设计（Interface-Driven Design）
type CartDAOInterface interface {
	InitTable()
	CreateCart(*models.Cart) (int64, error)
	DeleteCart(int64) error
	UpdateCart(*models.Cart) (int64, error)
	FindCartByID(Id int64) (*models.Cart, error)
	FindAll() ([]*models.Cart, error)
}

type CartDAO struct {
	db *gorm.DB
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

func (c *CartDAO) DeleteCart(Id int64) error {
	if err := c.db.Delete(&models.Cart{}, Id).Error; err != nil {
		return err
	}
	return nil
}

func (c *CartDAO) UpdateCart(cart *models.Cart) (int64, error) {
	res := c.db.Save(cart) // 如果没有包含主键，会调用Create方法 否则调用Update方法
	if err := res.Error; err != nil {
		return res.RowsAffected, err
	}
	return res.RowsAffected, res.Error
}

func (c *CartDAO) FindCartByID(Id int64) (*models.Cart, error) {
	cart := &models.Cart{}
	if err := c.db.First(cart, Id).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

func (c *CartDAO) FindAll() ([]*models.Cart, error) {
	var carts []*models.Cart
	if err := c.db.Find(&carts).Error; err != nil {
		return nil, err
	}
	return carts, nil
}
