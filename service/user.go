package service

import (
	"context"
	"errors"
	_ "errors"
	"fmt"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/GitEval/GitEval-Backend/pkg/llm"
	"github.com/google/go-github/v50/github"
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
	SaveUser(ctx context.Context, user model.User) error
	GetFollowingUsersJoinContact(ctx context.Context, id int64) ([]model.User, error)
	GetFollowersUsersJoinContact(ctx context.Context, id int64) ([]model.User, error)
}

type ContactDAOProxy interface {
	GetCountOfFollowing(ctx context.Context, id int64) (int64, error)
	GetCountOfFollowers(ctx context.Context, id int64) (int64, error)
	CreateContacts(ctx context.Context, contacts []model.FollowingContact) error
}
type DomainDAOProxy interface {
	Create(ctx context.Context, domain []model.Domain) error
	GetDomainById(ctx context.Context, id int64) ([]string, error)
}
type GithubProxy interface {
	GetFollowing(ctx context.Context, id int64) []model.User
	GetFollowers(ctx context.Context, id int64) []model.User
	CalculateScore(ctx context.Context, id int64, name string) float64
	GetAllRepositories(ctx context.Context, loginName string, userId int64) []*model.Repo
	GetClientFromMap(userID int64) (*github.Client, bool)
	GetAllEvents(ctx context.Context, username string, client *github.Client) ([]model.UserEvent, error)
}
type LLMProxy interface {
	GetArea(ctx context.Context, req llm.GetAreaRequest) (llm.GetAreaResponse, error)
	GetDomain(ctx context.Context, req llm.GetDomainRequest) (llm.GetDomainResponse, error)
	GetEvaluation(ctx context.Context, req llm.GetEvaluationRequest) (llm.GetEvaluationResponse, error)
}

type UserService struct {
	user    UserDAOProxy
	contact ContactDAOProxy
	domain  DomainDAOProxy
	tx      Transaction
	g       GithubProxy
	l       LLMProxy
}

func NewUserService(user UserDAOProxy, contact ContactDAOProxy, domain DomainDAOProxy, transaction Transaction, g GithubProxy, l LLMProxy) *UserService {
	return &UserService{
		user:    user,
		contact: contact,
		domain:  domain,
		tx:      transaction,
		g:       g,
		l:       l,
	}
}

// InitUser 存储user,同时搜索其following和followers,将他们也存入
func (s *UserService) InitUser(ctx context.Context, u model.User) (err error) {
	var (
		users = make([]model.User, 0)
	)

	following := s.g.GetFollowing(ctx, u.ID)
	users = append(users, following...)
	followers := s.g.GetFollowers(ctx, u.ID)
	users = append(users, followers...)

	var (
		followersLoc = make([]string, len(followers))
		followingLoc = make([]string, len(following))
	)

	// 获取followers和following的Loction
	// 顺便计算他们的分数
	for _, v := range followers {
		if v.Location != nil {
			followersLoc = append(followersLoc, *v.Location)
		}
		v.Score = s.g.CalculateScore(ctx, u.ID, v.LoginName)
	}

	for _, v := range following {
		if v.Location != nil {
			followingLoc = append(followingLoc, *v.Location)
		}
		v.Score = s.g.CalculateScore(ctx, u.ID, v.LoginName)
	}
	users = append(users, u)

	//得到关系
	followingContact := getContact(u.ID, following, Following)
	followersContact := getContact(u.ID, followers, Followers)

	//开启事务
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

	// 测试通过,花费时间大概10s
	go func() {
		ctx1 := context.Background()
		//得到用户的国籍,尝试存储这个用户的国籍
		Nation := s.generateNationality(ctx1, u.Bio, u.Company, u.Location, followersLoc, followingLoc)
		u.Nationality = Nation
		err := s.user.SaveUser(ctx1, u)
		if err != nil {
			return
		}
	}()

	// 测试通过,花费时间大概要到10s左右
	go func() {
		ctx2 := context.Background()
		var bio string
		if u.Bio != nil {
			bio = *u.Bio
		}
		//获取这个用户的主要技术领域
		userDomain := s.generateDomain(ctx2, u.LoginName, bio, u.ID)
		//将获取的结果转化成对应的model
		domains := StringToDomains(userDomain, u.ID)
		//存储domain
		if err := s.domain.Create(ctx2, domains); err != nil {
			return
		}
	}()

	return nil

}

