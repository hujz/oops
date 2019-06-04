package util

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "[config] ", log.Llongfile|log.LstdFlags|log.Lmicroseconds)

type TConfig struct {
	HttpAddr, ConsoleAddr, ShellDir, MetaDir, SpecsDir string
}

var Config *TConfig

func GetConfig() *TConfig {
	if Config == nil {
		m := build()
		config := &TConfig{}
		for k, v := range m {
			switch k {
			case "http":
				config.HttpAddr = v
			case "console":
				config.ConsoleAddr = v
			case "shell_dir":
				config.ShellDir = v
			case "meta_dir":
				config.MetaDir = v
			case "specs_dir":
				config.SpecsDir = v
			}
		}
		Config = config
	}
	return Config
}

func Build(content string) map[string]string {
	lines := strings.Split(content, "\n")
	m := make(map[string]string)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) >= 1 && l[len(l)-1] == '\r' {
			l = l[:len(l)-2]
		}
		if len(l) == 0 {
			continue
		}
		if l[0] == '#' {
			continue
		}
		n := strings.Index(l, " ")
		var k, v string
		if n != -1 {
			k, v = l[:n], strings.TrimSpace(l[n+1:])
		} else {
			k = l
		}
		m[k] = v
	}
	return m
}

func build() map[string]string {
	file, err := os.Open("config.conf")
	if err != nil {
		logger.Println(err)
		return nil
	}
	str, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Println(err)
		return nil
	}
	return Build(string(str))
}
