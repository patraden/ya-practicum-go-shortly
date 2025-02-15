package batcher

type Batch []*Operation

func (b Batch) SetError(err error) {
	for _, op := range b {
		op.SetError(err)
	}
}
