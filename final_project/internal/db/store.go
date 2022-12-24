package db

import (
	"context"

	"github.com/IB133/RPBD/final_project/internal/models"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	conn *sqlx.DB
}

func NewConnect(connString string) (*Service, error) {
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}
	// driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	// if err != nil {
	// 	return nil, err
	// }

	// m, err := migrate.NewWithDatabaseInstance(
	// 	"file://migrations",
	// 	"postgres", driver)
	// if err != nil {
	// 	panic(err)
	// }
	// m.Up()
	return &Service{
		db: &Store{conn: conn},
	}, nil
}

func (s *Store) addUSer(ctx context.Context, u *models.User) error {
	_, err := s.conn.ExecContext(ctx, `
	INSERT INTO users(email, login, password)
	VALUES ($1, $2, $3);
	`, u.Email, u.Login, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) getUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user []models.User
	err := s.conn.SelectContext(ctx, &user, `
	SELECT id, password
	FROM users
	WHERE email = $1;
	`, email)
	if err != nil {
		return models.User{}, err
	}
	return user[0], nil
}

func (s *Store) addParentComment(ctx context.Context, c *models.Comments) error {
	_, err := s.conn.ExecContext(ctx, `
	INSERT INTO comments(news_id, user_id, created_at, content)
	VALUES ($1, $2, $3, $4);
	`, c.News_id, c.User_id, c.Created_at, c.Content)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) addChildComment(ctx context.Context, c *models.Comments) error {
	_, err := s.conn.ExecContext(ctx, `
	INSERT INTO comments(news_id, user_id, created_at, reply_to, content)
	VALUES ($1, $2, $3, $4, $5);
	`, c.News_id, c.User_id, c.Created_at, c.Reply_to, c.Content)
	if err != nil {
		return err
	}
	return nil
}

type comm struct {
	Login      string
	Created_at string
	Content    string
	Status     bool
}

func (s *Store) getFirstLevlCommentList(ctx context.Context) ([]comm, error) {
	var ucom []comm
	err := s.conn.SelectContext(ctx, &ucom, `
	SELECT u.login, c.created_at, c.content, c.status
	FROM comments AS c, users AS u
	WHERE c.reply_to is null AND c.user_id = u.id
	ORDER BY c.created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	return ucom, nil
}

func (s *Store) getScecondLevelList(ctx context.Context, replyId int) ([]comm, error) {
	var ucom []comm
	err := s.conn.SelectContext(ctx, &ucom, `
	SELECT u.login, c.created_at, c.content, c.status
	FROM comments AS c, users AS u
	WHERE c.reply_to = $1 AND c.user_id = u.id
	ORDER BY c.created_at DESC;
	`, replyId)
	if err != nil {
		return nil, err
	}
	return ucom, nil
}

func (s *Store) addNewsByUser(ctx context.Context, n *models.News) error {
	_, err := s.conn.ExecContext(ctx, `
	INSERT INTO news(theme, title, content, user_id)
	VALUES ($1, $2, $3, $4);
	`, n.Theme, n.Title, n.Content, n.User_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) addNewsByModer(ctx context.Context, n *models.News) error {
	_, err := s.conn.ExecContext(ctx, `
	INSERT INTO news(theme, title, content, moder_id, posted, posted_at)
	VALUES ($1, $2, $3, $4, true, $5);
	`, n.Theme, n.Title, n.Content, n.Moder_id, n.Posted_at)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) acceptNews(ctx context.Context, postedDate string, moderId int, id int) error {
	_, err := s.conn.ExecContext(ctx, `
	UPDATE news
	SET posted = 'true', posted_at = $1, moder_id = $2
	WHERE id = $3;
	`, postedDate, moderId, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) declineNews(ctx context.Context, moderComm string, moderId int, id int) error {
	_, err := s.conn.ExecContext(ctx, `
	UPDATE news
	SET moder_comm = $1, moder_id = $2
	WHERE id = $3;
	`, moderComm, moderId, id)
	if err != nil {
		return err
	}
	return nil
}

type unpostednews struct {
	Login   string
	Id      int
	Theme   string
	Title   string
	Content string
}

func (s *Store) getUnpostedNews(ctx context.Context) ([]unpostednews, error) {
	var news []unpostednews
	err := s.conn.SelectContext(ctx, &news, `
	SELECT users.login, news.id, theme, title, content
	FROM news , users
	WHERE posted = false AND users.id = news.user_id;
	`)
	if err != nil {
		return nil, err
	}
	return news, nil
}

type usercomms struct {
	Title      string
	Content    string
	Created_at string
}

func (s *Store) getCommentsByUserId(ctx context.Context, userId int) ([]usercomms, error) {
	var comm []usercomms
	err := s.conn.SelectContext(ctx, &comm, `
	SELECT news.title, comments.content, created_at
	FROM comments
	INNER JOIN news ON comments.news_id = news.id
	WHERE comments.user_id = $1;
	`, userId)
	if err != nil {
		return nil, err
	}
	return comm, nil
}

func (s *Store) deleteComment(ctx context.Context, comId int) error {
	_, err := s.conn.ExecContext(ctx, `
	UPDATE comments
	SET status = false, content = 'Comment is delete'
	WHERE id = $1;
	`, comId)
	if err != nil {
		return err
	}
	return nil
}
