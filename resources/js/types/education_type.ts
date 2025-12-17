export type EducationType = 'None' | 'Primary' | 'Secondary' | 'Tertiary' | 'Postgraduate';

export const EDUCATION_TYPE_OPTIONS: { value: EducationType; label: string }[] = [
  { value: 'None', label: 'Education None' },
  { value: 'Primary', label: 'Education Primary' },
  { value: 'Secondary', label: 'Education Secondary' },
  { value: 'Tertiary', label: 'Education Tertiary' },
  { value: 'Postgraduate', label: 'Education Postgraduate' }
];
