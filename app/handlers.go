package app

import (
	"fmt"
	"net/http"
	"strconv"
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
type UpdateOrderStatusReq struct {
	Status string `json:"status"`
}

func (m *MarketPlaceAPIs) CreateOrder(c *gin.Context) {
	var newOrder NewOrderReq
	c.Bind(&newOrder)

	if !validID(newOrder.CustomerID, CUSTOMER_TABLENAME, m.DB) {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Invalid Customer ID provided"})
		return
	}

	if len(newOrder.Products) == 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Products cannot be empty. An order need to have atleast 1 product. Add a product and try again !"})
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
			gin.H{"status": http.StatusInternalServerError, "error": "Unable to insert order"})
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

func (m *MarketPlaceAPIs) GetOrder(c *gin.Context) {
	orderIdStr := c.Params.ByName("id")

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Order ID passed is not a valid number."})
		return
	}

	if !validID(uint(orderId), ORDERS_TABLENAME, m.DB) {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "error": "Invalid Order ID provided"})
		return
	}

	var orderProducts []OrderProduct

	query := fmt.Sprintf("SELECT * FROM %s where order_id=%d;", ORDER_PRODUCTS_TABLE_NAME, orderId)
	m.DB.Raw(query).Scan(&orderProducts)
	if len(orderProducts) == 0 {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "error": "No order with requested ID exists in the table. Invalid ID."})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{"status": http.StatusOK, "message": "Order Details Fetched Successfully!", "order": orderProducts})
}

func (m *MarketPlaceAPIs) UpdateOrderedProducts(c *gin.Context) {
	orderIdStr := c.Params.ByName("id")

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Order ID passed is not a valid number."})
		return
	}

	if !validID(uint(orderId), ORDERS_TABLENAME, m.DB) {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "error": "Invalid Order ID provided"})
		return
	}

	var updateOrderProducts NewOrderReq
	c.Bind(&updateOrderProducts)
	var failedProductUpdate []uint

	var orderProducts []OrderProduct

	query := fmt.Sprintf("SELECT * FROM %s where order_id=%d ORDER BY ID;", ORDER_PRODUCTS_TABLE_NAME, orderId)
	result := m.DB.Raw(query).Scan(&orderProducts)
	if len(orderProducts) == 0 || result.Error != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "error": "No order with requested ID exists in the table. Invalid ID."})
		return
	} else if len(orderProducts) != len(updateOrderProducts.Products) {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "v1 Update API supports only product details updation in the existing order. New products addition and existing products deletion from existing order will be supported in future API ver"})
		return
	} else {
		for i, orderProduct := range orderProducts {
			product := updateOrderProducts.Products[i]
			orderProduct.CustomerID = updateOrderProducts.CustomerID
			orderProduct.ProductID = product.ID
			orderProduct.UpdatedAt = time.Now()
			qResult := m.DB.Model(OrderProduct{}).Where("id = ?", orderProduct.ID).Updates(orderProduct)
			if qResult.Error != nil {
				failedProductUpdate = append(failedProductUpdate, orderProduct.ProductID)
			}
		}
	}

	if len(failedProductUpdate) > 0 {
		c.JSON(http.StatusBadGateway,
			gin.H{"status": http.StatusBadGateway, "message": fmt.Sprintf("Unableto update products with ID %v ", failedProductUpdate)})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{"status": http.StatusOK, "message": "Order Updated Successfully!"})

}

func (m *MarketPlaceAPIs) UpdateOrderStatus(c *gin.Context) {
	orderIdStr := c.Params.ByName("id")

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "error": "Order ID passed is not a valid number."})
		return
	}

	if !validID(uint(orderId), ORDERS_TABLENAME, m.DB) {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "error": "Invalid Order ID provided"})
		return
	}

	var order Order

	query := fmt.Sprintf("SELECT * FROM %s where id=%d;", ORDERS_TABLENAME, orderId)
	m.DB.Raw(query).First(&order)
	if (Order{}) == order {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "error": "No order with requested ID exists in the table. Invalid ID."})
		return
	}

	var updateOrder UpdateOrderStatusReq
	c.Bind(&updateOrder)

	result := m.DB.Model(&Order{}).Where("id = ?", orderId).Update("status", updateOrder.Status)
	if result.RowsAffected != 1 {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "message": "Order not updated!"})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{"status": http.StatusOK, "message": "Order Updated Successfully!"})

}
