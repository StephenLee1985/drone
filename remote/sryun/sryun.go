package sryun

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/drone/drone/model"
	"github.com/drone/drone/store"
	"github.com/ianschenck/envflag"
)

const (
	//	login    = "sysadmin"
	//	token    = "SRYUN-ABCD-999"
	//	email    = "sysadmin@dataman-inc.com"
	//	avatar   = "https://avatars3.githubusercontent.com/u/76609?v=3&s=460"
	fullName = "leonlee"
	name     = "docker-2048"
	repoLink = "https://omdev.riderzen.com:10080/leonlee/docker-2048.git"
	clone    = "https://omdev.riderzen.com:10080/leonlee/docker-2048.git"
	branch   = "master"
	//passwd   = "ppppp"
)

type Sryun struct {
	User     *model.User
	Password string
}

func Load(config string) *Sryun {
	log.Infoln("Loading sryun driver...")

	login := envflag.String("RC_SRY_USER", "sryadmin", "")
	password := envflag.String("RC_SRY_PWD", "sryun-pwd", "")
	token := envflag.String("RC_SRY_TOKEN", "EFDDF4D3-2EB9-400F-BA83-4A9D292A1170", "")
	email := envflag.String("RC_SRY_EMAIL", "sryadmin@dataman-inc.net", "")
	avatar := envflag.String("RC_SRY_AVATAR", "https://avatars3.githubusercontent.com/u/76609?v=3&s=460", "")

	user := model.User{}
	user.Token = *token
	user.Login = *login
	user.Email = *email
	user.Avatar = *avatar
	sryun := Sryun{
		User:     &user,
		Password: *password,
	}

	sryunJson, _ := json.Marshal(sryun)

	log.Infoln(string(sryunJson))
	return &sryun

}

// Login authenticates the session and returns the
// remote user details.
func (s *Sryun) Login(w http.ResponseWriter, r *http.Request) (*model.User, bool, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Infoln("got", username, "/", password)

	if username == s.User.Login && password == s.Password {
		return s.User, true, nil
	}
	return nil, false, errors.New("bad auth")

}

// Auth authenticates the session and returns the remote user
// login for the given token and secret
func (s *Sryun) Auth(token, secret string) (string, error) {
	return s.User.Login, nil

}

// Repo fetches the named repository from the remote system.
func (s *Sryun) Repo(u *model.User, owner, name string) (*model.Repo, error) {
	repo := &model.Repo{}
	repo.Owner = owner
	repo.FullName = fullName
	repo.Link = repoLink
	repo.IsPrivate = true
	repo.Clone = clone
	repo.Branch = branch
	repo.Avatar = s.User.Avatar
	repo.Kind = model.RepoGit

	return repo, nil
}

// Repos fetches a list of repos from the remote system.
func (s *Sryun) Repos(u *model.User) ([]*model.RepoLite, error) {
	repo := &model.RepoLite{
		Owner:    s.User.Login,
		Name:     name,
		FullName: fullName,
		Avatar:   s.User.Avatar,
	}
	return []*model.RepoLite{repo}, nil
}

// Perm fetches the named repository permissions from
// the remote system for the specified user.
func (s *Sryun) Perm(u *model.User, owner, repo string) (*model.Perm, error) {
	m := &model.Perm{
		Admin: true,
		Pull:  true,
		Push:  false,
	}

	return m, nil
}

// File fetches a file from the remote repository and returns in string
// format.
func (s *Sryun) File(u *model.User, r *model.Repo, b *model.Build, f string) ([]byte, error) {
	cfg := `
clone:
  skip_verify: true
build:
  image: alpine:latest
  commands:
    - echo 'done'
publish:
  docker:
    username: blackicebird
    password: youman
    email: blackicebird@126.com
    repo: blackicebird/hello-2048
    tag:
      - latest
    load: docker/hello-2048.tar
    save:
      destination: docker/hello-2048.tar
      tag: latest
cache:
  mount:
    - docker/hello-2048.tar	`

	return []byte(cfg), nil
}

// Status sends the commit status to the remote system.
// An example would be the GitHub pull request status.
func (s *Sryun) Status(u *model.User, r *model.Repo, b *model.Build, link string) error {
	return nil
}

// Netrc returns a .netrc file that can be used to clone
// private repositories from a remote system.
func (s *Sryun) Netrc(u *model.User, r *model.Repo) (*model.Netrc, error) {
	netrc := &model.Netrc{}
	return netrc, nil
}

// Activate activates a repository by creating the post-commit hook and
// adding the SSH deploy key, if applicable.
func (s *Sryun) Activate(u *model.User, r *model.Repo, k *model.Key, link string) error {
	return nil
}

// Deactivate removes a repository by removing all the post-commit hooks
// which are equal to link and removing the SSH deploy key.
func (s *Sryun) Deactivate(u *model.User, r *model.Repo, link string) error {
	return nil
}

// Hook parses the post-commit hook from the Request body
// and returns the required data in a standard format.
func (s *Sryun) Hook(r *http.Request) (*model.Repo, *model.Build, error) {
	//		var (
	//		err   error
	//		repo  *model.Repo
	//		build *model.Build
	//	)

	//	switch r.Header.Get("X-Gogs-Event") {
	//	case "push":
	//		var push *PushHook
	//		push, err = parsePush(r.Body)
	//		if err == nil {
	//			repo = repoFromPush(push)
	//			build = buildFromPush(push)
	//		}
	//	}
	//	return repo, build, err

	owner := r.FormValue("owner") //c.Query("owner")
	name := r.FormValue("name")   //c.Query("name")
	force := r.FormValue("force") //c.DefaultQuery("force", "false")
	if len(owner)&len(name) == 0 {
		return nil, nil, errors.New("bad args")
	}

	repo, err := store.GetRepoOwnerName(c, owner, name)
	if err != nil {
		return nil, nil, err
	}

	build, err := buildFromRepo(repo, force)
	if err != nil {
		return nil, nil, err
	}

	return repo, build, nil

}

func buildFromRepo(repo *model.Repo, force string) (*model.Build, error) {
	client, err := git.NewClient(repo.Clone, repo.Branch)
	if err != nil {
		return nil, err
	}
	var filter uint8
	if repo.AllowTag {
		filter = filter + git.FilterTags
	}
	if repo.AllowPush {
		filter = filter + git.FilterHeads
	}

	push, tag, err := client.LsRemote(filter, "")
	if err != nil {
		return nil, err
	}
	log.Println("push", push, "tag", tag)

	//build := &model.Build{
	//    Event:     model.EventPush,
	//    Commit:    hook.After,
	//    Ref:       hook.Ref,
	//    Link:      hook.Compare,
	//    Branch:    strings.TrimPrefix(hook.Ref, "refs/heads/"),
	//    Message:   hook.Commits[0].Message,
	//    Avatar:    avatar,
	//    Author:    hook.Sender.Login,
	//    Timestamp: time.Now().UTC().Unix(),
	//}

	return nil, nil
}
