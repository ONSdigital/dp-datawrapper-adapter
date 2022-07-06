package charts

type MongoStore struct{}

func (ms *MongoStore) GetCollectionID(chartID string) (string, error) {
	return "", nil
}
