package core

import (
	"github.com/efejjota/concourse-autopilot/concourse"
)

func GenerateAutopilot(apilotConfig TargetConfig, allPackConfigs []PackConfig) AutopilotPipeline {
	basePipeline := generateBasePipeline(apilotConfig, allPackConfigs)

	basePipeline.Groups = append(basePipeline.Groups, GeneratePackGroups(apilotConfig, allPackConfigs)...)

	basePipeline.Resource_types = append(basePipeline.Resource_types, GenerateRepositoryResourceTypes(apilotConfig.Repository)...)
	basePipeline.Resource_types = append(basePipeline.Resource_types, GeneratePackResourceTypes(apilotConfig, allPackConfigs)...)

	basePipeline.Resources = append(basePipeline.Resources, GeneratePackResources(apilotConfig, allPackConfigs)...)

	basePipeline.Jobs = append(basePipeline.Jobs, GeneratePackJobs(apilotConfig, allPackConfigs)...)
	return basePipeline
}

func generateBasePipeline(apilotConfig TargetConfig, allPackConfigs []PackConfig) AutopilotPipeline {
	return AutopilotPipeline{
		Target_name:    apilotConfig.Target_name,
		Concourse_url:  apilotConfig.Concourse.Url,
		Concourse_team: apilotConfig.Concourse.Team,
		Pipeline: concourse.Pipeline{
			Groups: []concourse.Group{
				{
					Name: "autopilot",
					Jobs: []string{"sync-pipelines"},
				},
			},
			Resource_types: []concourse.ResourceType{},
			Resources: []concourse.Resource{
				{
					Name:   "autopilot-repository",
					Type:   apilotConfig.Repository.Resource_type.Name,
					Source: apilotConfig.Repository.Source,
				},
			},
			Jobs: []concourse.Job{
				{
					Name: "sync-pipelines",
					Plan: []interface{}{
						concourse.Get{
							Get:     "autopilot-repository",
							Trigger: true,
						},
						concourse.Task{
							Task: "regenerate-pipeline",
							Config: concourse.TaskConfig{
								Platform: "linux",
								Inputs:   []concourse.Input{{Name: "autopilot-repository"}},
								Outputs:  []concourse.Output{{Name: "generated"}},
								Image_resource: concourse.AnonymousResource{
									Type: "registry-image",
									Source: map[string]interface{}{
										"repository": "efejjota/concourse-autopilot",
									},
								},
								Run: concourse.Command{
									Path: "/go/bin/autopilot",
									Args: []string{
										"gen",
										"autopilot-repository",
										apilotConfig.Path,
										"generated",
									},
								},
							},
						},
						concourse.SetPipeline{
							Set_pipeline: "self",
							File:         "generated/"+apilotConfig.Target_name+".yml",
						},
					},
				},
			},
		},
	}
}

func GenerateRepositoryResourceTypes(repository Repository) []concourse.ResourceType {
	if repository.Resource_type.Type == "" {
		return []concourse.ResourceType{}
	} else {
		return []concourse.ResourceType{
			{
				Name:   repository.Resource_type.Name,
				Type:   repository.Resource_type.Type,
				Source: repository.Resource_type.Source,
			},
		}
	}
}

func GeneratePackGroups(apilotConfig TargetConfig, allPackConfigs []PackConfig) []concourse.Group {
	allPacksGroups := []concourse.Group{}

	for _, packConfig := range allPackConfigs {
		currentPackGroup := concourse.Group{}
		currentPackGroup.Name = packConfig.Pack_name
		currentPackGroup.Jobs = []string{}

		for _, pipeline := range packConfig.Pipelines {
			currentPackGroup.Jobs = append(currentPackGroup.Jobs, "set-"+pipeline.Name)
		}
		allPacksGroups = append(allPacksGroups, currentPackGroup)
	}

	return allPacksGroups
}

func GeneratePackResourceTypes(apilotConfig TargetConfig, allPackConfigs []PackConfig) []concourse.ResourceType {
	allPacksResourceTypes := []concourse.ResourceType{}

	for _, packConfig := range allPackConfigs {
		currentPackResourceTypes := []concourse.ResourceType{}

		for _, repository := range packConfig.Repositories {
			currentRepoResourceType := GenerateRepositoryResourceTypes(repository.Repository)
			currentPackResourceTypes = append(currentPackResourceTypes, currentRepoResourceType...)
		}
		allPacksResourceTypes = append(allPacksResourceTypes, currentPackResourceTypes...)
	}

	return allPacksResourceTypes
}

func GeneratePackResources(apilotConfig TargetConfig, allPackConfigs []PackConfig) []concourse.Resource {
	allPacksResources := []concourse.Resource{}

	for _, packConfig := range allPackConfigs {
		currentPackResources := []concourse.Resource{}

		for _, repository := range packConfig.Repositories {
			currentRepoResource := concourse.Resource{
				Name:   repository.Name,
				Type:   repository.Resource_type.Name,
				Source: repository.Source,
			}
			currentPackResources = append(currentPackResources, currentRepoResource)
		}
		allPacksResources = append(allPacksResources, currentPackResources...)
	}

	return allPacksResources
}

func GeneratePackJobs(apilotConfig TargetConfig, allPackConfigs []PackConfig) []concourse.Job {
	allPacksJobs := []concourse.Job{}

	for _, packConfig := range allPackConfigs {
		currentPackJobs := []concourse.Job{}

		for _, pipeline := range packConfig.Pipelines {
			currentPipelineJob := concourse.Job{
				Name: "set-"+pipeline.Name,
				Plan: []interface{}{
					concourse.Get{
						Get:     pipeline.Repository,
						Trigger: true,
					},
					concourse.Get{
						Get:     "autopilot-repository",
						Trigger: true,
						Passed:  []string{"sync-pipelines"},
					},
					concourse.SetPipeline{
						Set_pipeline: pipeline.Name,
						File:         pipeline.Repository+"/"+pipeline.Manifest,
						Vars:         pipeline.Vars,
						Team:         packConfig.Team,
					},
				},
			}
			currentPackJobs = append(currentPackJobs, currentPipelineJob)
		}
		allPacksJobs = append(allPacksJobs, currentPackJobs...)
	}

	return allPacksJobs
}
