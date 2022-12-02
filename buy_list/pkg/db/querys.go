package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	Id       string
	Chat_id  int
	Username string
}

type BuyList struct {
	Weight    float32
	Prod_name string
	Buy_time  string
}

type Fridge struct {
	Status            string
	Prod_name         string
	Experitation_date time.Time
	Status_date       time.Time
}

type Connection struct {
	conn *sqlx.DB
}

func NewConnect(connString string) (*Connection, error) {
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal()
	}
	return &Connection{
		conn: conn,
	}, nil
}

func (c *Connection) AddProductToBuyList(user_id string, name string, weight string, time string) error {
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

func (c *Connection) GetBuyList(user_id string) ([]BuyList, error) {
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
	err := c.conn.GetContext(context.Background(), &u, "SELECT id FROM users WHERE name = $1", username)
	if err != nil {
		return u, fmt.Errorf("people not found: %w", err)
	}
	return u, nil
}

func (c *Connection) GetFridgeList(user_id string) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT prod_name, status, experitation_date
	FROM fridge
	WHERE user_id = $1 AND (status = 'opened' OR status = 'stored')
	ORDER BY prod_name DESC;
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

func (c *Connection) UpdateProductToCooked(user_id string, name string, date string) error {
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
	SET experitation_date= $1, status = 'opened'
	WHERE prod_name = $2 AND user_id = $3;
	`, date, name, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) AddUser(name string, chat_id int) error {
	_, err := c.conn.ExecContext(context.Background(), `
	INSERT INTO users(name, chat_id)
	VALUES ($1, $2);
	`, name, chat_id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) GetStoredProductsList(user_id string) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT prod_name, status, experitation_date
	FROM fridge
	WHERE user_id = $1 AND status = 'stored'
	ORDER BY prod_name DESC;
	`, user_id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return fridge, nil
}

func (c *Connection) GetStatsByDateDifference(user_id string, firstDate string, secondDate string) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT status
	FROM fridge
	WHERE user_id = $1 AND (status = 'used' OR status = 'dispose') AND status_date <= $2 AND status_date >= $3;
	`, user_id, firstDate, secondDate)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return fridge, nil
}
