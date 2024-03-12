package mask

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	DefaultMask = "*****"
)

type Masks struct {
	MaskChar string   `yaml:"maskChar"`
	Values   []string `yaml:"values"`
}

func (m *Masks) Add(r string) {
	if !m.contains(r) {
		m.Values = append(m.Values, r)
	}
}

func (m *Masks) contains(r string) bool {
	for _, val := range m.Values {
		if strings.EqualFold(val, r) {
			return true
		}
	}
	return false
}

func New() *Masks {
	return &Masks{
		MaskChar: DefaultMask,
		Values:   make([]string, 0),
	}
}

func LoadMasks(filepath string) *Masks {
	var m Masks
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return New()
	}

	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		return New()
	}

	return &m
}

func (m *Masks) Save(filepath string) {
	content, _ := yaml.Marshal(m)

	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Print("Err: Could not access configuration file")
		os.Exit(1)
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		fmt.Print("Err: Could not store configuration file")
		os.Exit(1)
	}
}

func (r *Masks) Remove(v string) bool {
	for i, existing := range r.Values {
		if strings.EqualFold(existing, v) {
			r.Values = append(r.Values[:i], r.Values[i+1:]...)
			return true
		}
	}
	return false
}

func (r *Masks) Compile() []*regexp.Regexp {
	var compiled []*regexp.Regexp
	for _, mask := range r.Values {
		re, err := regexp.Compile(fmt.Sprintf(`(?i)%s`, mask))
		if err != nil {
			fmt.Print("Err: invalid regex for mask", mask)
			continue
		}
		compiled = append(compiled, re)
	}
	return compiled
}
