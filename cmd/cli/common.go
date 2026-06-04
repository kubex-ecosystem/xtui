package cli

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"

	t "github.com/kubex-ecosystem/xtui/types"
	"github.com/spf13/pflag"
)

var (
	banners = []string{`
 __    __ ________          ______
|  \  |  \        \        |      \
| ▓▓  | ▓▓\▓▓▓▓▓▓▓▓__    __ \▓▓▓▓▓▓
 \▓▓\/  ▓▓  | ▓▓  |  \  |  \ | ▓▓
  >▓▓  ▓▓   | ▓▓  | ▓▓  | ▓▓ | ▓▓
 /  ▓▓▓▓\   | ▓▓  | ▓▓  | ▓▓ | ▓▓
|  ▓▓ \▓▓\  | ▓▓  | ▓▓__/ ▓▓_| ▓▓_
| ▓▓  | ▓▓  | ▓▓   \▓▓    ▓▓   ▓▓ \
 \▓▓   \▓▓   \▓▓    \▓▓▓▓▓▓ \▓▓▓▓▓▓`}
)

// GetDescriptions returns a map of descriptions for a command.
// This is useful for commands that have a lot of flags that can be passed as arguments.
func GetDescriptions(descriptionArg []string, _ bool) map[string]string {
	var description, banner string
	if descriptionArg != nil {
		if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
			description = descriptionArg[0]
		} else {
			description = descriptionArg[1]
		}
	} else {
		description = ""
	}
	bannerRandLen := len(banners)
	bannerRandIndex := rand.Intn(bannerRandLen)
	banner = banners[bannerRandIndex]
	return map[string]string{"banner": banner, "description": description}
}

// GetAvailableProperties returns a map of available properties.
// This is useful for commands that have a lot of flags that can be passed as arguments.
func GetAvailableProperties() map[string]string {
	return map[string]string{
		"property1": "value1",
		"property2": "value2",
	}
}

// AdaptArgsToProperties adapts a list of arguments to a list of properties.
// This is useful for commands that have a lot of flags that can be passed as arguments.
func AdaptArgsToProperties(args []string, properties map[string]string) []string {
	adaptedArgs := args
	for key, value := range properties {
		adaptedArgs = append(adaptedArgs, fmt.Sprintf("--%s=%s", key, value))
	}
	return adaptedArgs
}

// CreateFormConfig create a form config for a command and its flags.
func CreateFormConfig(commandName string, flags *pflag.FlagSet) t.FormConfig {
	var formFields []t.FormInputObject[any]

	flags.VisitAll(func(flag *pflag.Flag) {
		val := reflect.ValueOf(flag.Value).Interface()
		formFields = append(formFields, &t.Input[any]{
			Ph:                 flag.Name,
			Tp:                 reflect.TypeOf(flag.Value),
			Val:                &val,
			Req:                false,
			Min:                0,
			Max:                100,
			Err:                "",
			ValidationRulesVal: nil,
		})
	})

	return t.FormConfig{
		Title: fmt.Sprintf("Configure %s Command", commandName),
		FormFields: t.FormFields{
			Title:  fmt.Sprintf("Configure %s Command", commandName),
			Fields: formFields,
		},
	}
}
