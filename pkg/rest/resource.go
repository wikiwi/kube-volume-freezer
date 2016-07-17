/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package rest

// Resource is the interface of a REST Resource.
type Resource interface {
	Register(s *Server)
}
