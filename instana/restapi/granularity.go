package restapi

// Granularity custom type for an Alter Granularity
type Granularity int32

// Granularities custom type for a slice of Granularity
type Granularities []Granularity

// ToIntSlice Returns the corresponding int representations
func (granularities Granularities) ToIntSlice() []int {
	result := make([]int, len(granularities))
	for i, v := range granularities {
		result[i] = int(v)
	}
	return result
}

const (
	//Granularity60000 constant value for granularity of 1min(60 sec)
	Granularity60000 = Granularity(60000)
	//Granularity300000 constant value for granularity of 5min(300 sec)
	Granularity300000 = Granularity(300000)
	//Granularity600000 constant value for granularity of 10min(600 sec)
	Granularity600000 = Granularity(600000)
	//Granularity900000 constant value for granularity of 15min(900 sec)
	Granularity900000 = Granularity(900000)
	//Granularity1200000 constant value for granularity of 20min(1200 sec)
	Granularity1200000 = Granularity(1200000)
	//Granularity1800000 constant value for granularity of 30min(1800 sec)
	Granularity1800000 = Granularity(1800000)
)

// SupportedSmartAlertGranularities list of all supported Granularities
var SupportedSmartAlertGranularities = Granularities{Granularity60000, Granularity300000, Granularity600000, Granularity900000, Granularity1200000, Granularity1800000}
