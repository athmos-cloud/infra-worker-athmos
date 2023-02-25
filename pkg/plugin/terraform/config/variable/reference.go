package variable

import "fmt"

type Reference struct {
	Name   string `hcl:"name,label"`
	Source string `hcl:"source,optional"`
}
type ReferenceList []Reference

type Packaged struct {
	Variable        *Variable  `hcl:"variable,label"`
	ModuleReference *Reference `hcl:"module_reference,optional"`
	TFVar           *Reference `hcl:"reference,label"`
}

type PackagedList []Packaged

func (r *Reference) ToString() string {
	return fmt.Sprintf("%s = %s", r.Name, r.Source)
}
