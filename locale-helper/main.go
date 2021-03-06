/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"os/exec"
	"time"

	"sync"

	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

const (
	dbusServiceName = "com.deepin.api.LocaleHelper"
	dbusPath        = "/com/deepin/api/LocaleHelper"
	dbusInterface   = dbusServiceName
)

type Helper struct {
	service *dbusutil.Service
	mu      sync.Mutex
	running bool

	methods *struct {
		SetLocale      func() `in:"locale"`
		GenerateLocale func() `in:"locale"`
	}

	signals *struct {
		/**
		 * if failed, Success(false, reason), else Success(true, "")
		 **/
		Success struct {
			ok     bool
			reason string
		}
	}
}

var (
	logger = log.NewLogger(dbusServiceName)
)

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	service, err := dbusutil.NewSystemService()
	if err != nil {
		logger.Fatal("failed to new system service:", err)
	}

	hasOwner, err := service.NameHasOwner(dbusServiceName)
	if err != nil {
		logger.Fatal(err)
	}
	if hasOwner {
		logger.Fatalf("name %q already has the owner", dbusServiceName)
	}

	var h = &Helper{
		running: false,
		service: service,
	}
	err = service.Export(dbusPath, h)
	if err != nil {
		logger.Fatal("failed to export:", err)
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to request name:", err)
	}

	service.SetAutoQuitHandler(30*time.Second, h.canQuit)
	service.Wait()
}

func (*Helper) GetInterfaceName() string {
	return dbusInterface
}

func (h *Helper) canQuit() bool {
	h.mu.Lock()
	running := h.running
	h.mu.Unlock()
	return !running
}

func (h *Helper) doGenLocale() error {
	return exec.Command("/bin/sh", "-c", "locale-gen").Run()
}

// locales version <= 2.13
func (h *Helper) doGenLocaleWithParam(locale string) error {
	cmd := "locale-gen " + locale
	return exec.Command("/bin/sh", "-c", cmd).Run()
}
