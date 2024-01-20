package commands

import (
	"sync"

	"ash/internal/dto"
)

// 1. Internal actions (for key binds and manual startup)
// 2. Filesystems commands i.e. exec /usr/sbin/fdisk
// 3. Internal POSIX commands like 'cd'
type commandRouter struct {
	commandManagers []CommandManagerIface
}

func (r *commandRouter) AddNewCommandManager(newCommandManager CommandManagerIface) {
	r.commandManagers = append(r.commandManagers, newCommandManager)
}

func (r *commandRouter) SearchCommands(patterns ...dto.PatternIface) dto.CommandRouterSearchResult {
	var wg sync.WaitGroup
	res := NewCommandRouterSearchResult()

	resultChan := make(chan dto.CommandManagerSearchResult, r.getCommandManagerCount())
	defer close(resultChan)

	go func() {
		for r := range resultChan {
			res.AddResult(r)
			wg.Done()
		}
	}()

	for _, cm := range r.commandManagers {
		go cm.SearchCommands(resultChan, patterns...)
		wg.Add(len(patterns))
	}

	wg.Wait()

	return &res
}

func (r *commandRouter) getCommandManagerCount() int {
	return len(r.commandManagers)
}

func NewCommandRouter(commandManagers ...CommandManagerIface) commandRouter {
	c := commandRouter{
		commandManagers: make([]CommandManagerIface, len(commandManagers)),
	}

	for i, cm := range commandManagers {
		c.commandManagers[i] = cm
	}

	return c
}

type CommandManagerIface interface {
	SearchCommands(resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface)
}

func NewCommandRouterSearchResult() commandRouterSearchResult {
	c := commandRouterSearchResult{
		data: make(map[dto.PatternIface][]dto.CommandManagerSearchResult),
	}
	return c
}

type commandRouterSearchResult struct {
	data map[dto.PatternIface][]dto.CommandManagerSearchResult
}

func (c *commandRouterSearchResult) AddResult(searchResult dto.CommandManagerSearchResult) {
	c.data[searchResult.GetPattern()] = append(c.data[searchResult.GetPattern()], searchResult)
}

func (c *commandRouterSearchResult) GetDataByPattern(pattern dto.PatternIface) []dto.CommandManagerSearchResult {
	return c.data[pattern]
}
