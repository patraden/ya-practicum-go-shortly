package batcher

// Batch represents a sset of Batcher operations.
type Batch []*Operation

// SetError signals an error relating to all operation in the Batch.
func (b Batch) SetError(err error) {
	for _, op := range b {
		op.SetError(err)
	}
}
