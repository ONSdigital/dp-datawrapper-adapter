package charts

type Stub struct{}

func (ms *Stub) GetCollectionID(chartID string) (string, error) {
	return "", nil
}
