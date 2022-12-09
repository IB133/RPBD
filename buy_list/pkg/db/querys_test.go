package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateTestDatabase() (testcontainers.Container, *sqlx.DB) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "test",
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_USER":     "postgres",
		},
	}

	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		panic(err)
	}

	host, err := dbContainer.Host(context.Background())
	if err != nil {
		panic(err)
	}
	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		panic(err)
	}

	connString := fmt.Sprintf("postgres://postgres:pass@%v:%v/test?sslmode=disable", host, port.Port())
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:../../migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Up()
	return dbContainer, db
}

func TestAddUsers(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	err := query.addUser("fgdfgdf", 23232)
	if err != nil {
		t.Error(err)
	}
}

func TestAddProductToBuyList(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	err := query.addProductToBuyList(1, "fgdfgdf", "12", "2022-12-12")
	if err != nil {
		t.Error(err)
	}
}

func TestAddProductToFridge(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	err := query.addProductToFridge(1, "aboba", "2022-12-12")
	if err != nil {
		t.Error(err)
	}
}

func TestGetBuyList(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	_, err := query.getBuyList(10)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetUserByUsername(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("aboba", 1123213)
	_, err := query.getUserByUsername("aboba")
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetFridgeList(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	_, err := query.getFridgeList(1)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetUsedProductsList(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	_, err := query.getFridgeList(1)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestDeleteFromBuyList(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	query.addProductToBuyList(1, "aboba", "12", "2022-12-12")
	err := query.deleteFromBuyList(1, "aboba")
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestUpdateProductToCooked(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	query.addProductToFridge(1, "aboba", "2022-12-12")
	err := query.updateProductToCooked(1, "aboba", "2022-12-12")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateProductToDispose(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	query.addProductToFridge(1, "aboba", "2022-12-12")
	err := query.updateProductToDispose(1, "aboba", "2022-12-12")
	if err != nil {
		t.Error(err)
	}
}

func TestOpenProduct(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	query.addProductToFridge(1, "aboba", "2022-12-12")
	err := query.updateProductToDispose(1, "aboba", "2022-12-12")
	if err != nil {
		t.Error(err)
	}
}

func TestGetStoredProductsList(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	query.addProductToFridge(1, "aboba", "2022-12-12")
	_, err := query.getStoredProductsList(1)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetStatsByDateDifference(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	query.addProductToFridge(1, "aboba", "2022-12-12")
	_, err := query.getStatsByDateDifference(1, "2022-12-12", "2022-12-29")
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetBuyListForScheduler(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	query.addProductToBuyList(1, "aboba", "12", "2022-12-12")
	_, err := query.getBuyListForScheduler(1)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetFridgeListForScheduler(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	query.addProductToFridge(1, "aboba", "2022-12-12")
	_, err := query.getFridgeListForScheduler(1)
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetUsersList(t *testing.T) {
	container, conn := CreateTestDatabase()
	defer container.Terminate(context.Background())
	query := Connection{
		conn: conn,
	}
	query.addUser("fgdfgdf", 23232)
	_, err := query.getUsersList()
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}
