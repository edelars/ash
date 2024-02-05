package file_system

import (
	"os"
	"path/filepath"
	"strings"

	"ash/internal/commands"
	"ash/internal/dto"
)

const constManagerName = "Filesystem"

type fileSystemManager struct{}

func NewFileSystemManager() fileSystemManager {
	return fileSystemManager{}
}

func (m fileSystemManager) SearchCommands(iContext dto.InternalContextIface, resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	var data []dto.CommandIface
	paths := preparePathArr(iContext.GetEnv("$PATH"))

	defer func() {
		commandManager := commands.NewCommandManager(constManagerName, data...)
		commandManager.SearchCommands(iContext, resultChan, patterns...)
	}()

	if len(paths) == 0 {
		return
	}

	for _, pattern := range patterns {
		filesResults := getFileNamesInDirs(paths, pattern.IsPrecisionSearch(), getFileNamesInDir)
		for _, fileResItem := range filesResults {
			for _, fileName := range fileResItem.files {
				if pattern.IsPrecisionSearch() {
					if pattern.GetPattern() == fileName {
						data = append(data, NewSystemCommand(fileName))
					}
				} else {
					data = append(data, NewPseudoCommand(filepath.Join(fileResItem.dir, fileName)))
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
