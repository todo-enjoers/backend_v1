package config

type JWT struct {
	AccessTokenLifeTime  int    `config:"access_token_lifetime" toml:"access_token_lifetime"`
	RefreshTokenLifeTime int    `config:"refresh_token_lifetime" toml:"refresh_token_lifetime"`
	PublicKeyPath        string `config:"public-key-path" toml:"public_key_path"`
	PrivateKeyPath       string `config:"private-key-path" toml:"private_key_path"`
}
