package requirement

import (
	"os"
	"path/filepath"
)

var authentication Requirement = Requirement{
	Id: "authentication",
	Config: RequirementConfig{
		GetDefault: func() (string, error) {
			filename := "default.yaml"
			content, err := os.ReadFile(filepath.Join(resourceDir, "authentication", "config", filename))

			return string(content), err
		},
	},
}
