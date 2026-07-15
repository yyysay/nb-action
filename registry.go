package main

var actionRegistry = make(map[string]Action)

func Register(actions ...Action) {

	for _, action := range actions {
		actionRegistry[action.Name()] = action
	}
}

func GetAction(name string) (Action, bool) {

	action, ok := actionRegistry[name]

	return action, ok
}
