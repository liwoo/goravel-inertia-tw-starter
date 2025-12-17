export type MalawianStatusType = 'Citizen' | 'Permanent Resident' | 'Work Permit' | 'Student Visa' | 'Visitor';

export const MALAWIAN_STATUS_OPTIONS: { value: MalawianStatusType; label: string }[] = [
  { value: 'Citizen', label: 'Citizen' },
  { value: 'Permanent Resident', label: 'Permanent Resident' },
  { value: 'Work Permit', label: 'Work Permit Holder' },
  { value: 'Student Visa', label: 'Student Visa' },
  { value: 'Visitor', label: 'Visitor' },
];