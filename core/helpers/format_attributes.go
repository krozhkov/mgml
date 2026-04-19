package helpers

import (
	"github.com/elliotchance/orderedmap/v3"
	"github.com/krozhkov/mgml/core/types"
)

func FormatAttributes(attributes *orderedmap.OrderedMap[string, string], allowedAttributes map[string]string) *orderedmap.OrderedMap[string, string] {
	result := orderedmap.NewOrderedMap[string, string]()

	for attrName, val := range attributes.AllFromFront() {
		if allowedAttributes != nil {
			if unit, ok := allowedAttributes[attrName]; ok {
				TypeConstructor, err := types.InitializeType(unit)
				if err == nil {
					typeValue := TypeConstructor(val)

					// yarozhkov: do we really need it?
					result.Set(attrName, typeValue.Value())
					continue
				}
			}
		}

		result.Set(attrName, val)
	}

	return result
}
