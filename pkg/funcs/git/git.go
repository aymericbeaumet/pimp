// Package git contains Git helper functions (https://git-scm.com/)
package git

import (
	"text/template"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"GitBranches":       GitBranches,
		"GitLocalBranches":  GitLocalBranches,
		"GitOpen":           GitOpen,
		"GitReferences":     GitReferences,
		"GitRemoteBranches": GitRemoteBranches,
		"GitRemotes":        GitRemotes,
		"GitRoot":           GitRoot,
		"GitTags":           GitTags,
	}
}

type Reference struct {
	reference *plumbing.Reference
}

func (r Reference) String() string {
	return r.reference.Name().Short()
}

type Remote struct {
	remote *git.Remote
}

func (r Remote) String() string {
	return r.remote.Config().Name
}

type Tag struct {
	tag *object.Tag
}

func (t Tag) String() string {
	return t.tag.Name
}

type Repository struct {
	path       string
	repository *git.Repository
}

func (r Repository) String() string {
	return r.path
}

func (r Repository) Branches() ([]*Reference, error) {
	iter, err := r.repository.References()
	if err != nil {
		return nil, err
	}

	out := []*Reference{}
	if err := iter.ForEach(func(reference *plumbing.Reference) error {
		name := reference.Name()
		if name.IsBranch() || name.IsRemote() {
			out = append(out, &Reference{reference: reference})
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

func (r Repository) LocalBranches() ([]*Reference, error) {
	iter, err := r.repository.References()
	if err != nil {
		return nil, err
	}

	out := []*Reference{}
	if err := iter.ForEach(func(reference *plumbing.Reference) error {
		name := reference.Name()
		if name.IsBranch() {
			out = append(out, &Reference{reference: reference})
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

func (r Repository) References() ([]*Reference, error) {
	iter, err := r.repository.References()
	if err != nil {
		return nil, err
	}

	out := []*Reference{}
	if err := iter.ForEach(func(reference *plumbing.Reference) error {
		out = append(out, &Reference{reference: reference})
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

func (r Repository) RemoteBranches() ([]*Reference, error) {
	iter, err := r.repository.References()
	if err != nil {
		return nil, err
	}

	out := []*Reference{}
	if err := iter.ForEach(func(reference *plumbing.Reference) error {
		name := reference.Name()
		if name.IsRemote() {
			out = append(out, &Reference{reference: reference})
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

func (r Repository) Remotes() ([]*Remote, error) {
	remotes, err := r.repository.Remotes()
	if err != nil {
		return nil, err
	}

	out := []*Remote{}
	for _, remote := range remotes {
		out = append(out, &Remote{remote: remote})
	}

	return out, nil
}

func (r Repository) Root() string {
	return r.path
}

func (r Repository) Tags() ([]*Tag, error) {
	iter, err := r.repository.Tags()
	if err != nil {
		return nil, err
	}

	var out []*Tag
	if err := iter.ForEach(func(reference *plumbing.Reference) error {
		tag, err := r.repository.TagObject(reference.Hash())
		if err != nil {
			return err
		}
		out = append(out, &Tag{tag: tag})
		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}
