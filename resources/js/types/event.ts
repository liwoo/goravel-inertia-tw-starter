// TypeScript interfaces for Event entities and operations
import { BaseModel, PaginatedResult, ListRequest } from './crud';

// Enum types
export type EducationType = 'None' | 'Primary' | 'Secondary' | 'Tertiary' | 'Postgraduate';

export const EDUCATION_TYPE_OPTIONS: { value: EducationType; label: string }[] = [
  { value: 'None', label: 'Education None' },
  { value: 'Primary', label: 'Education Primary' },
  { value: 'Secondary', label: 'Education Secondary' },
  { value: 'Tertiary', label: 'Education Tertiary' },
  { value: 'Postgraduate', label: 'Education Postgraduate' }
];

export type GenderType = 'MALE' | 'FEMALE';

export const GENDER_TYPE_OPTIONS: { value: GenderType; label: string }[] = [
  { value: 'MALE', label: 'Gender Male' },
  { value: 'FEMALE', label: 'Gender Female' }
];

export type MalawianStatusType = 'Citizen' | 'Permanent Resident' | 'Work Permit' | 'Student Visa' | 'Visitor';

export const MALAWIAN_STATUS_TYPE_OPTIONS: { value: MalawianStatusType; label: string }[] = [
  { value: 'Citizen', label: 'Malawian Status Citizen' },
  { value: 'Permanent Resident', label: 'Malawian Status Permanent Resident' },
  { value: 'Work Permit', label: 'Malawian Status Work Permit' },
  { value: 'Student Visa', label: 'Malawian Status Student Visa' },
  { value: 'Visitor', label: 'Malawian Status Visitor' }
];

