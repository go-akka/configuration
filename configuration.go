package configuration

import (
	"io/ioutil"

	"github.com/go-akka/configuration/hocon"
)

func ParseString(text string, includeCallback hocon.IncludeCallback) *hocon.Config {
	root := hocon.Parse(text, includeCallback)
	return hocon.NewConfigFromRoot(root)
}

func LoadConfig(filename string) *hocon.Config {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return ParseString(string(data), defaultIncludeCallback)
}

func defaultIncludeCallback(filename string) *hocon.HoconRoot {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return hocon.Parse(string(data), defaultIncludeCallback)

}
