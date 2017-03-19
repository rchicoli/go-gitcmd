package gitcmd

import (
	"os"

	"github.com/libgit2/git2go"
)

// Clone a Repository
func (c *Commands) Clone(url, path string) (*git.Repository, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {

		if err := os.Mkdir(path, 0750); err == nil {
			return nil, err
		}
		if _, err := git.Clone(url, path, CloneOptions); err != nil {
			return nil, err
		}

	} else {
		return nil, err
	}

	return nil, nil
}
