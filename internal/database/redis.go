package main

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

func GetItems() ([]Item, error) {
	return []Item{}, nil
}

func AddItem(item Item) error {
	return nil
}

func DeleteItem(id string) error {
	return nil
}

func UpdateItem(item Item) error {
	return nil
}

func GetUser(username string) (User, error) {
	return User{}, nil
}

func AddUser(user User) error {
	return nil
}

func DeleteUser(username string) error {
	return nil
}

func UpdateUser(user User) error {
	return nil
}

func GetOrders() ([]Order, error) {
	return []Order{}, nil
}

func AddOrder(order Order) error {
	return nil
}

func DeleteOrder(id string) error {
	return nil
}

func UpdateOrder(order Order) error {
	return nil
}

func GetOrdersByUser(username string) ([]Order, error) {
	return []Order{}, nil
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
