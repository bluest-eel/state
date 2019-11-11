package marshaller

// Passthrough ...
type Passthrough struct{}

// NewPassthrough ...
func NewPassthrough() *Passthrough {
	return &Passthrough{}
}

// Marshal ...
func (m *Passthrough) Marshal(data *StateMetadata) (uint64, error) {
	return 1, nil
}

// Unmarshal ...
func (m *Passthrough) Unmarshal(data uint64) (StateMetadata, error) {
	return StateMetadata{}, nil
}
