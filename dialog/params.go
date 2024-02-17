package dialog

import (
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	views = []string{}

	//CurrentCommand is the command before assigning to variables
	CurrentCommand string
	//FinalCommand is the command after assigning to variables
	FinalCommand string
	// replace gocui
	// This matches most encountered patterns
	// Skips match if there is a whitespace at the end ex. <param='my >
	// Ignores <, > characters since they're used to match the pattern
	patternRegex = `<([^<>]*[^\s])>`
)

func insertParams(command string, params map[string]string) string {
	r := regexp.MustCompile(patternRegex)

	matches := r.FindAllStringSubmatch(command, -1)
	if len(matches) == 0 {
		return command
	}

	resultCommand := command
	for _, p := range matches {
		splitted := strings.Split(p[1], "=")
		resultCommand = strings.Replace(resultCommand, p[0], params[splitted[0]], -1)
	}

	return resultCommand
}

// SearchForParams returns variables from a command
func SearchForParams(command string) [][2]string {
	r := regexp.MustCompile(patternRegex)

	params := r.FindAllStringSubmatch(command, -1)
	if len(params) == 0 {
		return nil
	}

	extracted := map[string]string{}
	ordered_params := [][2]string{}
	for _, p := range params {
		splitted := strings.Split(p[1], "=")
		key := splitted[0]
		_, param_exists := extracted[key]

		// Set to empty if no value is provided and param is not already set
		if len(splitted) == 1 && !param_exists {
			extracted[key] = ""
		} else if len(splitted) > 1 {
			// Set the value instead if it is provided
			extracted[key] = splitted[1]
		}

		// Fill in the keys only if seen for the first time to track order
		if !param_exists {
			ordered_params = append(ordered_params, [2]string{key, ""})
		}
	}

	// Fill in the values
	for i, param := range ordered_params {
		pair := [2]string{param[0], extracted[param[0]]}
		ordered_params[i] = pair
	}
	return ordered_params
}

func evaluateParams(g *gocui.Gui, _ *gocui.View) error {
	paramsFilled := map[string]string{}
	for _, v := range views {
		view, _ := g.View(v)
		res := view.Buffer()
		res = strings.Replace(res, "\n", "", -1)
		paramsFilled[v] = strings.TrimSpace(res)
	}
	FinalCommand = insertParams(CurrentCommand, paramsFilled)
	return gocui.ErrQuit
}
