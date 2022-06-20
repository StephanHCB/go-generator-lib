package api

// Specifies what templates belong to a generator and what variables it needs to run.
//
// Will be read from a generator-*.yaml file in the root directory of the generator.
//
// The values of the variables as well as what generator to use come from a RenderSpec instead.
type GeneratorSpec struct {
	// The list of templates to render (if their condition evaluates to true)
	Templates []TemplateSpec `yaml:"templates"`

	// The list of available variables
	Variables map[string]VariableSpec `yaml:"variables"`
}

// Specifies a template to process, or a list to iterate over, if WithItems is nonempty (setting {{ item }} each run)
//
// Every field is evaluated as a template itself, so you can use variables in all fields.
//
// If Condition is set and evaluates to one of 'false', '0', 'no', the render run is skipped
type TemplateSpec struct {
	RelativeSourcePath string        `yaml:"source"`
	RelativeTargetPath string        `yaml:"target"`
	Condition          string        `yaml:"condition"`
	WithItems          []interface{} `yaml:"with_items"`
	JustCopy           bool          `yaml:"just_copy"`
}

// Specifies a variable that this generator uses, so it is made available in the templates.
//
// Actual values for an invocation of the generator are set in a RenderSpec, not the GeneratorSpec.
type VariableSpec struct {
	// Human readable description for the variable.
	Description string `yaml:"description"`

	// Regex validation pattern that the string representation (%v) of the value must match. No validation if left empty.
	ValidationPattern string `yaml:"pattern"`

	// Default value. If missing, the variable is considered required. Note that variables can have structured content.
	DefaultValue interface{} `yaml:"default"`
}
