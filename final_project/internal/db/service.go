package db

import (
	"context"
	"fmt"
	"time"

	"github.com/IB133/RPBD/final_project/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *Store
}

const (
	singingString = "qwesdfmksjow8457HDudf"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (s *Service) Registration(ctx context.Context, u *models.User) error {
	if err := u.ValidateRegistration(); err != nil {
		return err
	}
	if err := u.HashPassword(); err != nil {
		return err
	}
	if err := s.db.addUSer(ctx, u); err != nil {
		return err
	}
	return nil
}

func (s *Service) AuthAndGenerateToken(ctx context.Context, u *models.User) (string, error) {
	udb, err := s.db.getUserByEmail(ctx, u.Email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(udb.Password), []byte(u.Password)); err != nil {
		return "", fmt.Errorf("Incorrect password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: u.Id,
	})
	return token.SignedString([]byte(singingString))
}

func (s *Service) ParseToken(accesToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid singing method")
		}
		return []byte(singingString), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, fmt.Errorf("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil

}

func (s *Service) AddComment(ctx context.Context, c *models.Comments) error {
	if c.Reply_to == 0 {
		if err := s.db.addParentComment(ctx, c); err != nil {
			return err
		}
		return nil
	}
	if err := s.db.addChildComment(ctx, c); err != nil {
		return err

	}
	return nil
}

func (s *Service) GetComments(ctx context.Context, replyId int) (map[string]interface{}, error) {
	if replyId == 0 {
		com, err := s.db.getFirstLevlCommentList(ctx)
		if err != nil {
			return nil, err
		}
		mp := make(map[string]interface{})
		for _, v := range com {
			mp[v.Login] = fmt.Sprintf("%s, %s, %v", v.Content, v.Created_at, v.Status)
		}
		return mp, nil
	}
	com, err := s.db.getScecondLevelList(ctx, replyId)
	if err != nil {
		return nil, err
	}
	mp := make(map[string]interface{})
	for _, v := range com {
		mp[v.Login] = fmt.Sprintf("%s, %s, %v", v.Content, v.Created_at, v.Status)
	}
	return mp, nil
}

func (s *Service) AddNewsByUser(ctx context.Context, n *models.News) error {
	if err := n.ValidateNews(); err != nil {
		return err
	}
	if err := s.db.addNewsByUser(ctx, n); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddNewsByModer(ctx context.Context, n *models.News) error {
	if err := n.ValidateNews(); err != nil {
		return err
	}
	if err := s.db.addNewsByModer(ctx, n); err != nil {
		return err
	}
	return nil
}

func (s *Service) CheckoutNews(ctx context.Context, n *models.News) error {
	if !n.Posted {
		if err := n.ValidateModerComm(); err != nil {
			return err
		}

		if err := s.db.declineNews(ctx, n.Moder_comm, n.Moder_id, n.Id); err != nil {
			return err
		}
		return nil
	}
	if err := s.db.acceptNews(ctx, n.Posted_at, n.Moder_id, n.Id); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteComm(ctx context.Context, commId int) error {
	if err := s.db.deleteComment(ctx, commId); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUnpostedNewsList(ctx context.Context) (map[string]interface{}, error) {
	list, err := s.db.getUnpostedNews(ctx)
	if err != nil {
		return nil, err
	}
	mp := make(map[string]interface{})
	for _, v := range list {
		mp[v.Login] = fmt.Sprintf("%v, %s, %s, %s", v.Id, v.Title, v.Theme, v.Content)
	}
	return mp, nil
}

func (s *Service) GetUserCommentsList(ctx context.Context, userId int) (map[string]interface{}, error) {
	list, err := s.db.getCommentsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	mp := make(map[string]interface{})
	for _, v := range list {
		mp[v.Title] = fmt.Sprintf("%s, %s", v.Content, v.Created_at)
	}
	return mp, nil
}
