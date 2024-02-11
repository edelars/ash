package data_source

import (
	"sort"

	"ash/internal/dto"
)

type dataSourceImpl struct {
	originalData []dto.CommandManagerSearchResult
	keyMapping   map[rune]dto.CommandIface
}

type dataSourceImplItem struct {
	source string
	r      rune
	cmd    dto.CommandIface
}

func (ds *dataSourceImpl) GetCommand(r rune) dto.CommandIface {
	return ds.keyMapping[r]
}

// Return data, mainFieldMaxWid,descFieldMaxWid
func (ds *dataSourceImpl) GetData(avalaibleSpace, overheadLinesPerSource int) ([]dto.GetDataResult, int, int) {
	res := ds.initGetDataResult(avalaibleSpace, overheadLinesPerSource)

	var totalResCount int
	var r rune
	var mainFieldMaxWid, descFieldMaxWid int

	for _, v := range ds.originalData {
		if totalResCount >= len(res) {
			break
		}
		if v.Founded() == 0 {
			continue
		}
		res[totalResCount].SourceName = v.GetSourceName()

		originalCmds := sortSlice(v.GetCommands())

		for i := 0; i < len(res[totalResCount].Items) && i < len(originalCmds); i++ {
			r = ds.generateRune(r)
			res[totalResCount].Items[i].Name = originalCmds[i].GetName()
			res[totalResCount].Items[i].DisplayName = originalCmds[i].GetDisplayName()
			res[totalResCount].Items[i].ButtonRune = r
			res[totalResCount].Items[i].Description = originalCmds[i].GetDescription()
			ds.keyMapping[r] = originalCmds[i]

			if l := len(res[totalResCount].Items[i].DisplayName); l > mainFieldMaxWid {
				mainFieldMaxWid = l
			}
			if l := len(res[totalResCount].Items[i].Description); l > descFieldMaxWid {
				descFieldMaxWid = l
			}

		}
		totalResCount++
	}

	return res, mainFieldMaxWid, descFieldMaxWid
}

func (ds *dataSourceImpl) generateRune(i rune) rune {
	if i < 47 {
		return 48 // 0
	}

	if i > 47 && i < 57 { // 0 - 9
		return i + 1
	}

	if i == 57 {
		return 97
	}

	if i > 96 && i < 122 { // a - z
		return i + 1
	}

	if i == 122 {
		return 65
	}

	if i > 64 && i < 90 { // A - Z
		return i + 1
	}
	return 0
}

func (ds *dataSourceImpl) initGetDataResult(avalaibleSpace, overheadLinesPerSource int) []dto.GetDataResult {
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
	res := make([]dto.GetDataResult, totalFilledSources, totalFilledSources)

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
		newGetDataResult := dto.GetDataResult{
			SourceName: sr.GetSourceName(),
			Items:      make([]dto.GetDataResultItem, size),
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

func NewDataSource(sr []dto.CommandManagerSearchResult) dto.DataSource {
	ds := dataSourceImpl{
		originalData: sr,
		keyMapping:   make(map[rune]dto.CommandIface),
	}
	return &ds
}
