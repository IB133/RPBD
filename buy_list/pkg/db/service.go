package db

import (
	"fmt"
	"log"
	"time"

	"github.com/IB133/RPBD/buy_list/pkg/config"
)

//go:generate moq -out querys_mock_test.go . Querys
type Querys interface {
	AddProductToBuyList(userId int, name string, weight string, time string) error
	AddProductToFridge(userId int, name string, date string) error
	GetBuyList(userId int) ([]BuyList, error)
	GetUserByUsername(username string) (Users, error)
	GetFridgeList(userId int) ([]Fridge, error)
	GetUsedProductsList(userId int) ([]Fridge, error)
	DeleteFromBuyList(userId int, name string)
	UpdateProductToCooked(userId int, name string, date string) error
	UpdateProductToDispose(userId int, name string, date string) error
	OpenProduct(userId int, name string, date string) error
	AddUser(name string, chatId int64) error
	GetStoredProductsList(userId int) ([]Fridge, error)
	GetStatsByDateDifference(userId int, firstDate string, secondDate string) ([]Fridge, error)
	GetBuyListForScheduler(userId int) ([]BuyList, error)
	GetFridgeListForScheduler(userId int) ([]Fridge, error)
	GetUsersList() ([]Users, error)
}

type DB struct {
	Conn *Connection
}

func (d *DB) AddToBuyList(username string, prodName string, weight string, date string, mes config.Config) string {
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	err = d.Conn.addProductToBuyList(u.Id, prodName, weight, date)
	if err != nil {
		return mes.Default
	}
	return mes.Succesfull
}

func (d *DB) GetBuyList(username string, mes config.Config) string {
	var str string
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := d.Conn.getBuyList(u.Id)
	if err != nil {
		return mes.Default
	}
	if len(list) == 0 {
		return mes.NoRows
	}
	for _, v := range list {
		str += fmt.Sprintf("%s  %.0f\n", v.Prod_name, v.Weight)
	}
	return fmt.Sprintf("Выберите продукт из списка\n %s", str)
}

func (d *DB) UsedProducts(username string, mes config.Config) string {
	var str string
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := d.Conn.getUsedProductsList(u.Id)
	if err != nil {
		return mes.Default
	}
	for _, v := range list {
		str += fmt.Sprintf("%s  %s\n", v.Prod_name, v.Status)
	}
	return str
}

func (d *DB) AddProductToFridgeFromBuyList(prodName string, username string, date string, mes config.Config) string {
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	if err = d.Conn.addProductToFridge(u.Id, prodName, date); err != nil {
		return mes.Default
	}
	if err = d.Conn.deleteFromBuyList(u.Id, prodName); err != nil {
		return mes.Default
	}
	return mes.Succesfull
}

func (d *DB) AddProductToFridge(prodName string, username string, date string, mes config.Config) string {
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	if err = d.Conn.addProductToFridge(u.Id, prodName, date); err != nil {
		return mes.Default
	}
	return mes.Succesfull
}

func (d *DB) StoredProductList(username string, mes config.Config) string {
	var str string
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := d.Conn.getStoredProductsList(u.Id)
	if err != nil {
		return mes.Default
	}
	if len(list) == 0 {
		return mes.NoRows
	}
	for _, v := range list {
		str += fmt.Sprintf("%s  %s %s\n", v.Prod_name, "хранится", v.Experitation_date.Format("2006-01-02"))
	}
	return fmt.Sprintf("%s\n %s", mes.ProductOpen, str)
}

func (d *DB) OpenProduct(username string, prodName string, newDate string, mes config.Config) string {
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	if err = d.Conn.openProduct(u.Id, prodName, newDate); err != nil {
		return mes.Default
	}
	return mes.Succesfull
}

func (d *DB) FridgeList(username string, mes config.Config) string {
	var str string
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := d.Conn.getFridgeList(u.Id)
	if err != nil {
		return mes.Default
	}
	if len(list) == 0 {
		return mes.NoRows
	}
	for _, v := range list {
		switch v.Status {
		case "stored":
			v.Status = "хранится"
		case "opened":
			v.Status = "открыт"
		case "used":
			v.Status = "приготовлен"
		case "dispose":
			v.Status = "выкинут"
		}
		str += fmt.Sprintf("%s  %s %s\n", v.Prod_name, v.Status, v.Experitation_date.Format("2006-01-02"))
	}
	return str
}

func (d *DB) ChangeStatus(username string, prodName string, status string, mes config.Config) string {
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	switch status {
	case "приготовлен":
		if err = d.Conn.updateProductToCooked(u.Id, prodName, time.Now().Format("2006-01-02")); err != nil {
			return mes.Default
		}
	case "выкинут":
		if err = d.Conn.updateProductToDispose(u.Id, prodName, time.Now().Format("2006-01-02")); err != nil {
			return mes.Default
		}
	}
	return mes.Succesfull

}

func (d *DB) UsedProcutList(username string, mes config.Config) string {
	var str string
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := d.Conn.getUsedProductsList(u.Id)
	if err != nil {
		return mes.Default
	}
	if len(list) == 0 {
		return mes.NoRows
	}
	for _, v := range list {
		switch v.Status {
		case "used":
			v.Status = "приготовлен"
		case "dispose":
			v.Status = "выкинут"
		}
		str += fmt.Sprintf("%s  %s \n", v.Prod_name, v.Status)
	}
	return str
}

func (d *DB) GetStats(username string, fDate string, sDate string, mes config.Config) string {
	var cookedCount int
	var disposeCount int
	u, err := d.Conn.getUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := d.Conn.getStatsByDateDifference(u.Id, fDate, sDate)
	if err != nil {
		return mes.Default
	}
	if len(list) == 0 {
		return mes.NoRows
	}
	for _, v := range list {
		if v.Status == "used" {
			cookedCount++
			continue
		}
		disposeCount++
	}
	return fmt.Sprintf("Количество приготовленных продуктов: %v\nКоличество выкинутых продуктов: %v", cookedCount, disposeCount)
}

func (d *DB) SchedulerBuyList(userId int) []BuyList {
	list, err := d.Conn.getBuyListForScheduler(userId)
	if err != nil {
		log.Println(err)
	}
	return list
}

func (d *DB) SchedulerFridge(userId int) []Fridge {
	list, err := d.Conn.getFridgeListForScheduler(userId)
	if err != nil {
		log.Println(err)
	}
	return list
}

func (d *DB) UsersList() []Users {
	users, err := d.Conn.getUsersList()
	if err != nil {
		log.Println(err)
	}
	return users
}

func (d *DB) AddUser(username string, chatId int64, mes config.Config) string {
	_, err := d.Conn.getUserByUsername(username)
	if err != nil {
		err = d.Conn.addUser(username, chatId)
		if err != nil {
			return mes.Start
		}
	}
	return mes.Start

}
