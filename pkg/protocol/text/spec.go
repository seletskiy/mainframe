package text

import (
	"fmt"
	"image/color"
)

type Spec struct {
	Required []string
	Types    map[string]string
	Bound    map[string]interface{}
	Skipped  map[string]bool

	Opts struct {
		SkipUnknown bool
	}
}

func NewSpec() *Spec {
	return &Spec{
		Types:   map[string]string{},
		Bound:   map[string]interface{}{},
		Skipped: map[string]bool{},
	}
}

func (spec *Spec) Int(arg string, out interface{}) *Spec {
	spec.Types[arg] = "int"
	spec.Bound[arg] = out

	return spec
}

func (spec *Spec) Color(arg string, out interface{}) *Spec {
	spec.Types[arg] = "color"
	spec.Bound[arg] = out

	return spec
}

func (spec *Spec) String(arg string, out interface{}) *Spec {
	spec.Types[arg] = "string"
	spec.Bound[arg] = out

	return spec
}

func (spec *Spec) Bool(arg string, out interface{}) *Spec {
	spec.Types[arg] = "bool"
	spec.Bound[arg] = out

	return spec
}

func (spec *Spec) Require(arg string) *Spec {
	spec.Required = append(spec.Required, arg)

	return spec
}

func (spec *Spec) SkipUnknown() *Spec {
	spec.Opts.SkipUnknown = true

	return spec
}

func (spec *Spec) Skip(arg string) *Spec {
	spec.Skipped[arg] = true

	return spec
}

func (spec *Spec) Bind(args map[string]interface{}) error {
	for _, name := range spec.Required {
		if _, ok := args[name]; !ok {
			return ErrMissingArg(name)
		}
	}

	for name, _ := range args {
		if spec.Skipped[name] {
			continue
		}

		if _, ok := spec.Types[name]; !ok {
			if spec.Opts.SkipUnknown {
				continue
			}

			return ErrUnknownArg(name)
		}

		var (
			kind = spec.Types[name]
			out  = spec.Bound[name]
		)

		var value interface{}
		var ok bool

		switch kind {
		case "int":
			value, ok = args[name].(int)
		case "color":
			value, ok = args[name].(color.RGBA)
		case "string":
			value, ok = args[name].(string)
		case "bool":
			value, ok = args[name].(bool)
		default:
			panic(fmt.Sprintf("invalid arg type in spec: %q", kind))
		}

		if !ok {
			return ErrIncorrectType{
				Arg:   name,
				Value: args[name],
				Type:  kind,
			}
		}

		assign(value, out)
	}

	return nil
}

func assign(value interface{}, out interface{}) {
	switch out := out.(type) {
	case *int:
		*out = value.(int)
	case *color.RGBA:
		*out = value.(color.RGBA)
	case *string:
		*out = value.(string)
	case *bool:
		*out = value.(bool)

	case **int:
		value := value.(int)
		*out = &value
	case **color.RGBA:
		value := value.(color.RGBA)
		*out = &value
	case **string:
		value := value.(string)
		*out = &value
	case **bool:
		value := value.(bool)
		*out = &value
	}
}
