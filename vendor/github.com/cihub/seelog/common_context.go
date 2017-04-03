// Copyright (c) 2012 - Cloud Instruments Co., Ltd.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package seelog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	workingDir     = "/"
	stackCache     map[uintptr]*logContext
	stackCacheLock sync.RWMutex
)

func init() {
	wd, err := os.Getwd()
	if err == nil {
		workingDir = filepath.ToSlash(wd) + "/"
	}
	stackCache = make(map[uintptr]*logContext)
}

// Represents runtime caller context.
type LogContextInterface interface {
	// Caller's function name.
	Func() string
	// Caller's line number.
	Line() int
	// Caller's file short path (in slashed form).
	ShortPath() string
	// Caller's file full path (in slashed form).
	FullPath() string
	// Caller's file name (without path).
	FileName() string
	// True if the context is correct and may be used.
	// If false, then an error in context evaluation occurred and
	// all its other data may be corrupted.
	IsValid() bool
	// Time when log function was called.
	CallTime() time.Time
	// Custom context that can be set by calling logger.SetContext
	CustomContext() interface{}
}

// Returns context of the caller
func currentContext(custom interface{}) (LogContextInterface, error) {
	return specifyContext(1, custom)
}

func extractCallerInfo(skip int) (*logContext, error) {
	var stack [1]uintptr
	if runtime.Callers(skip+1, stack[:]) != 1 {
		return nil, errors.New("error  during runtime.Callers")
	}
	pc := stack[0]

	// do we have a cache entry?
	stackCacheLock.RLock()
	ctx, ok := stackCache[pc]
	stackCacheLock.RUnlock()
	if ok {
		return ctx, nil
	}

	// look up the details of the given caller
	funcInfo := runtime.FuncForPC(pc)
	if funcInfo == nil {
		return nil, errors.New("error during runtime.FuncForPC")
	}

	var shortPath string
	fullPath, line := funcInfo.FileLine(pc)
	if strings.HasPrefix(fullPath, workingDir) {
		shortPath = fullPath[len(workingDir):]
	} else {
		shortPath = fullPath
	}
	funcName := funcInfo.Name()
	if strings.HasPrefix(funcName, workingDir) {
		funcName = funcName[len(workingDir):]
	}

	ctx = &logContext{
		funcName:  funcName,
		line:      line,
		shortPath: shortPath,
		fullPath:  fullPath,
		fileName:  filepath.Base(fullPath),
	}

	// save the details in the cache; note that it's possible we might
	// have written an entry into the map in between the test above and
	// this section, but the behaviour is still correct
	stackCacheLock.Lock()
	stackCache[pc] = ctx
	stackCacheLock.Unlock()
	return ctx, nil
}

// Returns context of the function with placed "skip" stack frames of the caller
// If skip == 0 then behaves like currentContext
// Context is returned in any situation, even if error occurs. But, if an error
// occurs, the returned context is an error context, which contains no paths
// or names, but states that they can't be extracted.
func specifyContext(skip int, custom interface{}) (LogContextInterface, error) {
	callTime := time.Now()
	if skip < 0 {
		err := fmt.Errorf("can not skip negative stack frames")
		return &errorContext{callTime, err}, err
	}
	caller, err := extractCallerInfo(skip + 2)
	if err != nil {
		return &errorContext{callTime, err}, err
	}
	ctx := new(logContext)
	*ctx = *caller
	ctx.callTime = callTime
	ctx.custom = custom
	return ctx, nil
}

// Represents a normal runtime caller context.
type logContext struct {
	funcName  string
	line      int
	shortPath string
	fullPath  string
	fileName  string
	callTime  time.Time
	custom    interface{}
}

func (context *logContext) IsValid() bool {
	return true
}

func (context *logContext) Func() string {
	return context.funcName
}

func (context *logContext) Line() int {
	return context.line
}

func (context *logContext) ShortPath() string {
	return context.shortPath
}

func (context *logContext) FullPath() string {
	return context.fullPath
}

func (context *logContext) FileName() string {
	return context.fileName
}

func (context *logContext) CallTime() time.Time {
	return context.callTime
}

func (context *logContext) CustomContext() interface{} {
	return context.custom
}

// Represents an error context
type errorContext struct {
	errorTime time.Time
	err       error
}

func (errContext *errorContext) getErrorText(prefix string) string {
	return fmt.Sprintf("%s() error: %s", prefix, errContext.err)
}

func (errContext *errorContext) IsValid() bool {
	return false
}

func (errContext *errorContext) Line() int {
	return -1
}

func (errContext *errorContext) Func() string {
	return errContext.getErrorText("Func")
}

func (errContext *errorContext) ShortPath() string {
	return errContext.getErrorText("ShortPath")
}

func (errContext *errorContext) FullPath() string {
	return errContext.getErrorText("FullPath")
}

func (errContext *errorContext) FileName() string {
	return errContext.getErrorText("FileName")
}

func (errContext *errorContext) CallTime() time.Time {
	return errContext.errorTime
}

func (errContext *errorContext) CustomContext() interface{} {
	return nil
}
