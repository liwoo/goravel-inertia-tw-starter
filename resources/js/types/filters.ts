// Type definitions for the custom filter system

export type FilterType = 
  | 'string' 
  | 'number' 
  | 'date' 
  | 'datetime' 
  | 'boolean' 
  | 'enum' 
  | 'array';

export type FilterOperator = 
  // Common operators
  | 'equals'
  | 'not_equals'
  | 'is_null'
  | 'is_not_null'
  // String operators
  | 'contains'
  | 'not_contains'
  | 'starts_with'
  | 'ends_with'
  | 'is_empty'
  | 'is_not_empty'
  | 'regex_match'
  // Number operators
  | 'greater_than'
  | 'less_than'
  | 'greater_than_or_equal'
  | 'less_than_or_equal'
  | 'between'
  | 'not_between'
  // Date/DateTime operators
  | 'before'
  | 'after'
  | 'is_today'
  | 'is_yesterday'
  | 'is_this_week'
  | 'is_this_month'
  | 'is_this_year'
  | 'last_n_days'
  | 'next_n_days'
  // Boolean operators
  | 'is_true'
  | 'is_false'
  // Array operators
  | 'in'
  | 'not_in'
  | 'contains_any'
  | 'contains_all';

export type LogicOperator = 'AND' | 'OR';

export interface FilterDefinition {
  field: string;
  label: string;
  type: FilterType;
  operators: FilterOperator[];
  validation?: Record<string, any>;
  format?: string;
  enum_values?: string[];
  default?: any;
}

export interface FilterCondition {
  field: string;
  operator: FilterOperator;
  value: any;
  type?: FilterType;
}

export interface CompoundFilter {
  logic: LogicOperator;
  conditions: (FilterCondition | CompoundFilter)[];
}

export interface FilterMetadata {
  filters: FilterDefinition[];
  logic_operators: LogicOperator[];
  resource: string;
  searchable_fields: string[];
  sortable_fields: string[];
  filterable_fields: string[];
}

// Helper to get operator label
export function getOperatorLabel(operator: FilterOperator): string {
  const labels: Record<FilterOperator, string> = {
    // Common
    equals: 'Equals',
    not_equals: 'Not equals',
    is_null: 'Is empty',
    is_not_null: 'Is not empty',
    // String
    contains: 'Contains',
    not_contains: 'Does not contain',
    starts_with: 'Starts with',
    ends_with: 'Ends with',
    is_empty: 'Is empty',
    is_not_empty: 'Is not empty',
    regex_match: 'Matches pattern',
    // Number
    greater_than: 'Greater than',
    less_than: 'Less than',
    greater_than_or_equal: 'Greater than or equal',
    less_than_or_equal: 'Less than or equal',
    between: 'Between',
    not_between: 'Not between',
    // Date
    before: 'Before',
    after: 'After',
    is_today: 'Is today',
    is_yesterday: 'Is yesterday',
    is_this_week: 'This week',
    is_this_month: 'This month',
    is_this_year: 'This year',
    last_n_days: 'Last N days',
    next_n_days: 'Next N days',
    // Boolean
    is_true: 'Is true',
    is_false: 'Is false',
    // Array
    in: 'In',
    not_in: 'Not in',
    contains_any: 'Contains any',
    contains_all: 'Contains all',
  };
  return labels[operator] || operator;
}

// Helper to check if operator requires value
export function operatorRequiresValue(operator: FilterOperator): boolean {
  const noValueOperators: FilterOperator[] = [
    'is_null',
    'is_not_null',
    'is_empty',
    'is_not_empty',
    'is_true',
    'is_false',
    'is_today',
    'is_yesterday',
    'is_this_week',
    'is_this_month',
    'is_this_year'
  ];
  return !noValueOperators.includes(operator);
}

// Helper to check if operator requires two values (range)
export function operatorRequiresRange(operator: FilterOperator): boolean {
  return operator === 'between' || operator === 'not_between';
}

// Helper to check if operator requires array value
export function operatorRequiresArray(operator: FilterOperator): boolean {
  return ['in', 'not_in', 'contains_any', 'contains_all'].includes(operator);
}

// Helper to check if operator requires number value
export function operatorRequiresNumber(operator: FilterOperator): boolean {
  return operator === 'last_n_days' || operator === 'next_n_days';
}