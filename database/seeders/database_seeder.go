package seeders

type DatabaseSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *DatabaseSeeder) Signature() string {
	return "DatabaseSeeder"
}

// Run executes the seeder logic.
func (s *DatabaseSeeder) Run() error {
	// Run the RBAC seeder first (roles and permissions)
	rbacSeeder := &RBACSeeder{}
	if err := rbacSeeder.Run(); err != nil {
		return err
	}

	// Run the book seeder
	bookSeeder := &BookSeeder{}
	if err := bookSeeder.Run(); err != nil {
		return err
	}

	return nil
}
