package commands

import (
	"sort"
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

func (r *commandRouter) SearchCommands(iContext dto.InternalContextIface, patterns ...dto.PatternIface) dto.CommandRouterSearchResult {
	var wg sync.WaitGroup
	res := NewCommandRouterSearchResult()
	resultChan := make(chan dto.CommandManagerSearchResult, r.getCommandManagerCount())
	defer close(resultChan)

	wg.Add(len(patterns) * r.getCommandManagerCount())

	go func() {
		for rFromChan := range resultChan {
			if rFromChan.Founded() != 0 {
				res.AddResult(rFromChan)
			}
			wg.Done()
		}
	}()
	for _, cm := range r.commandManagers {
		go cm.SearchCommands(iContext, resultChan, patterns...)
	}

	wg.Wait()

	return res
}

func (r *commandRouter) getCommandManagerCount() int {
	return len(r.commandManagers)
}

func (r *commandRouter) GetSearchFunc() func(iContext dto.InternalContextIface, pattern dto.PatternIface) []dto.CommandManagerSearchResult {
	return func(iContext dto.InternalContextIface, pattern dto.PatternIface) []dto.CommandManagerSearchResult {
		return r.SearchCommands(iContext, pattern).GetDataByPattern(pattern)
	}
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
	SearchCommands(iContext dto.InternalContextIface, resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface)
}

func NewCommandRouterSearchResult() *commandRouterSearchResult {
	c := commandRouterSearchResult{
		data: make(map[dto.PatternIface][]dto.CommandManagerSearchResult),
	}
	return &c
}

type commandRouterSearchResult struct {
	sync.Mutex
	data map[dto.PatternIface][]dto.CommandManagerSearchResult
}

func (c *commandRouterSearchResult) AddResult(searchResult dto.CommandManagerSearchResult) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.data[searchResult.GetPattern()] = append(c.data[searchResult.GetPattern()], searchResult)
}

func (c *commandRouterSearchResult) GetDataByPattern(pattern dto.PatternIface) []dto.CommandManagerSearchResult {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	return sortCommandRouterSearchResult(c.data[pattern])
}

func sortCommandRouterSearchResult(cmsrs []dto.CommandManagerSearchResult) []dto.CommandManagerSearchResult {
	sort.Slice(cmsrs, func(i, j int) bool {
		return cmsrs[i].GetPriority() > cmsrs[j].GetPriority()
	})
	return cmsrs
}
