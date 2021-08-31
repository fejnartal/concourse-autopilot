package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/efejjota/concourse-autopilot/core"
	"gopkg.in/yaml.v2"
)

func Fly() {
	if len(os.Args) != 3 {
		fmt.Println("Wrong parameters")
	}
	glob := os.Args[2]
	files, _ := filepath.Glob(glob)

	flyPath, _ := exec.LookPath("fly")
	for _, filePath := range files {
		apilotPipeline, err := readPipeline(filePath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		flySetPipelineAutorepair(flyPath, apilotPipeline, filePath)
	}
}

func readPipeline(filePath string) (core.AutopilotPipeline, error) {
	var apilotPipeline core.AutopilotPipeline

	pipelineFile, err := os.ReadFile(filePath)
	if err != nil {
		return apilotPipeline, err
	}

	apilotPipeline = core.AutopilotPipeline{}
	err = yaml.Unmarshal([]byte(pipelineFile), &apilotPipeline)
	if err != nil {
		return apilotPipeline, err
	}

	return apilotPipeline, err
}

func flySync(flyPath string, autopilotPipeline core.AutopilotPipeline) error {
	cmd := exec.Command(flyPath, "sync", "-c", autopilotPipeline.Concourse_url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err
}

func flyLogin(flyPath string, autopilotPipeline core.AutopilotPipeline) error {
	cmd := exec.Command(flyPath, "-t", autopilotPipeline.Target_name, "login", "-c", autopilotPipeline.Concourse_url, "--team-name", autopilotPipeline.Concourse_team)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err
}

func flySetPipeline(flyPath string, autopilotPipeline core.AutopilotPipeline, filePath string) error {
	cmd := exec.Command(flyPath, "-t", autopilotPipeline.Target_name, "set-pipeline", "-p", "autopilot", "-c", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err
}

func flyEditTarget(flyPath string, autopilotPipeline core.AutopilotPipeline) error {
	cmd := exec.Command(flyPath, "-t", autopilotPipeline.Target_name, "edit-target", "--concourse-url", autopilotPipeline.Concourse_url, "--team-name", autopilotPipeline.Concourse_team)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err
}

func flySetPipelineAutorepair(flyPath string, autopilotPipeline core.AutopilotPipeline, filePath string) error {
	err := flyEditTarget(flyPath, autopilotPipeline)
	if err != nil {
		err = flyLogin(flyPath, autopilotPipeline)
		if err != nil {
			return err
		}
	}
	err = flySetPipeline(flyPath, autopilotPipeline, filePath)
	if err != nil {
		err = flyLogin(flyPath, autopilotPipeline)
		if err != nil {
			return err
		}
		err = flySync(flyPath, autopilotPipeline)
		if err != nil {
			return err
		}
	}
	return nil
}
