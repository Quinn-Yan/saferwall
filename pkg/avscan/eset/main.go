// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package eset

import (
	"regexp"
	"strings"

	"github.com/saferwall/saferwall/pkg/utils"
)

const (
	cmd = "/opt/eset/esets/sbin/esets_scan"
)

// ScanFile performs antivirus scan
func ScanFile(filePath string) (string, error) {

	// Execute the scanner with the given file path
	// --unsafe                  scan for potentially unsafe applications
	// --unwanted                scan for potentially unwanted applications
	// --clean-mode=MODE         use cleaning MODE for infected objects.
	// 							 Available options: none, standard (default),
	// 							 strict, rigorous, delete

	esetOut, err := utils.ExecCommand(cmd, "--unsafe", "--unwanted",
		"--clean-mode=NONE", filePath)

	// Exit codes:
	// 0    no threat found
	// 1    threat found and cleaned
	// 10   some files could not be scanned (may be threats)
	// 50   threat found
	// 100  error

	if err != nil && err.Error() != "exit status 1" && err.Error() != "exit status 50" {
		return "", err
	}

	// Scan started at:   Tue Jan  1 01:32:57 2019
	// name="/samples/aa.exx", threat="a variant of Win32/Injector.BIIZ trojan", action="", info=""

	detection := ""
	re := regexp.MustCompile(`threat="([\s\w/.]+)"`)
	l := re.FindStringSubmatch(esetOut)
	if len(l) > 0 {
		detection = l[1]
	}

	// Clean up detection name
	detection = strings.TrimPrefix(detection, "a variant of ")
	detection = strings.TrimSuffix(detection, "  potentially unwanted application")
	detection = strings.TrimSuffix(detection, "  potentially unsafe application")
	return detection, nil
}

// GetProgramVersion returns program version
func GetProgramVersion() (string, error) {

	// Execute the scanner with the given file path
	versionOut, err := utils.ExecCommand(cmd, "--version")
	if err != nil {
		return "", err
	}

	// Extract the version
	version := strings.Split(versionOut, " ")[2]
	version = strings.TrimSuffix(version, "\n")
	return version, nil
}