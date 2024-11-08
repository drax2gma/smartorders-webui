package db

import (
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func CloseRedisClient() {
	RedisClient.Close()
}

func GetCustomer(username string) (Customer, error) {
	return Customer{}, nil
}

func AddCustomer(user Customer) error {
	return nil
}

func DeleteCustomer(username string) error {
	return nil
}

func UpdateCustomer(user Customer) error {
	return nil
}

func GetOrders() ([]Item, error) {
	return []Item{}, nil
}

func AddOrder(order Item) error {
	return nil
}

func DeleteOrder(id string) error {
	return nil
}

func UpdateOrder(order Item) error {
	return nil
}

func GetOrdersByCustomer(username string) ([]Item, error) {
	return []Item{}, nil
}

func GetOrder(id string) (Order, error) {
	return Order{}, nil
}

func AddOrderItem(orderItem OrderItem) error {
	return nil
}

func DeleteOrderItem(id string) error {
	return nil
}

func UpdateOrderItem(orderItem OrderItem) error {
	return nil
}

func GetOrderItemsByOrder(id string) ([]OrderItem, error) {
	return []OrderItem{}, nil
}

func GetOrderItemsByUser(username string) ([]OrderItem, error) {
	return []OrderItem{}, nil
}
