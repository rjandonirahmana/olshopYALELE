package product

import "errors"

type serviceProduct struct {
	repo RepoProduct
}

type ServiceProductInt interface {
	GetProductCategory(name_category string) ([]Product, error)
	AddShoppingCart(customerid int) (int, error)
	InsertProductByCartID(customerid, productid, cartid int) (Product, error)
	GetListInShopCart(cartid int, customerid int) ([]ShopeCart, error)
	DeleteListOnshoppingCart(cartid, customerid, productid int) ([]ShopeCart, error)
	GetShopCartCustomer(customerid int) ([]Cart, error)
}

func NewService(repo RepoProduct) *serviceProduct {
	return &serviceProduct{repo: repo}
}

func (s *serviceProduct) GetProductCategory(name_category string) ([]Product, error) {

	products, err := s.repo.GetProductByCategoryName(name_category)
	if err != nil {
		return []Product{}, err
	}

	if len(products) == 0 {
		return []Product{}, errors.New("products is not found")
	}

	return products, nil
}

func (s *serviceProduct) AddShoppingCart(customerid int) (int, error) {
	id, _ := s.repo.GetLastID()

	id += 1

	err := s.repo.CreateCart(customerid, id)
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (s *serviceProduct) InsertProductByCartID(customerid, productid, cartid int) (Product, error) {

	product, err := s.repo.GetProductByID(productid)
	if err != nil {
		return Product{}, err
	}

	cart, err := s.repo.GetShopCartIDCustomer(customerid, cartid)
	if err != nil {
		return Product{}, err
	}

	check, _ := s.repo.CheckInshopCart(cartid, product.Name)
	if check <= 0 {
		err := s.repo.InsertShoppingCart(cart.ID, productid, product.Price, product.Name)
		if err != nil {
			return Product{}, err

		}

		return product, nil
	}

	err = s.repo.IncreaseQuantityInshopCart(cartid, productid)
	if err != nil {
		return Product{}, err
	}

	return product, nil

}

func (s *serviceProduct) GetListInShopCart(cartid int, customerid int) ([]ShopeCart, error) {

	_, err := s.repo.GetShopCartIDCustomer(customerid, cartid)
	if err != nil {
		return []ShopeCart{}, err
	}

	shopcart, err := s.repo.GetListCartByID(cartid)
	if err == nil && len(shopcart) == 0 {
		return []ShopeCart{}, nil
	}
	if err != nil {
		return []ShopeCart{}, err
	}

	return shopcart, nil
}

func (s *serviceProduct) DeleteListOnshoppingCart(cartid, customerid, productid int) ([]ShopeCart, error) {
	err := s.repo.DeleteProductInShopCart(cartid, customerid, productid)
	if err != nil {
		return []ShopeCart{}, err
	}

	productLeft, err := s.GetListInShopCart(cartid, customerid)
	if err != nil {
		return []ShopeCart{}, err
	}

	return productLeft, nil
}

func (s *serviceProduct) GetShopCartCustomer(customerid int) ([]Cart, error) {

	cart, err := s.repo.ShopCartCustomer(customerid)
	if err != nil {
		return []Cart{}, err
	}

	return cart, nil
}
