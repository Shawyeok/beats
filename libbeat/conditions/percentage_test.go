package conditions

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

func TestPercentage_Check(t *testing.T) {
	logger := logp.NewLogger("TestPercentage")
	tests := []struct {
		name  string
		p     *Percentage
		event ValuesMap
		want  bool
	}{
		{
			name: "zero-ratio",
			p: &Percentage{
				Field:        "traceId",
				Selector:     "type",
				DefaultRatio: 0.2,
				Groups: map[string]float32{
					"a": 0,
				},
				Logger: logger,
			},
			event: common.MapStr{
				"type":    "a",
				"traceId": "",
			},
			want: false,
		},
		{
			name: "one-ratio",
			p: &Percentage{
				Field:        "traceId",
				Selector:     "type",
				DefaultRatio: 0.2,
				Groups: map[string]float32{
					"a": 1,
				},
				Logger: logger,
			},
			event: common.MapStr{
				"type":    "a",
				"traceId": "",
			},
			want: true,
		},
		{
			name: "nil-selector",
			p: &Percentage{
				Field:        "traceId",
				Selector:     "type",
				DefaultRatio: 0.2,
				Groups: map[string]float32{
					"a": 0.2,
				},
				Logger: logger,
			},
			event: common.MapStr{
				"traceId": "",
			},
			want: false,
		},
		{
			name: "nil-field",
			p: &Percentage{
				Field:        "traceId",
				Selector:     "type",
				DefaultRatio: 0.2,
				Groups: map[string]float32{
					"a": 0.2,
				},
				Logger: logger,
			},
			event: common.MapStr{
				"type": "a",
			},
			want: false,
		},
		{
			name: "empty-field",
			p: &Percentage{
				Field:        "traceId",
				Selector:     "type",
				DefaultRatio: 0.2,
				Groups: map[string]float32{
					"a": 0.2,
				},
				Logger: logger,
			},
			event: common.MapStr{
				"type":    "a",
				"traceId": "",
			},
			want: false,
		},
		{
			name: "hit",
			p: &Percentage{
				Field:        "traceId",
				Selector:     "type",
				DefaultRatio: 0.2,
				Groups: map[string]float32{
					"a": 0.46,
				},
				Logger: logger,
			},
			event: common.MapStr{
				"type":    "a",
				"traceId": "cb6ca4c3e10648b09a75e96be6f8e32c",
			},
			want: true,
		},
		{
			name: "default-ratio-not-cover",
			p: &Percentage{
				Field:        "traceId",
				Selector:     "type",
				DefaultRatio: 0.2,
				Groups:       map[string]float32{},
				Logger:       logger,
			},
			event: common.MapStr{
				"type":    "a",
				"traceId": "cb6ca4c3e10648b09a75e96be6f8e32c",
			},
			want: false,
		},
		{
			name: "default-ratio-covered",
			p: &Percentage{
				Field:        "traceId",
				Selector:     "type",
				DefaultRatio: 0.46,
				Groups:       map[string]float32{},
				Logger:       logger,
			},
			event: common.MapStr{
				"type":    "a",
				"traceId": "cb6ca4c3e10648b09a75e96be6f8e32c",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Check(tt.event); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
