package migrations

import (
	"github.com/goravel/framework/facades"
)

type M20260112090000UpdateUsmeDistrictCodes struct{}

// Signature The unique signature for the migration.
func (r *M20260112090000UpdateUsmeDistrictCodes) Signature() string {
	return "20260112090000_update_usme_district_codes"
}

// Up Run the migrations.
func (r *M20260112090000UpdateUsmeDistrictCodes) Up() error {
	// Create the Luhn check digit calculation function
	_, err := facades.Orm().Query().Exec(`
		-- Function to calculate Luhn check digit for USME numbers
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
			-- Extract only numeric digits from the USME number (excluding the last check digit)
			digits := regexp_replace(usme_without_check, '[^0-9]', '', 'g');

			-- Convert to array of integers
			digit_array := ARRAY[]::INT[];
			FOR i IN 1..length(digits) LOOP
				digit_array := array_append(digit_array, CAST(substring(digits FROM i FOR 1) AS INT));
			END LOOP;

			-- Process digits from right to left
			FOR i IN REVERSE array_length(digit_array, 1)..1 LOOP
				digit := digit_array[i];
				-- Double every second digit from the right
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

			-- Calculate check digit
			check_digit := (10 - (sum_val % 10)) % 10;

			RETURN CAST(check_digit AS CHAR);
		END;
		$$ LANGUAGE plpgsql IMMUTABLE;
	`)
	if err != nil {
		return err
	}

	// Create the main update function
	_, err = facades.Orm().Query().Exec(`
		-- Function to update USME numbers with new district codes
		CREATE OR REPLACE FUNCTION update_usme_district_codes()
		RETURNS TABLE(
			sme_id BIGINT,
			old_usme VARCHAR,
			new_usme VARCHAR,
			old_district_code VARCHAR,
			new_district_code VARCHAR
		) AS $$
		DECLARE
			rec RECORD;
			parts TEXT[];
			old_code VARCHAR;
			new_code VARCHAR;
			new_usme_without_check VARCHAR;
			new_check_digit CHAR;
			updated_usme VARCHAR;
			update_count INT := 0;
		BEGIN
			-- District code mapping: OLD -> NEW
			-- Note: Order matters due to conflicts (e.g., MN->MH, MH->MZ)
			-- We use a temporary placeholder approach to handle conflicts

			-- Create temporary mapping table
			CREATE TEMP TABLE IF NOT EXISTS district_code_mapping (
				old_code VARCHAR(3),
				new_code VARCHAR(3),
				priority INT -- Lower number = process first
			) ON COMMIT DROP;

			-- Clear and populate mapping
			DELETE FROM district_code_mapping;

			-- Insert mappings with priority to handle conflicts
			-- Codes that could conflict are given lower priority (processed first with temp placeholder)
			INSERT INTO district_code_mapping (old_code, new_code, priority) VALUES
				-- Conflicting codes - process first with temp markers
				('MH', '__MZ__', 1),  -- Mzimba: MH -> MZ (temp)
				('MN', '__MH__', 1),  -- Mangochi: MN -> MH (temp, conflicts with old Mzimba)
				('MW', '__MN__', 1),  -- Mwanza: MW -> MN (temp, conflicts with old Mangochi)
				('NE', '__NN__', 1),  -- Neno: NE -> NN (temp)
				('NS', '__NE__', 1),  -- Nsanje: NS -> NE (temp, conflicts with old Neno)

				-- Non-conflicting codes
				('CT', 'CP', 2),      -- Chitipa
				('KR', 'KA', 2),      -- Karonga
				('LK', 'LA', 2),      -- Likoma
				('DE', 'DZ', 2),      -- Dedza
				('DO', 'DA', 2),      -- Dowa
				('KS', 'KU', 2),      -- Kasungu
				('LI', 'LL', 2),      -- Lilongwe
				('NK', 'KK', 2),      -- Nkhotakota
				('BA', 'BLK', 2),     -- Balaka
				('CR', 'CZ', 2),      -- Chiradzulu
				('MG', 'MHG', 2),     -- Machinga
				('NI', 'NS', 2),      -- Ntchisi
				('TH', 'TO', 2),      -- Thyolo
				('PH', 'PE', 2),      -- Phalombe
				('ZO', 'ZA', 2);      -- Zomba

			-- Process all SMEs with USME numbers
			FOR rec IN
				SELECT id, usme_number
				FROM smes
				WHERE usme_number IS NOT NULL
				  AND usme_number LIKE 'MW-%'
				  AND deleted_at IS NULL
				ORDER BY id
			LOOP
				-- Parse the USME number: MW-YYYY-DD-CT-NNNNNN-C
				parts := string_to_array(rec.usme_number, '-');

				IF array_length(parts, 1) >= 6 THEN
					old_code := parts[3]; -- District code is the 3rd part

					-- Check if this district code needs updating
					SELECT m.new_code INTO new_code
					FROM district_code_mapping m
					WHERE m.old_code = old_code
					ORDER BY m.priority
					LIMIT 1;

					IF new_code IS NOT NULL THEN
						-- Build new USME without check digit
						new_usme_without_check := parts[1] || '-' || parts[2] || '-' || new_code || '-' || parts[4] || '-' || parts[5];

						-- Calculate new Luhn check digit
						new_check_digit := calculate_luhn_check_digit(new_usme_without_check);

						-- Build complete new USME
						updated_usme := new_usme_without_check || '-' || new_check_digit;

						-- Return the mapping for logging
						sme_id := rec.id;
						old_usme := rec.usme_number;
						new_usme := updated_usme;
						old_district_code := old_code;
						new_district_code := new_code;
						RETURN NEXT;

						-- Update the record
						UPDATE smes SET usme_number = updated_usme WHERE id = rec.id;
						update_count := update_count + 1;
					END IF;
				END IF;
			END LOOP;

			-- Second pass: resolve temporary placeholders to final codes
			UPDATE smes SET usme_number = replace(usme_number, '__MZ__', 'MZ') WHERE usme_number LIKE '%__MZ__%';
			UPDATE smes SET usme_number = replace(usme_number, '__MH__', 'MH') WHERE usme_number LIKE '%__MH__%';
			UPDATE smes SET usme_number = replace(usme_number, '__MN__', 'MN') WHERE usme_number LIKE '%__MN__%';
			UPDATE smes SET usme_number = replace(usme_number, '__NN__', 'NN') WHERE usme_number LIKE '%__NN__%';
			UPDATE smes SET usme_number = replace(usme_number, '__NE__', 'NE') WHERE usme_number LIKE '%__NE__%';

			-- Recalculate check digits for all records that were updated with temp placeholders
			FOR rec IN
				SELECT id, usme_number
				FROM smes
				WHERE usme_number IS NOT NULL
				  AND usme_number LIKE 'MW-%'
				  AND deleted_at IS NULL
				  AND (
					  usme_number LIKE 'MW-%-MZ-%' OR
					  usme_number LIKE 'MW-%-MH-%' OR
					  usme_number LIKE 'MW-%-MN-%' OR
					  usme_number LIKE 'MW-%-NN-%' OR
					  usme_number LIKE 'MW-%-NE-%'
				  )
			LOOP
				parts := string_to_array(rec.usme_number, '-');
				IF array_length(parts, 1) >= 6 THEN
					-- Rebuild without check digit
					new_usme_without_check := parts[1] || '-' || parts[2] || '-' || parts[3] || '-' || parts[4] || '-' || parts[5];
					new_check_digit := calculate_luhn_check_digit(new_usme_without_check);
					updated_usme := new_usme_without_check || '-' || new_check_digit;

					IF updated_usme != rec.usme_number THEN
						UPDATE smes SET usme_number = updated_usme WHERE id = rec.id;
					END IF;
				END IF;
			END LOOP;

			RAISE NOTICE 'Updated % USME numbers with new district codes', update_count;
		END;
		$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		return err
	}

	// Create a dry-run function that shows what would be updated without making changes
	_, err = facades.Orm().Query().Exec(`
		-- Dry-run function to preview changes without updating
		CREATE OR REPLACE FUNCTION preview_usme_district_code_updates()
		RETURNS TABLE(
			sme_id BIGINT,
			sme_name VARCHAR,
			district VARCHAR,
			old_usme VARCHAR,
			new_usme VARCHAR,
			old_district_code VARCHAR,
			new_district_code VARCHAR
		) AS $$
		DECLARE
			rec RECORD;
			parts TEXT[];
			old_code VARCHAR;
			new_code VARCHAR;
			new_usme_without_check VARCHAR;
			new_check_digit CHAR;
			updated_usme VARCHAR;
		BEGIN
			-- District code mapping: OLD -> NEW
			CREATE TEMP TABLE IF NOT EXISTS district_code_mapping_preview (
				old_code VARCHAR(3),
				new_code VARCHAR(3)
			) ON COMMIT DROP;

			DELETE FROM district_code_mapping_preview;

			INSERT INTO district_code_mapping_preview (old_code, new_code) VALUES
				('CT', 'CP'),      -- Chitipa
				('KR', 'KA'),      -- Karonga
				('LK', 'LA'),      -- Likoma
				('MH', 'MZ'),      -- Mzimba
				('DE', 'DZ'),      -- Dedza
				('DO', 'DA'),      -- Dowa
				('KS', 'KU'),      -- Kasungu
				('LI', 'LL'),      -- Lilongwe
				('NK', 'KK'),      -- Nkhotakota
				('BA', 'BLK'),     -- Balaka
				('CR', 'CZ'),      -- Chiradzulu
				('MG', 'MHG'),     -- Machinga
				('NI', 'NS'),      -- Ntchisi
				('MN', 'MH'),      -- Mangochi
				('MW', 'MN'),      -- Mwanza
				('NS', 'NE'),      -- Nsanje
				('TH', 'TO'),      -- Thyolo
				('PH', 'PE'),      -- Phalombe
				('ZO', 'ZA'),      -- Zomba
				('NE', 'NN');      -- Neno

			FOR rec IN
				SELECT s.id, s.name, s.district, s.usme_number
				FROM smes s
				WHERE s.usme_number IS NOT NULL
				  AND s.usme_number LIKE 'MW-%'
				  AND s.deleted_at IS NULL
				ORDER BY s.id
			LOOP
				parts := string_to_array(rec.usme_number, '-');

				IF array_length(parts, 1) >= 6 THEN
					old_code := parts[3];

					SELECT m.new_code INTO new_code
					FROM district_code_mapping_preview m
					WHERE m.old_code = old_code
					LIMIT 1;

					IF new_code IS NOT NULL THEN
						new_usme_without_check := parts[1] || '-' || parts[2] || '-' || new_code || '-' || parts[4] || '-' || parts[5];
						new_check_digit := calculate_luhn_check_digit(new_usme_without_check);
						updated_usme := new_usme_without_check || '-' || new_check_digit;

						sme_id := rec.id;
						sme_name := rec.name;
						district := rec.district;
						old_usme := rec.usme_number;
						new_usme := updated_usme;
						old_district_code := old_code;
						new_district_code := new_code;
						RETURN NEXT;
					END IF;
				END IF;
			END LOOP;
		END;
		$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260112090000UpdateUsmeDistrictCodes) Down() error {
	_, err := facades.Orm().Query().Exec(`
		DROP FUNCTION IF EXISTS update_usme_district_codes();
		DROP FUNCTION IF EXISTS preview_usme_district_code_updates();
		DROP FUNCTION IF EXISTS calculate_luhn_check_digit(VARCHAR);
	`)
	return err
}
