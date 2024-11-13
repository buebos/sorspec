package requirement

import "path/filepath"

type RequirementConfig struct {
	/*
		Returns a yaml formatted string with the default configuration for the requirement.
	*/
	GetDefault func() (string, error)
}

type Requirement struct {
	Id        string
	Shorthand string
	Config    RequirementConfig
}

var resourceDir string = filepath.Join("resource", "requirement")

var Requirements map[string]*Requirement = map[string]*Requirement{
	"authentication": &authentication,
	"authorization":  &authorization,
}
