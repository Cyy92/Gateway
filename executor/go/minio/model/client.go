package model

type Credential struct {
	Authentications Authentication `yaml:"authentication"`
}

type Authentication struct {
	MinIO MinIO `yaml:"minio,omitempty"`
}

type MinIO struct {
	// minio endpoint
	EndPoint string `yaml:"ep"`

	// minio access account
	AccessKey string `yaml:"accesskey"`
	SecretKey string `yaml:"secretkey"`
}
