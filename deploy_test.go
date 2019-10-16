package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	testURL := "https://github.com/kelseyhightower/nocode.git"
	testDir, _ := ioutil.TempDir("", "DeployTest")
	defer os.RemoveAll(testDir)

	err := cloneRepo(testURL, testDir)
	if err != nil {
		t.Errorf("Failed to clone repo: %s", err)
	}

	files, err := ioutil.ReadDir(testDir)
	if len(files) == 0 {
		t.Errorf("Didn't get any files")
	}
}

func TestRunCmd(t *testing.T) {
	cmd := "echo hi"
	dir := "/"
	err := runCMD(cmd, dir)
	if err != nil {
		t.Errorf("Couldn't run command: %s", err)
	}
}

func TestLoadConfig(t *testing.T) {
	var dCfg deployConfig
	configFilePath := "test.yml"
	dCfg.loadConfig(configFilePath)
	if dCfg.BuildCMD != "echo build" || dCfg.InstallCMD != "echo install" || dCfg.RepositoryURL != "https://github.com/kelseyhightower/nocode.git" {
		t.Errorf("Didn't load config properly")
	}
}

func TestDeploy(t *testing.T) {
	var dCfg deployConfig
	configFilePath := "test.yml"
	dCfg.loadConfig(configFilePath)
	err := runDeploy(&dCfg)
	if err != nil {
		t.Errorf("Didn't deploy: %s", err)
	}
}
