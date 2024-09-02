package definition

import (
	"github.com/scorix/grib-go/pkg/grib2/drt"
	"github.com/scorix/grib-go/pkg/grib2/gdt"
	"github.com/scorix/grib-go/pkg/grib2/pdt"
)

// Section 0
//
// https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_sect0.shtml
//
// This section serves to identify the start of the record in a human readable form, indicate the total length of the message, and indicate the Edition number of GRIB used to construct or encode the message. For GRIB2, this section is always 16 octets long.
//
// don't edit
type Section0 struct {
	GribLiteral   [4]byte // "GRIB"
	Reserved      [2]byte
	Discipline    Discipline
	EditionNumber EditionNumber
	GribLength    uint64
}

// https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table0-0.shtml
type Discipline uint8

const (
	DisciplineMeteorologicalProducts         Discipline = 0
	DisciplineHydrologicalProducts           Discipline = 1
	DisciplineLandSurfaceProducts            Discipline = 2
	DisciplineSatelliteRemoteSensingProducts Discipline = 3
	DisciplineSpaceWeatherProducts           Discipline = 4
	// 5-9 Reserved
	DisciplineOceanographicProducts Discipline = 10
	// 11-19 Reserved
	DisciplineHealthAndSocioeconomicImpacts Discipline = 20
	// 21-191 Reserved
	// 192-254 Reserved For Local Use
	DisciplineMissing = 255
)

type EditionNumber uint8

const (
	EditionNumberGrib1 EditionNumber = 1
	EditionNumberGrib2 EditionNumber = 2
	EditionNumberGrib3 EditionNumber = 3
)

type Section1 struct {
	Section1FixedPart
	Reserved []byte // 22-N Reserved
}

// don't edit
type Section1FixedPart struct {
	Section1Length                  uint32        // Length of the section in octets (21 or N)
	NumberOfSection                 uint8         // 1 - Number of the section
	Center                          uint16        // Identification of originating/generating center: https://www.nco.ncep.noaa.gov/pmb/docs/on388/table0.html
	SubCenter                       uint16        // Identification of originating/generating subcenter: https://www.nco.ncep.noaa.gov/pmb/docs/on388/tablec.html
	TableVersion                    uint8         // GRIB master tables version number (currently 2): https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-0.shtml
	LocalTableVersion               uint8         // Version number of GRIB local tables used to augment Master Tables: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-1.shtml
	SignificanceOfReferenceTime     ReferenceTime // Significance of reference time: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-2.shtml
	Year                            uint16        // 4 digits
	Month                           uint8         // Month
	Day                             uint8         // Day
	Hour                            uint8         // Hour
	Minute                          uint8         // Minute
	Second                          uint8         // Second
	ProductionStatusOfProcessedData uint8         // Production Status of Processed data in the GRIB message: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-3.shtml
	TypeOfProcessedData             uint8         // Type of processed data in this GRIB message: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-4.shtml
}

type ReferenceTime uint8

const (
	ReferenceTimeAnalysis                ReferenceTime = 0
	ReferenceTimeStartOfForecast         ReferenceTime = 1
	ReferenceTimeVerifyingTimeOfForecast ReferenceTime = 2
	ReferenceTimeObservationTime         ReferenceTime = 3
	ReferenceTimeLocalTime               ReferenceTime = 4
	ReferenceTimeSimulationStart         ReferenceTime = 5
	// 6-191 Reserved
	// 192-254 Reserved For Local Use
	ReferenceTimeMissing ReferenceTime = 255
)

type Section2 struct {
	Section2FixedPart
	Local []byte // 6-N Local Use
}

// don't edit
type Section2FixedPart struct {
	Section2Length  uint32 // Length of the section in octets (N)
	NumberOfSection uint8  // 2 - Number of the section
}

type Section3 struct {
	Section3FixedPart
	GridDefinitionTemplate gdt.Template
}

// don't edit
type Section3FixedPart struct {
	Section3Length                   uint32 // Length of the section in octets (N)
	NumberOfSection                  uint8  // 3 - Number of the section
	SourceOfGridDefinition           uint8  // Source of grid definition
	NumberOfDataPoints               uint32
	NumberOfOctectsForNumberOfPoints uint8
	InterpretationOfNumberOfPoints   uint8
	GridDefinitionTemplateNumber     uint16
}

type Section4 struct {
	Section4FixedPart
	ProductDefinitionTemplate pdt.Template
}

// don't edit
type Section4FixedPart struct {
	Section4Length                  uint32 // Length of the section in octets (N)
	NumberOfSection                 uint8  // 4 - Number of the section
	NV                              uint16
	ProductDefinitionTemplateNumber uint16
}

type Section5 struct {
	Section5FixedPart
	DataRepresentationTemplate drt.Template
}

// don't edit
type Section5FixedPart struct {
	Section5Length                   uint32 // Length of the section in octets (N)
	NumberOfSection                  uint8  // 5 - Number of the section
	NumberOfValues                   uint32
	DataRepresentationTemplateNumber drt.TemplateNumber
}

type Section6 struct {
	Section6Length  uint32 // Length of the section in octets (N)
	NumberOfSection uint8  // 6 - Number of the section
	BitMapIndicator uint8
}

type Section7 struct {
	Section7FixedPart
	Data []byte
}

// don't edit
type Section7FixedPart struct {
	Section7Length  uint32 // Length of the section in octets (N)
	NumberOfSection uint8  // 7 - Number of the section
}

// don't edit
type Section8 struct {
	MagicNumber [4]byte
}
