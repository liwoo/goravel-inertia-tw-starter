package migrations

import (
	"time"
	"github.com/goravel/framework/facades"
	"players/app/models"
)

type RestoreBookPublishedDates struct {
}

// Signature The name and signature of the console command.
func (receiver *RestoreBookPublishedDates) Signature() string {
	return "20250724000001_restore_book_published_dates"
}

// Description The console command description.
func (receiver *RestoreBookPublishedDates) Description() string {
	return "Restore published dates for seeded books"
}

// Up Run the migrations.
func (receiver *RestoreBookPublishedDates) Up() error {
	// Map of ISBN to published date
	bookDates := map[string]string{
		// Classic Literature
		"978-0-06-112008-4": "1960-07-11",
		"978-0-452-28423-4": "1949-06-08",
		"978-0-14-143951-8": "1813-01-28",
		"978-0-7432-7356-5": "1925-04-10",
		"978-0-14-144114-6": "1847-10-16",
		// Science Fiction
		"978-0-441-17271-9": "1965-08-01",
		"978-0-553-29335-0": "1951-05-01",
		"978-0-441-56956-9": "1984-07-01",
		"978-0-345-39180-3": "1979-10-12",
		"978-0-812-55070-2": "1985-01-15",
		// Fantasy
		"978-0-547-92822-7": "1954-07-29",
		"978-0-439-70818-8": "1997-06-26",
		"978-0-553-10354-0": "1996-08-01",
		"978-0-7564-0474-1": "2007-03-27",
		"978-0-7653-2635-5": "2010-08-31",
		// Mystery/Thriller
		"978-0-307-49892-6": "2005-08-01",
		"978-0-307-58836-4": "2012-06-05",
		"978-0-385-50420-1": "2003-03-18",
		"978-0-06-207348-6": "1939-11-06",
		"978-0-394-75828-5": "1939-01-01",
		// Non-Fiction
		"978-0-06-231609-7": "2011-01-01",
		"978-0-399-59050-4": "2018-02-20",
		"978-1-4000-5217-2": "2010-02-02",
		"978-0-374-53355-7": "2011-10-25",
		"978-1-4000-6928-6": "2012-02-28",
		// Contemporary Fiction
		"978-1-59448-000-3": "2003-05-29",
		"978-0-15-100811-7": "2001-09-11",
		"978-0-375-83100-3": "2005-03-14",
		"978-0-735-21953-0": "2018-08-14",
		"978-1-501-16134-8": "2017-06-13",
		// Horror
		"978-0-307-74365-9": "1977-01-28",
		"978-0-486-41109-7": "1897-05-26",
		"978-0-486-28211-4": "1818-01-01",
		// Romance
		"978-0-446-60523-4": "1996-10-01",
		"978-0-14-312454-1": "2012-01-05",
		// Young Adult
		"978-0-439-02348-1": "2008-09-14",
		"978-0-525-47881-2": "2012-01-10",
		"978-0-06-202402-2": "2011-04-25",
		// Historical Fiction
		"978-0-449-21394-8": "1929-01-29",
		"978-0-451-16689-5": "1989-01-01",
		"978-0-399-15534-5": "2009-02-10",
		// Biography
		"978-1-451-64853-9": "2011-10-24",
		"978-0-316-54585-6": "1994-10-01",
		// Philosophy
		"978-0-486-29823-2": "0171-01-01",
		"978-1-59030-963-7": "0500-01-01",
		// Business
		"978-0-06-662099-2": "2001-10-16",
		"978-0-307-88789-4": "2011-09-13",
		// Technology
		"978-0-13-235088-4": "2008-08-01",
		"978-0-201-61622-4": "1999-10-30",
	}

	// Update each book
	for isbn, dateStr := range bookDates {
		parsedTime, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		
		facades.Orm().Query().
			Model(&models.Book{}).
			Where("isbn = ?", isbn).
			Update("published_at", parsedTime)
	}

	return nil
}

// Down Reverse the migrations.
func (receiver *RestoreBookPublishedDates) Down() error {
	// This is a data migration, no need to reverse
	return nil
}