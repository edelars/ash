package data_source

import (
	"sort"

	"ash/internal/dto"
)

type DataSource interface {
	GetCommand(r rune) dto.CommandIface
	GetData(avalaibleSpace, overheadLinesPerSource int) []GetDataResult
}

type GetDataResult struct {
	SourceName string
	Items      []GetDataResultItem
}

type GetDataResultItem struct {
	Name       string
	ButtonRune rune
}

type dataSourceImpl struct {
	originalData []SearchResultIface
	keyMapping   map[rune]dto.CommandIface
}

type SearchResultIface interface {
	GetSourceName() string
	GetCommands() []dto.CommandIface
	Founded() int
}

type dataSourceImplItem struct {
	source string
	r      rune
	cmd    dto.CommandIface
}

func (ds *dataSourceImpl) GetCommand(r rune) dto.CommandIface {
	return ds.keyMapping[r]
}

func (ds *dataSourceImpl) GetData(avalaibleSpace, overheadLinesPerSource int) []GetDataResult {
	res := ds.initGetDataResult(avalaibleSpace, overheadLinesPerSource)

	var totalResCount int
	var r rune

	for _, v := range ds.originalData {
		if totalResCount >= len(res) {
			break
		}
		if v.Founded() == 0 {
			continue
		}
		res[totalResCount].SourceName = v.GetSourceName()

		originalCmds := sortSlice(v.GetCommands())

		for i := 0; i < len(res[totalResCount].Items); i++ {
			r = ds.generateRune(r)
			res[totalResCount].Items[i].Name = originalCmds[i].GetName()
			res[totalResCount].Items[i].ButtonRune = r
			ds.keyMapping[r] = originalCmds[i]
		}
		totalResCount++
	}

	return res
}

func (ds *dataSourceImpl) generateRune(i rune) rune {
	if i > 96 && i < 122 {
		return i + 1
	} else {
		return 97
	}
}

func (ds *dataSourceImpl) initGetDataResult(avalaibleSpace, overheadLinesPerSource int) []GetDataResult {
	var totalFilledSources, additionalFreeSpace, overloadSourceCount, addSpacePerSource int

	for _, v := range ds.originalData {
		if c := v.Founded(); c > 0 {
			totalFilledSources++
		}
	}

	// if too much sources with data
	if totalFilledSources*(overheadLinesPerSource+1) > avalaibleSpace {
		totalFilledSources = avalaibleSpace / (overheadLinesPerSource + 1)
	}

	if totalFilledSources == 0 {
		return nil
	}

	spaceForEverySource := (avalaibleSpace - overheadLinesPerSource*totalFilledSources) / totalFilledSources
	res := make([]GetDataResult, totalFilledSources, totalFilledSources)

	for _, sr := range ds.originalData {
		if sr.Founded() == 0 {
			continue
		}
		cmdCount := len(sr.GetCommands())
		if cmdCount <= spaceForEverySource {
			additionalFreeSpace = additionalFreeSpace + spaceForEverySource - cmdCount
		} else {
			overloadSourceCount++
		}
	}

	// count additional space for every overload source
	if overloadSourceCount > 0 {
		addSpacePerSource = additionalFreeSpace / overloadSourceCount
	}

	var drCount int
	for _, sr := range ds.originalData {
		if sr.Founded() == 0 {
			continue
		}
		size := spaceForEverySource
		c := len(sr.GetCommands())
		if c > size && addSpacePerSource > 0 {
			size = size + addSpacePerSource
		} else if size > c {
			size = c
		}
		newGetDataResult := GetDataResult{
			SourceName: sr.GetSourceName(),
			Items:      make([]GetDataResultItem, size),
		}
		res[drCount] = newGetDataResult
		drCount++
	}
	return res
}

func sortSlice(cmds []dto.CommandIface) []dto.CommandIface {
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].GetMathWeight() > cmds[j].GetMathWeight()
	})
	return cmds
}

func NewDataSource(sr []SearchResultIface) dataSourceImpl {
	ds := dataSourceImpl{
		originalData: sr,
		keyMapping:   make(map[rune]dto.CommandIface),
	}
	return ds
}
