package config

// Database config
type Database struct {
	ServerDB       Mysql `yaml:"serverDB"`
	ApprovalInfoDB Mysql `yaml:"approvalInfoDB"`
	Redis          Redis `yaml:"redis"`
}

// Mysql config
type Mysql struct {
	Driver string `json:"driver" yaml:"driver"`
	Source string `json:"source" yaml:"source"`
}

// Redis config
type Redis struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
}
