package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	received   = "RECEIVED"
	inProgress = "IN_PROGRESS"
	shipped    = "SHIPPED"
	delivered  = "DELIVERED"
	cancelled  = "CANCELLED"
)

const (
	CREATE_ORDER_API = "CREATE_ORDER_API"
)

const (
	CUSTOMER_TABLENAME        = "customers"
	PODUCT_TABLENAME          = "products"
	ORDER_PRODUCTS_TABLE_NAME = "order_products"
	ORDERS_TABLENAME          = "orders"
)

type NewOrderReq struct {
	CustomerID uint      `json:"customer_id"`
	Products   []Product `json:"products"`
}

func (m *MarketPlaceAPIs) CreateOrder(c *gin.Context) {
	var newOrder NewOrderReq
	c.Bind(&newOrder)

	if !validID(newOrder.CustomerID, CUSTOMER_TABLENAME, m.DB) {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "Error": "Invalid Customer ID provided"})
		return
	}

	if len(newOrder.Products) == 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "Error": "Products cannot be empty. An order need to have atleast 1 product. Add a product and try again !"})
		return
	}

	order := &Order{
		CustomerID: newOrder.CustomerID,
		Status:     received,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result := m.DB.Create(&order)
	if result.Error != nil {
		apiErr(CREATE_ORDER_API, "unable to insert new order", result.Error)
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "Error": "Unable to insert order"})
		return
	}
	//Iterate over the request data 'products array' and for each product
	//in this new order, add an entry into the 'order_products' table.
	for _, product := range newOrder.Products {
		orderProduct := &OrderProduct{
			OrderID:    order.ID,
			ProductID:  product.ID,
			CustomerID: newOrder.CustomerID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		//Insert the order -> product mapping record in the 'order_products' table.
		productRes := m.DB.Create(orderProduct)
		if productRes.Error != nil {
			apiErr(CREATE_ORDER_API, "Add new order_product mapping failed in order_products table", productRes.Error)
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": http.StatusInternalServerError, "error": "Unable to insert the order in ordr product "})
			return

		}

	}

	c.JSON(http.StatusCreated,
		gin.H{"status": http.StatusCreated, "message": "Order Created Successfully!", "OrderID": order.ID})

}
