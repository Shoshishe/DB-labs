package main

import (
	"context"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//go:embed userid-name-email-createdat.csv
var usersCsv string

//go:embed categoryid-name.csv
var categoriesCsv string

//go:embed orderid-userid-createdat-status.csv
var userCategoriesCsv string

//go:embed orderitemid-orderid-productid-quantity-price.csv
var orderItemsCsv string

//go:embed productid-name-categoryid-price.csv
var productsCsv string

//go:embed orderid-userid-createdat-status.csv
var orderCsv string

type User struct {
	UserId    int       `bson:"user_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
}

type Category struct {
	CategoryId int    `bson:"category_id"`
	Name       string `bson:"name"`
}

type Product struct {
	ProductId  int     `bson:"product_id"`
	Name       string  `bson:"name"`
	CategoryId int     `bson:"category_id"`
	Price      float64 `bson:"price"`
}

type OrderItem struct {
	OrderItemId int     `bson:"order_item_id"`
	Quantity    int     `bson:"quantity"`
	Price       float64 `bson:"price"`
	Product     Product `bson:"product"`
}

type Order struct {
	OrderId   int         `bson:"order_id"`
	User      User        `bson:"owner"`
	Items     []OrderItem `bson:"items"`
	CreatedAt time.Time   `bson:"created_at"`
	Status    string      `bson:"status"`
}

func main() {
	client, _ := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	setupCategoriesCollection(client)
	setupUsersCollection(client)
	setupProductsCollection(client)
	setupCategoriesCollection(client)
	setupOrdersCollection(client)
}

func setupCategoriesCollection(client *mongo.Client) {
	err := client.Database("testing").CreateCollection(context.Background(), "categories")
	must(err)
	coll := client.Database("testing").Collection("categories")

	r := csv.NewReader(strings.NewReader(categoriesCsv))
	header, err := r.Read()
	must(err)
	headers := map[string]int{}
	for i, val := range header {
		headers[val] = i
	}

	categories := []Category{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		must(err)
		categoryId, err := strconv.Atoi(record[headers["category_id"]])
		must(err)
		name := record[headers["name"]]
		categories = append(categories, Category{CategoryId: categoryId, Name: name})
	}
	_, err = coll.InsertMany(context.Background(), categories)
	must(err)
}

func setupUsersCollection(client *mongo.Client) {
	err := client.Database("testing").CreateCollection(context.Background(), "users")
	must(err)
	coll := client.Database("testing").Collection("users")

	r := csv.NewReader(strings.NewReader(usersCsv))
	header, err := r.Read()
	must(err)
	headers := map[string]int{}
	for i, val := range header {
		headers[val] = i
	}
	users := []User{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		must(err)
		userId, err := strconv.Atoi(record[headers["user_id"]])
		must(err)
		name := record[headers["name"]]
		email := record[headers["email"]]
		createdAt, err := time.Parse("2006-01-02 15:04:05", record[headers["created_at"]])
		must(err)
		users = append(users, User{UserId: userId, Name: name, Email: email, CreatedAt: createdAt})
	}
	coll.InsertMany(context.Background(), users)
}

func setupProductsCollection(client *mongo.Client) {
	err := client.Database("testing").CreateCollection(context.Background(), "products")
	must(err)
	coll := client.Database("testing").Collection("products")

	r := csv.NewReader(strings.NewReader(productsCsv))
	header, err := r.Read()
	must(err)
	headers := map[string]int{}
	for i, val := range header {
		headers[val] = i
	}
	products := []Product{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		must(err)
		productId, err := strconv.Atoi(record[headers["product_id"]])
		must(err)
		name := record[headers["name"]]
		categoryId, err := strconv.Atoi(record[headers["category_id"]])
		must(err)
		price, err := strconv.ParseFloat(record[headers["price"]], 64)
		must(err)
		products = append(products, Product{ProductId: productId, Name: name, CategoryId: categoryId, Price: price})
	}
	coll.InsertMany(context.Background(), products)
}

func setupOrdersCollection(client *mongo.Client) {
	err := client.Database("testing").CreateCollection(context.Background(), "orders")
	must(err)
	coll := client.Database("testing").Collection("orders")
	users := client.Database("testing").Collection("users")
	products := client.Database("testing").Collection("products")

	r := csv.NewReader(strings.NewReader(orderCsv))
	header, err := r.Read()
	must(err)
	headers := map[string]int{}
	for i, val := range header {
		headers[val] = i
	}
	orders := map[int]*Order{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		must(err)
		orderId, err := strconv.Atoi(record[headers["order_id"]])
		must(err)
		userId, err := strconv.Atoi(record[headers["user_id"]])
		must(err)
		createdAt, err := time.Parse("2006-01-02 15:04:05", record[headers["created_at"]])
		must(err)
		status := record[headers["status"]]

		user := User{}
		err = users.FindOne(context.Background(), bson.D{{Key: "user_id", Value: fmt.Sprint(userId)}}).Decode(&user)
		must(err)

		orders[orderId] = &Order{
			OrderId:   orderId,
			CreatedAt: createdAt,
			Status:    status,
			User:      user,
		}
	}

	r = csv.NewReader(strings.NewReader(orderItemsCsv))
	header, err = r.Read()
	must(err)
	headers = map[string]int{}
	for i, val := range header {
		headers[val] = i
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		must(err)
		orderId, err := strconv.Atoi(record[headers["order_id"]])
		must(err)
		orderItemId, err := strconv.Atoi(record[headers["order_item_id"]])
		must(err)
		productId, err := strconv.Atoi(record[headers["product_id"]])
		must(err)
		quantity, err := strconv.ParseFloat(record[headers["quantity"]], 64)
		must(err)
		price, err := strconv.ParseFloat(record[headers["price"]], 64)
		must(err)

		product := Product{}
		err = products.FindOne(context.Background(), bson.D{{Key: "product_id", Value: fmt.Sprint(productId)}}).Decode(&product)
		must(err)

		orders[orderId].Items = append(orders[orderId].Items, OrderItem{Product: product, OrderItemId: orderItemId, Price: price, Quantity: int(quantity)})
	}
	coll.InsertMany(context.Background(), orders)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
