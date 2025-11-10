export type EducationType = 'None' | 'Primary' | 'Secondary' | 'Tertiary' | 'Postgraduate';

export const EDUCATION_OPTIONS: { value: EducationType; label: string }[] = [
  { value: 'None', label: 'None' },
  { value: 'Primary', label: 'Primary School' },
  { value: 'Secondary', label: 'Secondary School' },
  { value: 'Tertiary', label: 'Tertiary Education' },
  { value: 'Postgraduate', label: 'Postgraduate' },
];