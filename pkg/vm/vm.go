package vm

type VM interface {
	// Run runs a
	Run(code string, args ...interface{}) error
	RunFile(path string, args ...interface{}) error
	RunIn(code, entry string, args ...interface{}) error
	RunInFile(path, entry string, args ...interface{}) error
}

func init() {
}