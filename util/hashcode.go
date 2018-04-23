package util

//hashcode 有可能会重复，此处仅用于city，不会重复。重复概率最低的是只有小写字母,空格，_,-
//如需要不重复需自行验证
func Hashcode(str string) uint32 {
	var h uint32
	length := len(str)
	if length == 0 {
		return 0
	}

	r := []rune(str)
	// pp.Println(str, len(r))
	for i := 0; i < len(r); i++ {
		h = h<<5 - h + uint32(r[i]) //h*31+r[i]
	}

	return h
}
