package tcc

import "reflect"

// I not sure this method is common, so put it here. the reflec

type ReferencedService interface {
	Reference() string
}

// GetReference return the reference id of the service.
// If the service implemented the ReferencedService interface,
// it will call the Reference method. If not, it will
// return the struct name as the reference id.
func GetReference(service interface{}) string {
	if s, ok := service.(ReferencedService); ok {
		return s.Reference()
	}
	ref := ""
	sType := reflect.TypeOf(service)
	kind := sType.Kind()
	switch kind {
	case reflect.Struct:
		ref = sType.Name()
	case reflect.Ptr:
		sName := sType.Elem().Name()
		if sName != "" {
			ref = sName
		} else {
			ref = sType.Elem().Field(0).Name
		}
	}
	return ref
}