export type Nationality = 'Afghan' | 'Albanian' | 'Algerian' | 'American' | 'Andorran' | 'Angolan' | 'Anguillan' | 'Citizen of Antigua and Barbuda' | 'Argentine' | 'Armenian' | 'Australian' | 'Austrian' | 'Azerbaijani' | 'Bahamian' | 'Bahraini' | 'Bangladeshi' | 'Barbadian' | 'Belarusian' | 'Belgian' | 'Belizean' | 'Beninese' | 'Bermudian' | 'Bhutanese' | 'Bolivian' | 'Citizen of Bosnia and Herzegovina' | 'Botswanan' | 'Brazilian' | 'British' | 'British Virgin Islander' | 'Bruneian' | 'Bulgarian' | 'Burkinan' | 'Burmese' | 'Burundian' | 'Cambodian' | 'Cameroonian' | 'Canadian' | 'Cape Verdean' | 'Cayman Islander' | 'Central African' | 'Chadian' | 'Chilean' | 'Chinese' | 'Colombian' | 'Comoran' | 'Congolese (Congo)' | 'Congolese (DRC)' | 'Cook Islander' | 'Costa Rican' | 'Croatian' | 'Cuban' | 'Cymraes' | 'Cymro' | 'Cypriot' | 'Czech' | 'Danish' | 'Djiboutian' | 'Dominican' | 'Citizen of the Dominican Republic' | 'Dutch' | 'East Timorese' | 'Ecuadorean' | 'Egyptian' | 'Emirati' | 'English' | 'Equatorial Guinean' | 'Eritrean' | 'Estonian' | 'Ethiopian' | 'Faroese' | 'Fijian' | 'Filipino' | 'Finnish' | 'French' | 'Gabonese' | 'Gambian' | 'Georgian' | 'German' | 'Ghanaian' | 'Gibraltarian' | 'Greek' | 'Greenlandic' | 'Grenadian' | 'Guamanian' | 'Guatemalan' | 'Citizen of Guinea-Bissau' | 'Guinean' | 'Guyanese' | 'Haitian' | 'Honduran' | 'Hong Konger' | 'Hungarian' | 'Icelandic' | 'Indian' | 'Indonesian' | 'Iranian' | 'Iraqi' | 'Irish' | 'Israeli' | 'Italian' | 'Ivorian' | 'Jamaican' | 'Japanese' | 'Jordanian' | 'Kazakh' | 'Kenyan' | 'Kittitian' | 'Citizen of Kiribati' | 'Kosovan' | 'Kuwaiti' | 'Kyrgyz' | 'Lao' | 'Latvian' | 'Lebanese' | 'Liberian' | 'Libyan' | 'Liechtenstein citizen' | 'Lithuanian' | 'Luxembourger' | 'Macanese' | 'Macedonian' | 'Malagasy' | 'Malawian' | 'Malaysian' | 'Maldivian' | 'Malian' | 'Maltese' | 'Marshallese' | 'Martiniquais' | 'Mauritanian' | 'Mauritian' | 'Mexican' | 'Micronesian' | 'Moldovan' | 'Monegasque' | 'Mongolian' | 'Montenegrin' | 'Montserratian' | 'Moroccan' | 'Mosotho' | 'Mozambican' | 'Namibian' | 'Nauruan' | 'Nepalese' | 'New Zealander' | 'Nicaraguan' | 'Nigerian' | 'Nigerien' | 'Niuean' | 'North Korean' | 'Northern Irish' | 'Norwegian' | 'Omani' | 'Pakistani' | 'Palauan' | 'Palestinian' | 'Panamanian' | 'Papua New Guinean' | 'Paraguayan' | 'Peruvian' | 'Pitcairn Islander' | 'Polish' | 'Portuguese' | 'Prydeinig' | 'Puerto Rican' | 'Qatari' | 'Romanian' | 'Russian' | 'Rwandan' | 'Salvadorean' | 'Sammarinese' | 'Samoan' | 'Sao Tomean' | 'Saudi Arabian' | 'Scottish' | 'Senegalese' | 'Serbian' | 'Citizen of Seychelles' | 'Sierra Leonean' | 'Singaporean' | 'Slovak' | 'Slovenian' | 'Solomon Islander' | 'Somali' | 'South African' | 'South Korean' | 'South Sudanese' | 'Spanish' | 'Sri Lankan' | 'St Helenian' | 'St Lucian' | 'Stateless' | 'Sudanese' | 'Surinamese' | 'Swazi' | 'Swedish' | 'Swiss' | 'Syrian' | 'Taiwanese' | 'Tajik' | 'Tanzanian' | 'Thai' | 'Togolese' | 'Tongan' | 'Trinidadian' | 'Tristanian' | 'Tunisian' | 'Turkish' | 'Turkmen' | 'Turks and Caicos Islander' | 'Tuvaluan' | 'Ugandan' | 'Ukrainian' | 'Uruguayan' | 'Uzbek' | 'Vatican citizen' | 'Citizen of Vanuatu' | 'Venezuelan' | 'Vietnamese' | 'Vincentian' | 'Wallisian' | 'Welsh' | 'Yemeni' | 'Zambian' | 'Zimbabwean';

