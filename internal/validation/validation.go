package validation

type Rule func() error

type Rules []Rule

type Validator struct {
	Rules Rules
}

func (v *Validator) Add(rules ...Rules) {
	for i := range rules {
		v.Rules = append(v.Rules, rules[i]...)
	}
}

func (v *Validator) Validate() []string {
	var errors []string
	for _, rule := range v.Rules {
		if err := rule(); err != nil {
			errors = append(errors, err.Error())
		}
	}
	return errors
}
