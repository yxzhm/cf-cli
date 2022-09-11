package codegen

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type DefEntity struct {
	Header *DefHeader
	Bodies []*DefBody
}

type DefHeader struct {
	CmtCode string
	Path    string
	Parent  string
}

type DefBody struct {
	Tag           string
	Desc          string
	Path          string
	Required      bool
	Sign          bool
	Type          string
	IsEmbeddedSet bool
	EmbeddedSet   *DefEntity
}

var headerReg *regexp.Regexp
var bodyReg *regexp.Regexp

func init() {
	headerReg = regexp.MustCompile("CMTCODE:(\\S+)\\s+(\\S+)\\s+PATH:(\\S+)\\s+ParentCMTCODE:(\\S+)")
	bodyReg = regexp.MustCompile("(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+(\\S+)")
}

func ReadXmlDef(path string) (map[string]*DefEntity, error) {
	allEntities, err := readPlainXmlDef(path)
	if err != nil {
		return nil, err
	}

	if len(allEntities) == 0 {
		return nil, errors.New("the entity is empty")
	}

	return buildXmlDefTree(allEntities)
}

func readPlainXmlDef(path string) ([]*DefEntity, error) {
	textLines, err := readFileToStrings(path)
	if err != nil {
		return nil, err
	}
	if len(textLines) == 0 {
		return nil, errors.New("the Xml Def file is empty")
	}
	var allEntities []*DefEntity
	var entity *DefEntity = nil
	for _, line := range textLines {
		if strings.TrimSpace(line) == "" || isComment(line) {
			continue
		}
		if isHeader(line) {
			entity = newDefEntity()
			header, err := newDefHeader(line)
			if err != nil {
				return nil, err
			}
			entity.Header = header
			allEntities = append(allEntities, entity)
		}
		if isBody(line) {
			body, err := newDefBody(line)
			if err != nil {
				return nil, err
			}
			entity.Bodies = append(entity.Bodies, body)
		}
	}

	for _, entity := range allEntities {
		err := validEntity(entity)
		if err != nil {
			return nil, err
		}
	}

	return allEntities, nil
}

func validEntity(entity *DefEntity) error {
	allTags := make(map[string]int)
	for _, body := range entity.Bodies {
		if allTags[body.Tag] == 1 {
			return errors.New(fmt.Sprintf("the tag %s is dupldated in %s", body.Tag, entity.Header.CmtCode))
		} else {
			allTags[body.Tag] = 1
		}
	}
	return nil
}

func buildXmlDefTree(allEntities []*DefEntity) (map[string]*DefEntity, error) {
	rootEntities := make(map[string]*DefEntity)
	allCmtCodes := make(map[string]*DefEntity)
	for _, entity := range allEntities {
		if allCmtCodes[entity.Header.CmtCode] != nil {
			return nil, errors.New(fmt.Sprintf("the cmt code %s is duplicated", entity.Header.CmtCode))
		}
		allCmtCodes[entity.Header.CmtCode] = entity

		if entity.Header.Parent == "NULL" {
			rootEntities[entity.Header.CmtCode] = entity
		} else {
			rootEntity := allCmtCodes[entity.Header.Parent]
			for _, body := range rootEntity.Bodies {
				if body.Path == entity.Header.Path {
					body.IsEmbeddedSet = true
					body.EmbeddedSet = entity
				}
			}
		}
	}

	return rootEntities, nil
}

func isComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

func isHeader(line string) bool {
	return strings.HasPrefix(line, "CMTCODE:")
}

func isBody(line string) bool {
	return !isHeader(line)
}

func newDefBody(line string) (*DefBody, error) {
	allString := bodyReg.FindStringSubmatch(line)
	if len(allString) != 7 {
		return nil, errors.New(fmt.Sprintf("the %s body is invalid", line))
	}

	return &DefBody{Tag: allString[1],
			Desc:     allString[2],
			Path:     allString[3],
			Required: allString[4] == "M",
			Sign:     allString[5] == "y",
			Type:     allString[6]},
		nil
}

func newDefHeader(line string) (*DefHeader, error) {
	allString := headerReg.FindStringSubmatch(line)
	if len(allString) != 5 {
		return nil, errors.New(fmt.Sprintf("the %s header is invalid", line))
	}

	return &DefHeader{CmtCode: allString[1],
			Path:   allString[3],
			Parent: allString[4]},
		nil
}

func newDefEntity() *DefEntity {
	entity := DefEntity{}
	entity.Bodies = []*DefBody{}
	return &entity
}

func readFileToStrings(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var fileTextLines []string
	firstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		if firstLine {
			line = strings.Replace(line, "\uFEFF", "", -1)
			firstLine = false
		}
		fileTextLines = append(fileTextLines, strings.TrimSpace(line))
	}
	return fileTextLines, nil
}