export const NATIONALITY_OPTIONS: { value: Nationality; label: string }[] = [
  { value: 'Afghan', label: 'Afghan' },
  { value: 'Albanian', label: 'Albanian' },
  { value: 'Algerian', label: 'Algerian' },
  { value: 'American', label: 'American' },
  { value: 'Andorran', label: 'Andorran' },
  { value: 'Angolan', label: 'Angolan' },
  { value: 'Anguillan', label: 'Anguillan' },
  { value: 'Citizen of Antigua and Barbuda', label: 'Antigua And Barbuda' },
  { value: 'Argentine', label: 'Argentine' },
  { value: 'Armenian', label: 'Armenian' },
  { value: 'Australian', label: 'Australian' },
  { value: 'Austrian', label: 'Austrian' },
  { value: 'Azerbaijani', label: 'Azerbaijani' },
  { value: 'Bahamian', label: 'Bahamian' },
  { value: 'Bahraini', label: 'Bahraini' },
  { value: 'Bangladeshi', label: 'Bangladeshi' },
  { value: 'Barbadian', label: 'Barbadian' },
  { value: 'Belarusian', label: 'Belarusian' },
  { value: 'Belgian', label: 'Belgian' },
  { value: 'Belizean', label: 'Belizean' },
  { value: 'Beninese', label: 'Beninese' },
  { value: 'Bermudian', label: 'Bermudian' },
  { value: 'Bhutanese', label: 'Bhutanese' },
  { value: 'Bolivian', label: 'Bolivian' },
  { value: 'Citizen of Bosnia and Herzegovina', label: 'Bosnia And Herzegovina' },
  { value: 'Botswanan', label: 'Botswanan' },
  { value: 'Brazilian', label: 'Brazilian' },
  { value: 'British', label: 'British' },
  { value: 'British Virgin Islander', label: 'British Virgin Islander' },
  { value: 'Bruneian', label: 'Bruneian' },
  { value: 'Bulgarian', label: 'Bulgarian' },
  { value: 'Burkinan', label: 'Burkinan' },
  { value: 'Burmese', label: 'Burmese' },
  { value: 'Burundian', label: 'Burundian' },
  { value: 'Cambodian', label: 'Cambodian' },
  { value: 'Cameroonian', label: 'Cameroonian' },
  { value: 'Canadian', label: 'Canadian' },
  { value: 'Cape Verdean', label: 'Cape Verdean' },
  { value: 'Cayman Islander', label: 'Cayman Islander' },
  { value: 'Central African', label: 'Central African' },
  { value: 'Chadian', label: 'Chadian' },
  { value: 'Chilean', label: 'Chilean' },
  { value: 'Chinese', label: 'Chinese' },
  { value: 'Colombian', label: 'Colombian' },
  { value: 'Comoran', label: 'Comoran' },
  { value: 'Congolese (Congo)', label: 'Congolese Congo' },
  { value: 'Congolese (DRC)', label: 'Congolese D R C' },
  { value: 'Cook Islander', label: 'Cook Islander' },
  { value: 'Costa Rican', label: 'Costa Rican' },
  { value: 'Croatian', label: 'Croatian' },
  { value: 'Cuban', label: 'Cuban' },
  { value: 'Cymraes', label: 'Cymraes' },
  { value: 'Cymro', label: 'Cymro' },
  { value: 'Cypriot', label: 'Cypriot' },
  { value: 'Czech', label: 'Czech' },
  { value: 'Danish', label: 'Danish' },
  { value: 'Djiboutian', label: 'Djiboutian' },
  { value: 'Dominican', label: 'Dominican' },
  { value: 'Citizen of the Dominican Republic', label: 'Dominican Republic' },
  { value: 'Dutch', label: 'Dutch' },
  { value: 'East Timorese', label: 'East Timorese' },
  { value: 'Ecuadorean', label: 'Ecuadorean' },
  { value: 'Egyptian', label: 'Egyptian' },
  { value: 'Emirati', label: 'Emirati' },
  { value: 'English', label: 'English' },
  { value: 'Equatorial Guinean', label: 'Equatorial Guinean' },
  { value: 'Eritrean', label: 'Eritrean' },
  { value: 'Estonian', label: 'Estonian' },
  { value: 'Ethiopian', label: 'Ethiopian' },
  { value: 'Faroese', label: 'Faroese' },
  { value: 'Fijian', label: 'Fijian' },
  { value: 'Filipino', label: 'Filipino' },
  { value: 'Finnish', label: 'Finnish' },
  { value: 'French', label: 'French' },
  { value: 'Gabonese', label: 'Gabonese' },
  { value: 'Gambian', label: 'Gambian' },
  { value: 'Georgian', label: 'Georgian' },
  { value: 'German', label: 'German' },
  { value: 'Ghanaian', label: 'Ghanaian' },
  { value: 'Gibraltarian', label: 'Gibraltarian' },
  { value: 'Greek', label: 'Greek' },
  { value: 'Greenlandic', label: 'Greenlandic' },
  { value: 'Grenadian', label: 'Grenadian' },
  { value: 'Guamanian', label: 'Guamanian' },
  { value: 'Guatemalan', label: 'Guatemalan' },
  { value: 'Citizen of Guinea-Bissau', label: 'Guinea Bissau' },
  { value: 'Guinean', label: 'Guinean' },
  { value: 'Guyanese', label: 'Guyanese' },
  { value: 'Haitian', label: 'Haitian' },
  { value: 'Honduran', label: 'Honduran' },
  { value: 'Hong Konger', label: 'Hong Konger' },
  { value: 'Hungarian', label: 'Hungarian' },
  { value: 'Icelandic', label: 'Icelandic' },
  { value: 'Indian', label: 'Indian' },
  { value: 'Indonesian', label: 'Indonesian' },
  { value: 'Iranian', label: 'Iranian' },
  { value: 'Iraqi', label: 'Iraqi' },
  { value: 'Irish', label: 'Irish' },
  { value: 'Israeli', label: 'Israeli' },
  { value: 'Italian', label: 'Italian' },
  { value: 'Ivorian', label: 'Ivorian' },
  { value: 'Jamaican', label: 'Jamaican' },
  { value: 'Japanese', label: 'Japanese' },
  { value: 'Jordanian', label: 'Jordanian' },
  { value: 'Kazakh', label: 'Kazakh' },
  { value: 'Kenyan', label: 'Kenyan' },
  { value: 'Kittitian', label: 'Kittitian' },
  { value: 'Citizen of Kiribati', label: 'Kiribati' },
  { value: 'Kosovan', label: 'Kosovan' },
  { value: 'Kuwaiti', label: 'Kuwaiti' },
  { value: 'Kyrgyz', label: 'Kyrgyz' },
  { value: 'Lao', label: 'Lao' },
  { value: 'Latvian', label: 'Latvian' },
  { value: 'Lebanese', label: 'Lebanese' },
  { value: 'Liberian', label: 'Liberian' },
  { value: 'Libyan', label: 'Libyan' },
  { value: 'Liechtenstein citizen', label: 'Liechtenstein' },
  { value: 'Lithuanian', label: 'Lithuanian' },
  { value: 'Luxembourger', label: 'Luxembourger' },
  { value: 'Macanese', label: 'Macanese' },
  { value: 'Macedonian', label: 'Macedonian' },
  { value: 'Malagasy', label: 'Malagasy' },
  { value: 'Malawian', label: 'Malawian' },
  { value: 'Malaysian', label: 'Malaysian' },
  { value: 'Maldivian', label: 'Maldivian' },
  { value: 'Malian', label: 'Malian' },
  { value: 'Maltese', label: 'Maltese' },
  { value: 'Marshallese', label: 'Marshallese' },
  { value: 'Martiniquais', label: 'Martiniquais' },
  { value: 'Mauritanian', label: 'Mauritanian' },
  { value: 'Mauritian', label: 'Mauritian' },
  { value: 'Mexican', label: 'Mexican' },
  { value: 'Micronesian', label: 'Micronesian' },
  { value: 'Moldovan', label: 'Moldovan' },
  { value: 'Monegasque', label: 'Monegasque' },
  { value: 'Mongolian', label: 'Mongolian' },
  { value: 'Montenegrin', label: 'Montenegrin' },
  { value: 'Montserratian', label: 'Montserratian' },
  { value: 'Moroccan', label: 'Moroccan' },
  { value: 'Mosotho', label: 'Mosotho' },
  { value: 'Mozambican', label: 'Mozambican' },
  { value: 'Namibian', label: 'Namibian' },
  { value: 'Nauruan', label: 'Nauruan' },
  { value: 'Nepalese', label: 'Nepalese' },
  { value: 'New Zealander', label: 'New Zealander' },
  { value: 'Nicaraguan', label: 'Nicaraguan' },
  { value: 'Nigerian', label: 'Nigerian' },
  { value: 'Nigerien', label: 'Nigerien' },
  { value: 'Niuean', label: 'Niuean' },
  { value: 'North Korean', label: 'North Korean' },
  { value: 'Northern Irish', label: 'Northern Irish' },
  { value: 'Norwegian', label: 'Norwegian' },
  { value: 'Omani', label: 'Omani' },
  { value: 'Pakistani', label: 'Pakistani' },
  { value: 'Palauan', label: 'Palauan' },
  { value: 'Palestinian', label: 'Palestinian' },
  { value: 'Panamanian', label: 'Panamanian' },
  { value: 'Papua New Guinean', label: 'Papua New Guinean' },
  { value: 'Paraguayan', label: 'Paraguayan' },
  { value: 'Peruvian', label: 'Peruvian' },
  { value: 'Pitcairn Islander', label: 'Pitcairn Islander' },
  { value: 'Polish', label: 'Polish' },
  { value: 'Portuguese', label: 'Portuguese' },
  { value: 'Prydeinig', label: 'Prydeinig' },
  { value: 'Puerto Rican', label: 'Puerto Rican' },
  { value: 'Qatari', label: 'Qatari' },
  { value: 'Romanian', label: 'Romanian' },
  { value: 'Russian', label: 'Russian' },
  { value: 'Rwandan', label: 'Rwandan' },
  { value: 'Salvadorean', label: 'Salvadorean' },
  { value: 'Sammarinese', label: 'Sammarinese' },
  { value: 'Samoan', label: 'Samoan' },
  { value: 'Sao Tomean', label: 'Sao Tomean' },
  { value: 'Saudi Arabian', label: 'Saudi Arabian' },
  { value: 'Scottish', label: 'Scottish' },
  { value: 'Senegalese', label: 'Senegalese' },
  { value: 'Serbian', label: 'Serbian' },
  { value: 'Citizen of Seychelles', label: 'Seychelles' },
  { value: 'Sierra Leonean', label: 'Sierra Leonean' },
  { value: 'Singaporean', label: 'Singaporean' },
  { value: 'Slovak', label: 'Slovak' },
  { value: 'Slovenian', label: 'Slovenian' },
  { value: 'Solomon Islander', label: 'Solomon Islander' },
  { value: 'Somali', label: 'Somali' },
  { value: 'South African', label: 'South African' },
  { value: 'South Korean', label: 'South Korean' },
  { value: 'South Sudanese', label: 'South Sudanese' },
  { value: 'Spanish', label: 'Spanish' },
  { value: 'Sri Lankan', label: 'Sri Lankan' },
  { value: 'St Helenian', label: 'St Helenian' },
  { value: 'St Lucian', label: 'St Lucian' },
  { value: 'Stateless', label: 'Stateless' },
  { value: 'Sudanese', label: 'Sudanese' },
  { value: 'Surinamese', label: 'Surinamese' },
  { value: 'Swazi', label: 'Swazi' },
  { value: 'Swedish', label: 'Swedish' },
  { value: 'Swiss', label: 'Swiss' },
  { value: 'Syrian', label: 'Syrian' },
  { value: 'Taiwanese', label: 'Taiwanese' },
  { value: 'Tajik', label: 'Tajik' },
  { value: 'Tanzanian', label: 'Tanzanian' },
  { value: 'Thai', label: 'Thai' },
  { value: 'Togolese', label: 'Togolese' },
  { value: 'Tongan', label: 'Tongan' },
  { value: 'Trinidadian', label: 'Trinidadian' },
  { value: 'Tristanian', label: 'Tristanian' },
  { value: 'Tunisian', label: 'Tunisian' },
  { value: 'Turkish', label: 'Turkish' },
  { value: 'Turkmen', label: 'Turkmen' },
  { value: 'Turks and Caicos Islander', label: 'Turks And Caicos Islander' },
  { value: 'Tuvaluan', label: 'Tuvaluan' },
  { value: 'Ugandan', label: 'Ugandan' },
  { value: 'Ukrainian', label: 'Ukrainian' },
  { value: 'Uruguayan', label: 'Uruguayan' },
  { value: 'Uzbek', label: 'Uzbek' },
  { value: 'Vatican citizen', label: 'Vatican Citizen' },
  { value: 'Citizen of Vanuatu', label: 'Vanuatu' },
  { value: 'Venezuelan', label: 'Venezuelan' },
  { value: 'Vietnamese', label: 'Vietnamese' },
  { value: 'Vincentian', label: 'Vincentian' },
  { value: 'Wallisian', label: 'Wallisian' },
  { value: 'Welsh', label: 'Welsh' },
  { value: 'Yemeni', label: 'Yemeni' },
  { value: 'Zambian', label: 'Zambian' },
  { value: 'Zimbabwean', label: 'Zimbabwean' }
];

