# go-order-service



Sample API s 
http://localhost:8080/sv1/order/create

{
    "customer_id": 1,
    "products": [
        {
            "id": 1
        },
        {
            "id": 2
        }
    ]
}



Get Order 
http://localhost:8080/sv1/order/7



Update Status of an order

{
    "status":"DELIEVERED"
}


http://localhost:8080/sv1/order/updateOrderProducts/7
{
    "customer_id": 1,
    "products": [
        {
            "id": 3
        },
        {
            "id": 4
        }
    ]
}