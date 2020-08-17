package cmd

func bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func strikethrough(str string) string {
	return "\033[9m" + str + "\033[0m"
}
