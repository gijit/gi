// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements objsets.
//
// An objset is similar to a Scope but objset elements
// are identified by their unique id, instead of their
// object name.

package types

import "fmt"

// An objset is a set of objects identified by their unique id.
// The zero value for objset is a ready-to-use empty objset.
type objset map[string]Object // initialized lazily

// insert attempts to insert an object obj into objset s.
// If s already contains an alternative object alt with
// the same name, insert leaves s unchanged and returns alt.
// Otherwise it inserts obj and returns nil.
func (s *objset) insert(obj Object) Object {
	pp("objset.insert called with obj.Name()='%s', obj.Id()='%s'", obj.Name(), obj.Id())
	id := obj.Id()
	if alt := (*s)[id]; alt != nil {
		return alt
	}
	if *s == nil {
		*s = make(map[string]Object)
	}
	(*s)[id] = obj
	return nil
}

func (s *objset) del(obj Object) {
	if s == nil || *s == nil {
		return
	}
	id := obj.Id()
	delete(*s, id)
}

func (s *objset) exists(obj Object) bool {
	if *s == nil {
		return false
	}
	id := obj.Id()
	_, found := (*s)[id]
	return found
}

func (s *objset) replace(obj Object) (alt Object) {
	pp("objset.replace called with obj.Name()='%s', Id='%s'", obj.Name(), obj.Id())
	id := obj.Id()
	alt = (*s)[id]
	if alt != nil {
		// jea:
		pp("objset.replace is replacing a prior '%s'", id)
	}
	if *s == nil {
		*s = make(map[string]Object)
	}
	(*s)[id] = obj
	return
}

func (s *objset) String() string {
	var r string
	for i := range *s {
		r += fmt.Sprintf("objset[%v] = '%s'\n", i, (*s)[i])
	}
	return r
}
