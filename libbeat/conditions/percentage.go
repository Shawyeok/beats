package conditions

import (
	"fmt"
	"hash/crc32"
	"strconv"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

type Percentage struct {
	Field        string
	Selector     string
	DefaultRatio float32
	Groups       map[string]float32

	Logger *logp.Logger
}

type PercentageConfig struct {
	Field        string            `config:"field"`
	Selector     string            `config:"selector"`
	DefaultRatio string            `config:"default_ratio"`
	Groups       map[string]string `config:"groups"`
}

func NewPercentageCondition(config *PercentageConfig) (p Percentage, err error) {
	p = Percentage{}
	if config.DefaultRatio != "" {
		defaultRatio, err := strconv.ParseFloat(config.DefaultRatio, 32)
		if err != nil {
			return p, err
		}
		if defaultRatio < 0 || defaultRatio > 1 {
			return p, fmt.Errorf("percentage: default_ratio must be between zero to one")
		}
		p.DefaultRatio = float32(defaultRatio)
	}
	p.Groups = make(map[string]float32)
	for key, value := range config.Groups {
		ratio, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return p, err
		}
		p.Groups[key] = float32(ratio)
		if ratio < 0 || ratio > 1 {
			return p, fmt.Errorf("percentage: ratio must be between zero to one")
		}
	}
	if config.Field == "" {
		return p, fmt.Errorf("percentage: field cannot be empty")
	}
	if config.Selector == "" {
		return p, fmt.Errorf("percentage: selector cannot be empty")
	}
	p.Field = config.Field
	p.Selector = config.Selector
	p.Logger = logp.NewLogger("percentage")
	return p, nil
}

func (p Percentage) Check(event ValuesMap) bool {
	ratio := p.getRatio(event)
	if ratio == 1 {
		return true
	}
	if ratio == 0 {
		return false
	}
	value, err := event.GetValue(p.Field)
	if err != nil {
		if err != common.ErrKeyNotFound {
			p.Logger.Warn("unexpect err get value by field[%s]: %v", p.Field, err)
		}
		return false
	}
	s := fmt.Sprintf("%v", value)
	if s == "" {
		return false
	}
	d := float32(crc32.ChecksumIEEE([]byte(s))%100) / 100.0
	return d <= ratio
}

func (p Percentage) String() string {
	return fmt.Sprintf("percentage{field: %s, selector: %s, defaultRatio: %f, groups: %v}",
		p.Field, p.Selector, p.DefaultRatio, p.Groups)
}

func (p Percentage) getRatio(event ValuesMap) float32 {
	v, err := event.GetValue(p.Selector)
	if err != nil {
		if err != common.ErrKeyNotFound {
			p.Logger.Warn("unexpect err get value by selector[%s]: %v", p.Selector, err)
		}
		return p.DefaultRatio
	}
	s := fmt.Sprintf("%v", v)
	if ratio, ok := p.Groups[s]; ok {
		return ratio
	}
	return p.DefaultRatio
}
