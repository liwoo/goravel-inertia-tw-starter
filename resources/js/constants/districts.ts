// Malawi districts with their ISO 3166-2:MW codes
// These are used for UBI generation and location tracking

export interface District {
  name: string;
  code: string;
  region: 'Northern' | 'Central' | 'Southern';
}

export const DISTRICTS: readonly District[] = [
  // Northern Region
  { name: 'Chitipa', code: 'CP', region: 'Northern' },
  { name: 'Karonga', code: 'KA', region: 'Northern' },
  { name: 'Mzuzu', code: 'MZ', region: 'Northern' },
  { name: 'Nkhata Bay', code: 'NB', region: 'Northern' },
  { name: 'Rumphi', code: 'RU', region: 'Northern' },
  { name: 'Likoma', code: 'LA', region: 'Northern' },
  { name: 'Mzimba', code: 'MZ', region: 'Northern' },

  // Central Region
  { name: 'Dedza', code: 'DZ', region: 'Central' },
  { name: 'Dowa', code: 'DA', region: 'Central' },
  { name: 'Kasungu', code: 'KU', region: 'Central' },
  { name: 'Lilongwe', code: 'LL', region: 'Central' },
  { name: 'Mchinji', code: 'MC', region: 'Central' },
  { name: 'Nkhotakota', code: 'KK', region: 'Central' },
  { name: 'Ntcheu', code: 'NU', region: 'Central' },
  { name: 'Ntchisi', code: 'NS', region: 'Central' },
  { name: 'Salima', code: 'SA', region: 'Central' },

  // Southern Region
  { name: 'Balaka', code: 'BLK', region: 'Southern' },
  { name: 'Blantyre', code: 'BT', region: 'Southern' },
  { name: 'Chikwawa', code: 'CK', region: 'Southern' },
  { name: 'Chiradzulu', code: 'CZ', region: 'Southern' },
  { name: 'Machinga', code: 'MHG', region: 'Southern' },
  { name: 'Mangochi', code: 'MH', region: 'Southern' },
  { name: 'Mulanje', code: 'MJ', region: 'Southern' },
  { name: 'Mwanza', code: 'MN', region: 'Southern' },
  { name: 'Nsanje', code: 'NE', region: 'Southern' },
  { name: 'Thyolo', code: 'TO', region: 'Southern' },
  { name: 'Phalombe', code: 'PE', region: 'Southern' },
  { name: 'Zomba', code: 'ZA', region: 'Southern' },
  { name: 'Neno', code: 'NN', region: 'Southern' },
] as const;

// Type for district names
export type DistrictName = typeof DISTRICTS[number]['name'];

// Type for district codes
export type DistrictCode = typeof DISTRICTS[number]['code'];

// Type for regions
export type Region = 'Northern' | 'Central' | 'Southern';

// Helper function to get district by name
export function getDistrictByName(name: string): District | undefined {
  return DISTRICTS.find(d => d.name === name);
}

// Helper function to get district by code
export function getDistrictByCode(code: string): District | undefined {
  return DISTRICTS.find(d => d.code === code);
}

// Helper function to get region from district name
export function getRegionFromDistrict(districtName: string): Region | undefined {
  const district = getDistrictByName(districtName);
  return district?.region;
}

// Helper function to get all districts by region
export function getDistrictsByRegion(region: Region): District[] {
  return DISTRICTS.filter(d => d.region === region);
}

// Helper function to check if a string is a valid district name
export function isValidDistrictName(value: string): value is DistrictName {
  return DISTRICTS.some(d => d.name === value);
}

// Helper function to check if a string is a valid district code
export function isValidDistrictCode(value: string): value is DistrictCode {
  return DISTRICTS.some(d => d.code === value);
}

// Get all district names as an array
export const DISTRICT_NAMES = DISTRICTS.map(d => d.name);

// Get all district codes as an array
export const DISTRICT_CODES = DISTRICTS.map(d => d.code);
