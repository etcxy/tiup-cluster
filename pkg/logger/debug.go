// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pingcap-incubator/tiup-cluster/pkg/colorutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var debugBuffer *bytes.Buffer

func newDebugLogCore() zapcore.Core {
	debugBuffer = new(bytes.Buffer)
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	return zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(debugBuffer)), zapcore.DebugLevel)
}

// OutputDebugLog outputs debug log in the current working directory.
func OutputDebugLog() {
	if err := os.MkdirAll("./logs", 0755); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\nCreate debug logs directory failed %v.\n", err)
		return
	}

	// FIXME: Stupid go does not allow writing fraction seconds without a leading dot.
	fileName := time.Now().Format("./logs/tiup-cluster-debug-2006-01-02-15-04-05.log")
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		filePath = fileName
	}

	err = ioutil.WriteFile(filePath, debugBuffer.Bytes(), 0644)
	if err != nil {
		_, _ = colorutil.ColorWarningMsg.Fprint(os.Stderr, "\nWarn: Failed to write error debug log.\n")
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "\nVerbose debug logs has been written to %s.\n", colorutil.ColorKeyword.Sprint(filePath))
	}
	debugBuffer.Reset()
}
