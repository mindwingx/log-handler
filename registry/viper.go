package registry

import (
	src "github.com/mindwingx/log-handler"
	"github.com/mindwingx/log-handler/constants"
	registry "github.com/spf13/viper"
	"log"
)

type RegAbstraction interface {
	InitRegistry()
	Parse(interface{})
}

type Viper struct {
	*registry.Viper
}

func NewViper() RegAbstraction {
	return &Viper{registry.New()}
}

func (v *Viper) InitRegistry() {
	v.AddConfigPath(src.Root())
	v.SetConfigName(".env")
	v.SetConfigType(constants.EnvFile)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func (v *Viper) Parse(item interface{}) {
	err := v.Unmarshal(&item)
	if err != nil {
		log.Fatal(err)
	}
}
