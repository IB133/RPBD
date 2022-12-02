package db

import (
	"fmt"
	"time"

	"github.com/IB133/RPBD/buy_list/pkg/config"
)

//go:generate moq -out querys_test.go . Querys
type Querys interface {
	GetBuyList(username string, s Connection) (string, error)
}

func GetBuyList(username string, s Connection) (string, error) {
	var str string
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	list, err := s.GetBuyList(u.Id)
	if err != nil {
		return "", err
	}
	for _, v := range list {
		str += fmt.Sprintf("%s  %.0f\n", v.Prod_name, v.Weight)
	}
	return str, nil
}

func UsedProducts(username string, s Connection) (string, error) {
	var str string
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	list, err := s.GetUsedProductsList(u.Id)
	if err != nil {
		return "", err
	}
	for _, v := range list {
		str += fmt.Sprintf("%s  %s\n", v.Prod_name, v.Status)
	}
	return str, nil
}

func AddProductToFridgeFromBuyList(prodName string, username string, date string, s Connection) error {
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if err = s.AddProductToFridge(u.Id, prodName, date); err != nil {
		return err
	}
	if err = s.DeleteFromBuyList(u.Id, prodName); err != nil {
		return err
	}
	return nil
}

func AddProductToFridge(prodName string, username string, date string, s Connection) error {
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if err = s.AddProductToFridge(u.Id, prodName, date); err != nil {
		return err
	}
	return nil
}

func StoredProductList(username string, s Connection, mes config.Config) string {
	var str string
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := s.GetStoredProductsList(u.Id)
	if err != nil {
		return err.Error()
	}
	for _, v := range list {
		str += fmt.Sprintf("%s  %s %s\n", v.Prod_name, "хранится", v.Experitation_date.Format("2006-01-02"))
	}
	return fmt.Sprintf("Выберите продукт из списка\n %s", str)
}

func OpenProduct(username string, prodName string, newDate string, s Connection, mes config.Config) string {
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	if err = s.OpenProduct(u.Id, prodName, newDate); err != nil {
		return err.Error()
	}
	return mes.Succesfull
}

func FridgeList(username string, s Connection, mes config.Config) string {
	var str string
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := s.GetFridgeList(u.Id)
	if err != nil {
		return err.Error()
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

func ChangeStatus(username string, prodName string, status string, s Connection, mes config.Config) string {
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	switch status {
	case "приготовлен":
		if err = s.UpdateProductToCooked(u.Id, prodName, time.Now().Format("2006-01-02")); err != nil {
			return err.Error()
		}
	case "выкинут":
		if err = s.UpdateProductToDispose(u.Id, prodName, time.Now().Format("2006-01-02")); err != nil {
			return err.Error()
		}
	}
	return mes.Succesfull

}

func UsedProcutList(username string, s Connection, mes config.Config) string {
	var str string
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := s.GetUsedProductsList(u.Id)
	if err != nil {
		return err.Error()
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

func GetStats(username string, fDate string, sDate string, s Connection, mes config.Config) string {
	var cookedCount int
	var disposeCount int
	u, err := s.GetUserByUsername(username)
	if err != nil {
		return mes.UserNotFound
	}
	list, err := s.GetStatsByDateDifference(u.Id, fDate, sDate)
	if err != nil {
		return err.Error()
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
