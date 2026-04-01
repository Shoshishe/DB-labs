package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"maps"
	"math"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()
	setup()
	// fillRedis(rdb)
	// products := getByCategory(rdb, "Электроника")
	// print(products)
}

type Category = string

type Product struct {
	id       int      `redis:"id"`
	category Category `redis:"category"`
	price    float64  `redis:"price"`
	name     string   `redis:"name"`
}
type OrderId = int
type UserId = int
type CategoryId = int

type Order struct {
	status    string    `redis:"status"`
	createdAt time.Time `redis:"created_at"`
	userId    int       `redis:"user_id"`
	orderId   int       `redis:"order_id"`
	products  map[int]struct {
		product  Product
		quantity int
	} `redis:"products"`
}

var (
	userOrders       = map[UserId][]*Order{}
	orders           = map[OrderId]*Order{}
	products         = map[int]*Product{}
	categories       = map[int]Category{}
	categoryProducts = map[CategoryId][]*Product{}
)

func setup() {
	// userOrders := map[UserId][]*Order{}
	// orders := map[OrderId]*Order{}
	// products := map[int]*Product{}
	// categories := map[int]Category{}
	// categoryProducts := map[CategoryId][]*Product{}

	file, err := os.Open("../categoryid-name.csv")
	must(err)
	defer file.Close()
	r := csv.NewReader(file)
	header, err := r.Read()
	must(err)
	headers := map[string]int{}
	for i, val := range header {
		headers[val] = i
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		must(err)
		categoryId, err := strconv.Atoi(record[headers["category_id"]])
		must(err)
		categories[categoryId] = record[headers["name"]]
	}

	file, err = os.Open("../productid-name-categoryid-price.csv")
	must(err)
	defer file.Close()
	r = csv.NewReader(file)
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
		productId, err := strconv.Atoi(record[headers["product_id"]])
		must(err)
		categoryId, err := strconv.Atoi(record[headers["category_id"]])
		must(err)
		price, err := strconv.ParseFloat(record[headers["price"]], 64)
		products[productId] = &Product{id: productId, category: categories[categoryId], price: price, name: record[headers["name"]]}
		categoryProducts[categoryId] = append(categoryProducts[categoryId], products[productId])
	}

	file, err = os.Open("../orderid-userid-createdat-status.csv")
	must(err)
	defer file.Close()
	r = csv.NewReader(file)
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
		userId, err := strconv.Atoi(record[headers["user_id"]])
		must(err)
		createdAt, err := time.Parse("2006-01-02 15:04:05", record[headers["created_at"]])
		must(err)
		orders[orderId] = &Order{orderId: orderId, userId: userId, createdAt: createdAt, status: record[headers["status"]], products: map[int]struct {
			product  Product
			quantity int
		}{}}
		userOrders[userId] = append(userOrders[userId], orders[orderId])
	}

	file, err = os.Open("../orderitemid-orderid-productid-quantity-price.csv")
	must(err)
	defer file.Close()

	r = csv.NewReader(file)
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
		productId, err := strconv.Atoi(record[headers["product_id"]])
		must(err)
		quantity, err := strconv.Atoi(record[headers["quantity"]])
		must(err)

		curProduct := products[productId]
		orders[orderId].products[productId] = struct {
			product  Product
			quantity int
		}{product: *curProduct, quantity: quantity}
	}
}

