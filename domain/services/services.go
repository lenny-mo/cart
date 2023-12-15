package services

import (
	"github.com/lenny-mo/cart/domain/dao"
	"github.com/lenny-mo/cart/domain/models"
)

type CartServiceInterface interface {
	AddCart(*models.Cart) (int64, error)
	DeleteCart(int64) error
	UpdateCart(*models.Cart) (int64, error)
	FindCartByUserIDandSKUID(userid, skuid string) (*models.Cart, error)
	FindAll(int64) ([]*models.Cart, error)
}

type CartService struct {
	CartDAO dao.CartDAO
}

func NewCartService(dao *dao.CartDAO) CartServiceInterface {
	return &CartService{CartDAO: *dao}
}

func (c *CartService) AddCart(cart *models.Cart) (int64, error) {
	return c.CartDAO.CreateCart(cart)
}

func (c *CartService) DeleteCart(userId int64) error {
	return c.CartDAO.DeleteCart(userId)
}

func (c *CartService) UpdateCart(cart *models.Cart) (int64, error) {
	return c.CartDAO.UpdateCart(cart)
}

func (c *CartService) FindCartByUserIDandSKUID(userid, skuid string) (*models.Cart, error) {
	return c.CartDAO.FindCartByUserIDandSKUID(userid, skuid)
}

func (c *CartService) FindAll(UserId int64) ([]*models.Cart, error) {
	return c.CartDAO.FindAll(UserId)
}
