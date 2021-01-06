// Copyright 2020 The CDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cdk

// Imaginary Type Tag

type TypeTag interface {
	Tag() CTypeTag
	String() string
}

// Used to denote a concrete type identity
type CTypeTag string

func NewTypeTag(tag string) TypeTag {
	return CTypeTag(tag)
}

func (tag CTypeTag) Tag() CTypeTag {
	return tag
}

// Stringer interface implementation
func (tag CTypeTag) String() string {
	return string(tag)
}
