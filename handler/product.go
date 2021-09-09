package handler

import (
	"fmt"
	"graphql/customer"
	"graphql/product"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerProduct struct {
	service product.ServiceProductInt
}

func NewProductHandler(service product.ServiceProductInt) *HandlerProduct {
	return &HandlerProduct{service: service}
}

func (h *HandlerProduct) GetProductByCategory(c *gin.Context) {
	category := c.Request.FormValue("category")

	products, err := h.service.GetProductCategory(category)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success", http.StatusOK, "products", products)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) CreateShopCart(c *gin.Context) {

	customer := c.MustGet("currentCustomer").(customer.Customer)

	err := h.service.AddShoppingCart(customer.ID)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success", http.StatusOK, fmt.Sprintf("Create cart for customerid %d", customer.ID), nil)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) InsertToShopCart(c *gin.Context) {

	idChart, err := strconv.Atoi((c.Request.FormValue("id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	productID, err := strconv.Atoi((c.Request.FormValue("product_id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	customer := c.MustGet("currentCustomer").(customer.Customer)

	product, err := h.service.InsertProductByCartID(productID, idChart)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusUnprocessableEntity, "failed", err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := APIResponse("success", http.StatusOK, fmt.Sprintf("Create shopecart for customer with id %d", customer.ID), product)
	c.JSON(http.StatusOK, response)

}

func (h *HandlerProduct) GetListProductShopCart(c *gin.Context) {

	cartID, err := strconv.Atoi((c.Request.FormValue("id")))
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	customer := c.MustGet("currentCustomer").(customer.Customer)

	products, err := h.service.GetListShopCart(cartID, customer.ID)
	if err != nil {
		response := APIResponse(err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	if len(products) == 0 {
		response := APIResponse(fmt.Sprintf("list product cart id : %d with customer id %d not found", cartID, customer.ID), http.StatusNotFound, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := APIResponse("success", http.StatusOK, fmt.Sprintf("list product shop cart with id customer : %d and cart_id : %d", customer.ID, cartID), products)
	c.JSON(http.StatusOK, response)

}
