package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lampadovnikita/StorekeeperTask/pkg/config"
	"github.com/lampadovnikita/StorekeeperTask/pkg/data"
	"github.com/lampadovnikita/StorekeeperTask/pkg/database"
)

var notEnoughArgsError = errors.New("Not enough arguments provided")
var invalidArgsError = errors.New("The arguments are invalid")

func main() {

	orderIDs, err := getOrderIDs()
	if err == notEnoughArgsError {
		fmt.Println("Please provide al least one non-negative number as an argument!")
		return
	} else if err == invalidArgsError {
		fmt.Println("Invalid Arguments! There is no non-negative number among the arguments!")
		return
	} else if err != nil {
		panic(err)
	}

	pgcfg, err := config.GetPGConfig()
	if err != nil {
		panic(err)
	}

	pgxpool, err := database.NewPGXPool(context.Background(), pgcfg)
	if err != nil {
		panic(err)
	}

	pgstorage := database.NewPGStorage(pgxpool)

	gatheringInfo, err := pgstorage.GetGatheringInfo(orderIDs)
	if err != nil {
		panic(err)
	}

	printGatheringInfo(orderIDs, gatheringInfo)
}

func printGatheringInfo(orderIDs []int, info []data.GatheringInfo) {
	fmt.Println("=+=+=+=")

	ordersStr := strings.Replace(strings.Trim(fmt.Sprint(orderIDs), "[]"), " ", ",", -1)
	fmt.Println("Страница сборки заказов", ordersStr)
	fmt.Println()

	var prevRackName string
	for _, elem := range info {
		if prevRackName != elem.RackName {
			prevRackName = elem.RackName

			fmt.Println("===Стеллаж", elem.RackName)
		}

		fmt.Printf("%s (id=%d)\n", elem.ProductName, elem.ProductID)
		fmt.Printf("заказ %d, %d шт\n", elem.OrderID, elem.Amount)

		if len(elem.AdditionalRacks) > 0 {
			fmt.Println("доп стеллаж:", strings.Join(elem.AdditionalRacks, ","))
		}

		fmt.Println()
	}
}

func getOrderIDs() ([]int, error) {
	arguments := os.Args
	if len(arguments) == 1 {
		return nil, notEnoughArgsError
	}

	orderIDs := make([]int, 0)
	for i, arg := range arguments {
		if i == 0 {
			continue
		}

		orderID, err := strconv.Atoi(arg)
		if (err != nil) || (orderID < 0) {
			continue
		}

		orderIDs = append(orderIDs, orderID)
	}

	if len(orderIDs) == 0 {
		return nil, invalidArgsError
	}

	return orderIDs, nil
}
