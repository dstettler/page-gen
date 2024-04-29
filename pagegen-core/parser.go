package pagegencore

import (
	"slices"
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

type VariableType int

const (
	variableTypeDirectVal
	variableTypeArrayIndex
	variableTypeStructArrayVal
)

type Replacements struct {
	ReplacementStart, ReplacementEnd int
	ReplacedContent                  string
}

type Refname struct {
	RefString string
	RefParent HTMLTag
}

func GetRefnamesStringArray(refnames []Refname) []string {
	var arr []string = make([]string, 0)
	for i := range refnames {
		str := refnames[i].RefString
		arr = append(arr, str)
	}

	return arr
}

func IsTagnameCustomTag(tag string) bool {
	return (tag == "if-exists" ||
		tag == "if" ||
		tag == "for")
}

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

// Returns string, int, or float64
func RecursiveParseVar(interiorStr string, contents *ReaderContents, refnames []Refname) interface{} {
	for i := range refnames {
		if (interiorStr)
	}
}

// Returns end tag corresponding to topTag, and a string of the contents
func RecursiveParseTag(interiorStr string, contents *ReaderContents, topTag HTMLTag, refnames []Refname) (HTMLTag, string) {
	var buffer strings.Builder

	readingVariableName := false
	var varnameBuffer strings.Builder
	varnameStart := 0

	var replacements []Replacements

	for charIndex := range interiorStr {
		currentChar := interiorStr[charIndex]

		if currentChar == '<' {
			tag, found := GetHTMLTagFromString(interiorStr[charIndex:])
			if found && IsTagnameCustomTag(tag.TagName) {
				if tag.TagName == topTag.TagName && tag.TagType == htmlTagTypeCloseTag {
					return tag, buffer.String()
				} else if tag.TagType == htmlTagTypeOpenTag {
					if refname, err := tag.TagItems["refname"]; !err {
						if slices.Contains(GetRefnamesStringArray(refnames), refname.(string)) {
							panicStr := "Refname: " + refname.(string) + " already used!"
							panic(panicStr)
						}

						var newRef Refname
						newRef.RefString = refname.(string)
						newRef.RefParent = tag

						refnames = append(refnames, newRef)
					}

					RecursiveParseTag(interiorStr[tag.EndPos:], contents, tag, refnames)
				}
			}
		} else if currentChar == '{' {
			readingVariableName = true
			varnameStart = charIndex
		} else if currentChar == '}' {
			if strings.Contains(varnameBuffer.String(), ".") {
				strSlices := strings.Split(varnameBuffer.String(), ".")

				if len(strSlices) != 2 {
					panic("Variable must have only one string on either side of the '.'")
				}

				ident := strSlices[0]
				var replacement Replacements

				if content, err := contents.DirectVals[ident]; !err {
					switch content.(type) {
					case string:
						replacement.ReplacedContent = RecursiveParseVar(content.(string), contents, refnames)
					}
				}

				replacement.ReplacementEnd = charIndex
				replacement.ReplacementStart = varnameStart

				replacements = append(replacements, replacement)
			}

		}
	}
}

func RecursiveParseString(valStr string, contents *ReaderContents) string {
	var currentTag HTMLTag

	for charIndex := range valStr {
		currentChar := valStr[charIndex]

		if currentChar == '<' {
			tag, found := GetHTMLTagFromString(valStr[charIndex:])
			if found && IsTagnameCustomTag(tag.TagName) {
				// Skip tag content since we'll have already scanned these chars anyways
				charIndex = tag.EndPos
				currentTag = tag
			}
		}
	}
}

// Will returned parsed value
func ParseValContent(valKey string, contents *ReaderContents) interface{} {
	val := contents.DirectVals[valKey]
	switch val.(type) {
	// No need to recursive scan/parse contents if the val isn't a string :)
	case int:
	case float64:
		return val
	case string:
		return RecursiveParseString(val.(string), contents)
	}
}

func VariablesParser(contents *ReaderContents) {

}

func TemplateParser(contents *ReaderContents, templatePath string) string {
	return ""
}
