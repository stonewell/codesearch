package vcs_ignore

import (
	"log"
	"os"
	"errors"
	"path/filepath"
	gitignore "github.com/sabhiram/go-gitignore"
)

type VCSIgnore struct {
	parent *VCSIgnore
	path string
	patterns []*gitignore.GitIgnore
}

func NewVCSIgnore(path string, parent *VCSIgnore) *VCSIgnore {
	vcs_ignore := &VCSIgnore{
		parent: parent,
	}

	abs_path, err := filepath.Abs(path)

	if err != nil {
		log.Printf("convert abs failed: %s, %s", path, err)
		vcs_ignore.path = abs_path
	} else {
		vcs_ignore.path = path
	}

	vcs_ignore.LoadVCSIgnoreFilesInPath(vcs_ignore.path)

	return vcs_ignore
}

func (vcs_ignore *VCSIgnore)LoadVCSIgnoreFilesInPath(path string) {
	IGNORE_FILES := []string {
		".ignore",
			".gitignore",
			".git/info/exclude",
			".hgignore",
		}

	for _, f := range(IGNORE_FILES) {
		ignore_file := filepath.Join(path, f)

		log.Printf("try to load ignore file:%s", ignore_file)
		object, err := gitignore.CompileIgnoreFile(ignore_file)

		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Printf("load vcs ignore file %s failed: %s", path, err)
			}
			continue
		}

		log.Printf("load ignore file:%s success", ignore_file)
		vcs_ignore.patterns = append(vcs_ignore.patterns, object)
	}
}

func (vcs_ignore *VCSIgnore)ShouldIgnorePath(path string) bool {
	abs_path, err := filepath.Abs(path)

	if err != nil {
		log.Printf("convert abs failed: %s, %s", path, err)
		abs_path = path
	}

	rel_path, err := filepath.Rel(vcs_ignore.path, abs_path)

	if err != nil {
		log.Printf("try to get rel path failed:%s, %s", vcs_ignore.path, path)
		rel_path = path
	}

	for _, pattern := range(vcs_ignore.patterns) {
		if pattern.MatchesPath(rel_path) {
			return true
		}
	}

	if vcs_ignore.parent != nil {
		return vcs_ignore.parent.ShouldIgnorePath(path)
	}

	return false
}
