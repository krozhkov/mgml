package validator

var globalDependencies = make(map[string]map[string]struct{})

func RegisterDependency(name string, deps []string) {
	upsertDependency(globalDependencies, name, deps)
}

func upsertDependency(dst map[string]map[string]struct{}, name string, deps []string) {
	if existing, ok := dst[name]; ok {
		for _, dep := range deps {
			if _, ok := existing[dep]; !ok {
				existing[dep] = struct{}{}
			}
		}
	} else {
		s := make(map[string]struct{}, len(deps))
		for _, value := range deps {
			s[value] = struct{}{}
		}
		dst[name] = s
	}
}
