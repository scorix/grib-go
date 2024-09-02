package definition

// don't edit
type SimplePacking struct {
	R    float32 // Reference value (R) (IEEE 32-bit floating-point value)
	E    uint16  // Binary scale factor
	D    uint16  // Decimal scale factor
	Bits uint8   // Number of bits used for each packed value for simple packing, or for each group reference value for complex packing or spatial differencing
	Type uint8   // Type of original field values: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table5-1.shtml
}

// don't edit
type ComplexPacking struct {
	SimplePacking

	GroupSplittingMethodUsed   uint8
	MissingValueManagementUsed uint8
	PrimaryMissingSubstitute   uint32
	SecondaryMissingSubstitute uint32
	NumberOfGroups             uint32
	GroupWidths                uint8
	GroupWidthsBits            uint8
	GroupLengthsReference      uint32
	GroupLengthIncrement       uint8
	GroupLastLength            uint32
	GroupScaledLengthsBits     uint8
}

// don't edit
type ComplexPackingAndSpatialDifferencing struct {
	ComplexPacking

	SpatialOrderDifference uint8
	OctetsNumber           uint8
}
