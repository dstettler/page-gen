package pagegencore

import (
	"log/slog"
	"slices"
	"strconv"

	"github.com/goccy/go-yaml"
)

type ReaderContents struct {
	DirectVals      map[string]interface{}
	SimpleArrays    map[string][]interface{}
	StructureArrays map[string][]map[string]interface{}
}

// Uses the contents of `structure` to fill missing fields in the structure arrays read from the contents file
// This will panic if a default value is not found for a field that is missing
func FillDefaults(vals map[string][]map[string]interface{}, defaults map[string]interface{}, structure map[string][]string) map[string][]map[string]interface{} {
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
	var valsMap map[string]interface{} = make(map[string]interface{})
	var arrsMap map[string][]interface{} = make(map[string][]interface{})
	var arrStructsMap map[string]([]map[string]interface{}) = make(map[string][]map[string]interface{})
	var defaults map[string]interface{} = make(map[string]interface{})

	// Build the contents of valsMap and structures for each array/array structure and value encountered
	for key, entry := range out {
		switch entryContent := entry.(type) {
		case string:
			slog.Debug("Got var: ", "key", key, "entry", entry)
			valsMap[key] = entry.(string)
		case int:
			slog.Debug("Got var: ", "key", key, "entry", entry)
			valsMap[key] = entry.(int)
		case float64:
			slog.Debug("Got var: ", "key", key, "entry", entry)
			valsMap[key] = entry.(float64)

		// We can't immediately check for the type of interface{}, so we can only check at this level for now
		// We know that the only valid types of arrays are []string or []map[string]string, however
		case []interface{}:
			slog.Debug("Got array: ", "var", key)

			// Manually do logic for defaults all at once, since it's simple enough
			if key == "defaults" {
				for index := range entryContent {
					varContent := entryContent[index]
					for k, v := range varContent.(map[string]interface{}) {
						switch v.(type) {
						case string:
							defaults[k] = v.(string)
						case int:
							defaults[k] = v.(int)
						case float64:
							defaults[k] = v.(float64)
						}

					}
				}

				continue
			}

			// For each array found
			for arrValIndex := range entryContent {
				arrValContent := entryContent[arrValIndex]
				switch arrValContent.(type) {
				case string:
					slog.Debug("Item ", "index", arrValIndex, "content", arrValContent)
					arrsMap[key] = append(arrsMap[key], arrValContent.(string))
				case int:
					arrsMap[key] = append(arrsMap[key], arrValContent.(int))
				case float64:
					arrsMap[key] = append(arrsMap[key], arrValContent.(float64))

				case map[string]interface{}:
					slog.Debug("Got structured arr: ", "index", arrValIndex)

					var structMap map[string]interface{} = make(map[string]interface{})

					// Read each key/value pair in the structure and append to the overall struct array
					for k, v := range arrValContent.(map[string]interface{}) {
						if !slices.Contains(structures[key], k) {
							structures[key] = append(structures[key], k)
						}

						slog.Debug("Got val: ", "key", k, "val", v)
						switch v.(type) {
						case string:
							structMap[k] = v.(string)
						case int:
							structMap[k] = v.(int)
						case float64:
							structMap[k] = v.(float64)
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
