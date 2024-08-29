package grib

import "time"

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

type Section1 struct {
	section1
	reserved []byte // 22-N Reserved
}

type section1 struct {
	Length                      uint32        // Length of the section in octets (21 or N)
	SectionNumber               uint8         // 1 - Number of the section
	Center                      uint16        // Identification of originating/generating center: https://www.nco.ncep.noaa.gov/pmb/docs/on388/table0.html
	SubCenter                   uint16        // Identification of originating/generating subcenter: https://www.nco.ncep.noaa.gov/pmb/docs/on388/tablec.html
	MasterTableVersion          uint8         // GRIB master tables version number (currently 2): https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-0.shtml
	LocalTableVersion           uint8         // Version number of GRIB local tables used to augment Master Tables: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-1.shtml
	SignificanceOfReferenceTime ReferenceTime // Significance of reference time: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-2.shtml
	Year                        uint16        // 4 digits
	Month                       uint8         // Month
	Day                         uint8         // Day
	Hour                        uint8         // Hour
	Minute                      uint8         // Minute
	Second                      uint8         // Second
	DataProductionStatus        uint8         // Production Status of Processed data in the GRIB message: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-3.shtml
	DataType                    uint8         // Type of processed data in this GRIB message: https://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc/grib2_table1-4.shtml
}

func (s *Section1) GetSectionNumber() int {
	return int(s.SectionNumber)
}

func (s *Section1) GetTime() time.Time {
	return time.Date(int(s.Year), time.Month(s.Month), int(s.Day), int(s.Hour), int(s.Minute), int(s.Second), 0, time.UTC)
}

func (s *Section1) GetMasterTableVersion() int {
	return int(s.MasterTableVersion)
}

func (s *Section1) GetLocalTableVersion() int {
	return int(s.LocalTableVersion)
}
