package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251020113904AddCascadeDeleteToSmeRelationships struct{}

// Signature The unique signature for the migration.
func (r *M20251020113904AddCascadeDeleteToSmeRelationships) Signature() string {
	return "20251020113904_add_cascade_delete_to_sme_relationships"
}

// Up Run the migrations.
func (r *M20251020113904AddCascadeDeleteToSmeRelationships) Up() error {
	// Drop and recreate foreign keys with CASCADE DELETE for primary_business_owner
	if err := facades.Schema().Table("primary_business_owner", func(table schema.Blueprint) {
		table.DropForeign("primary_business_owner_sme_id_foreign")
	}); err != nil {
		return err
	}

	if err := facades.Schema().Table("primary_business_owner", func(table schema.Blueprint) {
		table.Foreign("sme_id").References("id").On("smes").CascadeOnDelete()
	}); err != nil {
		return err
	}

	// Drop and recreate foreign keys with CASCADE DELETE for additional_business_members
	if err := facades.Schema().Table("additional_business_members", func(table schema.Blueprint) {
		table.DropForeign("additional_business_members_sme_id_foreign")
	}); err != nil {
		return err
	}

	if err := facades.Schema().Table("additional_business_members", func(table schema.Blueprint) {
		table.Foreign("sme_id").References("id").On("smes").CascadeOnDelete()
	}); err != nil {
		return err
	}

	// Drop and recreate foreign keys with CASCADE DELETE for business_formalisation
	if err := facades.Schema().Table("business_formalisation", func(table schema.Blueprint) {
		table.DropForeign("business_formalisation_sme_id_foreign")
	}); err != nil {
		return err
	}

	if err := facades.Schema().Table("business_formalisation", func(table schema.Blueprint) {
		table.Foreign("sme_id").References("id").On("smes").CascadeOnDelete()
	}); err != nil {
		return err
	}

	// Drop and recreate foreign keys with CASCADE DELETE for business_employee_summary
	if err := facades.Schema().Table("business_employee_summary", func(table schema.Blueprint) {
		table.DropForeign("business_employee_summary_sme_id_foreign")
	}); err != nil {
		return err
	}

	if err := facades.Schema().Table("business_employee_summary", func(table schema.Blueprint) {
		table.Foreign("sme_id").References("id").On("smes").CascadeOnDelete()
	}); err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251020113904AddCascadeDeleteToSmeRelationships) Down() error {
	// Revert back to foreign keys without CASCADE DELETE

	// Drop CASCADE DELETE foreign keys and recreate without CASCADE for primary_business_owner
	if err := facades.Schema().Table("primary_business_owner", func(table schema.Blueprint) {
		table.DropForeign("primary_business_owner_sme_id_foreign")
	}); err != nil {
		return err
	}

	if err := facades.Schema().Table("primary_business_owner", func(table schema.Blueprint) {
		table.Foreign("sme_id").References("id").On("smes")
	}); err != nil {
		return err
	}

	// Drop CASCADE DELETE foreign keys and recreate without CASCADE for additional_business_members
	if err := facades.Schema().Table("additional_business_members", func(table schema.Blueprint) {
		table.DropForeign("additional_business_members_sme_id_foreign")
	}); err != nil {
		return err
	}

	if err := facades.Schema().Table("additional_business_members", func(table schema.Blueprint) {
		table.Foreign("sme_id").References("id").On("smes")
	}); err != nil {
		return err
	}

	// Drop CASCADE DELETE foreign keys and recreate without CASCADE for business_formalisation
	if err := facades.Schema().Table("business_formalisation", func(table schema.Blueprint) {
		table.DropForeign("business_formalisation_sme_id_foreign")
	}); err != nil {
		return err
	}

	if err := facades.Schema().Table("business_formalisation", func(table schema.Blueprint) {
		table.Foreign("sme_id").References("id").On("smes")
	}); err != nil {
		return err
	}

	// Drop CASCADE DELETE foreign keys and recreate without CASCADE for business_employee_summary
	if err := facades.Schema().Table("business_employee_summary", func(table schema.Blueprint) {
		table.DropForeign("business_employee_summary_sme_id_foreign")
	}); err != nil {
		return err
	}

	if err := facades.Schema().Table("business_employee_summary", func(table schema.Blueprint) {
		table.Foreign("sme_id").References("id").On("smes")
	}); err != nil {
		return err
	}

	return nil
}
