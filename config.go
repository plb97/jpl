// Copyright (c) 2017 plb97.
// All rights reserved.
// Use of this source code is governed by a CeCILL-B_V1
// (BSD-style) license that can be found in the
// LICENCE (French) or LICENSE (English) file.
package jpl

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"fmt"
	"os"
)

type Config struct {
	InputDir  string    `json:"inputDir"`
	OutputDir string    `json:"outputDir"`
	TchebNmax int       `json:"tchebNmax"`
	SQL       SqlConfig `json:"sqlConfig"`
}
func (c Config) String() string {
	return fmt.Sprintf("config.json\n" +
		"{\n" +
		"\t\"inputDir\":\"%v\",\n" +
		"\t\"outputDir\":\"%v\",\n" +
		"\t\"tchebNmax\":%v,\n" +
		"\t\"sqlConfig\":%v" +
		"}",c.InputDir,c.OutputDir,c.TchebNmax,c.SQL)
}

func LoadConfig() Config {
	configPtr := flag.String("config", "./config.json", "config.json.txt file path")
	flag.Parse()
	config := LoadConfigFile(*configPtr)
	return config
}

func LoadConfigFile(path string) Config {
	var config Config
	file, err := ioutil.ReadFile(path)
	if err != nil {
		const (
			tcheb = 14
			input  = "/input"
			output = "/output"
			host = "127.0.0.1"
			port = 3306
			user = "test"
			pwd = "test"
			db = "test"
		)
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		config = Config{InputDir:wd+input, OutputDir:wd+output,TchebNmax:tcheb,SQL:SqlConfig{Host:host,Port:port,User:user,Pwd:pwd,Db:db}}
		fmt.Println(config.String())
	} else {
		err = json.Unmarshal(file, &config)
		if err != nil {
			panic(err)
		}
	}
	return config
}

var Configuration Config = LoadConfig()

type SqlConfig struct {
	Host string `json:"hostname"`
	Port int    `json:"port"`
	User string `json:"username"`
	Pwd  string `json:"password"`
	Db   string `json:"database"`
}
func (c SqlConfig) String() string {
	return fmt.Sprintf("\t{\n" +
		"\t\t\"hostname\":\"%v\",\n" +
		"\t\t\"port\":%v,\n" +
		"\t\t\"username\":\"%v\",\n" +
		"\t\t\"password\":\"%v\",\n" +
		"\t\t\"database\":\"%v\"\n" +
		"\t}\n",c.Host,c.Port,c.User,c.Pwd,c.Db)
}
