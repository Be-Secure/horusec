// Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package testutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/onsi/ginkgo"
)

type HorusecCmd string

func (h HorusecCmd) String() string {
	return string(h)
}

const (
	VersionCmd  HorusecCmd = "version"
	StartCmd    HorusecCmd = "start"
	GenerateCmd HorusecCmd = "generate"
)

func GinkgoGetHorusecCmd(horusecCmd HorusecCmd) *exec.Cmd {
	bin := ginkgoBuildHorusecBinary()
	args := setLogLevelArgsToHorusecCmd(horusecCmd.String())
	return exec.Command(bin, args...)
}

func GinkgoGetHorusecCmdWithFlags(cmdArg HorusecCmd, flags map[string]string) *exec.Cmd {
	bin := ginkgoBuildHorusecBinary()
	args := setLogLevelArgsToHorusecCmd(cmdArg.String())
	for flag, value := range flags {
		args = append(args, fmt.Sprintf("%s=%s", flag, value))
	}
	return exec.Command(bin, args...)
}

func ginkgoBuildHorusecBinary(customArgs ...string) string {
	binary := filepath.Join(os.TempDir(), getBinaryNameBySystem())
	args := []string{
		"build",
		`-ldflags=-X 'github.com/ZupIT/horusec/cmd/app/version.Version=vTest'`,
		fmt.Sprintf("-o=%s", binary), filepath.Join(RootPath, "cmd", "app"),
	}
	args = append(args, customArgs...)
	cmd := exec.Command("go", args...)
	err := cmd.Run()
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Error on build Horusec binary for e2e test %v", err))
	}
	return binary
}

func setLogLevelArgsToHorusecCmd(horusecCmd ...string) []string {
	return append(horusecCmd, fmt.Sprintf("%s=%s", "--log-level", "debug"))
}

func getBinaryNameBySystem() string {
	binaryName := "horusec-e2e"
	if runtime.GOOS == "windows" {
		binaryName = fmt.Sprintf("%s.exe", binaryName)
	}
	return binaryName
}