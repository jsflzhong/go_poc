// Code generated by "stringer -type=FormatCode"; DO NOT EDIT.

package pgwirebase

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FormatText-0]
	_ = x[FormatBinary-1]
}

const _FormatCode_name = "FormatTextFormatBinary"

var _FormatCode_index = [...]uint8{0, 10, 22}

func (i FormatCode) String() string {
	if i >= FormatCode(len(_FormatCode_index)-1) {
		return "FormatCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FormatCode_name[_FormatCode_index[i]:_FormatCode_index[i+1]]
}
