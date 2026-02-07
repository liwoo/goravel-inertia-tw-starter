export type CitizenshipStatusType = 'Citizen' | 'Permanent Resident' | 'Work Permit' | 'Student Visa' | 'Visitor';

export const CITIZENSHIP_STATUS_OPTIONS: { value: CitizenshipStatusType; label: string }[] = [
  { value: 'Citizen', label: 'Citizen' },
  { value: 'Permanent Resident', label: 'Permanent Resident' },
  { value: 'Work Permit', label: 'Work Permit Holder' },
  { value: 'Student Visa', label: 'Student Visa' },
  { value: 'Visitor', label: 'Visitor' },
];

// Backward compatibility aliases
export type MalawianStatusType = CitizenshipStatusType;
export const MALAWIAN_STATUS_OPTIONS = CITIZENSHIP_STATUS_OPTIONS;
