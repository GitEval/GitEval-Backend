package service

import (
	"context"
	"errors"
	_ "errors"
	"fmt"
	llmv1 "github.com/GitEval/GitEval-Backend/client/gen"
	"github.com/GitEval/GitEval-Backend/errs"
	"github.com/GitEval/GitEval-Backend/model"
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
	SearchUser(ctx context.Context, nation *string, domain string, page int, pageSize int) ([]model.User, error)
}

type ContactDAOProxy interface {
	GetCountOfFollowing(ctx context.Context, id int64) (int64, error)
	GetCountOfFollowers(ctx context.Context, id int64) (int64, error)
	CreateContacts(ctx context.Context, contacts []model.FollowingContact) error
}

type DomainDAOProxy interface {
	Create(ctx context.Context, domain []model.Domain) error
	GetDomainById(ctx context.Context, id int64) ([]string, error)
	Delete(ctx context.Context, id int64) error
}

type GithubProxy interface {
	GetFollowing(ctx context.Context, id int64) []model.User
	GetFollowers(ctx context.Context, id int64) []model.User
	CalculateScore(ctx context.Context, id int64, name string) float64
	GetAllRepositories(ctx context.Context, loginName string, userId int64) []*model.Repo
	GetClientFromMap(userID int64) (*github.Client, bool)
	GetAllUserEvents(ctx context.Context, username string, client *github.Client) ([]model.UserEvent, error)
}

type UserService struct {
	user    UserDAOProxy
	contact ContactDAOProxy
	domain  DomainDAOProxy
	tx      Transaction
	g       GithubProxy
	l       llmv1.LLMServiceClient
}

func NewUserService(user UserDAOProxy, contact ContactDAOProxy, domain DomainDAOProxy, transaction Transaction, g GithubProxy, l llmv1.LLMServiceClient) *UserService {
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
	followers := s.g.GetFollowers(ctx, u.ID)
	var (
		followersLoc = make([]string, len(followers))
		followingLoc = make([]string, len(following))
	)

	// 获取followers和following的Location
	// 顺便计算他们的分数
	for i := range followers {
		followersLoc = append(followersLoc, followers[i].Location)
		followers[i].Score = s.g.CalculateScore(ctx, u.ID, followers[i].LoginName)
	}
	users = append(users, following...)

	for i := range following {
		followingLoc = append(followingLoc, following[i].Location)
		following[i].Score = s.g.CalculateScore(ctx, u.ID, following[i].LoginName)
	}
	users = append(users, followers...)

	//将user二次存入,这个地方主要是为了能够保证每次用户上号这个评分都能更新
	u.Score = s.g.CalculateScore(ctx, u.ID, u.LoginName)
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
		//获取这个用户的主要技术领域
		userDomain := s.generateDomain(ctx2, u.LoginName, u.Bio, u.ID)
		//将获取的结果转化成对应的model
		domains := StringToDomains(userDomain, u.ID)
		//先删除之前的记录,这个地方不够优雅
		err = s.domain.Delete(ctx2, u.ID)
		if err != nil {
			return
		}
		//存储domain
		if err = s.domain.Create(ctx2, domains); err != nil {
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
		UserID:    user.ID,
		UserName:  user.LoginName,
		AvatarURL: user.AvatarURL,
		Score:     user.Score,
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
		return "", errs.LoginFailErr
	}

	events, err := s.g.GetAllUserEvents(ctx, user.LoginName, client)
	if err != nil {
		return "", err
	}

	var userEvents []*llmv1.UserEvent
	for _, event := range events {
		userEvents = append(userEvents, &llmv1.UserEvent{
			Repo: &llmv1.RepoInfo{
				Name:             event.Repo.Name,
				Description:      event.Repo.Description,
				StargazersCount:  int32(event.Repo.StargazersCount),
				ForksCount:       int32(event.Repo.ForksCount),
				CreatedAt:        event.Repo.CreatedAt,
				SubscribersCount: int32(event.Repo.SubscribersCount),
			},
			CommitCount:      int32(event.PushCount),
			IssuesCount:      int32(event.IssuesCount),
			PullRequestCount: int32(event.PullRequestCount),
		})
	}

	//此处允许获取值为空而不报错,因为可能用户没有成功获取领域就直接开始做评价了
	domains, _ := s.domain.GetDomainById(ctx, user.ID)
	evaluation, err := s.l.GetEvaluation(ctx, &llmv1.GetEvaluationRequest{
		Bio:               user.Bio,
		Followers:         int32(len(followers)),
		Following:         int32(len(following)),
		TotalPrivateRepos: int32(user.TotalPrivateRepos),
		TotalPublicRepos:  int32(user.PublicRepos),
		UserEvents:        userEvents,
		Domains:           domains,
	})
	if err != nil {
		return "", err
	}

	return evaluation.Evaluation, nil
}

