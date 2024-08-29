package grib

type Section2 struct {
	section2
	local []byte // 6-N Local Use
}

type section2 struct {
	Length        uint32 // Length of the section in octets (N)
	SectionNumber uint8  // 2 - Number of the section
}

func (s *Section2) GetSectionNumber() int {
	return int(s.SectionNumber)
}
