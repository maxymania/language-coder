/*
      Copyright(C) 2015 Simon Schmidt

      This Source Code Form is subject to the
      terms of the Mozilla Public License, v.
      2.0. If a copy of the MPL was not
      distributed with this file, You can
      obtain one at
      http://mozilla.org/MPL/2.0/.
*/

/*
 Low Level Representation of program code.
*/
package lowrepr

const (
	OC_NOOP = iota
	OC_RETURN
	OC_LIT
	OC_JUMP
	OC_JUMP_IF
	OC_JUMP_UNLESS
	OC_MOVE
	// the base number for user commands
	OC_BASE = 0x20
)

type ValueType struct{
	Type int `xml:"type,attr,omitempty"`
	TypeEx string `xml:"type-ex,attr,omitempty"`
}
func (s ValueType) VSpec(n string) ValueSpec { return ValueSpec{s.Type,s.TypeEx,n} }

type ValueSpec struct{
	Type int `xml:"type,attr,omitempty"`
	TypeEx string `xml:"type-ex,attr,omitempty"`
	Content string `xml:",chardata"`
}
func (s ValueSpec) VType() ValueType { return ValueType{s.Type,s.TypeEx} }
func (s *ValueSpec) SetVType(v ValueType) {
	s.Type = v.Type
	s.TypeEx = v.TypeEx
}


type Op struct{
	Code int `xml:"code,attr,omitempty"`
	Label string `xml:"label,attr,omitempty"`
	Target string `xml:"target,attr,omitempty"`
	Content string `xml:",chardata"`
	ContentSmall string `xml:"content,attr,omitempty"`
	ContentConcise string `xml:"content,omitempty"`
	In []ValueSpec `xml:"in"`
	Out []ValueSpec `xml:"out"`
	Closure *Function `xml:"closure"`
}

type Function struct{
	Params []ValueSpec `xml:"param"`
	Locals []ValueSpec `xml:"local"`
	Ops []Op `xml:"op"`
}


