/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

// Package version represents the current version of the project.
package version

// Version is the current version of the kube-volume-freezer.
// Update this whenever making a new release.
// The version is of the format Major.Minor.Patch[-Stage.No[+]]
// Increment major number for new feature additions and behavioral changes.
// Increment minor number for bug fixes and performance enhancements.
// Increment patch number for critical fixes to existing releases.
// Stage is added during development and is one of (pre-alpha|alpha|beta|rc)
// accompanied by its iteration number. The trailing '+' sign is only removed
// during a release.
var Version = "0.2.0-alpha.1+"
