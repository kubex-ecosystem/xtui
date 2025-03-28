package cli

import (
	"fmt"
	t "github.com/faelmori/xtui/types"
	"github.com/spf13/pflag"
	"os"
	"reflect"
	"strings"
)

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

	banner = `
 ___    ___  _________   ___  ___   ___     
|\  \  /  /||\___   ___\|\  \|\  \ |\  \    
\ \  \/  / /\|___ \  \_|\ \  \\\  \\ \  \   
 \ \    / /      \ \  \  \ \  \\\  \\ \  \  
  /     \/        \ \  \  \ \  \\\  \\ \  \ 
 /  /\   \         \ \__\  \ \_______\\ \__\
/__/ /\ __\         \|__|   \|_______| \|__|
|__|/ \|__|                                 
`
	return map[string]string{"banner": banner, "description": description}
}

func getAvailableProperties() map[string]string {
	return map[string]string{
		"property1": "value1",
		"property2": "value2",
	}
}

func adaptArgsToProperties(args []string, properties map[string]string) []string {
	adaptedArgs := args
	for key, value := range properties {
		adaptedArgs = append(adaptedArgs, fmt.Sprintf("--%s=%s", key, value))
	}
	return adaptedArgs
}

func createFormConfig(commandName string, flags *pflag.FlagSet) t.FormConfig {
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
