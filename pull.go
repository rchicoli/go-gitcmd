package gitcmd

import (
	"errors"
	"fmt"

	"github.com/libgit2/git2go"
)

func (c *Commands) Pull(workingDir, branchName string) error {
	// branch, err := r.Branch()
	// if err != nil {
	// 	return err
	// }
	// name, err := branch.Name()
	// if err != nil {
	// 	return err
	// }

	repo, err := git.OpenRepository(workingDir)
	if err != nil {
		return err
	}

	// Get remote master
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + branchName)
	if err != nil {
		return err
	}

	remoteBranchID := remoteBranch.Target()
	// Get annotated commit
	annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
	if err != nil {
		return err
	}

	// Do the merge analysis
	// mergeHeads := make([]*git.AnnotatedCommit, 1)
	// mergeHeads[0] = annotatedCommit
	// analysis, _, err := repo.MergeAnalysis(mergeHeads)
	analysis, _, err := repo.MergeAnalysis([]*git.AnnotatedCommit{annotatedCommit})
	if err != nil {
		return err
	}
	fmt.Println("test", analysis)

	// Get repo head
	head, err := repo.Head()
	if err != nil {
		return err
	}

	// http://www.rubydoc.info/github/libgit2/rugged/Rugged/Repository:merge_analysis
	if analysis&git.MergeAnalysisUpToDate != 0 {
		// The given commit is reachable from HEAD,
		// meaning HEAD is up-to-date and no merge needs to be performed.
		return nil

	} else if analysis&git.MergeAnalysisNormal != 0 {
		// A "normal" merge is possible,
		// both HEAD and the given commit have diverged from their common ancestor.
		// The divergent commits must be merged.

		// Just merge changes
		if err := repo.Merge([]*git.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
			return err
		}
		// Check for conflicts
		index, err := repo.Index()
		if err != nil {
			return err
		}

		if index.HasConflicts() {
			return errors.New("Conflicts encountered. Please resolve them.")
		}

		// Make the merge commit
		sig, err := repo.DefaultSignature()
		if err != nil {
			return err
		}

		// Get Write Tree
		treeId, err := index.WriteTree()
		if err != nil {
			return err
		}

		tree, err := repo.LookupTree(treeId)
		if err != nil {
			return err
		}

		localCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}

		remoteCommit, err := repo.LookupCommit(remoteBranchID)
		if err != nil {
			return err
		}

		repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)

		// Clean up
		repo.StateCleanup()

	} else if analysis&git.MergeAnalysisFastForward != 0 {
		// The given commit is a fast-forward from HEAD and no merge needs to be performed.
		// HEAD can simply be set to the given commit.

		// Get remote tree
		remoteTree, err := repo.LookupTree(remoteBranchID)
		if err != nil {
			return err
		}

		// Checkout
		if err := repo.CheckoutTree(remoteTree, nil); err != nil {
			return err
		}

		branchRef, err := repo.References.Lookup("refs/heads/" + branchName)
		if err != nil {
			return err
		}

		// Point branch to the object
		branchRef.SetTarget(remoteBranchID, "")
		if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Unexpected merge analysis result %d", analysis)
	}

	return nil
}
