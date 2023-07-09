package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
)

type Order struct {
	Products []Product
	Customer Customer
}

type Product struct {
	Name     string
	Price    float32
	Quantity int
}

type Customer struct {
	Name string
}

func generateOrder(orderChan chan Order, wg *sync.WaitGroup, productsNum uint) {
	orderChan <- Order{
		Products: generateProducts(productsNum),
		Customer: generateCustomer(),
	}
	wg.Done()
}

func generateCustomer() Customer {
	return Customer{
		Name: "Customer_" + string(rune(rand.Intn(100))),
	}
}

func generateProducts(pNum uint) []Product {
	var products []Product

	for i := 0; i < int(pNum); i++ {
		products = append(products, Product{
			Name:     fmt.Sprintf("Customer_%d", i),
			Price:    float32(rand.Intn(100)) + 1,
			Quantity: rand.Intn(10) + 1,
		})
	}

	return products
}

func totalPriceCalc(totalPriceChan chan string, orderChan chan Order, wg *sync.WaitGroup) {
	var totalPrice float32

	order := <-orderChan

	for _, val := range order.Products {
		totalPrice += val.Price
	}

	totalPriceChan <- fmt.Sprintf("Total sum for order: %.2f", totalPrice)
	wg.Done()
}

func main() {
	orderChannel := make(chan Order)
	totalPriceChannel := make(chan string)
	defer close(orderChannel)
	defer close(totalPriceChannel)

	num := flag.Uint("num_of_products", 1, "The number of products which the program should generate")
	flag.Parse()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go generateOrder(orderChannel, &wg, *num)

	wg.Add(1)
	go totalPriceCalc(totalPriceChannel, orderChannel, &wg)

	totalPrice := <-totalPriceChannel
	fmt.Println(totalPrice)

	wg.Wait()
}
