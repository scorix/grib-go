package gdt

type ScanningMode interface {
	GetScanMode() int8

	GetGridPointLL(i, j int) (lat, lon float32)
	GetGridPointFromLL(lat float32, lon float32) (i, j, n int)
}

type ScanningMode0000 struct {
	Ni                          int32 `json:"ni"`
	Nj                          int32 `json:"nj"`
	LatitudeOfFirstGridPoint    int32 `json:"latitudeOfFirstGridPoint"`
	LongitudeOfFirstGridPoint   int32 `json:"longitudeOfFirstGridPoint"`
	ResolutionAndComponentFlags int8  `json:"resolutionAndComponentFlags"`
	LatitudeOfLastGridPoint     int32 `json:"latitudeOfLastGridPoint"`
	LongitudeOfLastGridPoint    int32 `json:"longitudeOfLastGridPoint"`
	IDirectionIncrement         int32 `json:"iDirectionIncrement"`
	JDirectionIncrement         int32 `json:"jDirectionIncrement"`
	N                           int32 `json:"n,omitempty"`
	getGridIndexFunc            func(lat, lon float32) (i, j, n int)
	getGridPointByIndexFunc     func(i, j int) (lat, lon float32)
}

func (sm *ScanningMode0000) GetGridPointLL(i, j int) (lat, lon float32) {
	return sm.getGridPointByIndexFunc(i, j)
}

func (sm *ScanningMode0000) GetGridPointFromLL(lat float32, lon float32) (i, j, n int) {
	return sm.getGridIndexFunc(lat, lon)
}

func (sm *ScanningMode0000) GetScanMode() int8 {
	return 0
}
