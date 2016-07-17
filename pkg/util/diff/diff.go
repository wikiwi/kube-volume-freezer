/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

// Package diff provides diff capabilites.
package diff

import (
	"github.com/wikiwi/kube-volume-freezer/third_party/forked/pretty"
)

// ObjectDiff returns the diff of given objects.
var ObjectDiff = pretty.Diff

// ObjectDiffDerivative is the same as ObjectDiff but ignores empty fields of the first object.
var ObjectDiffDerivative = pretty.DiffDerivative
