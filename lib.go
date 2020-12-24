package gitstore

import (
	"fmt"
	"strings"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func Read(branch, filename string) (string, error) {
	target := fmt.Sprintf("refs/gitstore/%s:%s", branch, filename)
	return git("show", target)
}

func git(args ...string) (string, error) {
	out, err := exec.Command("git", args...).Output()
	return strings.TrimSpace(string(out)), err
}

func Write(branch, filename, value string) error {
	refbranch := "refs/gitstore/"
	ref := refbranch + branch
	refbranchpath := ".git/" + refbranch
	refpath := ".git/" + ref
	if e := os.MkdirAll(refbranchpath, os.ModePerm); e != nil {
		return errors.Wrap(e, "could not ensure refs path exists")
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "gitstore-")
	if err != nil {
		return errors.Wrap(err, "could not create a temporary file to work with")
	}
	defer os.Remove(tmpFile.Name())
	if _, e := tmpFile.Write([]byte(value)); e != nil {
		return errors.Wrap(e, "could not write data to temporary file")
	}

	hashed_value, err := git("hash-object", "-w", tmpFile.Name())
	if err != nil {
		return errors.Wrap(err, "could not hash the data")
	}

	hashed_stage, err := git("update-index", "--add", "--cacheinfo", fmt.Sprintf("100644,%s,%s", string(hashed_value), filename))
	if err != nil {
		return errors.Wrap(err, "could not stage the hashed data")
	}

	hashed_tree, err := git("write-tree", string(hashed_stage))
	if err != nil {
		return errors.Wrap(err, "could not write hashed stage to a tree")
	}

	if _, err := os.Stat(refpath); err == nil {

		hashed_parent, err := ioutil.ReadFile(refpath)
		if err != nil {
			return errors.Wrap(err, "could not read parent commit")
		}
		hashed_commit, err := git("commit-tree", string(hashed_tree), "-m", "update", "-p", string(hashed_parent))
		if err != nil {
			return errors.Wrap(err, "could not write commit ontop of the parent")
		}
		_, err = git("update-ref", ref, string(hashed_commit))
		return errors.Wrap(err, "could not update the branch ref to the new commit")
	} else if os.IsNotExist(err) {

		hashed_commit, err := git("commit-tree", string(hashed_tree), "-m", "init")
		if err != nil {
			return errors.Wrap(err, "could not create the first commit")
		}
		err = ioutil.WriteFile(refpath, []byte(hashed_commit), 0644)
		return errors.Wrap(err, "failed to write commit hash to branch ref")

	} else {
		return errors.Wrap(err, "i am confused by the branch ref existance state")
	}
}
