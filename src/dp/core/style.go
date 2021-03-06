package core

type codeStyle struct {
	UseDefine       bool
	UsePrefixInc    bool
	UseArrayFromOne bool
}

func NewCodeStyle() *codeStyle {
	ret := new(codeStyle)
	ret.UseArrayFromOne = false
	ret.UseDefine = false
	ret.UsePrefixInc = true
	return ret
}
