package grib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

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
	section  section1
	reserved []byte // 22-N Reserved
}

func (s *Section1) SectionLength() int {
	return s.section.SectionLength()
}

func (s *Section1) SectionNumber() int {
	return s.section.SectionNumber()
}

func (s *Section1) ReadFrom(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &s.section); err != nil {
		return fmt.Errorf("binary read: %w", err)
	}

	var reserved bytes.Buffer

	if _, err := io.Copy(&reserved, r); err != nil {
		return fmt.Errorf("copy reserved: %w", err)
	}

	s.reserved = reserved.Bytes()

	return nil
}

func (s *Section1) GetTime(loc *time.Location) time.Time {
	return time.Date(int(s.section.Year), time.Month(s.section.Month), int(s.section.Day), int(s.section.Hour), int(s.section.Minute), int(s.section.Second), 0, loc)
}

type section1 struct {
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

func (s *section1) SectionLength() int {
	return int(s.Section1Length)
}

func (s *section1) SectionNumber() int {
	return int(s.NumberOfSection)
}
