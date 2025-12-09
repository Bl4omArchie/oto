package oto

type OtoMap struct {
	Name string
	JobMap map[string]OtoMap
}


func NewMap(name string, jobMap map[string]OtoMap) *OtoMap {
	return &OtoMap{
		Name: name,
		JobMap: jobMap,
	}
}
