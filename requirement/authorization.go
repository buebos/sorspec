package requirement

var authorization Requirement = Requirement{
	Id: "authorization",
	Config: RequirementConfig{
		GetDefault: func() (string, error) {
			return "method: JWT", nil
		},
	},
}
