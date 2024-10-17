package definition

// don't edit
type SimplePacking struct {
	R float32 `json:"r"` // Reference value (R) (IEEE 32-bit floating-point value)
	B uint16  `json:"b"` // Binary scale factor
	D uint16  `json:"d"` // Decimal scale factor
	L uint8   `json:"l"` // Number of bits used for each packed value for simple packing, or for each group reference value for complex packing or spatial differencing
	T uint8   `json:"t"` // Type of original field values: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table5-1.shtml
}

// don't edit
type ComplexPacking struct {
	SimplePacking

	GroupSplittingMethodUsed   uint8  `json:"groupSplittingMethodUsed"`
	MissingValueManagementUsed uint8  `json:"missingValueManagementUsed"`
	PrimaryMissingSubstitute   uint32 `json:"primaryMissingSubstitute"`
	SecondaryMissingSubstitute uint32 `json:"secondaryMissingSubstitute"`
	NumberOfGroups             uint32 `json:"numberOfGroups"`
	GroupWidths                uint8  `json:"groupWidths"`
	GroupWidthsBits            uint8  `json:"groupWidthsBits"`
	GroupLengthsReference      uint32 `json:"groupLengthsReference"`
	GroupLengthIncrement       uint8  `json:"groupLengthIncrement"`
	GroupLastLength            uint32 `json:"groupLastLength"`
	GroupScaledLengthsBits     uint8  `json:"groupScaledLengthsBits"`
}

// don't edit
type ComplexPackingAndSpatialDifferencing struct {
	ComplexPacking

	SpatialOrderDifference uint8 `json:"spatialOrderDifference"`
	OctetsNumber           uint8 `json:"octetsNumber"`
}

/*
Section 5 - Template 41 : Grid point data - Portable Network Graphics (PNG) format

Octets	Key	Type	Content
12-15	referenceValue	ieeefloat	Reference value (R) (IEEE 32-bit floating-point value)
16-17	binaryScaleFactor	signed	Binary scale factor (E)
18-19	decimalScaleFactor	signed	Decimal scale factor (D)
20	bitsPerValue	unsigned	Number of bits required to hold the resulting scaled and referenced data values. (i.e. The depth of the image.) (see Note 2)
21	typeOfOriginalFieldValues	codetable	Type of original field values (see Code Table 5.1)
*/
type PNG SimplePacking
