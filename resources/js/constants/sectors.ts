// Sector options for Malawian SMEs
// These match the sectors defined in database/seeders/config_seeder.go

export const SECTOR_OPTIONS = [
  { value: 'Agriculture', label: 'Agriculture' },
  { value: 'Manufacturing', label: 'Manufacturing' },
  { value: 'Services', label: 'Services' },
  { value: 'Trade', label: 'Trade' },
  { value: 'Tourism and Hospitality', label: 'Tourism and Hospitality' },
  { value: 'Construction', label: 'Construction' },
  { value: 'Mining', label: 'Mining' },
  { value: 'Energy', label: 'Energy' },
  { value: 'Transport', label: 'Transport' },
  { value: 'ICT and Technology', label: 'ICT and Technology' },
] as const;

export const SECTOR_NAMES = SECTOR_OPTIONS.map(s => s.value);

export type Sector = typeof SECTOR_OPTIONS[number]['value'];
