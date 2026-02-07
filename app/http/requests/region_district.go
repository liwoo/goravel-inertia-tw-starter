package requests

// Region represents the three administrative regions of Malawi
type Region string

const (
	RegionNorthern Region = "Northern Region"
	RegionCentral  Region = "Central Region"
	RegionSouthern Region = "Southern Region"
)

// District represents the 28 administrative districts of Malawi
type District string

// Northern Region Districts (6 districts)
const (
	DistrictChitipa   District = "Chitipa"
	DistrictKaronga   District = "Karonga"
	DistrictLikoma    District = "Likoma"
	DistrictMzimba    District = "Mzimba"
	DistrictNkhataBay District = "Nkhata Bay"
	DistrictRumphi    District = "Rumphi"
)

// Central Region Districts (9 districts)
const (
	DistrictDedza      District = "Dedza"
	DistrictDowa       District = "Dowa"
	DistrictKasungu    District = "Kasungu"
	DistrictLilongwe   District = "Lilongwe"
	DistrictMchinji    District = "Mchinji"
	DistrictNkhotakota District = "Nkhotakota"
	DistrictNtcheu     District = "Ntcheu"
	DistrictNtchisi    District = "Ntchisi"
	DistrictSalima     District = "Salima"
)

// Southern Region Districts (13 districts)
const (
	DistrictBalaka     District = "Balaka"
	DistrictBlantyre   District = "Blantyre"
	DistrictChikwawa   District = "Chikwawa"
	DistrictChiradzulu District = "Chiradzulu"
	DistrictMachinga   District = "Machinga"
	DistrictMangochi   District = "Mangochi"
	DistrictMulanje    District = "Mulanje"
	DistrictMwanza     District = "Mwanza"
	DistrictNeno       District = "Neno"
	DistrictNsanje     District = "Nsanje"
	DistrictPhalombe   District = "Phalombe"
	DistrictThyolo     District = "Thyolo"
	DistrictZomba      District = "Zomba"
)

// DistrictToRegionMap maps each district to its respective region
var DistrictToRegionMap = map[District]Region{
	// Northern Region
	DistrictChitipa:   RegionNorthern,
	DistrictKaronga:   RegionNorthern,
	DistrictLikoma:    RegionNorthern,
	DistrictMzimba:    RegionNorthern,
	DistrictNkhataBay: RegionNorthern,
	DistrictRumphi:    RegionNorthern,

	// Central Region
	DistrictDedza:      RegionCentral,
	DistrictDowa:       RegionCentral,
	DistrictKasungu:    RegionCentral,
	DistrictLilongwe:   RegionCentral,
	DistrictMchinji:    RegionCentral,
	DistrictNkhotakota: RegionCentral,
	DistrictNtcheu:     RegionCentral,
	DistrictNtchisi:    RegionCentral,
	DistrictSalima:     RegionCentral,

	// Southern Region
	DistrictBalaka:     RegionSouthern,
	DistrictBlantyre:   RegionSouthern,
	DistrictChikwawa:   RegionSouthern,
	DistrictChiradzulu: RegionSouthern,
	DistrictMachinga:   RegionSouthern,
	DistrictMangochi:   RegionSouthern,
	DistrictMulanje:    RegionSouthern,
	DistrictMwanza:     RegionSouthern,
	DistrictNeno:       RegionSouthern,
	DistrictNsanje:     RegionSouthern,
	DistrictPhalombe:   RegionSouthern,
	DistrictThyolo:     RegionSouthern,
	DistrictZomba:      RegionSouthern,
}

// GetRegionByDistrict returns the region for a given district
func GetRegionByDistrict(district District) (Region, bool) {
	region, exists := DistrictToRegionMap[district]
	return region, exists
}

// GetDistrictsByRegion returns all districts in a given region
func GetDistrictsByRegion(region Region) []District {
	var districts []District
	for district, r := range DistrictToRegionMap {
		if r == region {
			districts = append(districts, district)
		}
	}
	return districts
}

// GetAllRegions returns all regions
func GetAllRegions() []Region {
	return []Region{
		RegionNorthern,
		RegionCentral,
		RegionSouthern,
	}
}

// GetAllDistricts returns all districts
func GetAllDistricts() []District {
	return []District{
		// Northern Region
		DistrictChitipa,
		DistrictKaronga,
		DistrictLikoma,
		DistrictMzimba,
		DistrictNkhataBay,
		DistrictRumphi,

		// Central Region
		DistrictDedza,
		DistrictDowa,
		DistrictKasungu,
		DistrictLilongwe,
		DistrictMchinji,
		DistrictNkhotakota,
		DistrictNtcheu,
		DistrictNtchisi,
		DistrictSalima,

		// Southern Region
		DistrictBalaka,
		DistrictBlantyre,
		DistrictChikwawa,
		DistrictChiradzulu,
		DistrictMachinga,
		DistrictMangochi,
		DistrictMulanje,
		DistrictMwanza,
		DistrictNeno,
		DistrictNsanje,
		DistrictPhalombe,
		DistrictThyolo,
		DistrictZomba,
	}
}

// ValidateDistrict checks if a district name is valid
func ValidateDistrict(district string) bool {
	_, exists := DistrictToRegionMap[District(district)]
	return exists
}

// ValidateRegion checks if a region name is valid
func ValidateRegion(region string) bool {
	r := Region(region)
	return r == RegionNorthern || r == RegionCentral || r == RegionSouthern
}
