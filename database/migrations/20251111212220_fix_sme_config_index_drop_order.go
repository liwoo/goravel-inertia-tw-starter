package migrations

type M20251111212220FixSmeConfigIndexDropOrder struct{}

// Signature The unique signature for the migration.
func (r *M20251111212220FixSmeConfigIndexDropOrder) Signature() string {
	return "20251111212220_fix_sme_config_index_drop_order"
}

// Up Run the migrations.
func (r *M20251111212220FixSmeConfigIndexDropOrder) Up() error {
	return nil
}

// Down Reverse the migrations.
func (r *M20251111212220FixSmeConfigIndexDropOrder) Down() error {
	// This migration was a fix for index drop order, nothing to reverse
	// Attempting to drop the index would cause errors since it may not exist
	return nil
}
