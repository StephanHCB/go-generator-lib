package api

// All Information needed by a render run.
//
// The idea is that this is read from a generator-<name>.yaml file in the target directory so runs can be repeated.
//
// You can also have this library render a file with the default values for a given generator.
type RenderSpec struct {
	// Name of the generator to use (determines yaml file to read for generator spec). The main one should be called
	// main.yaml, and GeneratorName should be set to "main".
	GeneratorName string `yaml:"generator"`

	// Assign variable "key" value "value". All values are evaluated as templates until they no longer change,
	// so the value of one variable can refer to other variables, even if using their default values.
	Variables map[string]string `yaml:"parameters"`
}
