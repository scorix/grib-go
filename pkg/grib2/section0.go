package grib

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

// Section 0
//
// https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_sect0.shtml
//
// This section serves to identify the start of the record in a human readable form, indicate the total length of the message, and indicate the Edition number of GRIB used to construct or encode the message. For GRIB2, this section is always 16 octets long.
type Section0 struct {
	GribLiteral   [4]byte // "GRIB"
	Reserved      [2]byte
	Discipline    Discipline
	EditionNumber EditionNumber
	Length        uint64
}

// Edition number
//
// - 1 for GRIB1
// - 2 for GRIB2
func (s *Section0) GetEditionNumber() int {
	return int(s.EditionNumber)
}

// Discipline (From [Table 0.0](https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table0-0.shtml))
func (s *Section0) GetDiscipline() int {
	return int(s.Discipline)
}

// Total length of GRIB message in octets (All sections)
func (s *Section0) GetGribLength() int {
	return int(s.Length)
}
