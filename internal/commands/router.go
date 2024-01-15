package commands

import "sync"

// 1. Internal commands
// 2. Filesystems commands i.e. exec /usr/sbin/fdisk
// 3. Internal POSIX commands like 'cd'
type CommandRouter struct {
	commandManagers []CommandManagerIface
}

func (r *CommandRouter) AddNewCommandManager(newCommandManager CommandManagerIface) {
	r.commandManagers = append(r.commandManagers, newCommandManager)
}

func (r *CommandRouter) SearchCommands(patterns ...PatternIface) commandRouterSearchResult {
	var wg sync.WaitGroup
	res := NewCommandRouterSearchResult()

	resultChan := make(chan CommandManagerSearchResult, r.getCommandManagerCount())
	defer close(resultChan)

	go func() {
		for r := range resultChan {
			// if _, ok := res[r.GetSourceName()]; ok {
			// } else {
			// }
			res.AddResult(r)
			wg.Done()
		}
	}()

	for _, cm := range r.commandManagers {
		go cm.SearchCommands(resultChan, patterns...)
		wg.Add(len(patterns))
	}

	wg.Wait()

	return res
}

func (r *CommandRouter) getCommandManagerCount() int {
	return len(r.commandManagers)
}

func NewCommandRouter(commandManagers ...CommandManagerIface) CommandRouter {
	c := CommandRouter{
		commandManagers: make([]CommandManagerIface, len(commandManagers)),
	}

	for i, cm := range commandManagers {
		c.commandManagers[i] = cm
	}

	return c
}

type CommandManagerIface interface {
	SearchCommands(resultChan chan CommandManagerSearchResult, patterns ...PatternIface)
}

type CommandManagerSearchResult interface {
	GetSourceName() string
	GetCommands() []CommandIface
	GetPattern() PatternIface
}

func NewCommandRouterSearchResult() commandRouterSearchResult {
	c := commandRouterSearchResult{
		data: make(map[PatternIface][]CommandManagerSearchResult),
	}
	return c
}

type commandRouterSearchResult struct {
	data map[PatternIface][]CommandManagerSearchResult
}

func (c *commandRouterSearchResult) AddResult(searchResult CommandManagerSearchResult) {
	c.data[searchResult.GetPattern()] = append(c.data[searchResult.GetPattern()], searchResult)
}

func (c *commandRouterSearchResult) GetDataByPattern(pattern PatternIface) []CommandManagerSearchResult {
	return c.data[pattern]
}
