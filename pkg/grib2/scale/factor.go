package scale

type Factor uint16

// https://codes.ecmwf.int/grib/format/grib2/regulations/
// 92.1.5 If applicable, negative values shall be indicated by setting the most significant bit to “1”.
func (s Factor) Int16() int16 {
	negtive := s&0x8000 > 0
	i := int16(s & 0x7fff)

	if negtive {
		return -i
	}

	return i
}
