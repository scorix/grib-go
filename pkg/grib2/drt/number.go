package drt

type DataRepresentationTemplateNumber uint16

const (
	GridPointDataSimplePacking                        DataRepresentationTemplateNumber = 0
	MatrixValueAtGridPointSimplePacking               DataRepresentationTemplateNumber = 1
	GridPointDataComplexPacking                       DataRepresentationTemplateNumber = 2
	GridPointDataComplexPackingAndSpatialDifferencing DataRepresentationTemplateNumber = 3
	GridPointDataIEEEFloatingPointData                DataRepresentationTemplateNumber = 4
	// 5-39 Reserved
	GridPointDataJPEG2000CodeStreamFormat DataRepresentationTemplateNumber = 40
	GridPointDataPNG                      DataRepresentationTemplateNumber = 41
	GridPointDataCCSDS                    DataRepresentationTemplateNumber = 42
	// 43-49 Reserved
	SpectralDataSimplePacking  DataRepresentationTemplateNumber = 50
	SpectralDataComplexPacking DataRepresentationTemplateNumber = 51
	// 52 Reserved
	SpectralDataComplexPackinForLimitedAreaModels DataRepresentationTemplateNumber = 53
	// 54-60 Reserved
	GridPointDataSimplePackingWithLogarithmPreProcessing DataRepresentationTemplateNumber = 61
	// 62-199 Reserved
	RunLengthPackingWithLevelValues DataRepresentationTemplateNumber = 200
	// 201-49151 Reserved
	// 49152-65534 Reserved For Local Use
	DataRepresentationTemplateNumberMissing DataRepresentationTemplateNumber = 255
)
