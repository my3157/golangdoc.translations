// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package driver implements the core pprof functionality. It can be
// parameterized with a flag implementation, fetch and symbolize
// mechanisms.

// Package driver implements the core pprof functionality. It can be
// parameterized with a flag implementation, fetch and symbolize mechanisms.
package driver

import (
	"bytes"
	"cmd/internal/pprof/commands"
	"cmd/internal/pprof/plugin"
	"cmd/internal/pprof/profile"
	"cmd/internal/pprof/report"
	"cmd/internal/pprof/tempfile"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PProf acquires a profile, and symbolizes it using a profile
// manager. Then it generates a report formatted according to the
// options selected through the flags package.

// PProf acquires a profile, and symbolizes it using a profile manager. Then it
// generates a report formatted according to the options selected through the
// flags package.
func PProf(flagset plugin.FlagSet, fetch plugin.Fetcher, sym plugin.Symbolizer, obj plugin.ObjTool, ui plugin.UI, overrides commands.Commands) error