export type Region = 'Northern Region' | 'Central Region' | 'Southern Region';

export const REGION_OPTIONS: { value: Region; label: string }[] = [
  { value: 'Northern Region', label: 'Northern' },
  { value: 'Central Region', label: 'Central' },
  { value: 'Southern Region', label: 'Southern' }
];

export type District = 'Chitipa' | 'Karonga' | 'Likoma' | 'Mzimba' | 'Nkhata Bay' | 'Rumphi' | 'Dedza' | 'Dowa' | 'Kasungu' | 'Lilongwe' | 'Mchinji' | 'Nkhotakota' | 'Ntcheu' | 'Ntchisi' | 'Salima' | 'Balaka' | 'Blantyre' | 'Chikwawa' | 'Chiradzulu' | 'Machinga' | 'Mangochi' | 'Mulanje' | 'Mwanza' | 'Neno' | 'Nsanje' | 'Phalombe' | 'Thyolo' | 'Zomba';

export const DISTRICT_OPTIONS: { value: District; label: string }[] = [
  { value: 'Chitipa', label: 'Chitipa' },
  { value: 'Karonga', label: 'Karonga' },
  { value: 'Likoma', label: 'Likoma' },
  { value: 'Mzimba', label: 'Mzimba' },
  { value: 'Nkhata Bay', label: 'Nkhata Bay' },
  { value: 'Rumphi', label: 'Rumphi' },
  { value: 'Dedza', label: 'Dedza' },
  { value: 'Dowa', label: 'Dowa' },
  { value: 'Kasungu', label: 'Kasungu' },
  { value: 'Lilongwe', label: 'Lilongwe' },
  { value: 'Mchinji', label: 'Mchinji' },
  { value: 'Nkhotakota', label: 'Nkhotakota' },
  { value: 'Ntcheu', label: 'Ntcheu' },
  { value: 'Ntchisi', label: 'Ntchisi' },
  { value: 'Salima', label: 'Salima' },
  { value: 'Balaka', label: 'Balaka' },
  { value: 'Blantyre', label: 'Blantyre' },
  { value: 'Chikwawa', label: 'Chikwawa' },
  { value: 'Chiradzulu', label: 'Chiradzulu' },
  { value: 'Machinga', label: 'Machinga' },
  { value: 'Mangochi', label: 'Mangochi' },
  { value: 'Mulanje', label: 'Mulanje' },
  { value: 'Mwanza', label: 'Mwanza' },
  { value: 'Neno', label: 'Neno' },
  { value: 'Nsanje', label: 'Nsanje' },
  { value: 'Phalombe', label: 'Phalombe' },
  { value: 'Thyolo', label: 'Thyolo' },
  { value: 'Zomba', label: 'Zomba' }
];

