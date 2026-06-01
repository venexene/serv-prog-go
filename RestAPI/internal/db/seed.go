package db

import (
	"gravity-game-store/internal/entity"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(d *gorm.DB, log *logrus.Logger) error {
	var n int64
	d.Model(&entity.Author{}).Count(&n)
	if n > 0 {
		log.Info("DB already seeded, skip")
		return nil
	}

	log.Info("Seeding...")

	steps := []func(*gorm.DB) error{
		seedUsers,
		seedAuthors,
		seedGames,
		seedGameAuthors,
		seedOrderStatuses,
		seedShippingMethods,
		seedCustomers,
		seedAddresses,
		seedCustomerAddresses,
		seedCustOrders,
		seedOrderLines,
	}

	for _, fn := range steps {
		if err := fn(d); err != nil {
			return err
		}
	}

	log.Info("Seed done")
	return nil
}

func seedUsers(d *gorm.DB) error {
	users := []entity.User{
		{Username: "admin", Password: hash("admin123"), Role: "admin"},
		{Username: "user", Password: hash("user123"), Role: "user"},
	}
	return d.Create(&users).Error
}

func seedAuthors(d *gorm.DB) error {
	authors := []entity.Author{
		{Name: "FromSoftware"},
		{Name: "Nintendo EAD"},
		{Name: "Valve Corporation"},
		{Name: "CD Projekt Red"},
		{Name: "Rockstar Games"},
		{Name: "Bethesda Game Studios"},
		{Name: "Blizzard Entertainment"},
		{Name: "BioWare"},
		{Name: "Naughty Dog"},
		{Name: "Mojang Studios"},
		{Name: "Epic Games"},
		{Name: "Square Enix"},
		{Name: "Capcom"},
		{Name: "Konami"},
		{Name: "Supergiant Games"},
		{Name: "id Software"},
		{Name: "Larian Studios"},
		{Name: "ConcernedApe"},
		{Name: "Team Cherry"},
		{Name: "George R.R. Martin"},
	}
	return d.Create(&authors).Error
}

func seedGames(d *gorm.DB) error {
	games := []entity.Game{
		{Title: "Demon's Souls", Genre: "Action RPG", Platform: "PS3", PublisherID: 1, ReleaseDate: "2009-02-05", NumPlayers: 1},
		{Title: "Dark Souls", Genre: "Action RPG", Platform: "PC", PublisherID: 2, ReleaseDate: "2011-09-22", NumPlayers: 4},
		{Title: "Elden Ring", Genre: "Action RPG", Platform: "PC", PublisherID: 2, ReleaseDate: "2022-02-25", NumPlayers: 4},
		{Title: "Super Mario 64", Genre: "Platformer", Platform: "N64", PublisherID: 3, ReleaseDate: "1996-06-23", NumPlayers: 1},
		{Title: "The Legend of Zelda: Ocarina of Time", Genre: "Action-Adventure", Platform: "N64", PublisherID: 3, ReleaseDate: "1998-11-21", NumPlayers: 1},
		{Title: "The Legend of Zelda: Breath of the Wild", Genre: "Action-Adventure", Platform: "Switch", PublisherID: 3, ReleaseDate: "2017-03-03", NumPlayers: 1},
		{Title: "Half-Life 2", Genre: "FPS", Platform: "PC", PublisherID: 4, ReleaseDate: "2004-11-16", NumPlayers: 32},
		{Title: "Portal 2", Genre: "Puzzle", Platform: "PC", PublisherID: 4, ReleaseDate: "2011-04-19", NumPlayers: 2},
		{Title: "Dota 2", Genre: "MOBA", Platform: "PC", PublisherID: 4, ReleaseDate: "2013-07-09", NumPlayers: 10},
		{Title: "The Witcher 3: Wild Hunt", Genre: "Action RPG", Platform: "PC", PublisherID: 5, ReleaseDate: "2015-05-19", NumPlayers: 1},
		{Title: "Cyberpunk 2077", Genre: "Action RPG", Platform: "PC", PublisherID: 5, ReleaseDate: "2020-12-10", NumPlayers: 1},
		{Title: "The Witcher 2: Assassins of Kings", Genre: "Action RPG", Platform: "PC", PublisherID: 5, ReleaseDate: "2011-05-17", NumPlayers: 1},
		{Title: "Grand Theft Auto V", Genre: "Action-Adventure", Platform: "PC", PublisherID: 6, ReleaseDate: "2013-09-17", NumPlayers: 30},
		{Title: "Red Dead Redemption 2", Genre: "Action-Adventure", Platform: "PC", PublisherID: 6, ReleaseDate: "2018-10-26", NumPlayers: 32},
		{Title: "Grand Theft Auto: San Andreas", Genre: "Action-Adventure", Platform: "PS2", PublisherID: 6, ReleaseDate: "2004-10-26", NumPlayers: 1},
		{Title: "The Elder Scrolls V: Skyrim", Genre: "Action RPG", Platform: "PC", PublisherID: 7, ReleaseDate: "2011-11-11", NumPlayers: 1},
		{Title: "Fallout 4", Genre: "Action RPG", Platform: "PC", PublisherID: 7, ReleaseDate: "2015-11-10", NumPlayers: 1},
		{Title: "Doom (2016)", Genre: "FPS", Platform: "PC", PublisherID: 7, ReleaseDate: "2016-05-13", NumPlayers: 16},
		{Title: "World of Warcraft", Genre: "MMORPG", Platform: "PC", PublisherID: 8, ReleaseDate: "2004-11-23", NumPlayers: 40},
		{Title: "StarCraft II", Genre: "RTS", Platform: "PC", PublisherID: 8, ReleaseDate: "2010-07-27", NumPlayers: 8},
		{Title: "Diablo III", Genre: "Action RPG", Platform: "PC", PublisherID: 8, ReleaseDate: "2012-05-15", NumPlayers: 4},
		{Title: "Mass Effect 2", Genre: "Action RPG", Platform: "PC", PublisherID: 9, ReleaseDate: "2010-01-26", NumPlayers: 1},
		{Title: "Dragon Age: Origins", Genre: "RPG", Platform: "PC", PublisherID: 9, ReleaseDate: "2009-11-03", NumPlayers: 1},
		{Title: "Star Wars: Knights of the Old Republic", Genre: "RPG", Platform: "PC", PublisherID: 9, ReleaseDate: "2003-07-15", NumPlayers: 1},
		{Title: "The Last of Us", Genre: "Action-Adventure", Platform: "PS3", PublisherID: 10, ReleaseDate: "2013-06-14", NumPlayers: 8},
		{Title: "Uncharted 4: A Thief's End", Genre: "Action-Adventure", Platform: "PS4", PublisherID: 10, ReleaseDate: "2016-05-10", NumPlayers: 10},
		{Title: "The Last of Us Part II", Genre: "Action-Adventure", Platform: "PS4", PublisherID: 10, ReleaseDate: "2020-06-19", NumPlayers: 1},
		{Title: "Minecraft", Genre: "Sandbox", Platform: "PC", PublisherID: 11, ReleaseDate: "2011-11-18", NumPlayers: 30},
		{Title: "Minecraft Dungeons", Genre: "Dungeon Crawler", Platform: "PC", PublisherID: 11, ReleaseDate: "2020-05-26", NumPlayers: 4},
		{Title: "Fortnite", Genre: "Battle Royale", Platform: "PC", PublisherID: 12, ReleaseDate: "2017-07-25", NumPlayers: 100},
		{Title: "Final Fantasy VII", Genre: "JRPG", Platform: "PS1", PublisherID: 13, ReleaseDate: "1997-01-31", NumPlayers: 1},
		{Title: "Chrono Trigger", Genre: "JRPG", Platform: "SNES", PublisherID: 13, ReleaseDate: "1995-03-11", NumPlayers: 1},
		{Title: "Kingdom Hearts", Genre: "Action RPG", Platform: "PS2", PublisherID: 13, ReleaseDate: "2002-03-28", NumPlayers: 1},
		{Title: "Resident Evil 4", Genre: "Survival Horror", Platform: "GameCube", PublisherID: 14, ReleaseDate: "2005-01-11", NumPlayers: 1},
		{Title: "Monster Hunter: World", Genre: "Action RPG", Platform: "PC", PublisherID: 14, ReleaseDate: "2018-01-26", NumPlayers: 4},
		{Title: "Devil May Cry 5", Genre: "Action", Platform: "PC", PublisherID: 14, ReleaseDate: "2019-03-08", NumPlayers: 3},
		{Title: "Metal Gear Solid", Genre: "Stealth", Platform: "PS1", PublisherID: 15, ReleaseDate: "1998-09-03", NumPlayers: 1},
		{Title: "Metal Gear Solid V: The Phantom Pain", Genre: "Stealth", Platform: "PC", PublisherID: 15, ReleaseDate: "2015-09-01", NumPlayers: 16},
		{Title: "Silent Hill 2", Genre: "Survival Horror", Platform: "PS2", PublisherID: 15, ReleaseDate: "2001-09-24", NumPlayers: 1},
		{Title: "Bastion", Genre: "Action RPG", Platform: "PC", PublisherID: 16, ReleaseDate: "2011-07-20", NumPlayers: 1},
		{Title: "Hades", Genre: "Roguelike", Platform: "PC", PublisherID: 16, ReleaseDate: "2020-09-17", NumPlayers: 1},
		{Title: "Transistor", Genre: "Action RPG", Platform: "PC", PublisherID: 16, ReleaseDate: "2014-05-20", NumPlayers: 1},
		{Title: "Doom Eternal", Genre: "FPS", Platform: "PC", PublisherID: 7, ReleaseDate: "2020-03-20", NumPlayers: 3},
		{Title: "Quake III Arena", Genre: "FPS", Platform: "PC", PublisherID: 17, ReleaseDate: "1999-12-02", NumPlayers: 32},
		{Title: "Baldur's Gate 3", Genre: "RPG", Platform: "PC", PublisherID: 18, ReleaseDate: "2023-08-03", NumPlayers: 4},
		{Title: "Divinity: Original Sin 2", Genre: "RPG", Platform: "PC", PublisherID: 18, ReleaseDate: "2017-09-14", NumPlayers: 4},
		{Title: "Stardew Valley", Genre: "Simulation", Platform: "PC", PublisherID: 19, ReleaseDate: "2016-02-26", NumPlayers: 4},
		{Title: "Hollow Knight", Genre: "Metroidvania", Platform: "PC", PublisherID: 20, ReleaseDate: "2017-02-24", NumPlayers: 1},
		{Title: "Hollow Knight: Silksong", Genre: "Metroidvania", Platform: "PC", PublisherID: 20, ReleaseDate: "2025-12-31", NumPlayers: 1},
	}
	return d.Create(&games).Error
}

func seedGameAuthors(d *gorm.DB) error {
	rels := []entity.GameAuthor{
		{GameID: 1, AuthorID: 1}, {GameID: 2, AuthorID: 1}, {GameID: 3, AuthorID: 1},
		{GameID: 4, AuthorID: 2}, {GameID: 5, AuthorID: 2}, {GameID: 6, AuthorID: 2},
		{GameID: 7, AuthorID: 3}, {GameID: 8, AuthorID: 3}, {GameID: 9, AuthorID: 3},
		{GameID: 10, AuthorID: 4}, {GameID: 11, AuthorID: 4}, {GameID: 12, AuthorID: 4},
		{GameID: 13, AuthorID: 5}, {GameID: 14, AuthorID: 5}, {GameID: 15, AuthorID: 5},
		{GameID: 16, AuthorID: 6}, {GameID: 17, AuthorID: 6}, {GameID: 18, AuthorID: 6},
		{GameID: 19, AuthorID: 7}, {GameID: 20, AuthorID: 7}, {GameID: 21, AuthorID: 7},
		{GameID: 22, AuthorID: 8}, {GameID: 23, AuthorID: 8}, {GameID: 24, AuthorID: 8},
		{GameID: 25, AuthorID: 9}, {GameID: 26, AuthorID: 9}, {GameID: 27, AuthorID: 9},
		{GameID: 28, AuthorID: 10}, {GameID: 29, AuthorID: 10},
		{GameID: 30, AuthorID: 11},
		{GameID: 31, AuthorID: 12}, {GameID: 32, AuthorID: 12}, {GameID: 33, AuthorID: 12},
		{GameID: 34, AuthorID: 13}, {GameID: 35, AuthorID: 13}, {GameID: 36, AuthorID: 13},
		{GameID: 37, AuthorID: 14}, {GameID: 38, AuthorID: 14}, {GameID: 39, AuthorID: 14},
		{GameID: 40, AuthorID: 15}, {GameID: 41, AuthorID: 15}, {GameID: 42, AuthorID: 15},
		{GameID: 43, AuthorID: 16}, {GameID: 44, AuthorID: 16},
		{GameID: 45, AuthorID: 17}, {GameID: 46, AuthorID: 17},
		{GameID: 47, AuthorID: 18},
		{GameID: 48, AuthorID: 19}, {GameID: 49, AuthorID: 19},
		{GameID: 3, AuthorID: 20},
		{GameID: 43, AuthorID: 6},
	}
	return d.Create(&rels).Error
}

func seedOrderStatuses(d *gorm.DB) error {
	return d.Create(&[]entity.OrderStatus{
		{StatusID: 1, StatusValue: "Order Received"},
		{StatusID: 2, StatusValue: "Pending Delivery"},
		{StatusID: 3, StatusValue: "Delivered"},
		{StatusID: 4, StatusValue: "Cancelled"},
	}).Error
}

func seedShippingMethods(d *gorm.DB) error {
	return d.Create(&[]entity.ShippingMethod{
		{MethodID: 1, MethodName: "Standard", Cost: 5.99},
		{MethodID: 2, MethodName: "Express", Cost: 12.99},
		{MethodID: 3, MethodName: "Overnight", Cost: 24.99},
		{MethodID: 4, MethodName: "Pickup", Cost: 0.00},
	}).Error
}

func seedCustomers(d *gorm.DB) error {
	cust := []entity.Customer{
		{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob.johnson@example.com"},
		{FirstName: "Alice", LastName: "Williams", Email: "alice.williams@example.com"},
		{FirstName: "Charlie", LastName: "Brown", Email: "charlie.brown@example.com"},
		{FirstName: "Diana", LastName: "Miller", Email: "diana.miller@example.com"},
		{FirstName: "Edward", LastName: "Davis", Email: "edward.davis@example.com"},
		{FirstName: "Fiona", LastName: "Garcia", Email: "fiona.garcia@example.com"},
		{FirstName: "George", LastName: "Rodriguez", Email: "george.rodriguez@example.com"},
		{FirstName: "Hannah", LastName: "Wilson", Email: "hannah.wilson@example.com"},
		{FirstName: "Ivan", LastName: "Martinez", Email: "ivan.martinez@example.com"},
		{FirstName: "Julia", LastName: "Anderson", Email: "julia.anderson@example.com"},
		{FirstName: "Kevin", LastName: "Taylor", Email: "kevin.taylor@example.com"},
		{FirstName: "Laura", LastName: "Thomas", Email: "laura.thomas@example.com"},
		{FirstName: "Mike", LastName: "Hernandez", Email: "mike.hernandez@example.com"},
		{FirstName: "Nina", LastName: "Moore", Email: "nina.moore@example.com"},
		{FirstName: "Oscar", LastName: "Jackson", Email: "oscar.jackson@example.com"},
		{FirstName: "Paula", LastName: "Martin", Email: "paula.martin@example.com"},
		{FirstName: "Quinn", LastName: "Lee", Email: "quinn.lee@example.com"},
		{FirstName: "Rachel", LastName: "Perez", Email: "rachel.perez@example.com"},
	}
	return d.Create(&cust).Error
}

func seedAddresses(d *gorm.DB) error {
	return d.Create(&[]entity.Address{
		{StreetNumber: "123", StreetName: "Main St", City: "New York", CountryID: 1},
		{StreetNumber: "456", StreetName: "Oak Ave", City: "Los Angeles", CountryID: 1},
		{StreetNumber: "789", StreetName: "Pine Rd", City: "Chicago", CountryID: 1},
		{StreetNumber: "321", StreetName: "Elm St", City: "Houston", CountryID: 1},
		{StreetNumber: "654", StreetName: "Maple Dr", City: "Phoenix", CountryID: 1},
		{StreetNumber: "10", StreetName: "Downing Street", City: "London", CountryID: 2},
		{StreetNumber: "22", StreetName: "Baker Street", City: "London", CountryID: 2},
		{StreetNumber: "5", StreetName: "Avenue des Champs-Élysées", City: "Paris", CountryID: 3},
		{StreetNumber: "100", StreetName: "Unter den Linden", City: "Berlin", CountryID: 4},
		{StreetNumber: "15", StreetName: "Nevsky Prospekt", City: "Saint Petersburg", CountryID: 5},
		{StreetNumber: "8", StreetName: "Ginza", City: "Tokyo", CountryID: 6},
		{StreetNumber: "42", StreetName: "Rambla de Catalunya", City: "Barcelona", CountryID: 7},
		{StreetNumber: "77", StreetName: "Via del Corso", City: "Rome", CountryID: 8},
		{StreetNumber: "200", StreetName: "Yonge Street", City: "Toronto", CountryID: 9},
		{StreetNumber: "33", StreetName: "George Street", City: "Sydney", CountryID: 10},
		{StreetNumber: "55", StreetName: "Orchard Road", City: "Singapore", CountryID: 11},
		{StreetNumber: "144", StreetName: "Ocean Drive", City: "Miami", CountryID: 1},
		{StreetNumber: "900", StreetName: "Sunset Blvd", City: "Los Angeles", CountryID: 1},
		{StreetNumber: "350", StreetName: "Fifth Avenue", City: "New York", CountryID: 1},
		{StreetNumber: "1", StreetName: "Infinite Loop", City: "Cupertino", CountryID: 1},
	}).Error
}

func seedCustomerAddresses(d *gorm.DB) error {
	return d.Create(&[]entity.CustomerAddress{
		{CustomerID: 1, AddressID: 1, StatusID: 1},
		{CustomerID: 2, AddressID: 2, StatusID: 1},
		{CustomerID: 3, AddressID: 3, StatusID: 1},
		{CustomerID: 4, AddressID: 4, StatusID: 1},
		{CustomerID: 5, AddressID: 5, StatusID: 1},
		{CustomerID: 6, AddressID: 6, StatusID: 1},
		{CustomerID: 7, AddressID: 7, StatusID: 1},
		{CustomerID: 8, AddressID: 8, StatusID: 1},
		{CustomerID: 9, AddressID: 9, StatusID: 1},
		{CustomerID: 10, AddressID: 10, StatusID: 1},
		{CustomerID: 11, AddressID: 11, StatusID: 1},
		{CustomerID: 12, AddressID: 12, StatusID: 1},
		{CustomerID: 13, AddressID: 13, StatusID: 1},
		{CustomerID: 14, AddressID: 14, StatusID: 1},
		{CustomerID: 15, AddressID: 15, StatusID: 1},
		{CustomerID: 16, AddressID: 16, StatusID: 1},
		{CustomerID: 17, AddressID: 17, StatusID: 1},
		{CustomerID: 18, AddressID: 18, StatusID: 1},
		{CustomerID: 19, AddressID: 19, StatusID: 1},
		{CustomerID: 20, AddressID: 20, StatusID: 1},
		{CustomerID: 1, AddressID: 17, StatusID: 2},
		{CustomerID: 5, AddressID: 18, StatusID: 2},
		{CustomerID: 10, AddressID: 19, StatusID: 2},
		{CustomerID: 15, AddressID: 20, StatusID: 2},
	}).Error
}

func seedCustOrders(d *gorm.DB) error {
	return d.Create(&[]entity.CustOrder{
		{CustomerID: 1, ShippingMethodID: 1, DestAddressID: 1},
		{CustomerID: 1, ShippingMethodID: 2, DestAddressID: 17},
		{CustomerID: 2, ShippingMethodID: 1, DestAddressID: 2},
		{CustomerID: 2, ShippingMethodID: 3, DestAddressID: 2},
		{CustomerID: 3, ShippingMethodID: 2, DestAddressID: 3},
		{CustomerID: 3, ShippingMethodID: 1, DestAddressID: 3},
		{CustomerID: 4, ShippingMethodID: 4, DestAddressID: 4},
		{CustomerID: 5, ShippingMethodID: 1, DestAddressID: 5},
		{CustomerID: 5, ShippingMethodID: 2, DestAddressID: 18},
		{CustomerID: 6, ShippingMethodID: 1, DestAddressID: 6},
		{CustomerID: 7, ShippingMethodID: 3, DestAddressID: 7},
		{CustomerID: 8, ShippingMethodID: 2, DestAddressID: 8},
		{CustomerID: 9, ShippingMethodID: 1, DestAddressID: 9},
		{CustomerID: 10, ShippingMethodID: 1, DestAddressID: 10},
		{CustomerID: 10, ShippingMethodID: 3, DestAddressID: 19},
		{CustomerID: 11, ShippingMethodID: 2, DestAddressID: 11},
		{CustomerID: 12, ShippingMethodID: 1, DestAddressID: 12},
		{CustomerID: 13, ShippingMethodID: 4, DestAddressID: 13},
		{CustomerID: 14, ShippingMethodID: 1, DestAddressID: 14},
		{CustomerID: 15, ShippingMethodID: 3, DestAddressID: 15},
		{CustomerID: 15, ShippingMethodID: 1, DestAddressID: 20},
		{CustomerID: 16, ShippingMethodID: 1, DestAddressID: 16},
		{CustomerID: 17, ShippingMethodID: 2, DestAddressID: 17},
		{CustomerID: 17, ShippingMethodID: 1, DestAddressID: 17},
		{CustomerID: 18, ShippingMethodID: 4, DestAddressID: 18},
		{CustomerID: 19, ShippingMethodID: 2, DestAddressID: 19},
		{CustomerID: 19, ShippingMethodID: 1, DestAddressID: 19},
		{CustomerID: 20, ShippingMethodID: 1, DestAddressID: 20},
		{CustomerID: 20, ShippingMethodID: 3, DestAddressID: 20},
	}).Error
}

func seedOrderLines(d *gorm.DB) error {
	return d.Create(&[]entity.OrderLine{
		{OrderID: 1, GameID: 1, Price: 49.99}, {OrderID: 1, GameID: 2, Price: 39.99},
		{OrderID: 2, GameID: 3, Price: 59.99}, {OrderID: 2, GameID: 6, Price: 59.99},
		{OrderID: 3, GameID: 7, Price: 9.99}, {OrderID: 3, GameID: 8, Price: 9.99},
		{OrderID: 4, GameID: 10, Price: 39.99}, {OrderID: 4, GameID: 11, Price: 59.99}, {OrderID: 4, GameID: 12, Price: 19.99},
		{OrderID: 5, GameID: 13, Price: 29.99}, {OrderID: 5, GameID: 14, Price: 59.99}, {OrderID: 5, GameID: 15, Price: 14.99},
		{OrderID: 6, GameID: 16, Price: 39.99}, {OrderID: 6, GameID: 17, Price: 29.99},
		{OrderID: 7, GameID: 19, Price: 14.99}, {OrderID: 7, GameID: 20, Price: 29.99},
		{OrderID: 8, GameID: 22, Price: 19.99}, {OrderID: 8, GameID: 23, Price: 19.99},
		{OrderID: 9, GameID: 25, Price: 29.99}, {OrderID: 9, GameID: 26, Price: 39.99},
		{OrderID: 10, GameID: 28, Price: 26.95}, {OrderID: 10, GameID: 30, Price: 0.00},
		{OrderID: 11, GameID: 31, Price: 24.99}, {OrderID: 11, GameID: 32, Price: 19.99},
		{OrderID: 12, GameID: 34, Price: 19.99}, {OrderID: 12, GameID: 35, Price: 39.99},
		{OrderID: 13, GameID: 37, Price: 14.99}, {OrderID: 13, GameID: 38, Price: 29.99},
		{OrderID: 14, GameID: 40, Price: 14.99}, {OrderID: 14, GameID: 41, Price: 24.99},
		{OrderID: 15, GameID: 18, Price: 19.99}, {OrderID: 15, GameID: 43, Price: 39.99},
		{OrderID: 16, GameID: 45, Price: 59.99}, {OrderID: 16, GameID: 46, Price: 44.99},
		{OrderID: 17, GameID: 47, Price: 14.99}, {OrderID: 17, GameID: 48, Price: 14.99},
		{OrderID: 18, GameID: 33, Price: 29.99}, {OrderID: 18, GameID: 36, Price: 39.99},
		{OrderID: 19, GameID: 21, Price: 39.99}, {OrderID: 19, GameID: 19, Price: 14.99},
		{OrderID: 20, GameID: 24, Price: 9.99}, {OrderID: 20, GameID: 22, Price: 19.99},
		{OrderID: 21, GameID: 3, Price: 59.99}, {OrderID: 22, GameID: 27, Price: 49.99},
		{OrderID: 23, GameID: 42, Price: 19.99}, {OrderID: 24, GameID: 29, Price: 19.99},
		{OrderID: 25, GameID: 39, Price: 39.99}, {OrderID: 26, GameID: 44, Price: 14.99},
		{OrderID: 27, GameID: 9, Price: 0.00}, {OrderID: 28, GameID: 5, Price: 39.99},
		{OrderID: 29, GameID: 4, Price: 49.99},
	}).Error
}

func hash(pw string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic("hash: " + err.Error())
	}
	return string(h)
}