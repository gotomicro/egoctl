// Copyright 2013 bee authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/gotomicro/egoctl/cmd"
	_ "github.com/gotomicro/egoctl/cmd/migrate"
	_ "github.com/gotomicro/egoctl/cmd/pb"
	_ "github.com/gotomicro/egoctl/cmd/run"
	_ "github.com/gotomicro/egoctl/cmd/version"
	_ "github.com/gotomicro/egoctl/cmd/web"
	"github.com/gotomicro/egoctl/internal/config"
)

func main() {
	config.LoadConfig()
	err := cmd.RootCommand.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}
