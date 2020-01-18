/*
 * MIT License
 *
 * Copyright (c) 2020 Dan <https://github.com/delivrance>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"nekobin/limiter"
	"time"
)

type Nekobin struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`

	MaxTitleLength   int `yaml:"max_title_length"`
	MaxAuthorLength  int `yaml:"max_author_length"`
	MaxContentLength int `yaml:"max_content_length"`
}

type Database struct {
	URI string `yaml:"uri"`

	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type Documents struct {
	Get  []limiter.Limit `yaml:"get"`
	Post []limiter.Limit `yaml:"post"`
}

type Limits struct {
	Documents Documents `yaml:"documents"`
}

type Config struct {
	Nekobin  Nekobin  `yaml:"nekobin"`
	Database Database `yaml:"database"`
	Limits   Limits   `yaml:"limits"`
}

func Load(path string) *Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &Config{}

	err = yaml.UnmarshalStrict(file, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// YAML time values are kept in seconds for convenience.
	// Convert them here to nanoseconds because that's what Limiter needs.
	{
		for i, get := 0, cfg.Limits.Documents.Get; i < len(get); i++ {
			get[i].Period *= time.Second
		}

		for i, post := 0, cfg.Limits.Documents.Post; i < len(post); i++ {
			post[i].Period *= time.Second
		}
	}

	return cfg
}
