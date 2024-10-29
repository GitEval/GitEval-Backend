package service

import (
	"context"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/GitEval/GitEval-Backend/pkg/github"
)

type AuthService interface {
	Login(ctx context.Context) (url string, err error)
	CallBack(ctx context.Context, code string) (userId int64, err error)
}

type authService struct {
	userDAO   model.UserDAO
	githubAPI github.GitHubAPI
}

func NewAuthService(userDAO model.UserDAO) AuthService {
	return &authService{userDAO: userDAO}
}

func (s *authService) Login(ctx context.Context) (url string, err error) {
	githubAPI := github.NewGitHubAPI()
	url = githubAPI.GetLoginUrl()
	return url, nil
}

func (s *authService) CallBack(ctx context.Context, code string) (userId int64, err error) {
	client, err := s.githubAPI.GetClientByCode(code)
	if err != nil {
		return 0, err
	}

	userInfo, err := s.githubAPI.GetUserInfo(ctx, client)
	if err != nil {
		return 0, err
	}

	// 根据用户 ID 查找用户
	user, err := s.userDAO.GetUserByID(userInfo.GetID())
	if err != nil {
		return 0, err
	}

	// 如果用户不存在，创建新用户,如果存在
	if user == nil {
		user = &model.User{
			ID:                      userInfo.GetID(),
			Login:                   userInfo.GetLogin(),
			NodeID:                  userInfo.GetNodeID(),
			AvatarURL:               userInfo.GetAvatarURL(),
			HTMLURL:                 userInfo.GetHTMLURL(),
			GravatarID:              userInfo.GetGravatarID(),
			Name:                    userInfo.GetName(),
			Company:                 userInfo.GetCompany(),
			Blog:                    userInfo.GetBlog(),
			Location:                userInfo.GetLocation(),
			Email:                   userInfo.GetEmail(),
			Hireable:                userInfo.GetHireable(),
			Bio:                     userInfo.GetBio(),
			TwitterUsername:         userInfo.GetTwitterUsername(),
			PublicRepos:             userInfo.GetPublicRepos(),
			PublicGists:             userInfo.GetPublicGists(),
			Followers:               userInfo.GetFollowers(),
			Following:               userInfo.GetFollowing(),
			CreatedAt:               userInfo.GetCreatedAt().Time,
			UpdatedAt:               userInfo.GetUpdatedAt().Time,
			SuspendedAt:             userInfo.GetSuspendedAt().Time,
			Type:                    userInfo.GetType(),
			SiteAdmin:               userInfo.GetSiteAdmin(),
			TotalPrivateRepos:       userInfo.GetTotalPrivateRepos(),
			OwnedPrivateRepos:       userInfo.GetOwnedPrivateRepos(),
			PrivateGists:            userInfo.GetPrivateGists(),
			DiskUsage:               userInfo.GetDiskUsage(),
			Collaborators:           userInfo.GetCollaborators(),
			TwoFactorAuthentication: userInfo.GetTwoFactorAuthentication(),
			LdapDn:                  userInfo.GetLdapDn(),
			URL:                     userInfo.GetURL(),
			EventsURL:               userInfo.GetEventsURL(),
			FollowingURL:            userInfo.GetFollowingURL(),
			FollowersURL:            userInfo.GetFollowersURL(),
			GistsURL:                userInfo.GetGistsURL(),
			OrganizationsURL:        userInfo.GetOrganizationsURL(),
			ReceivedEventsURL:       userInfo.GetReceivedEventsURL(),
			ReposURL:                userInfo.GetReposURL(),
			StarredURL:              userInfo.GetStarredURL(),
			SubscriptionsURL:        userInfo.GetSubscriptionsURL(),
			Permissions:             userInfo.GetPermissions(),
			RoleName:                userInfo.GetRoleName(),
		}

		// 创建用户
		err = s.userDAO.CreateUser(user)
		if err != nil {
			return 0, err
		}
		//存储用户到
		s.githubAPI.SetClient(user.ID, client)
	}

	return user.ID, nil
}
