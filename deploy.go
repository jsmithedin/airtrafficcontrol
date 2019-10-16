package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
)

type DeployConfig struct {
	repositoryURL string
	buildCMD      string
	installCMD    string
}

func (c *DeployConfig) loadConfig(path string) *DeployConfig {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Print(err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Print(err)
	}

	return c
}

func runDeploy(deployConfig *DeployConfig) error {
	// Create temp dir
	workingDir := os.TempDir()
	defer os.RemoveAll(workingDir)

	// Check out repo
	err := cloneRepo(deployConfig.repositoryURL, workingDir)
	if err != nil {
		return err
	}

	// Build
	err = runCMD(deployConfig.buildCMD, workingDir)
	if err != nil {
		return err
	}

	// Install
	err = runCMD(deployConfig.installCMD, workingDir)
	if err != nil {
		return err
	}

	return nil
}

func runCMD(cmd string, dir string) error {
	_ = os.Chdir(dir)
	run := exec.Command(cmd)
	_, err := run.Output()
	if err != nil {
		return err
	}

	return nil
}

func cloneRepo(url string, targetDirectory string) error {
	_, err := git.PlainClone(targetDirectory, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	return nil
}
