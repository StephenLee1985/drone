package sryun

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/drone/drone/model"
	"github.com/drone/drone/remote"
	"github.com/drone/drone/store"
)

const (
	fullName = "leonlee"
	name     = "docker-2048"
	repoLink = "https://omdev.riderzen.com:10080/leonlee/docker-2048.git"
	clone    = "https://omdev.riderzen.com:10080/leonlee/docker-2048.git"
	branch   = "master"
)

// Opts defines configuration options.
type Opts struct {
	Login        string
	Password     string
	Token        string
	Email        string
	Avatar       string
	Workspace    string
	ScriptName   string
	SecName      string
	Registry     string
	Insecure     bool
	Storage      string
	PluginPrefix string

	//	login := env.String("RC_SRY_USER", "sryadmin")
	//	password := env.String("RC_SRY_PWD", "sryun-pwd")
	//	token := env.String("RC_SRY_TOKEN", "EFDDF4D3-2EB9-400F-BA83-4A9D292A1170")
	//	email := env.String("RC_SRY_EMAIL", "sryadmin@dataman-inc.net")
	//	avatar := env.String("RC_SRY_AVATAR", "https://avatars3.githubusercontent.com/u/76609?v=3&s=460")
	//	workspace := env.String("RC_SRY_WORKSPACE", "/var/lib/drone/ws/")
	//	scriptName := env.String("RC_SRY_SCRIPT", ".sryci.yaml")
	//	secName := env.String("RC_SRY_SEC", ".sryci.sec")
	//	registry := env.String("RC_SRY_REG_HOST", "")
	//	insecure := env.Bool("RC_SRY_REG_INSECURE", false)
	//	storage := env.String("DOCKER_STORAGE", "aufs")
	//	pluginPrefix := env.String("PLUGIN_PREFIX", "")
}

//Sryun modelss
type Sryun struct {
	User         *model.User
	Password     string
	Workspace    string
	ScriptName   string
	SecName      string
	Registry     string
	Insecure     bool
	Storage      string
	PluginPrefix string
	//store        store.Store
}

//type client struct {
//	URL         string
//	Machine     string
//	Username    string
//	Password    string
//	PrivateMode bool
//	SkipVerify  bool
//}

// New returns a Remote implementation that integrates with Gogs, an open
// source Git service written in Go. See https://gogs.io/
func New(opts Opts) (remote.Remote, error) {

	user := model.User{}
	user.Token = opts.Token
	user.Login = opts.Login
	user.Email = opts.Email
	user.Avatar = opts.Avatar

	sryun := Sryun{
		User:         &user,
		Password:     opts.Password,
		Workspace:    opts.Workspace,
		ScriptName:   opts.ScriptName,
		SecName:      opts.SecName,
		Registry:     opts.Registry,
		Storage:      opts.Storage,
		Insecure:     opts.Insecure,
		PluginPrefix: opts.PluginPrefix,
		//store:        store,
	}

	sryunJSON, _ := json.Marshal(sryun)
	log.Infoln(string(sryunJSON))

	log.Infoln("loaded sryun remote driver")

	return &sryun, nil

}

// Login authenticates the session and returns the
// remote user details.
func (sry *Sryun) Login(res http.ResponseWriter, req *http.Request) (*model.User, error) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	log.Infoln("got", username, "/", password)

	if username == sry.User.Login && password == sry.Password {
		return sry.User, nil
	}
	return nil, errors.New("bad auth")
}

// Auth authenticates the session and returns the remote user
// login for the given token and secret
func (sry *Sryun) Auth(token, secret string) (string, error) {
	return sry.User.Login, nil
}

// Teams fetches a list of team memberships from the remote system.
func (sry *Sryun) Teams(u *model.User) ([]*model.Team, error) {
	return nil, nil
}

// Repo fetches the named repository from the remote system.
func (sry *Sryun) Repo(u *model.User, owner, name string) (*model.Repo, error) {
	repo := &model.Repo{}
	repo.FullName = fmt.Sprintf("%s/%s", owner, name)
	repo.IsPrivate = true
	repo.Avatar = sry.User.Avatar
	repo.Kind = model.RepoGit
	repo.AllowPull = true
	repo.AllowDeploy = true
	repo.IsTrusted = true

	if !repo.AllowTag && !repo.AllowPush {
		repo.AllowPush = true
	}
	if len(repo.Branch) < 1 {
		repo.Branch = "master"
	}

	return repo, nil
}

