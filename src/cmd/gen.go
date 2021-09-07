package cmd

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/efejjota/concourse-autopilot/core"
	"gopkg.in/yaml.v2"
)

func Gen() {
	var repoDir, configFilesGlob, outDir string
	switch len(os.Args) {
	case 3:
		repoDir = os.Args[2]
		configFilesGlob = ".autopilot/targets/*.yml"
		outDir = path.Join(repoDir, ".generated")
	case 4:
		repoDir = os.Args[2]
		configFilesGlob = os.Args[3]
		outDir = path.Join(repoDir, ".generated")
	case 5:
		repoDir = os.Args[2]
		configFilesGlob = os.Args[3]
		outDir = os.Args[4]
	}

	configFiles, _ := filepath.Glob(path.Join(repoDir, configFilesGlob))

	for _, configPath := range configFiles {
		targetConfig, err := readTargetConfig(configPath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		// ****************************************************************
		// SMELL!!!! TODO - Improve the way we handle Path for target files
		// Currently we store the Path as if it was part of the target file
		// ****************************************************************
		targetConfig.Path, _ = filepath.Rel(repoDir, configPath)
		allPacksConfig := readAllPacksConfig(targetConfig, repoDir)

		autopilotPipeline, err := generatePipeline(targetConfig, allPacksConfig)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		err = writePipeline(autopilotPipeline, outDir)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
}

func readTargetConfig(configPath string) (core.TargetConfig, error) {
	var targetConfig core.TargetConfig

	configYaml, err := os.ReadFile(configPath)
	if err != nil {
		return targetConfig, err
	}

	targetConfig = core.TargetConfig{}
	err = yaml.Unmarshal([]byte(configYaml), &targetConfig)
	if err != nil {
		return targetConfig, err
	}

	return targetConfig, nil
}

func readPackConfig(configPath string) (core.PackConfig, error) {
	var packConfig core.PackConfig

	configYaml, err := os.ReadFile(configPath)
	if err != nil {
		return packConfig, err
	}

	packConfig = core.PackConfig{}
	err = yaml.Unmarshal([]byte(configYaml), &packConfig)
	if err != nil {
		return packConfig, err
	}

	return packConfig, nil
}

func generatePipeline(targetConfig core.TargetConfig, allPackConfigs []core.PackConfig) (core.AutopilotPipeline, error) {
	generatedPipeline := core.GenerateAutopilot(targetConfig, allPackConfigs)
	return generatedPipeline, nil
}

func writePipeline(generatedPipeline core.AutopilotPipeline, outDir string) error {
	generatedPipelineYaml, err := yaml.Marshal(generatedPipeline)
	if err != nil {
		return err
	}

	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(outDir, generatedPipeline.Target_name+".yml"), generatedPipelineYaml, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func readAllPacksConfig(targetConfig core.TargetConfig, repoDir string) []core.PackConfig {
	allPacksConfig := []core.PackConfig{}
	for _, pack := range targetConfig.Packs {
		packFiles, _ := filepath.Glob(path.Join(repoDir, pack.Path))

		for _, packFile := range packFiles {
			packConfig, err := readPackConfig(packFile)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			// **************************************************************
			// SMELL!!!! TODO - Improve the way we handle Team for pack files
			// Currently we store the team as if it was part of the pack file
			// **************************************************************
			packConfig.Team = pack.Team
			allPacksConfig = append(allPacksConfig, packConfig)
		}
	}
	return allPacksConfig
}
