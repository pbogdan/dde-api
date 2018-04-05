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

// Gtk/Icon/Cursor theme scanner.
package themes

import (
	"os"
	"path"
	"pkg.deepin.io/dde/api/themes/scanner"
	xdg "pkg.deepin.io/lib/xdg/basedir"
)

// Check whether 'theme' in 'list'
func IsThemeInList(theme string, list []string) bool {
	name := path.Base(theme)
	for _, l := range list {
		if path.Base(l) == name {
			return true
		}
	}
	return false
}

// List gtk theme in system.
//
// Scan '/usr/share/themes' and '$HOME/.themes'
func ListGtkTheme() []string {
	sysdirs := []string{}
	for _, dir := range xdg.GetSystemDataDirs() {
		sysdirs = append(sysdirs, path.Join(dir, "themes"))
	}
	return doListTheme(
		[]string{
			path.Join(xdg.GetUserDataDir(), "themes"),
			path.Join(os.Getenv("HOME"), ".themes"),
		},
		sysdirs,
		scanner.ListGtkTheme)
}

// List icon theme in system.
//
// Scan '/usr/share/icons' and '$HOME/.icons'
func ListIconTheme() []string {
	sysdirs := []string{}
	for _, dir := range xdg.GetSystemDataDirs() {
		sysdirs = append(sysdirs, path.Join(dir, "icons"))
	}
	return doListTheme(
		[]string{
			path.Join(xdg.GetUserDataDir(), "icons"),
			path.Join(os.Getenv("HOME"), ".icons"),
		},
		sysdirs,
		scanner.ListIconTheme)
}

// List cursor theme in system.
//
// Scan '/usr/share/icons' and '$HOME/.icons'
func ListCursorTheme() []string {
	sysdirs := []string{}
	for _, dir := range xdg.GetSystemDataDirs() {
		sysdirs = append(sysdirs, path.Join(dir, "icons"))
	}
	return doListTheme(
		[]string{
			path.Join(xdg.GetUserDataDir(), "icons"),
			path.Join(os.Getenv("HOME"), ".icons"),
		},
		sysdirs,
		scanner.ListCursorTheme)
}

func doListTheme(local, sys []string, scanner func(string) ([]string, error)) []string {
	list := scanThemeDirs(local, scanner)
	sysList := scanThemeDirs(sys, scanner)
	return mergeThemeList(list, sysList)
}

func scanThemeDirs(dirs []string, scanner func(string) ([]string, error)) []string {
	var list []string
	for _, d := range dirs {
		tmp, err := scanner(d)
		if err != nil {
			continue
		}
		list = append(list, tmp...)
	}
	return list
}

func mergeThemeList(src, target []string) []string {
	if len(target) == 0 {
		return src
	}

	for _, t := range target {
		if IsThemeInList(t, src) {
			continue
		}
		src = append(src, t)
	}
	return src
}
