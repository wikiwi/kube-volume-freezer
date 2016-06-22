package reflect

import (
	reflect "k8s.io/kubernetes/third_party/forked/reflect"
)

var def Equalities

// Equalities is a wrapper around Equalities. The functionality comes from
// Kubernetes which itself forks from the official go std library.
type Equalities struct {
	reflect.Equalities
}

// EqualitiesOrDie panics on errors for convience.
func EqualitiesOrDie(funcs ...interface{}) Equalities {
	e := Equalities{reflect.Equalities{}}
	if err := e.AddFuncs(funcs...); err != nil {
		panic(err)
	}
	return e
}

// DeepEqual performs reflect.DeepEqual without equality functions defined.
func DeepEqual(a interface{}, b interface{}) bool {
	return def.DeepEqual(a, b)
}

// DeepDerivative is similar to DeepEqual except that unset fields in a1 are
// ignored (not compared). This allows us to focus on the fields that matter to
// the semantic comparison.
//
// The unset fields include a nil pointer and an empty string.
// No equality functions are used.
func DeepDerivative(a interface{}, b interface{}) bool {
	return def.DeepDerivative(a, b)
}
