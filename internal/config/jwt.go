package config

type JWT struct {
	AccessTokenLifeTime  int    `config:"access_token_lifetime"`
	RefreshTokenLifeTime int    `config:"refresh_token_lifetime"`
	PublicKeyPath        string `config:"public-key-path"`
	PrivateKeyPath       string `config:"private-key-path"`
}
