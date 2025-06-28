package seeders

import (
	"github.com/goravel/framework/facades"
	"players/app/models"
)

type BookSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *BookSeeder) Signature() string {
	return "BookSeeder"
}

// Run executes the seeder logic.
func (s *BookSeeder) Run() error {
	books := []models.Book{
		// Classic Literature
		{
			Title:       "To Kill a Mockingbird",
			Author:      "Harper Lee",
			ISBN:        "978-0-06-112008-4",
			Description: "A gripping, heart-wrenching, and wholly remarkable tale of coming-of-age in a South poisoned by virulent prejudice.",
			Price:       14.99,
			Status:      "AVAILABLE",
			PublishedAt: "1960-07-11",
		},
		{
			Title:       "1984",
			Author:      "George Orwell",
			ISBN:        "978-0-452-28423-4",
			Description: "A dystopian social science fiction novel and cautionary tale about the dangers of totalitarianism.",
			Price:       13.99,
			Status:      "BORROWED",
			PublishedAt: "1949-06-08",
		},
		{
			Title:       "Pride and Prejudice",
			Author:      "Jane Austen",
			ISBN:        "978-0-14-143951-8",
			Description: "A romantic novel of manners written by Jane Austen. It follows the character development of Elizabeth Bennet.",
			Price:       12.99,
			Status:      "AVAILABLE",
			PublishedAt: "1813-01-28",
		},
		{
			Title:       "The Great Gatsby",
			Author:      "F. Scott Fitzgerald",
			ISBN:        "978-0-7432-7356-5",
			Description: "A 1925 novel written by American author F. Scott Fitzgerald that follows a cast of characters living in West Egg.",
			Price:       15.99,
			Status:      "MAINTENANCE",
			PublishedAt: "1925-04-10",
		},
		{
			Title:       "Jane Eyre",
			Author:      "Charlotte Brontë",
			ISBN:        "978-0-14-144114-6",
			Description: "A bildungsroman which follows the experiences of its eponymous heroine.",
			Price:       11.99,
			Status:      "AVAILABLE",
			PublishedAt: "1847-10-16",
		},

		// Science Fiction
		{
			Title:       "Dune",
			Author:      "Frank Herbert",
			ISBN:        "978-0-441-17271-9",
			Description: "Set in the distant future amidst a feudal interstellar society in which various noble houses control planetary fiefs.",
			Price:       16.99,
			Status:      "AVAILABLE",
			PublishedAt: "1965-08-01",
		},
		{
			Title:       "Foundation",
			Author:      "Isaac Asimov",
			ISBN:        "978-0-553-29335-0",
			Description: "A cycle of five interrelated short stories, first published as a single book in 1951.",
			Price:       14.99,
			Status:      "BORROWED",
			PublishedAt: "1951-05-01",
		},
		{
			Title:       "Neuromancer",
			Author:      "William Gibson",
			ISBN:        "978-0-441-56956-9",
			Description: "A 1984 science fiction novel. It is one of the best-known works in the cyberpunk genre.",
			Price:       13.99,
			Status:      "AVAILABLE",
			PublishedAt: "1984-07-01",
		},
		{
			Title:       "The Hitchhiker's Guide to the Galaxy",
			Author:      "Douglas Adams",
			ISBN:        "978-0-345-39180-3",
			Description: "A comedy science fiction series created by Douglas Adams.",
			Price:       12.99,
			Status:      "AVAILABLE",
			PublishedAt: "1979-10-12",
		},
		{
			Title:       "Ender's Game",
			Author:      "Orson Scott Card",
			ISBN:        "978-0-812-55070-2",
			Description: "A 1985 military science fiction novel. Set at an unspecified date in Earth's future.",
			Price:       15.99,
			Status:      "BORROWED",
			PublishedAt: "1985-01-15",
		},

		// Fantasy
		{
			Title:       "The Lord of the Rings: The Fellowship of the Ring",
			Author:      "J.R.R. Tolkien",
			ISBN:        "978-0-547-92822-7",
			Description: "The first volume in The Lord of the Rings. It is preceded by The Hobbit.",
			Price:       18.99,
			Status:      "AVAILABLE",
			PublishedAt: "1954-07-29",
		},
		{
			Title:       "Harry Potter and the Philosopher's Stone",
			Author:      "J.K. Rowling",
			ISBN:        "978-0-439-70818-8",
			Description: "A fantasy novel written by British author J. K. Rowling. The first novel in the Harry Potter series.",
			Price:       17.99,
			Status:      "BORROWED",
			PublishedAt: "1997-06-26",
		},
		{
			Title:       "A Game of Thrones",
			Author:      "George R.R. Martin",
			ISBN:        "978-0-553-10354-0",
			Description: "The first novel in A Song of Ice and Fire, a series of fantasy novels by American author George R. R. Martin.",
			Price:       19.99,
			Status:      "AVAILABLE",
			PublishedAt: "1996-08-01",
		},
		{
			Title:       "The Name of the Wind",
			Author:      "Patrick Rothfuss",
			ISBN:        "978-0-7564-0474-1",
			Description: "A heroic fantasy novel written by American author Patrick Rothfuss. It is the first book in the ongoing trilogy The Kingkiller Chronicle.",
			Price:       16.99,
			Status:      "MAINTENANCE",
			PublishedAt: "2007-03-27",
		},
		{
			Title:       "The Way of Kings",
			Author:      "Brandon Sanderson",
			ISBN:        "978-0-7653-2635-5",
			Description: "An epic fantasy novel written by American author Brandon Sanderson and the first book in The Stormlight Archive series.",
			Price:       21.99,
			Status:      "AVAILABLE",
			PublishedAt: "2010-08-31",
		},

		// Mystery/Thriller
		{
			Title:       "The Girl with the Dragon Tattoo",
			Author:      "Stieg Larsson",
			ISBN:        "978-0-307-49892-6",
			Description: "A psychological thriller novel. It is the first book of the Millennium series.",
			Price:       15.99,
			Status:      "BORROWED",
			PublishedAt: "2005-08-01",
		},
		{
			Title:       "Gone Girl",
			Author:      "Gillian Flynn",
			ISBN:        "978-0-307-58836-4",
			Description: "A thriller novel. The story is told from the point of view of husband Nick Dunne and his wife Amy Dunne.",
			Price:       16.99,
			Status:      "AVAILABLE",
			PublishedAt: "2012-06-05",
		},
		{
			Title:       "The Da Vinci Code",
			Author:      "Dan Brown",
			ISBN:        "978-0-385-50420-1",
			Description: "A mystery thriller novel. It is the second novel to include the character Robert Langdon.",
			Price:       14.99,
			Status:      "AVAILABLE",
			PublishedAt: "2003-03-18",
		},
		{
			Title:       "And Then There Were None",
			Author:      "Agatha Christie",
			ISBN:        "978-0-06-207348-6",
			Description: "A mystery novel. It was first published in the United Kingdom by the Collins Crime Club.",
			Price:       13.99,
			Status:      "BORROWED",
			PublishedAt: "1939-11-06",
		},
		{
			Title:       "The Big Sleep",
			Author:      "Raymond Chandler",
			ISBN:        "978-0-394-75828-5",
			Description: "A hardboiled crime novel. It has been adapted for film twice, in 1946 and again in 1978.",
			Price:       12.99,
			Status:      "MAINTENANCE",
			PublishedAt: "1939-01-01",
		},

		// Non-Fiction
		{
			Title:       "Sapiens: A Brief History of Humankind",
			Author:      "Yuval Noah Harari",
			ISBN:        "978-0-06-231609-7",
			Description: "A book by Yuval Noah Harari, first published in Hebrew in Israel in 2011.",
			Price:       18.99,
			Status:      "AVAILABLE",
			PublishedAt: "2011-01-01",
		},
		{
			Title:       "Educated",
			Author:      "Tara Westover",
			ISBN:        "978-0-399-59050-4",
			Description: "A memoir by American historian and author Tara Westover.",
			Price:       17.99,
			Status:      "BORROWED",
			PublishedAt: "2018-02-20",
		},
		{
			Title:       "The Immortal Life of Henrietta Lacks",
			Author:      "Rebecca Skloot",
			ISBN:        "978-1-4000-5217-2",
			Description: "A non-fiction book by American author Rebecca Skloot.",
			Price:       16.99,
			Status:      "AVAILABLE",
			PublishedAt: "2010-02-02",
		},
		{
			Title:       "Thinking, Fast and Slow",
			Author:      "Daniel Kahneman",
			ISBN:        "978-0-374-53355-7",
			Description: "A 2011 book by psychologist Daniel Kahneman.",
			Price:       19.99,
			Status:      "AVAILABLE",
			PublishedAt: "2011-10-25",
		},
		{
			Title:       "The Power of Habit",
			Author:      "Charles Duhigg",
			ISBN:        "978-1-4000-6928-6",
			Description: "A book by Charles Duhigg, a New York Times reporter, published in February 2012.",
			Price:       15.99,
			Status:      "MAINTENANCE",
			PublishedAt: "2012-02-28",
		},

		// Contemporary Fiction
		{
			Title:       "The Kite Runner",
			Author:      "Khaled Hosseini",
			ISBN:        "978-1-59448-000-3",
			Description: "The debut novel by Afghan-American author Khaled Hosseini.",
			Price:       14.99,
			Status:      "AVAILABLE",
			PublishedAt: "2003-05-29",
		},
		{
			Title:       "Life of Pi",
			Author:      "Yann Martel",
			ISBN:        "978-0-15-100811-7",
			Description: "A Canadian philosophical novel by Yann Martel published in 2001.",
			Price:       13.99,
			Status:      "BORROWED",
			PublishedAt: "2001-09-11",
		},
		{
			Title:       "The Book Thief",
			Author:      "Markus Zusak",
			ISBN:        "978-0-375-83100-3",
			Description: "A 2005 historical novel by Australian author Markus Zusak.",
			Price:       15.99,
			Status:      "AVAILABLE",
			PublishedAt: "2005-03-14",
		},
		{
			Title:       "Where the Crawdads Sing",
			Author:      "Delia Owens",
			ISBN:        "978-0-735-21953-0",
			Description: "A 2018 novel by American zoologist Delia Owens.",
			Price:       16.99,
			Status:      "BORROWED",
			PublishedAt: "2018-08-14",
		},
		{
			Title:       "The Seven Husbands of Evelyn Hugo",
			Author:      "Taylor Jenkins Reid",
			ISBN:        "978-1-501-16134-8",
			Description: "A novel by American author Taylor Jenkins Reid and published in 2017.",
			Price:       14.99,
			Status:      "AVAILABLE",
			PublishedAt: "2017-06-13",
		},

		// Horror
		{
			Title:       "The Shining",
			Author:      "Stephen King",
			ISBN:        "978-0-307-74365-9",
			Description: "A horror novel by American author Stephen King.",
			Price:       15.99,
			Status:      "MAINTENANCE",
			PublishedAt: "1977-01-28",
		},
		{
			Title:       "Dracula",
			Author:      "Bram Stoker",
			ISBN:        "978-0-486-41109-7",
			Description: "An 1897 Gothic horror novel by Irish author Bram Stoker.",
			Price:       11.99,
			Status:      "AVAILABLE",
			PublishedAt: "1897-05-26",
		},
		{
			Title:       "Frankenstein",
			Author:      "Mary Shelley",
			ISBN:        "978-0-486-28211-4",
			Description: "An 1818 novel written by English author Mary Shelley.",
			Price:       10.99,
			Status:      "BORROWED",
			PublishedAt: "1818-01-01",
		},

		// Romance
		{
			Title:       "The Notebook",
			Author:      "Nicholas Sparks",
			ISBN:        "978-0-446-60523-4",
			Description: "A 1996 romantic novel by American novelist Nicholas Sparks.",
			Price:       13.99,
			Status:      "AVAILABLE",
			PublishedAt: "1996-10-01",
		},
		{
			Title:       "Me Before You",
			Author:      "Jojo Moyes",
			ISBN:        "978-0-14-312454-1",
			Description: "A romance novel written by Jojo Moyes.",
			Price:       14.99,
			Status:      "BORROWED",
			PublishedAt: "2012-01-05",
		},

		// Young Adult
		{
			Title:       "The Hunger Games",
			Author:      "Suzanne Collins",
			ISBN:        "978-0-439-02348-1",
			Description: "A 2008 dystopian novel by American writer Suzanne Collins.",
			Price:       12.99,
			Status:      "AVAILABLE",
			PublishedAt: "2008-09-14",
		},
		{
			Title:       "The Fault in Our Stars",
			Author:      "John Green",
			ISBN:        "978-0-525-47881-2",
			Description: "A novel by John Green. It is his fourth solo novel, and sixth novel overall.",
			Price:       13.99,
			Status:      "BORROWED",
			PublishedAt: "2012-01-10",
		},
		{
			Title:       "Divergent",
			Author:      "Veronica Roth",
			ISBN:        "978-0-06-202402-2",
			Description: "A novel in the Divergent trilogy by Veronica Roth.",
			Price:       14.99,
			Status:      "MAINTENANCE",
			PublishedAt: "2011-04-25",
		},

		// Historical Fiction
		{
			Title:       "All Quiet on the Western Front",
			Author:      "Erich Maria Remarque",
			ISBN:        "978-0-449-21394-8",
			Description: "A novel by Erich Maria Remarque, a German veteran of World War I.",
			Price:       12.99,
			Status:      "AVAILABLE",
			PublishedAt: "1929-01-29",
		},
		{
			Title:       "The Pillars of the Earth",
			Author:      "Ken Follett",
			ISBN:        "978-0-451-16689-5",
			Description: "A historical novel by Welsh author Ken Follett published in 1989.",
			Price:       17.99,
			Status:      "BORROWED",
			PublishedAt: "1989-01-01",
		},
		{
			Title:       "The Help",
			Author:      "Kathryn Stockett",
			ISBN:        "978-0-399-15534-5",
			Description: "A 2009 novel by American author Kathryn Stockett.",
			Price:       15.99,
			Status:      "AVAILABLE",
			PublishedAt: "2009-02-10",
		},

		// Biography
		{
			Title:       "Steve Jobs",
			Author:      "Walter Isaacson",
			ISBN:        "978-1-451-64853-9",
			Description: "An authorized biography of Steve Jobs, the co-founder and longtime chief executive officer of Apple Inc.",
			Price:       19.99,
			Status:      "AVAILABLE",
			PublishedAt: "2011-10-24",
		},
		{
			Title:       "Long Walk to Freedom",
			Author:      "Nelson Mandela",
			ISBN:        "978-0-316-54585-6",
			Description: "An autobiographical work written by South African President Nelson Mandela.",
			Price:       18.99,
			Status:      "MAINTENANCE",
			PublishedAt: "1994-10-01",
		},

		// Philosophy
		{
			Title:       "Meditations",
			Author:      "Marcus Aurelius",
			ISBN:        "978-0-486-29823-2",
			Description: "A series of personal writings by Marcus Aurelius, Roman Emperor from 161 to 180 AD.",
			Price:       9.99,
			Status:      "AVAILABLE",
			PublishedAt: "0171-01-01",
		},
		{
			Title:       "The Art of War",
			Author:      "Sun Tzu",
			ISBN:        "978-1-59030-963-7",
			Description: "An ancient Chinese military treatise dating from the Late Spring and Autumn Period.",
			Price:       8.99,
			Status:      "BORROWED",
			PublishedAt: "0500-01-01",
		},

		// Business
		{
			Title:       "Good to Great",
			Author:      "Jim Collins",
			ISBN:        "978-0-06-662099-2",
			Description: "A management book by Jim C. Collins that describes how companies transition from being good companies to great companies.",
			Price:       17.99,
			Status:      "AVAILABLE",
			PublishedAt: "2001-10-16",
		},
		{
			Title:       "The Lean Startup",
			Author:      "Eric Ries",
			ISBN:        "978-0-307-88789-4",
			Description: "A book by Eric Ries describing his proposed lean startup strategy for startup companies.",
			Price:       16.99,
			Status:      "BORROWED",
			PublishedAt: "2011-09-13",
		},

		// Technology
		{
			Title:       "Clean Code",
			Author:      "Robert C. Martin",
			ISBN:        "978-0-13-235088-4",
			Description: "A handbook of agile software craftsmanship by Robert C. Martin.",
			Price:       24.99,
			Status:      "AVAILABLE",
			PublishedAt: "2008-08-01",
		},
		{
			Title:       "The Pragmatic Programmer",
			Author:      "David Thomas",
			ISBN:        "978-0-201-61622-4",
			Description: "A book about computer programming and software engineering, written by David Thomas and Andrew Hunt.",
			Price:       23.99,
			Status:      "MAINTENANCE",
			PublishedAt: "1999-10-30",
		},
	}

	// Insert all books into database
	for _, book := range books {
		if err := facades.Orm().Query().Create(&book); err != nil {
			return err
		}
	}

	return nil
}