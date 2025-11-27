package migrations

import (
	"github.com/goravel/framework/facades"
)

type M20251020113904AddCascadeDeleteToSmeRelationships struct{}

// Signature The unique signature for the migration.
func (r *M20251020113904AddCascadeDeleteToSmeRelationships) Signature() string {
	return "20251020113904_add_cascade_delete_to_sme_relationships"
}

// Up Run the migrations.
func (r *M20251020113904AddCascadeDeleteToSmeRelationships) Up() error {
	// Use raw SQL for PostgreSQL compatibility
	// Drop existing foreign keys and recreate with CASCADE DELETE

	// primary_business_owner
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE primary_business_owner DROP CONSTRAINT IF EXISTS primary_business_owner_sme_id_foreign",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE primary_business_owner ADD CONSTRAINT primary_business_owner_sme_id_foreign FOREIGN KEY (sme_id) REFERENCES smes(id) ON DELETE CASCADE",
	); err != nil {
		return err
	}

	// additional_business_members
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE additional_business_members DROP CONSTRAINT IF EXISTS additional_business_members_sme_id_foreign",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE additional_business_members ADD CONSTRAINT additional_business_members_sme_id_foreign FOREIGN KEY (sme_id) REFERENCES smes(id) ON DELETE CASCADE",
	); err != nil {
		return err
	}

	// business_formalisation
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_formalisation DROP CONSTRAINT IF EXISTS business_formalisation_sme_id_foreign",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_formalisation ADD CONSTRAINT business_formalisation_sme_id_foreign FOREIGN KEY (sme_id) REFERENCES smes(id) ON DELETE CASCADE",
	); err != nil {
		return err
	}

	// business_employee_summary
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_employee_summary DROP CONSTRAINT IF EXISTS business_employee_summary_sme_id_foreign",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_employee_summary ADD CONSTRAINT business_employee_summary_sme_id_foreign FOREIGN KEY (sme_id) REFERENCES smes(id) ON DELETE CASCADE",
	); err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251020113904AddCascadeDeleteToSmeRelationships) Down() error {
	// Revert back to foreign keys without CASCADE DELETE using raw SQL

	// primary_business_owner
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE primary_business_owner DROP CONSTRAINT IF EXISTS primary_business_owner_sme_id_foreign",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE primary_business_owner ADD CONSTRAINT primary_business_owner_sme_id_foreign FOREIGN KEY (sme_id) REFERENCES smes(id)",
	); err != nil {
		return err
	}

	// additional_business_members
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE additional_business_members DROP CONSTRAINT IF EXISTS additional_business_members_sme_id_foreign",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE additional_business_members ADD CONSTRAINT additional_business_members_sme_id_foreign FOREIGN KEY (sme_id) REFERENCES smes(id)",
	); err != nil {
		return err
	}

	// business_formalisation
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_formalisation DROP CONSTRAINT IF EXISTS business_formalisation_sme_id_foreign",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_formalisation ADD CONSTRAINT business_formalisation_sme_id_foreign FOREIGN KEY (sme_id) REFERENCES smes(id)",
	); err != nil {
		return err
	}

	// business_employee_summary
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_employee_summary DROP CONSTRAINT IF EXISTS business_employee_summary_sme_id_foreign",
	); err != nil {
		return err
	}
	if _, err := facades.Orm().Query().Exec(
		"ALTER TABLE business_employee_summary ADD CONSTRAINT business_employee_summary_sme_id_foreign FOREIGN KEY (sme_id) REFERENCES smes(id)",
	); err != nil {
		return err
	}

	return nil
}
