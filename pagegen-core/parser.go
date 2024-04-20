package pagegencore

import (
	"strconv"
	"strings"
	"unicode"
)

type HTMLTagType int

const (
	htmlTagTypeOpenTag  HTMLTagType = 0
	htmlTagTypeCloseTag HTMLTagType = 1
)

type HTMLTag struct {
	TagName          string
	TagItems         map[string]interface{}
	TagType          HTMLTagType
	StartPos, EndPos int
}

type HTMLReaderState int

const (
	htmlReaderStateSpace                 HTMLReaderState = -1
	htmlReaderStatePreReading            HTMLReaderState = 0
	htmlReaderStateReadingTagName        HTMLReaderState = 1
	htmlReaderStateReadingTagItemName    HTMLReaderState = 2
	htmlReaderStateReadingTagItemContent HTMLReaderState = 3
)

// Returns HTMLTag and true if tag found, or an empty HTMLTag and false otherwise
func GetHTMLTagFromString(content string) (HTMLTag, bool) {
	var buffer strings.Builder
	var name string
	var items map[string]interface{} = make(map[string]interface{})
	var currentItemName string
	isEndTag := false
	startPos := -1

	readingState := htmlReaderStatePreReading
	previousState := readingState

	for strIndex := range content {
		currentChar := content[strIndex]

		if readingState == htmlReaderStateSpace && currentChar == '>' {
			if startPos < strIndex && startPos != -1 {
				var tag HTMLTag
				tag.TagName = name
				tag.StartPos = startPos
				tag.EndPos = strIndex
				tag.TagItems = items
				tag.TagType = htmlTagTypeOpenTag
				return tag, true
			}

			// Tag has no content- reset vars and continue searching for one with something
			buffer.Reset()
			name = ""
			items = make(map[string]interface{})

			startPos = -1

			continue
		} else if currentChar == '>' && readingState == htmlReaderStateReadingTagName && isEndTag {
			var tag HTMLTag
			tag.TagName = buffer.String()
			tag.StartPos = startPos
			tag.EndPos = strIndex
			tag.TagItems = items
			tag.TagType = htmlTagTypeCloseTag
			return tag, true
		}

		if readingState == htmlReaderStateSpace && !unicode.IsSpace(rune(currentChar)) && currentChar != '=' {
			readingState = previousState + 1
		}

		switch readingState {
		case htmlReaderStatePreReading:
			if currentChar == '<' {
				startPos = strIndex

				previousState = readingState
				readingState = htmlReaderStateSpace
			}

		case htmlReaderStateReadingTagName:
			if unicode.IsSpace(rune(currentChar)) {
				name = buffer.String()

				previousState = readingState
				readingState = htmlReaderStateSpace
				buffer.Reset()
			} else if currentChar == '/' {
				isEndTag = true
			} else {
				buffer.WriteByte(currentChar)
			}

		case htmlReaderStateReadingTagItemName:
			if currentChar == '=' || unicode.IsSpace(rune(currentChar)) {
				currentItemName = buffer.String()

				previousState = readingState
				readingState = htmlReaderStateSpace
				buffer.Reset()
			} else {
				buffer.WriteByte(currentChar)
			}

		case htmlReaderStateReadingTagItemContent:
			if currentChar == '"' && (buffer.Len() != 0 || content[strIndex-1] == '"') {
				if i, err := strconv.Atoi(buffer.String()); err == nil {
					items[currentItemName] = i
				} else if f, err := strconv.ParseFloat(buffer.String(), 64); err == nil {
					items[currentItemName] = f
				} else {
					items[currentItemName] = buffer.String()
				}

				buffer.Reset()
				currentItemName = ""
				previousState = htmlReaderStateReadingTagName
				readingState = htmlReaderStateSpace

			} else if currentChar == '"' && buffer.Len() == 0 {
				continue
			} else {
				buffer.WriteByte(currentChar)
			}
		}
	}

	return HTMLTag{}, false
}

// Will returned parsed value
func ParseValContent(valKey string, contents *ReaderContents) string {
	return ""
}

func VariablesParser(contents *ReaderContents) {

}

func TemplateParser(contents *ReaderContents, templatePath string) string {
	return ""
}
