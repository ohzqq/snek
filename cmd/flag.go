package cmd

type Flag struct {
	Name       string `yaml:"Name"`
	Shorthand  string `yaml:"Shorthand"`
	Usage      string `yaml:"Usage"`
	Type       string `yaml:"Type"`
	Var        string `yaml:"Var"`
	Value      string `yaml:"Value"`
	Persistent bool   `yaml:"Persistent"`
	Viper      bool   `yaml:"Viper"`
}
