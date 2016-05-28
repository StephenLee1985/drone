package sryun

import (
	"net/http"

	"github.com/drone/drone/model"
	//"github.com/drone/drone/remote"
)

type Sryun struct {
	User *model.User
}

func Load(config string) *Sryun {
	sryun := Sryun{}
	return &sryun

}
func (s *Sryun) Login(w http.ResponseWriter, r *http.Request) (*model.User, bool, error) {
	return nil, true, nil
}

// Auth authenticates the session and returns the remote user
// login for the given token and secret
func (s *Sryun) Auth(token, secret string) (string, error) {
	return "", nil

}

// Repo fetches the named repository from the remote system.
func (s *Sryun) Repo(u *model.User, owner, repo string) (*model.Repo, error) {
	return nil, nil
}

// Repos fetches a list of repos from the remote system.
func (s *Sryun) Repos(u *model.User) ([]*model.RepoLite, error) {
	return nil, nil
}

// Perm fetches the named repository permissions from
// the remote system for the specified user.
func (s *Sryun) Perm(u *model.User, owner, repo string) (*model.Perm, error) {
	return nil, nil
}

// File fetches a file from the remote repository and returns in string
// format.
func (s *Sryun) File(u *model.User, r *model.Repo, b *model.Build, f string) ([]byte, error) {
	return nil, nil
}

// Status sends the commit status to the remote system.
// An example would be the GitHub pull request status.
func (s *Sryun) Status(u *model.User, r *model.Repo, b *model.Build, link string) error {
	return nil
}

// Netrc returns a .netrc file that can be used to clone
// private repositories from a remote system.
func (s *Sryun) Netrc(u *model.User, r *model.Repo) (*model.Netrc, error) {
	return nil, nil
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
	return nil, nil, nil
}
