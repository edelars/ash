package file_system

import (
	"os"
	"path/filepath"
	"strings"

	"ash/internal/commands"
	"ash/internal/dto"
)

const constManagerName = "Filesystem"

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
		commandManager := commands.NewCommandManager(constManagerName, 100, data...)
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
				c := NewPseudoCommand(item, m.inputSet)
				c.SetMathWeight(100)
				// c.SetDisplayName(filepath.Join(fileResItem.dir, fileName))
				data = append(data, c)
			}
		}

		for _, fileResItem := range filesResults {
			for _, fileName := range fileResItem.files {
				if pattern.IsPrecisionSearch() { // exec
					if pattern.GetPattern() == fileName {
						data = append(data, NewSystemCommand(fileName))
					}
				} else { // autocomplete
					c := NewPseudoCommand(fileName, m.inputSet)
					c.SetDisplayName(filepath.Join(fileResItem.dir, fileName))
					data = append(data, c)
				}
			}
		}
	}
}

func prepareData(filesResults []filesResult, getData func(dir string, skipDirs bool) []string, patterns ...dto.PatternIface) {
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

func getFileNamesInDir(dir string, skipDirs bool) (res []string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return res
	}

	for _, v := range files {
		if v.IsDir() && skipDirs {
			continue
		}
		res = append(res, v.Name())
	}
	return res
}

func getFileNamesInDirs(dirs []string, skipDirs bool, searchFunc func(dir string, skipDirs bool) []string) (res []filesResult) {
	for _, v := range dirs {
		res = append(res, filesResult{dir: v, files: searchFunc(v, skipDirs)})
	}
	return res
}

type filesResult struct {
	dir   string
	files []string
}

// var/log /var/log /var/lo var/lo
func findPathsByPrefixPath(pattern string) (paths []string) {
	pattern = filepath.Clean(pattern)
	if _, err := os.Lstat(pattern); err == nil {
		paths = append(paths, pattern)
	}

	var res []string
	if _, err := os.Lstat(pattern); err == nil {
		res = append(res, getFileNamesInDir(pattern, false)...)
	} else {
		pattern = filepath.Dir(pattern)
		res = append(res, getFileNamesInDir(pattern, false)...)
	}
	for _, v := range res {
		paths = append(paths, filepath.Join(pattern, v))
	}
	return paths
}
