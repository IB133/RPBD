package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Users struct {
	Id      int
	Chat_id int64 `db:"chat_id"`
	Name    string
}

type BuyList struct {
	User_id   int
	Weight    float32
	Prod_name string
	Buy_time  string
}

type Fridge struct {
	User_id           int
	Status            string
	Prod_name         string
	Experitation_date time.Time
	Status_date       time.Time
}

type Connection struct {
	conn *sqlx.DB
}

func NewConnect(connString string) (*DB, error) {
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Up()
	return &DB{
		Conn: &Connection{conn: conn},
	}, nil
}

func (c *Connection) addProductToBuyList(userId int, name string, weight string, time string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	INSERT INTO buy_list(user_id, prod_name, weight, buy_time)
	VALUES($1, $2, $3, $4);
	`, userId, name, weight, time)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) addProductToFridge(userId int, name string, date string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	INSERT INTO fridge(user_id, prod_name, experitation_date)
	VALUES($1, $2, $3);
	`, userId, name, date)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) getBuyList(userId int) ([]BuyList, error) {
	var buyList []BuyList
	err := c.conn.SelectContext(context.Background(), &buyList, `
	SELECT prod_name, weight, buy_time
	FROM buy_list
	WHERE user_id = $1;
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return buyList, nil
}

func (c *Connection) getUserByUsername(username string) (Users, error) {
	u := Users{}
	err := c.conn.GetContext(context.Background(), &u, "SELECT id FROM users WHERE name = $1", username)
	if err != nil {
		return u, fmt.Errorf("people not found: %w", err)
	}
	return u, nil
}

func (c *Connection) getFridgeList(userId int) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT prod_name, status, experitation_date
	FROM fridge
	WHERE user_id = $1 AND (status = 'opened' OR status = 'stored')
	ORDER BY prod_name;
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return fridge, nil
}

func (c *Connection) getUsedProductsList(userId int) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT prod_name, status
	FROM fridge
	WHERE user_id = $1 AND (status = 'used' OR status = 'dispose');
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return fridge, nil
}

func (c *Connection) deleteFromBuyList(userId int, name string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	DELETE FROM buy_list
	WHERE user_id = $1 AND prod_name = $2;
	`, userId, name)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) updateProductToCooked(userId int, name string, date string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	UPDATE fridge 
	SET status = 'used',  status_date= $1
	WHERE prod_name = $2 AND user_id = $3;
	`, date, name, userId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) updateProductToDispose(userId int, name string, date string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	UPDATE fridge 
	SET status = 'dispose',  status_date= $1
	WHERE prod_name = $2 AND user_id = $3;
	`, date, name, userId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) openProduct(userId int, name string, date string) error {
	_, err := c.conn.ExecContext(context.Background(), `
	UPDATE fridge 
	SET experitation_date= $1, status = 'opened'
	WHERE prod_name = $2 AND user_id = $3;
	`, date, name, userId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) addUser(name string, chatId int64) error {
	_, err := c.conn.ExecContext(context.Background(), `
	INSERT INTO users(name, chat_id)
	VALUES ($1, $2);
	`, name, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) getStoredProductsList(userId int) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT prod_name, status, experitation_date
	FROM fridge
	WHERE user_id = $1 AND status = 'stored'
	ORDER BY prod_name DESC;
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return fridge, nil
}

func (c *Connection) getStatsByDateDifference(userId int, firstDate string, secondDate string) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT status
	FROM fridge
	WHERE user_id = $1 AND (status = 'used' OR status = 'dispose') AND status_date <= $2 AND status_date >= $3;
	`, userId, firstDate, secondDate)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return fridge, nil
}

func (c *Connection) getBuyListForScheduler(userId int) ([]BuyList, error) {
	var buyList []BuyList
	err := c.conn.SelectContext(context.Background(), &buyList, `
	SELECT *
	FROM buy_list
	WHERE buy_time::date = CURRENT_DATE AND user_id = $1;
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return buyList, nil
}

func (c *Connection) getFridgeListForScheduler(userId int) ([]Fridge, error) {
	var fridge []Fridge
	err := c.conn.SelectContext(context.Background(), &fridge, `
	SELECT prod_name, status
	FROM fridge
	WHERE experitation_date = CURRENT_DATE AND user_id = $1 AND (status = 'opened' OR status = 'stored');
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return fridge, nil
}

func (c *Connection) getUsersList() ([]Users, error) {
	var users []Users
	err := c.conn.SelectContext(context.Background(), &users, `
	SELECT *
	FROM users;
	`)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return users, nil
}
