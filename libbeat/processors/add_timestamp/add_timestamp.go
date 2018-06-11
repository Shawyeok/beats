package add_timestamp

import (
	"time"

	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/processors"
)

type addTimestamp struct {
	TargetField string
}

func init() {
	processors.RegisterPlugin("add_timestamp", newAddTimestamp)
}

func newAddTimestamp(c *common.Config) (processors.Processor, error) {
	config := struct {
		Target string `config:"target"`
	}{
		Target: "beat.timestamp",
	}

	err := c.Unpack(&config)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unpack the add_timestamp configuration")
	}

	var tim addTimestamp

	tim.TargetField = config.Target

	return tim, nil
}

func (t addTimestamp) Run(event *beat.Event) (*beat.Event, error) {
	ts := time.Now()
	event.PutValue(t.TargetField, ts)
	return event, nil
}

func (t addTimestamp) String() string {
	return "add_timestamp=[target=" + t.TargetField + "]"
}
