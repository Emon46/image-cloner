package controllers

func IncludedNamespace(namespace string) bool {
	if namespace == ExcludedNameSpace {
		return false
	}
	return true
}
