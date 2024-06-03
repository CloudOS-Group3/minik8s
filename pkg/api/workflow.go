package api

type Workflow struct {
	// Metadata: name, namespace, uuid
	Metadata ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	// Graph: the graph of the workflow
	Graph Graph `json:"graph,omitempty" yaml:"graph,omitempty"`
	// Trigger: the trigger of the workflow
	Trigger TriggerType `json:"triggerType,omitempty" yaml:"triggerType,omitempty"`
}

type Graph struct {
	// Function: name, namespace, (uuid)
	Function ObjectMeta `json:"function,omitempty" yaml:"function,omitempty"`
	// Rule: a switch-case rule with corresponding successor
	Rule NodeRule `json:"rule,omitempty" yaml:"rule,omitempty"`
}

type NodeRule struct {
	// Case: the case of the switch-case rule
	Case []Case `json:"case,omitempty" yaml:"case,omitempty"`
	// Default: the default successor
	Default *Graph `json:"default,omitempty" yaml:"default,omitempty"`
}

type Case struct {
	// Expression: rules
	Expression []Expression `json:"expression,omitempty" yaml:"expression,omitempty"`
	// Successor: the successor of the case
	Successor *Graph `json:"successor,omitempty" yaml:"successor,omitempty"`
}

type Expression struct {
	// Variable: the variable in the expression
	Variable string `json:"variable,omitempty" yaml:"variable,omitempty"`
	// Opt: the operator in the expression
	Opt string `json:"opt,omitempty" yaml:"opt,omitempty"`
	// Value: eg. Variable  Equal  Value
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
	// Type: the type of the value
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

const (
	// Equal
	Equal = "="
	// NotEqual
	NotEqual = "!="
	// Greater
	Greater = ">"
	// Less
	Less = "<"
	// GreaterEqual
	GreaterEqual = ">="
	// LessEqual
	LessEqual = "<="
)

const (
	// Int
	Int = "int"
	// Float
	Float = "float"
	// String
	String = "string"
	// Bool
	Bool = "bool"
)

type WorkflowResult struct {
	// Metadata: name, namespace, uuid
	Metadata ObjectMeta `json:"metadata,omitempty"`
	// Result: the result of the workflow
	Result []string `json:"result,omitempty"`
	// EndTime: the status of the workflow
	EndTime string `json:"status,omitempty"`
	// InvokeTime: the time of the workflow invoked
	InvokeTime string `json:"invokeTime,omitempty"`
}
