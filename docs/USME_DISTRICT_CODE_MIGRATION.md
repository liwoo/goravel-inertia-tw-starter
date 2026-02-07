# USME District Code Migration Guide

This document explains how to update existing USME numbers to use the new district codes.

## Background

The USME (Unique SME/MSME) number format is: `MW-YYYY-DD-CT-NNNNNN-C`

Where:
- `MW` = Country code (Malawi)
- `YYYY` = Registration year
- `DD` = District code (2-3 letters)
- `CT` = Category code
- `NNNNNN` = Sequential number
- `C` = Luhn check digit

The district codes have been updated to better align with national standards:

| District | Old Code | New Code |
|----------|----------|----------|
| Chitipa | CT | CP |
| Karonga | KR | KA |
| Likoma | LK | LA |
| Mzimba | MH | MZ |
| Dedza | DE | DZ |
| Dowa | DO | DA |
| Kasungu | KS | KU |
| Lilongwe | LI | LL |
| Nkhotakota | NK | KK |
| Ntchisi | NI | NS |
| Balaka | BA | BLK |
| Chiradzulu | CR | CZ |
| Machinga | MG | MHG |
| Mangochi | MN | MH |
| Mwanza | MW | MN |
| Nsanje | NS | NE |
| Thyolo | TH | TO |
| Phalombe | PH | PE |
| Zomba | ZO | ZA |
| Neno | NE | NN |

## Prerequisites

1. Database backup - **Always create a backup before running migrations**
2. Application access to the database
3. Goravel CLI or direct database access

## Migration Steps

### Step 1: Run the Migration

The migration creates the necessary PostgreSQL functions. Run it using Goravel's migrate command:

```bash
# From the project root directory
go run . artisan migrate
```

Or if using the compiled binary:

```bash
./books-database artisan migrate
```

### Step 2: Preview Changes (Recommended)

Before making any changes, preview what will be updated:

```sql
-- Connect to your database and run:
SELECT * FROM preview_usme_district_code_updates();
```

This will show a table with:
- `sme_id`: The SME record ID
- `sme_name`: Business name
- `district`: District name
- `old_usme`: Current USME number
- `new_usme`: What the new USME will be
- `old_district_code`: Current district code
- `new_district_code`: New district code

Example output:
```
 sme_id |    sme_name     | district  |        old_usme         |         new_usme         | old_district_code | new_district_code
--------+-----------------+-----------+-------------------------+--------------------------+-------------------+-------------------
      1 | ABC Trading     | Balaka    | MW-2025-BA-MI-000001-5  | MW-2025-BLK-MI-000001-8  | BA                | BLK
      2 | XYZ Services    | Kasungu   | MW-2025-KS-SM-000001-3  | MW-2025-KU-SM-000001-6   | KS                | KU
```

### Step 3: Execute the Update

Once you've reviewed the preview and confirmed the changes are correct:

```sql
-- This will update all USME numbers and return a log of changes
SELECT * FROM update_usme_district_codes();
```

The function will:
1. Parse each USME number
2. Replace old district codes with new ones
3. Recalculate the Luhn check digit
4. Update the database record
5. Return a log of all changes made

### Step 4: Verify the Update

After running the update, verify the changes:

```sql
-- Count records by district code pattern
SELECT
    substring(usme_number FROM 9 FOR 3) AS district_code,
    COUNT(*) AS count
FROM smes
WHERE usme_number IS NOT NULL
  AND deleted_at IS NULL
GROUP BY substring(usme_number FROM 9 FOR 3)
ORDER BY count DESC;

-- Check for any old codes that might remain
SELECT id, name, district, usme_number
FROM smes
WHERE usme_number LIKE 'MW-%-BA-%'  -- Old Balaka code
   OR usme_number LIKE 'MW-%-KS-%'  -- Old Kasungu code
   OR usme_number LIKE 'MW-%-LI-%'  -- Old Lilongwe code
   -- Add other old codes as needed
;
```

## Rollback

If you need to revert the migration (removes only the functions, not the data changes):

```bash
go run . artisan migrate:rollback
```

**Note:** Rolling back the migration will only remove the PostgreSQL functions. It will NOT revert the USME number changes. To revert data changes, you must restore from backup.

## Manual SQL Execution

If you prefer to run the SQL directly without using the migration framework:

### Create Functions Manually

```sql
-- 1. Create Luhn check digit function
CREATE OR REPLACE FUNCTION calculate_luhn_check_digit(usme_without_check VARCHAR)
RETURNS CHAR AS $$
DECLARE
    digits TEXT;
    digit_array INT[];
    i INT;
    sum_val INT := 0;
    digit INT;
    doubled INT;
    check_digit INT;
BEGIN
    digits := regexp_replace(usme_without_check, '[^0-9]', '', 'g');
    digit_array := ARRAY[]::INT[];
    FOR i IN 1..length(digits) LOOP
        digit_array := array_append(digit_array, CAST(substring(digits FROM i FOR 1) AS INT));
    END LOOP;
    FOR i IN REVERSE array_length(digit_array, 1)..1 LOOP
        digit := digit_array[i];
        IF (array_length(digit_array, 1) - i) % 2 = 0 THEN
            doubled := digit * 2;
            IF doubled > 9 THEN
                doubled := doubled - 9;
            END IF;
            sum_val := sum_val + doubled;
        ELSE
            sum_val := sum_val + digit;
        END IF;
    END LOOP;
    check_digit := (10 - (sum_val % 10)) % 10;
    RETURN CAST(check_digit AS CHAR);
END;
$$ LANGUAGE plpgsql IMMUTABLE;
```

### Single District Update (Example)

If you only need to update a specific district:

```sql
-- Update Balaka from BA to BLK
UPDATE smes
SET usme_number =
    substring(usme_number FROM 1 FOR 8) ||
    'BLK' ||
    substring(usme_number FROM 11)
WHERE usme_number LIKE 'MW-%-BA-%'
  AND deleted_at IS NULL;

-- Then recalculate check digits for affected records
-- (Run the full update function or manually update check digits)
```

## Troubleshooting

### Function Not Found

If you get "function does not exist" errors:
```sql
-- Check if functions exist
SELECT routine_name
FROM information_schema.routines
WHERE routine_type = 'FUNCTION'
  AND routine_name LIKE '%usme%';
```

If they don't exist, run the migration again or create them manually.

### Check Digit Validation Failed

The Luhn check digit ensures data integrity. If a USME number fails validation after update:

```sql
-- Recalculate check digit for a specific record
SELECT
    id,
    usme_number,
    substring(usme_number FROM 1 FOR position('-' IN reverse(usme_number)) * -1 + length(usme_number)) ||
    '-' ||
    calculate_luhn_check_digit(
        substring(usme_number FROM 1 FOR position('-' IN reverse(usme_number)) * -1 + length(usme_number))
    ) AS corrected_usme
FROM smes
WHERE id = <your_sme_id>;
```

### Conflicting District Codes

Some district codes were swapped (e.g., MN → MH, MH → MZ). The update function handles this by:
1. Using temporary placeholders during the first pass
2. Resolving placeholders to final codes in the second pass
3. Recalculating all check digits

## Support

For issues with this migration, contact the development team or refer to the main application documentation.
