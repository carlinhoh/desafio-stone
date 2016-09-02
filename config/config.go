package config

type Config struct {
	Token		string		`json:"token"`
	Port		string		`json:"port"`
	MySqlConfig	MySqlConfig	`json:"mysql"`
}

type MySqlConfig struct {
	Address		string		`json:"address"`
	User		string		`json:"user"`
	Password	string		`json:"password"`
	Schema		string		`json:"schema"`
}

var Settings Config = Config{}