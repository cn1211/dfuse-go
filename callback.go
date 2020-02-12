package dfuse

import (
	"github.com/chenyihui555/dfuse-go/entity"
	jsoniter "github.com/json-iterator/go"
)

type Callback struct {
	msgBytes []byte
}

func newCallback() *Callback {
	return &Callback{msgBytes: make([]byte, 0)}
}

func (c *Callback) TableSnapshot() (*entity.TableSnapshotResp, error) {
	snapshot := entity.TableSnapshotResp{}
	if err := jsoniter.Unmarshal(c.msgBytes, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (c *Callback) TableDelta() (*entity.TableDeltaResp, error) {
	delta := entity.TableDeltaResp{}
	if err := jsoniter.Unmarshal(c.msgBytes, &delta); err != nil {
		return nil, err
	}
	return &delta, nil
}

func (c *Callback) TransactionLifecycle() (*entity.TransactionLifecycleResp, error) {
	lifecycleResp := entity.TransactionLifecycleResp{}
	if err := jsoniter.Unmarshal(c.msgBytes, &lifecycleResp); err != nil {
		return nil, err
	}
	return &lifecycleResp, nil
}

func (c *Callback) Progress() (*entity.ProgressResp, error) {
	progressResp := entity.ProgressResp{}
	if err := jsoniter.Unmarshal(c.msgBytes, &progressResp); err != nil {
		return nil, err
	}
	return &progressResp, nil
}

func (c *Callback) Listening() (*entity.ListeningResp, error) {
	listenResp := entity.ListeningResp{}
	if err := jsoniter.Unmarshal(c.msgBytes, &listenResp); err != nil {
		return nil, err
	}
	return &listenResp, nil
}

func (c *Callback) Error() (*entity.ErrorResp, error) {
	errResp := entity.ErrorResp{}
	if err := jsoniter.Unmarshal(c.msgBytes, &errResp); err != nil {
		return nil, err
	}
	return &errResp, nil
}

// TODO
func (c *Callback) ActionTraces() (interface{}, error) {

	return nil, nil
}

// TODO
func (c *Callback) HeadInfo() (interface{}, error) {

	return nil, nil
}
