package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	Id       string
	Chat_id  int
	Username string
}

type BuyList struct {
	Weight      float32
	ProductName string
	BuyTime     string
}

type Fridge struct {
	Status             string
	ProductName        string
	ExperetitationDate string
}

type Connection struct {
	conn *sqlx.DB
}

func NewConnect() (*Connection, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/?sslmode=disable", os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBHOST"), os.Getenv("DBPORT"))
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal()
	}
	return &Connection{
		conn: conn,
	}, nil
}

func (c *Connection) AddProductToBuyList(user_id string, name string, weight float32, time int) error {
	_, err := c.conn.ExecContext(context.Background(), `
	INSERT INTO buy_list(user_id, prod_name, weight, buy_time)
	VALUES($1, $2, $3, $4);
	`, user_id, name, weight, time)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) AddProductToFridge(user_id string, name string, date string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	INSERT INTO fridge(user_id, prod_name, experitation_date)
	VALUES($1, $2, $3);
	`, user_id, name, date)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) GetProductsFromBuyList(user_id string) ([]BuyList, error) {
	var buyList []BuyList
	err := c.conn.SelectContext(context.Background(), &buyList, `
	SELECT prod_name, weight, buy_time
	FROM buy_list
	WHERE user_id = $1;
	`, user_id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return buyList, nil
}

func (c *Connection) GetUserByUsername(username string) (User, error) {
	u := User{}
	err := c.conn.GetContext(context.Background(), &u, "SELECT id FROM user WHERE username = $1", username)
	if err != nil {
		return u, fmt.Errorf("people not found: %w", err)
	}
	return u, nil
}

func (c *Connection) GetProductsList(user_id string) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT prod_name, status, experitation_date
	FROM fridge
	WHERE user_id = $1
	ORDER BUY product_name DESC;
	`, user_id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return fridge, nil
}

func (c *Connection) GetUsedProductsList(user_id string) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT prod_name, status
	FROM fridge
	WHERE user_id = $1 AND (status = 'used' OR status = 'dispose');
	`, user_id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return fridge, nil
}

func (c *Connection) DeleteFromBuyList(user_id string, name string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	DELETE FROM buy_list
	WHERE user_id = $1 AND prod_name = $2;
	`, user_id, name)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) UpdateProductToCoocked(user_id string, name string, date string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	UPDATE fridge 
	SET status = 'used',  status_date= $1
	WHERE prod_name = $2 AND user_id = $3;
	`, date, name, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) UpdateProductToDispose(user_id string, name string, date string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	UPDATE fridge 
	SET status = 'dispose',  status_date= $1
	WHERE prod_name = $2 AND user_id = $3;
	`, date, name, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) OpenProduct(user_id string, name string, date string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	UPDATE fridge 
	SET status_date= $1
	WHERE prod_name = $2 AND user_id = $3;
	`, date, name, user_id)
	if err != nil {
		return err
	}
	return nil
}
