package service

import (
	"context"
	"github.com/GitEval/GitEval-Backend/model"
	"log"
)

const (
	Following = iota
	Followers
)

// 有关user的服务

type UserDAOProxy interface {
	CreateUsers(ctx context.Context, user []model.User) error
	GetUserByID(ctx context.Context, id int64) (model.User, error)
	GetFollowingUsersJoinContact(ctx context.Context, id int64) ([]model.User, error)
	GetFollowersUsersJoinContact(ctx context.Context, id int64) ([]model.User, error)
}
type ContactDAOProxy interface {
	GetCountOfFollowing(ctx context.Context, id int64) (int64, error)
	GetCountOfFollowers(ctx context.Context, id int64) (int64, error)
	CreateContacts(ctx context.Context, contacts []model.FollowingContact) error
}
type GithubProxy interface {
	GetFollowing(ctx context.Context, id int64) []model.User
	GetFollowers(ctx context.Context, id int64) []model.User
}

type UserService struct {
	user    UserDAOProxy
	contact ContactDAOProxy
	tx      Transaction
	g       GithubProxy
}

func NewUserService(user UserDAOProxy, contact ContactDAOProxy, transaction Transaction, g GithubProxy) *UserService {
	return &UserService{
		user:    user,
		contact: contact,
		tx:      transaction,
		g:       g,
	}
}

// InitUser 存储user,同时搜索其following和followers,将他们也存入
func (s *UserService) InitUser(ctx context.Context, u model.User) (err error) {
	var (
		users = make([]model.User, 0)
	)
	users = append(users, u)
	following := s.g.GetFollowing(ctx, u.ID)
	followers := s.g.GetFollowers(ctx, u.ID)
	followingContact := getContact(u.ID, following, Following)
	followersContact := getContact(u.ID, followers, Followers)
	err = s.tx.InTx(ctx, func(ctx context.Context) error {
		if err := s.user.CreateUsers(ctx, users); err != nil {
			return err
		}
		if err := s.contact.CreateContacts(ctx, followingContact); err != nil {
			return err
		}
		if err := s.contact.CreateContacts(ctx, followersContact); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("Init user failed")
		return err
	}
	return nil

}

// GetUserById 从ID获取用户信息
func (s *UserService) GetUserById(ctx context.Context, id int64) (model.User, error) {
	return s.user.GetUserByID(ctx, id)
}

// 从users中得到相应的关系
func getContact(Id int64, users []model.User, follow int) []model.FollowingContact {
	var (
		contact = make([]model.FollowingContact, len(users))
	)
	if follow == Following {
		for k, user := range users {
			contact[k].Subject = Id
			contact[k].Object = user.ID
		}
	}
	if follow == Followers {
		for k, user := range users {
			contact[k].Subject = user.ID
			contact[k].Object = Id
		}
	}
	return contact
}
