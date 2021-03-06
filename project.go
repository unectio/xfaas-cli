/////////////////////////////////////////////////////////////////////////////////
//
// Copyright (C) 2019-2020, Unectio Inc, All Right Reserved.
//
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
/////////////////////////////////////////////////////////////////////////////////

package main

import (
	"github.com/unectio/api"
	"github.com/unectio/api/apilet"
)

var prcols = apilet.Projects

type elementPr struct{ *api.ProjectImage }

func (pe elementPr) id() string {
	return string(pe.ProjectImage.Id)
}

func (pe elementPr) name() string {
	return pe.ProjectImage.Name
}

func (pe elementPr) short() string {
	return ""
}

func (pe elementPr) long() []*field {
	return []*field{
		{
			name: "Role",
			data: pe.ProjectImage.Role,
		},
		{
			name: "UserData",
			data: pe.ProjectImage.UserData,
		},
	}
}

func projectList() {
	var prs []*api.ProjectImage

	makeReq(prcols.List(), &prs)

	for _, pr := range prs {
		showListElement(elementPr{pr})
	}
}

func projectInfo(name *string) {
	var pr api.ProjectImage

	makeReq(prcols.Info(*name), &pr)

	showInfoElement(elementPr{&pr})
}
