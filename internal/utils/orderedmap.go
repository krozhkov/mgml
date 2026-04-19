package utils

import "github.com/elliotchance/orderedmap/v3"

func ToEntries[K comparable, V any](source *orderedmap.OrderedMap[K, V]) []*orderedmap.Element[K, V] {
	if source == nil || source.Len() == 0 {
		return nil
	}

	result := make([]*orderedmap.Element[K, V], 0, source.Len())

	for key := range source.Keys() {
		result = append(result, source.GetElement(key))
	}

	return result
}
