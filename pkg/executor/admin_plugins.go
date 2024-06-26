// Copyright 2019 PingCAP, Inc.
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

package executor

import (
	"context"

	"github.com/pingcap/tidb/pkg/domain"
	"github.com/pingcap/tidb/pkg/executor/internal/exec"
	"github.com/pingcap/tidb/pkg/planner/core"
	"github.com/pingcap/tidb/pkg/plugin"
	"github.com/pingcap/tidb/pkg/util/chunk"
)

// AdminPluginsExec indicates AdminPlugins executor.
type AdminPluginsExec struct {
	exec.BaseExecutor
	Action  core.AdminPluginsAction
	Plugins []string
}

// Next implements the Executor Next interface.
func (e *AdminPluginsExec) Next(context.Context, *chunk.Chunk) error {
	switch e.Action {
	case core.Enable:
		return e.changeDisableFlagAndFlush(false)
	case core.Disable:
		return e.changeDisableFlagAndFlush(true)
	}
	return nil
}

func (e *AdminPluginsExec) changeDisableFlagAndFlush(disabled bool) error {
	dom := domain.GetDomain(e.Ctx())
	for _, pluginName := range e.Plugins {
		err := plugin.ChangeDisableFlagAndFlush(dom, pluginName, disabled)
		if err != nil {
			return err
		}
	}
	return nil
}