func (s *UserService) GetNationByUserId(ctx context.Context, userId int64) (string, error) {
	user, err := s.user.GetUserByID(ctx, userId)
	if err != nil {
		log.Println("get user failed")
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

	var (
		followersLoc = make([]string, len(followers))
		followingLoc = make([]string, len(following))
	)

	// 获取followers和following的Loction
	for _, v := range followers {
		followersLoc = append(followersLoc, v.Location)
	}

	for _, v := range following {
		followingLoc = append(followingLoc, v.Location)
	}

	Nation := s.generateNationality(ctx, user.Bio, user.Company, user.Location, followersLoc, followingLoc)
	user.Nationality = Nation
	err = s.user.SaveUser(ctx, user)
	if err != nil {
		return "", err
	}

	return Nation, nil
}

func (s *UserService) GetDomainByUserId(ctx context.Context, userId int64) ([]string, error) {
	user, err := s.user.GetUserByID(ctx, userId)
	if err != nil {
		log.Println("get user failed")
		return nil, err
	}

	userDomain := s.generateDomain(ctx, user.LoginName, user.Bio, user.ID)
	//将获取的结果转化成对应的model
	domains := StringToDomains(userDomain, user.ID)
	//先删除之前的记录
	err = s.domain.Delete(ctx, user.ID)
	if err != nil {
		return []string{}, err
	}

	//存储domain
	if err := s.domain.Create(ctx, domains); err != nil {
		return []string{}, err
	}

	var resp []string
	for _, domain := range domains {
		resp = append(resp, domain.Domain)
	}
	return resp, nil

}

func (s *UserService) SearchUser(ctx context.Context, nation *string, domain string, page int, pageSize int) ([]model.User, error) {
	return s.user.SearchUser(ctx, nation, domain, page, pageSize)
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
		leaderboard[k].UserName = user.LoginName
		leaderboard[k].AvatarURL = user.AvatarURL
		leaderboard[k].Score = user.Score
	}
	return leaderboard
}

func removeTheSame(s []model.Leaderboard) []model.Leaderboard {
	var (
		result = make([]model.Leaderboard, 0)
		mp     = make(map[int64]model.Leaderboard)
	)

	for _, v := range s {
		mp[v.UserID] = v
	}
	for _, v := range mp {
		result = append(result, v)
	}
	return result
}

// 生成国籍
func (s *UserService) generateNationality(ctx context.Context, bio, company, location string, followerLoc, followingloc []string) string {
	res, err := s.l.GetArea(ctx, &llmv1.GetAreaRequest{
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
	nation := fmt.Sprintf("%s|(trust:%.2f)", res.Area, res.Confidence)
	return nation
}

func (s *UserService) generateDomain(ctx context.Context, LoginName, bio string, userId int64) []string {
	repos := s.g.GetAllRepositories(ctx, LoginName, userId)
	if len(repos) == 0 {
		return nil
	}

	// 使用 make 来预分配切片大小，提升性能
	r := make([]*llmv1.Repo, 0, len(repos))
	for _, v := range repos {
		repo := &llmv1.Repo{
			Name:     v.Name,
			Language: v.Language,
			Readme:   v.Readme,
		}
		r = append(r, repo)
	}

	domains, err := s.l.GetDomain(ctx, &llmv1.GetDomainRequest{
		Repos: r,
		Bio:   bio,
	})
	if err != nil {
		log.Println(errors.New("failed to get domain"))
		return nil
	}

	// 添加置信度并格式化输出
	resp := make([]string, 0, len(domains.Domains))
	for _, domain := range domains.Domains {
		resp = append(resp, fmt.Sprintf("%s|(trust:%.2f)", domain.Domain, domain.Confidence))
	}
	return resp
}
