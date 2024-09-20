package model

type Credential struct {
	Authentications Authentication `yaml:"authentication"`
}

type Authentication struct {
	MongoDB MongoDB `yaml:"mongodb,omitempty"`
}

type MongoDB struct {
	// mongodb uri
	URI string `yaml:"uri"`

	// mongodb account
	UserName string `yaml:"user"`
	PassWord string `yaml:"pw"`

	// mongodb database name
	Database string `yaml:"db"`
}
