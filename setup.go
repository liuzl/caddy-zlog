// Copyright 2015 Light Code Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zlog

import (
	"flag"
	"os"
	"path/filepath"
	"sync"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

type Config map[string]string

func init() {
	caddy.RegisterPlugin("zlog", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new mime middleware instance.
func setup(c *caddy.Controller) error {

	if !flag.Parsed() {
		flag.Parse()
	}

	configs, err := parse(c)
	if err != nil {
		return err
	}
	var once sync.Once

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		if log_dir, ok := configs["log_dir"]; ok {
			return WithLog(ZLog{Next: next}, log_dir, once)
		}
		log_dir := filepath.Join(filepath.Dir(os.Args[0]), "zerolog")
		return WithLog(ZLog{Next: next}, log_dir, once)
	})

	return nil
}

func parse(c *caddy.Controller) (Config, error) {
	configs := Config{}

	for c.Next() {
		// At least one extension is required

		args := c.RemainingArgs()
		switch len(args) {
		case 2:
			configs[args[0]] = args[1]
		case 1:
			return configs, c.ArgErr()
		case 0:
			for c.NextBlock() {
				ext := c.Val()
				if !c.NextArg() {
					return configs, c.ArgErr()
				}
				configs[ext] = c.Val()
			}
		}

	}

	return configs, nil
}
