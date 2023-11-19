package services

import (
	"cart/domain/dao"
	"cart/domain/models"
)

type CartServiceInterface interface {
	AddCart(*models.Cart) (int64, error)
	DeleteCart(int64) error
	UpdateCart(*models.Cart) (int64, error)
	FindCartByID(int64) (*models.Cart, error)
	FindAll() ([]*models.Cart, error)
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

func (c *CartService) DeleteCart(Id int64) error {
	return c.CartDAO.DeleteCart(Id)
}

func (c *CartService) UpdateCart(cart *models.Cart) (int64, error) {
	return c.CartDAO.UpdateCart(cart)
}

func (c *CartService) FindCartByID(Id int64) (*models.Cart, error) {
	return c.CartDAO.FindCartByID(Id)
}

func (c *CartService) FindAll() ([]*models.Cart, error) {
	return c.CartDAO.FindAll()
}

