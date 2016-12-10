package configuration

import (
	"github.com/go-akka/configuration/hocon"
	"io/ioutil"
)

func ParseString(text string, includeCallback hocon.IncludeCallback) hocon.Config {
	res := hocon.Parse(text, includeCallback)
	return hocon.NewConfigFromRoot(res)
}

func LoadConfig(filename string) hocon.Config {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return ParseString(string(data), nil)
}
