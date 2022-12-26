package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/IB133/RPBD/final_project/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
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

func clean(ctx context.Context, t *testing.T, conn *sqlx.DB, container testcontainers.Container) {
	t.Cleanup(func() {
		_, err := conn.ExecContext(ctx, `
		BEGIN;
		DELETE FROM users;
		DELETE FROM person;
		DELETE FROM comments;
		DELETE FROM news;
		COMMIT;
		`)
		if err != nil {
			t.Fatalf("Cant delete from tables")
		}
		container.Terminate(ctx)
	})
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	err := query.addUSer(ctx, &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	})
	assert.Nil(t, err)
	clean(ctx, t, conn, container)
}

func TestAuth(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	usr, err := query.getUserByEmail(ctx, "qwe@web.com")
	assert.Equal(t, usr.Email, m.Email)
	clean(ctx, t, conn, container)
}

func TestAddParentComment(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	c := &models.Comments{
		News_id:    1,
		User_id:    1,
		Created_at: "2022-12-23",
		Content:    "coment",
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	err = query.addParentComment(ctx, c)
	assert.Nil(t, err)
	clean(ctx, t, conn, container)
}

func TestAddChildComment(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	c := &models.Comments{
		News_id:    1,
		User_id:    1,
		Created_at: "2022-12-23",
		Content:    "coment",
		Reply_to:   1,
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	err = query.addChildComment(ctx, c)
	assert.Nil(t, err)
	clean(ctx, t, conn, container)
}

func TestGetFirsLvlComments(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	c := &models.Comments{
		News_id:    1,
		User_id:    1,
		Created_at: "2022-12-23",
		Content:    "coment",
	}
	type comm struct {
		Login      string
		Created_at string
		Content    string
		Status     bool
	}
	rc := comm{
		Login:      "login1",
		Created_at: "2022-12-23T00:00:00Z",
		Content:    "coment",
		Status:     true,
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	err = query.addParentComment(ctx, c)
	assert.Nil(t, err)
	list, err := query.getFirstLevlCommentList(ctx)
	assert.Nil(t, err)
	assert.Equal(t, list[0], rc)
	clean(ctx, t, conn, container)
}

func TestGetSecondLvlList(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	c := &models.Comments{
		News_id:    1,
		User_id:    1,
		Created_at: "2022-12-23",
		Content:    "coment",
		Reply_to:   1,
	}
	type comm struct {
		Login      string
		Created_at string
		Content    string
		Status     bool
	}
	rc := comm{
		Login:      "login1",
		Created_at: "2022-12-23T00:00:00Z",
		Content:    "coment",
		Status:     true,
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	err = query.addChildComment(ctx, c)
	assert.Nil(t, err)
	list, err := query.getScecondLevelList(ctx, 1)
	assert.Nil(t, err)
	assert.Equal(t, list[0], rc)
	clean(ctx, t, conn, container)
}

func TestAddNewsByUser(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	clean(ctx, t, conn, container)
}

// Тоже самое что и юзер, только добавляет модер
func TestAddNewsByModer(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	clean(ctx, t, conn, container)
}

func TestAcceptNews(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	err = query.acceptNews(ctx, "2022-12-24", 1, 1)
	assert.Nil(t, err)
	list, err := query.getUnpostedNews(ctx)
	assert.Nil(t, err)
	assert.Nil(t, list)
	clean(ctx, t, conn, container)
}

func TestDeclineNews(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	err = query.declineNews(ctx, "bad news", 1, 1)
	assert.Nil(t, err)
	list, err := query.getUnpostedNews(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	clean(ctx, t, conn, container)
}

func TestGetUnpostedNews(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	list, err := query.getUnpostedNews(ctx)
	assert.Nil(t, err)
	assert.Equal(t, list[0].Login, m.Login)
	clean(ctx, t, conn, container)
}

func TestGetUsersComments(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	c := &models.Comments{
		News_id:    1,
		User_id:    1,
		Created_at: "2022-12-12",
		Content:    "cont",
	}
	type usercomms struct {
		Title      string
		Content    string
		Created_at string
	}
	usc := &usercomms{
		Title:      "title",
		Content:    "cont",
		Created_at: "2022-12-12",
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	err = query.addParentComment(ctx, c)
	list, err := query.getCommentsByUserId(ctx, 1)
	assert.Equal(t, list[0], usc)
	clean(ctx, t, conn, container)
}

func TestDeleteComments(t *testing.T) {
	ctx := context.Background()
	container, conn := CreateTestDatabase()
	query := Store{
		conn: conn,
	}
	m := &models.User{
		Email:    "qwe@web.com",
		Password: "12345567890",
		Login:    "login1",
	}
	n := &models.News{
		Title:   "title",
		Theme:   "tema",
		Content: "content",
		User_id: 1,
	}
	c := &models.Comments{
		News_id:    1,
		User_id:    1,
		Created_at: "2022-12-12",
		Content:    "cont",
	}
	err := query.addUSer(ctx, m)
	assert.Nil(t, err)
	err = query.addNewsByUser(ctx, n)
	assert.Nil(t, err)
	err = query.addParentComment(ctx, c)
	err = query.deleteComment(ctx, 1)
	assert.Nil(t, err)
	clean(ctx, t, conn, container)
}

// func TestGetBuyListForScheduler(t *testing.T) {
// 	container, conn := CreateTestDatabase()
// 	defer container.Terminate(context.Background())
// 	query := Connection{
// 		conn: conn,
// 	}
// 	query.addUser("fgdfgdf", 23232)
// 	query.addProductToBuyList(1, "aboba", "12", "2022-12-12")
// 	_, err := query.getBuyListForScheduler(1)
// 	if err != nil && err != sql.ErrNoRows {
// 		t.Error(err)
// 	}
// }

// func TestGetFridgeListForScheduler(t *testing.T) {
// 	container, conn := CreateTestDatabase()
// 	defer container.Terminate(context.Background())
// 	query := Connection{
// 		conn: conn,
// 	}
// 	query.addUser("fgdfgdf", 23232)
// 	query.addProductToFridge(1, "aboba", "2022-12-12")
// 	_, err := query.getFridgeListForScheduler(1)
// 	if err != nil && err != sql.ErrNoRows {
// 		t.Error(err)
// 	}
// }

// func TestGetUsersList(t *testing.T) {
// 	container, conn := CreateTestDatabase()
// 	defer container.Terminate(context.Background())
// 	query := Connection{
// 		conn: conn,
// 	}
// 	query.addUser("fgdfgdf", 23232)
// 	_, err := query.getUsersList()
// 	if err != nil && err != sql.ErrNoRows {
// 		t.Error(err)
// 	}
// }
