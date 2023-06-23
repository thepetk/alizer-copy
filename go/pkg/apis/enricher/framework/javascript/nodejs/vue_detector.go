/*******************************************************************************
* Copyright (c) 2022 Red Hat, Inc.
* Distributed under license by Red Hat, Inc. All rights reserved.
* This program is made available under the terms of the
* Eclipse Public License v2.0 which accompanies this distribution,
* and is available at http://www.eclipse.org/legal/epl-v20.html
*
* Contributors:
* Red Hat, Inc.
******************************************************************************/

package enricher

import (
	"context"
	"regexp"

	"github.com/mike-hoang/alizer-copy/go/pkg/apis/model"
	"github.com/mike-hoang/alizer-copy/go/pkg/utils"
)

type VueDetector struct{}

func (v VueDetector) GetSupportedFrameworks() []string {
	return []string{"Vue"}
}

// DoFrameworkDetection uses a tag to check for the framework name
func (v VueDetector) DoFrameworkDetection(language *model.Language, config string) {
	if hasFramework(config, "vue") {
		language.Frameworks = append(language.Frameworks, "Vue")
	}
}

// DoPortsDetection searches for the port in package.json, .env file, and vue.config.js
func (v VueDetector) DoPortsDetection(component *model.Component, ctx *context.Context) {
	regexes := []string{`--port (\d*)`, `PORT=(\d*)`}
	// check if --port or PORT is set in start script in package.json
	port := getPortFromStartScript(component.Path, regexes)
	if utils.IsValidPort(port) {
		component.Ports = []int{port}
	}

	// check if --port or PORT is set in dev script in package.json
	port = getPortFromDevScript(component.Path, regexes)
	if utils.IsValidPort(port) {
		component.Ports = []int{port}
	}

	// check if port is set on .env file
	port = utils.GetPortValueFromEnvFile(component.Path, `PORT=(\d*)`)
	if utils.IsValidPort(port) {
		component.Ports = []int{port}
		return
	}

	//check inside the vue.config.js file
	bytes, err := utils.ReadAnyApplicationFile(component.Path, []model.ApplicationFileInfo{
		{
			Dir:  "",
			File: "vue.config.js",
		},
	}, ctx)
	if err != nil {
		return
	}
	re := regexp.MustCompile(`port:\s*(\d+)*`)
	component.Ports = utils.FindAllPortsSubmatch(re, string(bytes), 1)
}
