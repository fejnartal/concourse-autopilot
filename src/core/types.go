package core

import (
	"github.com/efejjota/concourse-autopilot/concourse"
)

type AutopilotPipeline struct {
	Target_name        string
	Concourse_url      string
	Concourse_team     string
	concourse.Pipeline `yaml:",inline"`
}

type TargetConfig struct {
	Target_name string
	Concourse   struct {
		Url  string
		Team string
	}
	Repository Repository
	Packs      []struct {
		Team string
		Path string
	} `yaml:",omitempty"`
	// ****************************************************************
	// SMELL!!!! TODO - Improve the way we handle Path for target files
	// Currently we store the Path as if it was part of the target file
	// ****************************************************************
	Path string
}

type PackConfig struct {
	Pack_name string

	Repositories []struct {
		Name       string
		Repository `yaml:",inline"`
	}
	Pipelines []struct {
		Name       string
		Manifest   string
		Repository string
		Vars       map[string]interface{} `yaml:",omitempty"`
	}
	// **************************************************************
	// SMELL!!!! TODO - Improve the way we handle Team for pack files
	// Currently we store the team as if it was part of the pack file
	// **************************************************************
	Team string
}

type Repository struct {
	Resource_type struct {
		Name   string
		Type   string                 `yaml:",omitempty"`
		Source map[string]interface{} `yaml:",omitempty"`
	}
	Source map[string]interface{}
}
