package main

import (
	"log/slog"
	"slices"
	"strconv"

	"github.com/goccy/go-yaml"
)

type ReaderContents struct {
	DirectVals      map[string]string
	SimpleArrays    map[string][]string
	StructureArrays map[string][]map[string]string
}

// Uses the contents of `structure` to fill missing fields in the structure arrays read from the contents file
// This will panic if a default value is not found for a field that is missing
func FillDefaults(vals map[string][]map[string]string, defaults map[string]string, structure map[string][]string) map[string][]map[string]string {
	// For each structure in the contents
	for arrayVarname, arrayTopLevelVal := range vals {
		refArray := structure[arrayVarname]

		// For each item in that structure
		for arrItemIndex := range arrayTopLevelVal {
			arrItem := arrayTopLevelVal[arrItemIndex]

			// For each reference item in the structure ref
			for refIndex := range refArray {
				refVal := refArray[refIndex]
				if _, exists := arrItem[refVal]; !exists {
					itemLocator := arrayVarname + "." + refVal
					if defaultVal, defaultExists := defaults[itemLocator]; defaultExists {
						vals[arrayVarname][arrItemIndex][refVal] = defaultVal
					} else {
						panicStr := "No default value for " + itemLocator
						panic(panicStr)
					}
				}
			}
		}
	}

	return vals
}

// Returns read content file
// Keys will be one of three things:
// - A direct variable key
// - An array key in the form of Array[i].ArrayVal (default values will be already added)
// - An array key in the form of Array[i] (for arrays with only one unnamed value)
func ContentReader(contentsPath string) ReaderContents {
	var out map[string]interface{}
	err := yaml.Unmarshal([]byte(ReadFileToString(contentsPath)), &out)
	CheckErr(err)

	var structures map[string][]string = make(map[string][]string)
	var valsMap map[string]string = make(map[string]string)
	var arrsMap map[string][]string = make(map[string][]string)
	var arrStructsMap map[string]([]map[string]string) = make(map[string][]map[string]string)
	var defaults map[string]string = make(map[string]string)

	// Build the contents of valsMap and structures for each array/array structure and value encountered
	for key, entry := range out {
		switch entryContent := entry.(type) {
		case string:
			slog.Debug("Got var: ", "key", key, "entry", entry)
			valsMap[key] = entry.(string)

		case []interface{}:
			slog.Debug("Got array: ", "var", key)

			if key == "defaults" {
				for index := range entryContent {
					varContent := entryContent[index]
					for k, v := range varContent.(map[string]interface{}) {
						defaults[k] = v.(string)
					}
				}

				continue
			}

			for arrValIndex := range entryContent {
				arrValContent := entryContent[arrValIndex]
				switch arrValContent.(type) {
				case string:
					slog.Debug("Item ", "index", arrValIndex, "content", arrValContent)
					arrsMap[key] = append(arrsMap[key], arrValContent.(string))

				case map[string]interface{}:
					slog.Debug("Got structured arr: ", "index", arrValIndex)

					var structMap map[string]string = make(map[string]string)

					for k, v := range arrValContent.(map[string]interface{}) {
						if !slices.Contains(structures[key], k) {
							structures[key] = append(structures[key], k)
						}

						slog.Debug("Got val: ", "key", k, "val", v)
						switch v.(type) {
						case string:
							structMap[k] = v.(string)
						default:
							errorString := "Invalid value for key [" + k + "] in index [" + strconv.Itoa(arrValIndex) + "] of struct array [" + key + "]"
							panic(errorString)
						}
					}

					arrStructsMap[key] = append(arrStructsMap[key], structMap)
				}
			}
		default:
			slog.Debug("Invalid format of entry! Receieved type: ", entryContent, " for entry ", key)
		}
	}

	slog.Debug("Read vals: ", "val map", valsMap)
	slog.Debug("Read arrays: ", "arrays", arrsMap)
	slog.Debug("Read struct arrays: ", "stuct arrays", arrStructsMap)
	slog.Debug("Read structure templates: ", "structure map", structures)
	slog.Debug("Read defaults: ", "defaults map", defaults)

	arrStructsMap = FillDefaults(arrStructsMap, defaults, structures)

	slog.Debug("Finalized vals: ", "val map", valsMap)
	slog.Debug("Finalized arrays: ", "arrays", arrsMap)
	slog.Debug("Finalized struct arrays: ", "stuct arrays", arrStructsMap)

	var finalContents ReaderContents
	finalContents.DirectVals = valsMap
	finalContents.SimpleArrays = arrsMap
	finalContents.StructureArrays = arrStructsMap

	return finalContents
}