func fillRedis(rdb *redis.Client) {
	file, err := os.Open("../userid-name-email-createdat.csv")
	must(err)
	defer file.Close()
	r := csv.NewReader(file)
	header, err := r.Read()
	must(err)
	headers := map[string]int{}
	for i, val := range header {
		headers[val] = i
	}

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

		rdb.HSet(context.Background(), fmt.Sprintf("users:%d", userId), []string{"name", name, "created_at",
			createdAt.Format(time.RFC3339), "email", email})
		var serializable []any
		for _, order := range userOrders[userId] {
			serializable = append(serializable, order.orderId)
		}
		rdb.RPush(context.Background(), fmt.Sprintf("user:%d:orders", userId), serializable...)
	}

	for _, order := range orders {
		rdb.HSet(context.Background(), fmt.Sprintf("order:%d", order.orderId), order)
	}

	for _, product := range products {
		rdb.HSet(context.Background(), fmt.Sprintf("product:%d", product.id), product)
	}

	for categoryId, products := range categoryProducts {
		var serializable []any
		for _, product := range products {
			serializable = append(serializable, product.id)
		}
		rdb.RPush(context.Background(), fmt.Sprintf("category:%s:products", categories[categoryId]), serializable...)
	}

	productsBought := map[int]float64{}
	for _, order := range orders {
		for _, product := range order.products {
			productsBought[product.product.id] = product.product.price * float64(product.quantity)
		}
	}
	sortedPairs := []struct {
		id   int
		sold float64
	}{}
	for k, v := range productsBought {
		sortedPairs = append(sortedPairs, struct {
			id   int
			sold float64
		}{id: k, sold: v})
	}

	slices.SortFunc(sortedPairs, func(a struct {
		id   int
		sold float64
	}, b struct {
		id   int
		sold float64
	}) int {
		return int(a.sold - b.sold)
	})

	for _, pair := range sortedPairs {
		rdb.ZAdd(context.Background(), "products:by_sales", redis.Z{Score: pair.sold, Member: pair.id})
	}

	productsVal := slices.Collect(maps.Values(products))
	slices.SortFunc(productsVal, func(a *Product, b *Product) int {
		return int(a.price - b.price)
	})
	for _, product := range products {
		rdb.ZAdd(context.Background(), "products:by_price", redis.Z{Score: product.price, Member: product.id})
	}

	for userId, userOrders := range userOrders {
		for _, order := range userOrders {
			for _, product := range order.products {
				rdb.RPush(context.Background(), fmt.Sprintf("user:%d:purchased", userId), product.product.id)
			}
		}
	}
}

func getLatestOrders(rdb *redis.Client, userId int) []Order {
	orderIds, err := rdb.LRange(context.Background(), fmt.Sprintf("user:%d:orders", userId), -10, -1).Result()
	must(err)
	returned := []Order{}
	for _, orderId := range orderIds {
		id, err := strconv.Atoi(orderId)
		must(err)
		returned = append(returned, *orders[id])
	}
	return returned
}

func getTopProductsByRevenue(rdb *redis.Client) []Product {
	productIds, err := rdb.ZRangeArgs(context.Background(), redis.ZRangeArgs{Key: "products:by_sales", Start: -10, Stop: -1}).Result()
	must(err)
	returned := []Product{}
	for _, productId := range productIds {
		id, err := strconv.Atoi(productId)
		must(err)
		returned = append(returned, *products[id])
	}
	return returned
}

func getTopNProductsByRevenue(rdb *redis.Client, N int) []Product {
	count, err := rdb.Exists(context.Background(), fmt.Sprintf("cache:top_products:by_sales:%d", N)).Result()
	if count != 0 {
		productIds, err := rdb.LRange(context.Background(), fmt.Sprintf("cache:top_products:by_sales:%d", N), int64(-N), -1).Result()
		if err == nil {
			returned := []Product{}
			for _, productId := range productIds {
				id, err := strconv.Atoi(productId)
				must(err)
				returned = append(returned, *products[id])
			}
			return returned
		}
	}
	productIds, err := rdb.ZRangeArgs(context.Background(), redis.ZRangeArgs{Key: "products:by_sales", Start: 0, Stop: math.Inf(1), ByScore: true, Count: 12}).Result()
	must(err)
	returned := []Product{}
	for _, productId := range productIds {
		id, err := strconv.Atoi(productId)
		must(err)
		returned = append(returned, *products[id])
	}
	var appended []any
	for _, returned := range returned {
		appended = append(appended, returned.id)
	}
	pipe := rdb.Pipeline()
	pipe.LPush(context.Background(), fmt.Sprintf("cache:top_products:by_sales:%d", N), appended...)
	pipe.Expire(context.Background(), fmt.Sprintf("cache:top_products:by_sales:%d", N), time.Minute)
	_, err = pipe.Exec(context.Background())
	must(err)
	return returned
}

func getByCategory(rdb *redis.Client, category Category) []Product {
	productIds, err := rdb.LRange(context.Background(), fmt.Sprintf("category:%s:products", category), 0, -1).Result()
	must(err)
	returned := []Product{}
	for _, productId := range productIds {
		id, err := strconv.Atoi(productId)
		must(err)
		returned = append(returned, *products[id])
	}
	return returned
}

func getByPrice(rdb *redis.Client, price float64) []Product {
	productIds, err := rdb.ZRangeArgs(context.Background(), redis.ZRangeArgs{Key: "products:by_price", ByScore: true, Start: price, Stop: math.Inf(1)}).Result()
	must(err)
	returned := []Product{}
	for _, productId := range productIds {
		id, err := strconv.Atoi(productId)
		must(err)
		returned = append(returned, *products[id])
	}
	return returned
}

func must(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
