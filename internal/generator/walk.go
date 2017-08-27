// Copyright 2016 Andreas Pannewitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	aDot = "."
)

func nameLessExt(path string) string {
	name := filepath.Base(path)
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// pathIs - given some path, returns a bool
type pathIs func(path string) bool

type pathDo func(path string) error

// matchBool - iff flag
func matchBool(flag bool) pathIs {
	return func(path string) bool {
		return flag
	}
}

// matchFunc - iff filename matches any pattern
func matchFunc(pattern ...string) pathIs {
	return func(path string) (matched bool) {
		for i := range pattern {
			matched, _ = filepath.Match(pattern[i], filepath.Base(path)) // ignore errors
			if matched {
				break
			}
		}
		return matched
	}
}

func pathPiler(match pathIs, pile *Pile) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if match(path) {
			pile.Pile(path)
		}
		return nil
	}
}

func namePiler(match pathIs, pile *Pile) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if match(path) {
			pile.Pile(nameLessExt(path))
		}
		return nil
	}
}

func ifFlagSkipDirWf(match pathIs) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if match(path) {
			return nil
		}
		return filepath.SkipDir
	}
}

func (t *toDo) isDirWf(dirWf, filWf filepath.WalkFunc) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		switch {
		case t.ctx.Err() != nil:
			return filepath.SkipDir
		case info.IsDir():
			return dirWf(path, info, err)
		default:
			return filWf(path, info, err)
		}
	}
}