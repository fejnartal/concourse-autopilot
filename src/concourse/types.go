package concourse

type Group struct {
	Name string
	Jobs []string
}

type Pipeline struct {
	Groups         []Group `yaml:",omitempty"`
	Resources      []Resource
	Resource_types []ResourceType `yaml:",omitempty"`
	Jobs           []Job
}

type Job struct {
	Name string
	Plan []interface{}
}

type AnonymousResource struct {
	Type   string
	Source map[string]interface{}
}

type Resource struct {
	Name   string
	Type   string
	Source map[string]interface{}
}

type ResourceType struct {
	Name   string
	Type   string
	Source map[string]interface{}
}

type Get struct {
	Get     string
	Passed  []string               `yaml:",omitempty"`
	Params  map[string]interface{} `yaml:",omitempty"`
	Trigger bool                   `yaml:",omitempty"`
}

type Put struct {
	Put    string
	Params map[string]interface{} `yaml:",omitempty"`
}

type Task struct {
	Task   string
	Config TaskConfig
}

type SetPipeline struct {
	Set_pipeline string
	File         string
	Vars         map[string]interface{} `yaml:",omitempty"`
	Team         string                 `yaml:",omitempty"`
}

type TaskConfig struct {
	Platform       string
	Image_resource AnonymousResource
	Inputs         []Input  `yaml:",omitempty"`
	Outputs        []Output `yaml:",omitempty"`
	Run            Command
}

type Input struct {
	Name     string
	Path     string `yaml:",omitempty"`
	Optional bool   `yaml:",omitempty"`
}

type Output struct {
	Name string
	Path string `yaml:",omitempty"`
}

type Command struct {
	Path string
	Args []string `yaml:",omitempty"`
}
