package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

const repoURL = "https://robot:foo@gitea.vrutkovs.eu/vadim/cluster-state-stage.git"
const gitAuthorEmail = "robot@vrutkovs.eu"
const gitAuthorName = "Vadim's robot"

// Clone the repo and return the path
func cloneRepo() (string, error) {
	var output string
	t, err := ioutil.TempDir("", "git")
	fmt.Println("Using temp dir", t)
	if err != nil {
		return "", err
	}
	fmt.Println("Cloning", repoURL)
	if output, err := exec.Command("git", "clone", repoURL, t).CombinedOutput(); err != nil {
		return "", fmt.Errorf("git repo clone error: %v. output: %s", err, string(output))
	}
	fmt.Print(string(output))
	return t, nil
}

// Remove all files in repo
func cleanRepo(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.Name() == ".git" {
			continue
		}
		err := os.RemoveAll(f.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// Config git repo
func gitConfig(path string) error {
	fmt.Println("Configuring git user")
	cmd := exec.Command("git", "config", "user.name", gitAuthorName)
	cmd.Dir = path
	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git config user.name error: %v. output: %s", err, string(o))
	}
	cmd = exec.Command("git", "config", "user.email", gitAuthorEmail)
	cmd.Dir = path
	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git config user.name error: %v. output: %s", err, string(o))
	}
	return nil
}

// Make a git commit
func gitCommit(path string) error {
	var output string

	fmt.Println("Committing changes")
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = path
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add error: %v. output: %s", err, string(output))
	}

	t := time.Now()
	date := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	cmd = exec.Command("git", "commit", "-am", fmt.Sprintf("Cluster state on %s", date))
	cmd.Dir = path
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit error: %v. output: %s", err, string(output))
	}
	fmt.Print(string(output))
	return nil
}

// Push git commit
func gitPush(path string) error {
	fmt.Println("Pushing changes")
	cmd := exec.Command("git", "push", "-f")
	cmd.Dir = path
	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push error: %v. output: %s", err, string(o))
	}
	fmt.Println("Changes pushed")
	return nil
}
