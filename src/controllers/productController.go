package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product

	database.DB.Find(&products)

	return c.JSON(products)
}
func CreateProducts(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}
	database.DB.Create(&product)

	return c.JSON(product)
}

func GetProducts(c *fiber.Ctx) error {
	var product models.Product

	id, _ := strconv.Atoi(c.Params("id"))

	product.Id = uint(id)

	database.DB.Find(&product)

	return c.JSON(product)

}

func UpdateProducts(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}

	product.Id = uint(id)
	if err := c.BodyParser(&product); err != nil {
		return err
	}
	database.DB.Model(&product).Updates(&product)
	return c.JSON(product)

}
func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}

	product.Id = uint(id)
	database.DB.Delete(&product)
	return nil
}

func ProductsFrontend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()
	result, err := database.Cache.Get(ctx, "products_frontend").Result()

	if err != nil {
		fmt.Println(err.Error())
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)
		if err != nil {
			panic(err)
		}
		if errKey := database.Cache.Set(ctx, "products_frontend", bytes, 30*time.Minute); err != nil {
			panic(errKey)
		}

	} else {
		json.Unmarshal([]byte(result), &products)
	}

	return c.JSON(products)
}
func ProductsBackend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()
	result, err := database.Cache.Get(ctx, "products_backend").Result()

	if err != nil {
		fmt.Println(err.Error())
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)
		if err != nil {
			panic(err)
		}
		if errKey := database.Cache.Set(ctx, "products_backend", bytes, 30*time.Minute); err != nil {
			panic(errKey)
		}

	} else {
		json.Unmarshal([]byte(result), &products)
	}

	var searchedProducts []models.Product

	if s := c.Query("s"); s != "" {
		lower := strings.ToLower(s)

		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lower) || strings.Contains(strings.ToLower(product.Description), lower) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		searchedProducts = products
	}
	if sortParam := c.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)

		if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	return c.JSON(searchedProducts)
}
