// Package pretty pretty is a fork of github.com/kr/pretty to add
// DiffDerivative support. License can be found under the LICENSE file.
package pretty

import (
	"fmt"
	"io"
	"reflect"

	"github.com/kr/pretty"
)

type sbuf []string

func (s *sbuf) Write(b []byte) (int, error) {
	*s = append(*s, string(b))
	return len(b), nil
}

// Diff returns a slice where each element describes
// a difference between a and b.
func Diff(a, b interface{}) (desc []string) {
	diffWriter{w: (*sbuf)(&desc)}.diff(reflect.ValueOf(a), reflect.ValueOf(b))
	return desc
}

// DiffDerivative is like Diff except it ignores unset fields from a.
func DiffDerivative(a, b interface{}) (desc []string) {
	diffWriter{w: (*sbuf)(&desc)}.diffDerivative(reflect.ValueOf(a), reflect.ValueOf(b))
	return desc
}

type diffWriter struct {
	w io.Writer
	l string // label
}

func (w diffWriter) printf(f string, a ...interface{}) {
	var l string
	if w.l != "" {
		l = w.l + ": "
	}
	fmt.Fprintf(w.w, l+f, a...)
}

func (w diffWriter) diff(av, bv reflect.Value) {
	if !av.IsValid() && bv.IsValid() {
		w.printf("nil != %#v", bv.Interface())
		return
	}
	if av.IsValid() && !bv.IsValid() {
		w.printf("%#v != nil", av.Interface())
		return
	}
	if !av.IsValid() && !bv.IsValid() {
		return
	}

	at := av.Type()
	bt := bv.Type()
	if at != bt {
		w.printf("%v != %v", at, bt)
		return
	}

	// numeric types, including bool
	if at.Kind() < reflect.Array {
		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%#v != %#v", a, b)
		}
		return
	}

	switch at.Kind() {
	case reflect.String:
		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%q != %q", a, b)
		}
	case reflect.Ptr:
		switch {
		case av.IsNil() && !bv.IsNil():
			w.printf("nil != %v", bv.Interface())
		case !av.IsNil() && bv.IsNil():
			w.printf("%v != nil", av.Interface())
		case !av.IsNil() && !bv.IsNil():
			w.diff(av.Elem(), bv.Elem())
		}
	case reflect.Struct:
		for i := 0; i < av.NumField(); i++ {
			w.relabel(at.Field(i).Name).diff(av.Field(i), bv.Field(i))
		}
	case reflect.Slice:
		lenA := av.Len()
		lenB := bv.Len()
		if lenA != lenB {
			w.printf("%s[%d] != %s[%d]", av.Type(), lenA, bv.Type(), lenB)
			break
		}
		for i := 0; i < lenA; i++ {
			w.relabel(fmt.Sprintf("[%d]", i)).diff(av.Index(i), bv.Index(i))
		}
	case reflect.Map:
		ak, both, bk := keyDiff(av.MapKeys(), bv.MapKeys())
		for _, k := range ak {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.printf("%q != (missing)", av.MapIndex(k))
		}
		for _, k := range both {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.diff(av.MapIndex(k), bv.MapIndex(k))
		}
		for _, k := range bk {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.printf("(missing) != %q", bv.MapIndex(k))
		}
	case reflect.Interface:
		w.diff(reflect.ValueOf(av.Interface()), reflect.ValueOf(bv.Interface()))
	default:
		if !reflect.DeepEqual(av.Interface(), bv.Interface()) {
			w.printf("%# v != %# v", pretty.Formatter(av.Interface()), pretty.Formatter(bv.Interface()))
		}
	}
}

func (w diffWriter) diffDerivative(av, bv reflect.Value) {
	if !av.IsValid() && bv.IsValid() {
		w.printf("nil != %#v", bv.Interface())
		return
	}
	if av.IsValid() && !bv.IsValid() {
		w.printf("%#v != nil", av.Interface())
		return
	}
	if !av.IsValid() && !bv.IsValid() {
		return
	}

	at := av.Type()
	bt := bv.Type()
	if at != bt {
		w.printf("%v != %v", at, bt)
		return
	}

	// numeric types, including bool
	if at.Kind() < reflect.Array {
		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%#v != %#v", a, b)
		}
		return
	}

	switch at.Kind() {
	case reflect.String:
		if av.Len() == 0 {
			return
		}

		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%q != %q", a, b)
		}
	case reflect.Ptr:
		if av.IsNil() {
			return
		}
		switch {
		case !av.IsNil() && bv.IsNil():
			w.printf("%v != nil", av.Interface())
		case !av.IsNil() && !bv.IsNil():
			w.diffDerivative(av.Elem(), bv.Elem())
		}
	case reflect.Struct:
		for i := 0; i < av.NumField(); i++ {
			w.relabel(at.Field(i).Name).diffDerivative(av.Field(i), bv.Field(i))
		}
	case reflect.Slice:
		if av.IsNil() || av.Len() == 0 {
			return
		}

		lenA := av.Len()
		lenB := bv.Len()
		if lenA != lenB {
			w.printf("%s[%d] != %s[%d]", av.Type(), lenA, bv.Type(), lenB)
			break
		}
		for i := 0; i < lenA; i++ {
			w.relabel(fmt.Sprintf("[%d]", i)).diffDerivative(av.Index(i), bv.Index(i))
		}
	case reflect.Map:
		if av.IsNil() || av.Len() == 0 {
			return
		}
		ak, both, bk := keyDiff(av.MapKeys(), bv.MapKeys())
		for _, k := range ak {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.printf("%q != (missing)", av.MapIndex(k))
		}
		for _, k := range both {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.diffDerivative(av.MapIndex(k), bv.MapIndex(k))
		}
		for _, k := range bk {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.printf("(missing) != %q", bv.MapIndex(k))
		}
	case reflect.Interface:
		if av.IsNil() {
			return
		}

		w.diffDerivative(reflect.ValueOf(av.Interface()), reflect.ValueOf(bv.Interface()))
	default:
		if !reflect.DeepEqual(av.Interface(), bv.Interface()) {
			w.printf("%# v != %# v", pretty.Formatter(av.Interface()), pretty.Formatter(bv.Interface()))
		}
	}
}

func (d diffWriter) relabel(name string) (d1 diffWriter) {
	d1 = d
	if d.l != "" && name[0] != '[' {
		d1.l += "."
	}
	d1.l += name
	return d1
}

func keyDiff(a, b []reflect.Value) (ak, both, bk []reflect.Value) {
	for _, av := range a {
		inBoth := false
		for _, bv := range b {
			if reflect.DeepEqual(av.Interface(), bv.Interface()) {
				inBoth = true
				both = append(both, av)
				break
			}
		}
		if !inBoth {
			ak = append(ak, av)
		}
	}
	for _, bv := range b {
		inBoth := false
		for _, av := range a {
			if reflect.DeepEqual(av.Interface(), bv.Interface()) {
				inBoth = true
				break
			}
		}
		if !inBoth {
			bk = append(bk, bv)
		}
	}
	return
}
