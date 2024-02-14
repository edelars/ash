package file_system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ash/internal/commands"
	"ash/internal/dto"
)

const (
	constManagerName = "Filesystem"
	constDirDisplay  = "dir "
	constFileDisplay = "file"
)

type fileSystemManager struct {
	inputSet func(r []rune)
}

func NewFileSystemManager(inputSet func(r []rune)) fileSystemManager {
	return fileSystemManager{inputSet: inputSet}
}

// PrecisionSearch = true:
//
//	search for exec function
//	"ls" - ls, "ct" - nil
//
// PrecisionSearch = false:
//
//	search for autocomplete
//	Examples:
//	"ls" - ls , "ct" - cat, "/var/l" - var/log, "var/lo" - var/log
//	.
func (m *fileSystemManager) SearchCommands(iContext dto.InternalContextIface, resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	var data []dto.CommandIface
	paths := preparePathArr(iContext.GetEnv("$PATH"))

	defer func() {
		commandManager := commands.NewCommandManager(constManagerName, 100, false, data...)
		commandManager.SearchCommands(iContext, resultChan, patterns...)
	}()

	if len(paths) == 0 {
		return
	}

	for _, pattern := range patterns {
		filesResults := getFileNamesInDirs(paths, false, getFileNamesInDir)
		filesResults = append(filesResults, getFileNamesInDirs([]string{iContext.GetCurrentDir()}, true, getFileNamesInDir)...)

		if !pattern.IsPrecisionSearch() {
			for _, item := range findPathsByPrefixPath(pattern.GetPattern()) {
				s := constFileDisplay
				if item.isDir {
					s = constDirDisplay
				}
				c := commands.NewPseudoCommand(item.name, m.inputSet, generateDescription(s, item.info), item.name)
				c.SetMathWeight(100)
				data = append(data, c)
			}
		}

		for _, fileResItem := range filesResults {
			for _, fInfo := range fileResItem.files {
				if pattern.IsPrecisionSearch() { // exec
					if pattern.GetPattern() == fInfo.name {
						data = append(data, NewSystemCommand(fInfo.name, ""))
					}
				} else { // autocomplete
					displayName := filepath.Join(fileResItem.dir, fInfo.name)
					c := commands.NewPseudoCommand(fInfo.name, m.inputSet, generateDescription(constFileDisplay, fInfo.info), displayName)
					data = append(data, c)
				}
			}
		}
	}
}

func preparePathArr(path string) (res []string) {
	for _, v := range strings.Split(path, ":") {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			res = append(res, v)
		}
	}
	return res
}

func getFileNamesInDir(dir string, skipDirs bool) (res []fileInfo) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return res
	}

	for _, v := range files {
		if v.IsDir() && skipDirs {
			continue
		}
		var info string
		p, err := v.Info()
		if err == nil {
			info = p.Mode().String()
		}
		res = append(res, fileInfo{v.IsDir(), v.Name(), info})
	}
	return res
}

func getFileNamesInDirs(dirs []string, skipDirs bool, searchFunc func(dir string, skipDirs bool) []fileInfo) (res []filesResult) {
	for _, v := range dirs {
		res = append(res, filesResult{dir: v, files: searchFunc(v, skipDirs)})
	}
	return res
}

type filesResult struct {
	dir   string
	files []fileInfo
}

type fileInfo struct {
	isDir bool
	name  string
	info  string
}

// var/log /var/log /var/lo var/lo
func findPathsByPrefixPath(pattern string) (paths []fileInfo) {
	pattern = filepath.Clean(pattern)
	if i, err := os.Lstat(pattern); err == nil {
		paths = append(paths, fileInfo{i.IsDir(), pattern, i.Mode().String()})
	}

	if _, err := os.Lstat(pattern); err != nil {
		pattern = filepath.Dir(pattern)
	}

	r := getFileNamesInDir(pattern, false)
	for _, v := range r {
		paths = append(paths, fileInfo{v.isDir, filepath.Join(pattern, v.name), v.info})
	}

	return paths
}

func generateDescription(constStr, info string) string {
	if info == "" {
		return constStr
	} else {
		return fmt.Sprintf("%s %s", constStr, info)
	}
}