// GetDomains 返回用户的领域（基于主要使用的语言）
// 接受userId，返回用户的领域
func (s *UserService) GetDomains(ctx context.Context, userId int64) []string {
	domains, err := s.domain.GetDomainById(ctx, userId)
	if err != nil {
		return nil
	}
	return domains
}

// GetUserById 从ID获取用户信息
func (s *UserService) GetUserById(ctx context.Context, id int64) (model.User, error) {
	return s.user.GetUserByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, u model.User) error {
	err := s.user.SaveUser(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

// 生成国籍
func (s *UserService) generateNationality(ctx context.Context, bio, company, location *string, followerLoc, followingloc []string) string {
	res, err := s.l.GetArea(ctx, llm.GetAreaRequest{
		Bio:            bio,
		Company:        company,
		Location:       location,
		FollowerAreas:  followerLoc,
		FollowingAreas: followingloc,
	})
	if err != nil {
		log.Println(errors.New("failed to get Nationality"))
		return ""
	}
	//添加置信度
	nation := fmt.Sprintf("%s(trust:%.2f)", res.Area, res.Confidence)
	return nation
}

func (s *UserService) generateDomain(ctx context.Context, LoginName, bio string, userId int64) []string {
	repos := s.g.GetAllRepositories(ctx, LoginName, userId)
	if len(repos) == 0 {
		return nil
	}
	var r = make([]llm.Repo, len(repos))
	for k, v := range repos {
		r[k] = llm.Repo{
			Name:         v.Name,
			MainLanguage: v.Language,
			Readme:       v.Readme,
		}
	}
	resp, err := s.l.GetDomain(ctx, llm.GetDomainRequest{
		Repos: r,
		Bio:   bio,
	})
	if err != nil {
		log.Println(errors.New("failed to get domain"))
		return nil
	}
	return resp.Domain
}

func StringToDomains(lang []string, id int64) []model.Domain {
	var (
		domainsModel = make([]model.Domain, len(lang))
	)
	for k, v := range lang {
		domainsModel[k].UserID = id
		domainsModel[k].Domain = v
	}
	return domainsModel
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

func (s *UserService) GetEvaluation(ctx context.Context, userId int64) (string, error) {
	user, err := s.user.GetUserByID(ctx, userId)
	if err != nil {
		return "", err
	}
	followers, err := s.user.GetFollowersUsersJoinContact(ctx, userId)
	if err != nil {
		return "", err
	}
	following, err := s.user.GetFollowingUsersJoinContact(ctx, userId)
	if err != nil {
		return "", err
	}
	client, ok := s.g.GetClientFromMap(userId)
	if !ok {
		return "", err
	}
	events, err := s.g.GetAllEvents(ctx, user.LoginName, client)
	if err != nil {
		return "", err
	}
	var userEvents []llm.UserEvent
	for _, event := range events {
		userEvents = append(userEvents, llm.UserEvent{
			Repo: &llm.RepoInfo{
				Description:      event.Repo.Description,
				StargazersCount:  event.Repo.StargazersCount,
				ForksCount:       event.Repo.ForksCount,
				CreatedAt:        event.Repo.CreatedAt,
				SubscribersCount: event.Repo.SubscribersCount,
			},
			CommitCount:      event.CommitCount,
			IssuesCount:      event.IssuesCount,
			PullRequestCount: event.PullRequestCount,
		})
	}

	domains, err := s.domain.GetDomainById(ctx, user.ID)
	if err != nil {
		return "", err
	}

	evaluation, err := s.l.GetEvaluation(ctx, llm.GetEvaluationRequest{
		Bio:               *user.Bio,
		Followers:         len(followers),
		Following:         len(following),
		TotalPrivateRepos: user.TotalPrivateRepos,
		TotalPublicRepos:  user.TotalPrivateRepos,
		UserEvents:        userEvents,
		Domains:           domains,
	})
	if err != nil {
		return "", err
	}
	return evaluation.Evaluation, nil
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