// Core Event interface matching the backend model
export interface Event extends BaseModel {
  title: string;
  description: string;
  date: string;
  end_date: string;
  venue: string;
  partners: string[];
  district: string;
  attending_smes: number[];
  attending_sme_details?: { id: number; name: string }[];
  notes?: string;
}

// Event creation data (matches EventCreateRequest)
export interface EventCreateData {
  title: string;
  description: string;
  date: string;
  end_date: string;
  venue: string;
  partners: string[];
  district: string;
  attending_smes: number[];
  notes?: string;
}

// Event update data (matches EventUpdateRequest - all optional)
export interface EventUpdateData {
  title?: string;
  description?: string;
  date?: string;
  end_date?: string;
  venue?: string;
  partners?: string[];
  district?: string;
  attending_smes?: number[];
  notes?: string;
}

// Event list response (matches service GetList response)
export interface EventListResponse extends PaginatedResult<Event> { }

// Event list request (extends base ListRequest with event-specific filters)
export interface EventListRequest extends ListRequest {
  // Add your custom filters here
}

// Form validation types
export interface EventFormErrors {
  title?: string;
  description?: string;
  date?: string;
  end_date?: string;
  venue?: string;
  partners?: string;
  district?: string;
  attending_smes?: string;
  notes?: string;
  general?: string;
}

// Event statistics (if provided by backend)
export interface EventStats {
  totalevents: number;
  // Add your custom stats here
}
