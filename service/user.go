package service

import (
	"context"
	"github.com/GitEval/GitEval-Backend/model"
	"log"
	"sort"
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
	CalculateScore(ctx context.Context, id int64, name string) float64
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
	users = append(users, following...)
	followers := s.g.GetFollowers(ctx, u.ID)
	users = append(users, followers...)
	//获取分数
	for _, v := range users {
		v.Score = s.g.CalculateScore(ctx, u.ID, v.LoginName)
	}
	//得到关系
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

// GetLeaderboard 获取排行榜
func (s *UserService) GetLeaderboard(ctx context.Context, userId int64) ([]model.Leaderboard, error) {
	var (
		leaderboard = make([]model.Leaderboard, 0)
		err         error
	)
	user, err := s.user.GetUserByID(ctx, userId)
	if err != nil {
		log.Println("get user failed")
		return nil, err
	}
	leaderboard = append(leaderboard, model.Leaderboard{
		UserID: user.ID,
		Score:  user.Score,
	})
	//获取following
	followings, err := s.user.GetFollowingUsersJoinContact(ctx, userId)
	if err != nil {
		log.Println("get following failed")
		return nil, err
	}
	//获取followers
	followers, err := s.user.GetFollowersUsersJoinContact(ctx, userId)
	if err != nil {
		log.Println("get followers failed")
		return nil, err
	}
	leaderboard = append(leaderboard, getLeaderboard(followings)...)
	leaderboard = append(leaderboard, getLeaderboard(followers)...)
	//去重
	leaderboard = removeTheSame(leaderboard)
	//从大到小排序
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})
	return leaderboard, nil
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
func getLeaderboard(users []model.User) []model.Leaderboard {
	var (
		leaderboard = make([]model.Leaderboard, len(users))
	)
	for k, user := range users {
		leaderboard[k].UserID = user.ID
		leaderboard[k].Score = user.Score
	}
	return leaderboard
}
func removeTheSame(s []model.Leaderboard) []model.Leaderboard {
	var (
		result = make([]model.Leaderboard, 0)
		mp     = make(map[int64]float64)
	)

	for _, v := range s {
		mp[v.UserID] = v.Score
	}
	for k, v := range mp {
		result = append(result, model.Leaderboard{UserID: k, Score: v})
	}
	return result
}
