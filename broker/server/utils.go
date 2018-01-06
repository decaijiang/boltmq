// Copyright 2017 luoji

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package server

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"

	"github.com/pquerna/ffjson/ffjson"
)

// create2WriteFile 创建文件
func CreateFile(filePath string) (bool, error) {
	err := ensureDir(filePath)
	if err != nil {
		return false, err
	}

	_, err = os.Create(filePath)
	if err != nil {
		return false, err
	}

	return true, nil
}

func String2File(data []byte, flPath string) error {
	filePath := filepath.FromSlash(flPath)
	tmpFilePath := filePath + ".tmp"
	bakFilePath := filePath + ".bak"

	err := create2WriteFile(data, tmpFilePath)
	if err != nil {
		return err
	}

	oldData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = create2WriteFile(oldData, bakFilePath)
	if err != nil {
		return err
	}

	// 删除原文件
	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	// 重命临时文件
	return os.Rename(tmpFilePath, filePath)
}

func create2WriteFile(data []byte, filePath string) error {
	err := ensureDir(filePath)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func ensureDir(filePath string) error {
	dir := path.Dir(filePath)
	exist, err := pathExists(dir)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	return os.MkdirAll(dir, os.ModePerm)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// min int64 的最小值
// Author rongzhihong
// Since 2017/9/18
func min(a, b int64) int64 {
	if a >= b {
		return b
	}
	return a
}

var blankReg = regexp.MustCompile(`\S+?`)

// IsBlank 是否为空
// Author: rongzhihong, <rongzhihong@gome.com.cn>
// Since: 2017/9/19
func IsBlank(content string) bool {
	if blankReg.FindString(content) != "" {
		return false
	}
	return true
}

// Encode Json Encode
// Author: rongzhihong, <rongzhihong@gome.com.cn>
// Since: 2017/9/19
func Encode(v interface{}) ([]byte, error) {
	return ffjson.Marshal(v)
}

// Decode Json Decode
// Author: rongzhihong, <rongzhihong@gome.com.cn>
// Since: 2017/9/19
func Decode(data []byte, v interface{}) error {
	return ffjson.Unmarshal(data, v)
}

// callShell 执行命令
// Author rongzhihong
// Since 2017/9/8
func CallShell(shellString string) error {
	process := exec.Command(shellString)
	err := process.Start()
	if err != nil {
		return err
	}

	err = process.Wait()
	if err != nil {
		return err
	}

	return nil
}
