package prelude

import (
	"reflect"
	"sort"
	"strings"
)

func Sort(input interface{}) interface{} {
	value := reflect.ValueOf(input)

	if sortable, ok := value.Interface().(sort.Interface); ok {
		sort.Stable(sortable)
	} else if value.Kind() == reflect.Slice {
		sort.SliceStable(value.Interface(), func(i, j int) bool {
			return strings.Compare(
				ToString(value.Index(i)),
				ToString(value.Index(j)),
			) < 0
		})
	}

	return input
}