// RepoSryun fetches the named repository from the remote system.
func (sry *Sryun) RepoSryun(u *model.User, owner, name string, repo *model.Repo) (*model.Repo, error) {
	repo.FullName = fmt.Sprintf("%s/%s", owner, name)
	repo.IsPrivate = true
	repo.Avatar = sry.User.Avatar
	repo.Kind = model.RepoGit
	repo.AllowPull = true
	repo.AllowDeploy = true
	repo.IsTrusted = true

	if !repo.AllowTag && !repo.AllowPush {
		repo.AllowPush = true
	}
	if len(repo.Branch) < 1 {
		repo.Branch = "master"
	}

	return repo, nil
}

// Repos fetches a list of repos from the remote system.
func (sry *Sryun) Repos(u *model.User) ([]*model.RepoLite, error) {
	repo := &model.RepoLite{
		Owner:    sry.User.Login,
		Name:     name,
		FullName: fullName,
		Avatar:   sry.User.Avatar,
	}
	return []*model.RepoLite{repo}, nil
}

// Perm fetches the named repository permissions from
// the remote system for the specified user.
func (sry *Sryun) Perm(u *model.User, owner, repo string) (*model.Perm, error) {
	m := &model.Perm{
		Admin: true,
		Pull:  true,
		Push:  false,
	}

	return m, nil
}

// File fetches a file from the remote repository and returns in string
// format.
func (sry *Sryun) File(u *model.User, repo *model.Repo, build *model.Build, f string) ([]byte, error) {

	keys, err := sry.store.Keys().Get(repo)
	if err != nil {
		return nil, nil, err
	}
	workDir := fmt.Sprintf("%d_%s_%s", repo.ID, repo.Owner, repo.Name)
	client, err := git.NewClient(sry.Workspace, workDir, repo.Clone, repo.Branch, keys.Private)
	if err != nil {
		return nil, nil, err
	}
	err = client.FetchRef(build.Ref)
	if err != nil {
		return nil, nil, err
	}
	script, err := client.ShowFile(build.Commit, sry.ScriptName)
	if err != nil {
		return nil, nil, err
	}
	sec, err := client.ShowFile(build.Commit, sry.SecName)
	if err != nil {
		sec = nil
	}

	log.Infoln("old script\n", string(script))
	script, err = yaml.GenScript(repo, build, script, sry.Insecure, sry.Registry, sry.Storage, sry.PluginPrefix)
	if err != nil {
		return nil, nil, err
	}

	log.Infoln("script\n", string(script))

	return script, sec, nil
}

// Status sends the commit status to the remote system.
// An example would be the GitHub pull request status.
func (sry *Sryun) Status(u *model.User, r *model.Repo, b *model.Build, link string) error {}

// Netrc returns a .netrc file that can be used to clone
// private repositories from a remote system.
func (sry *Sryun) Netrc(u *model.User, r *model.Repo) (*model.Netrc, error) {}

// Activate activates a repository by creating the post-commit hook.
func (sry *Sryun) Activate(u *model.User, r *model.Repo, link string) error {}

// Deactivate deactivates a repository by removing all previously created
// post-commit hooks matching the given link.
func (sry *Sryun) Deactivate(u *model.User, r *model.Repo, link string) error {}

// Hook parses the post-commit hook from the Request body and returns the
// required data in a standard format.
func (sry *Sryun) Hook(r *http.Request) (*model.Repo, *model.Build, error) {
	params := poller.Params{}
	err := c.Bind(&params)
	if err != nil {
		log.Errorln("bad params")
		return nil, nil, err
	}
	log.Infoln("hook params", params)

	repo, err := sry.store.Repos().GetName(params.Owner + "/" + params.Name)
	if err != nil {
		return nil, nil, err
	}

	push, tag, err := sry.retrieveUpdate(repo)
	if err != nil {
		log.Errorln("retrieve update failed", err)
		return nil, nil, ErrBadRetrieve
	}
	log.Infoln("getting build", repo.ID, "-", branch)
	lastBuild, err := sry.store.Builds().GetLast(repo, branch)
	if err != nil {
		log.Infoln("no build found", err)
	}
	if lastBuild != nil {
		log.Infof("lastBuild %q", *lastBuild)
	}
	build, err := formBuild(lastBuild, repo, push, tag, params.Force)
	if err != nil {
		return nil, nil, err
	}

	return repo, build, nil
}
