package harvestclient

type Arguments map[string]string

func DefaultArgs() Arguments {
	return make(Arguments)
}
