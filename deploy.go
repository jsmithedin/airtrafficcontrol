package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
)

type deployConfig struct {
	RepositoryURL string `yaml:"url"`
	BuildCMD      string `yaml:"build"`
	InstallCMD    string `yaml:"install"`
}

func (c *deployConfig) loadConfig(path string) *deployConfig {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.UnmarshalStrict(yamlFile, c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func runDeploy(dCfg *deployConfig) error {
	log.Println("Starting deploy")

	// Create temp dir
	workingDir, _ := ioutil.TempDir("/tmp", "airtrafficcontrol")
	defer os.RemoveAll(workingDir)
	log.Println(fmt.Sprintf("Created %s", workingDir))

	// Check out repo
	err := cloneRepo(dCfg.RepositoryURL, workingDir)
	if err != nil {
		return err
	}

	// Build
	err = runCMD(dCfg.BuildCMD, workingDir)
	if err != nil {
		return err
	}

	// Install
	err = runCMD(dCfg.InstallCMD, workingDir)
	if err != nil {
		return err
	}

	return nil
}

func runCMD(cmd string, dir string) error {
	log.Println(cmd)
	_ = os.Chdir(dir)
	run := exec.Command("bash", "-c", cmd)
	out, err := run.CombinedOutput()
	if err != nil {
		log.Printf("%s\n", string(out))
		log.Println(err)
		return err
	}

	log.Println(string(out))

	return nil
}

func cloneRepo(url string, targetDirectory string) error {
	log.Println("Cloning repo")
	_, err := git.PlainClone(targetDirectory, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
