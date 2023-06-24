package app

import (
	"fmt"
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
	CreateOrderAPI = "CREATE_ORDER_API"
)

type NewOrderReq struct {
	CustomerID uint      `json:"customer_id"`
	Products   []Product `json:"products"`
}

func (m *MarketPlaceAPIs) CreateOrder(c *gin.Context) {
	var newOrder NewOrderReq
	c.Bind(&newOrder)

	order := &Order{
		CustomerID: newOrder.CustomerID,
		Status:     received,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result := m.DB.Create(&order)
	if result.Error != nil {
		apiErr(CreateOrderAPI, "unable to insert new order", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "Error": "Unable to insert order"})
		return
	}
	fmt.Println(" order id ", order.ID)
	//Iterate over the request data 'products array' and for each product
	//in this new order, add an entry into the 'order_products' table.
	for _, product := range newOrder.Products {
		fmt.Println("product ", product.ID)
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
			apiErr(CreateOrderAPI, "Add new order_product mapping failed in order_products table", productRes.Error)
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": http.StatusInternalServerError, "error": "Unable to insert the order in ordr product "})
			return

		}

	} //End of for-loop

	c.JSON(http.StatusCreated,
		gin.H{"status": http.StatusCreated, "message": "Order Created Successfully!", "resourceId": order.ID})

}
