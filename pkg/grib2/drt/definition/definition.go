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
