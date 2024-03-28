package config

type QtumDataSource struct {
	NetType     string `json:"net_type" yaml:"netType"`
	Url         string `json:"url" yaml:"url"`
	AccessToken string `json:"access_token" yaml:"accessToken"`
}
