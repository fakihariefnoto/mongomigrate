package orchestrator

import (
	"fmt"
	"strconv"

	combinator "github.com/fakihariefnoto/mongomigrate/controller/combinator"
	"github.com/pkg/errors"
)

var (
	combinatorPkg combinator.Pkg
)

func Init() {
	combinatorPkg = combinator.New()
}

func ShowApp() {

	welcome := `
## Welcome to mongo migrator
## this is still beta version,
## so, you still should have to insert config manual
## and set query getter and updatter
	`

	fmt.Println(welcome)

	for true {
		var choice int64
		var choiceStr string
		printAction()
		_, err := fmt.Scanf("%s", &choiceStr)
		if err != nil {
			fmt.Println("Please input number")
			continue
		}

		choice, err = strconv.ParseInt(choiceStr, 10, 64)

		if choice == 5 {
			break
		}

		err = callAction(choice)
		if err != nil {
			fmt.Println("Wrong choice! -> ", err)
			continue
		}
	}

}

func printAction() {
	options := `
# 1. Set Receiver. (disabled)
# 2. Set Updatter. (disabled)
# 3. Summary?
# 4. Execution!
# 5. GET OUT 
	`
	fmt.Printf("%v\n # So, what's your command : ", options)

}

func callAction(choice int64) error {
	switch choice {
	case 1:
		return errors.New("Edit your config file to configure receiver, this options is disabled!")
	case 2:
		return errors.New("Edit your config file to configure updater, this options is disabled!")
	case 3:

	case 4:
		combinatorPkg.Runner()
	default:
		return errors.New("Choose the right number, please!")
	}
	return nil
}
