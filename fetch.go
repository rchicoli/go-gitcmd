package gitcmd

import "github.com/libgit2/git2go"

func (c *Commands) Fetch(workingDir, remoteName string) error {

	repo, err := git.OpenRepository(workingDir)
	if err != nil {
		return err
	}

	remoteURL, err := repo.Remotes.Lookup(remoteName)
	if err != nil {
		return err
	}

	// Fetch changes from remote
	if err := remoteURL.Fetch([]string{}, &FetchOptions, ""); err != nil {
		return err
	}
	return nil

}
