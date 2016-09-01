package main

type Config struct {
	Token		string		`json:"token"`
	MySqlConfig	MySqlConfig	`json:"mysql"`
}

type MySqlConfig struct {
	Address		string		`json:"address"`
	User		string		`json:"user"`
	Password	string		`json:"password"`
	Schema		string		`json:"schema"`
}
