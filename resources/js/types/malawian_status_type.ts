export type MalawianStatusType = 'Citizen' | 'Permanent Resident' | 'Work Permit' | 'Student Visa' | 'Visitor';

export const MALAWIAN_STATUS_TYPE_OPTIONS: { value: MalawianStatusType; label: string }[] = [
  { value: 'Citizen', label: 'Malawian Status Citizen' },
  { value: 'Permanent Resident', label: 'Malawian Status Permanent Resident' },
  { value: 'Work Permit', label: 'Malawian Status Work Permit' },
  { value: 'Student Visa', label: 'Malawian Status Student Visa' },
  { value: 'Visitor', label: 'Malawian Status Visitor' }
];
