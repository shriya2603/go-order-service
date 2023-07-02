# Order Service

Order service is services which has a APIs to create , update, update order status of an order. 

## Features 
- Create Order: Lets the user to create a order with help of product and customer details 
- Update Order Status: Lets the user to update the status of the order 
- Update Order product : Lets the user to update the product id 
- Observality : Added api hit metric, elapsed time of APIs metric 

## Observality 
- Grafana Screenshot
![Grafana Screen shot ](screenshots/grafana_screenshot.png)

- APIs Hit metrics Screenshots
![APIs Hit metrics](screenshots/Apis_hit_metrics.png)


### Reference : 
- https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
- https://promlabs.com/promql-cheat-sheet/

## Sample API s 
http://localhost:8080/v1/order/create

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
http://localhost:8080/v1/order/1

{
    "status":"DELIEVERED"
}


http://localhost:8080/v1/order/updateOrderProducts/7
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




## TODO 
- [ ] Add delete api 
- [ ] Add any update to created order status in kafka topic 
- [x] Use Grafana dashboard to visualize all the prometheus metrics 
    - For uptime metric : uptime > 0
    - For Api counter metric : rate(<api_name>[$__rate_interval])
    - For elapse time metric : order_service_apis_elapsed_time with label as apiname 
- [ ] Add a docker compose file 