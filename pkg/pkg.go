package pkg

import (
	"github.com/GitEval/GitEval-Backend/pkg/github"
	"github.com/GitEval/GitEval-Backend/pkg/github/expireMap"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	github.NewGitHubAPI,
	expireMap.NewExpireMap, //github
	//llm.NewLLMClient,       //llm
)
