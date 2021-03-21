package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/zyahrial/blantik-be/api/models"
)

var users = []models.User{
	models.User{
		Nickname: "Khaerus Zyahrial",
		Email:    "khaerus@gmail.com",
		Password: "password",
		Kecamatan:    "Pekalongan Barat",
		Kota:    "Pekalongan",
		Provinsi:    "Jawa Tengah",
		Negara:    "Indonesia",
		Latitude:    "23676327423as",
		Longitude:	"16512363645as",
		Level:	0,
		Alamat:    "Jl. Urip Sumoharjo no.171 02/04 Desa Pringlangu",
	},
	models.User{
		Nickname: "Sheza",
		Email:    "sheza@gmail.com",
		Password: "password",
		Kecamatan:    "Medan Satria",
		Kota:    "Bekasi",
		Provinsi:    "Jawa Barat",
		Negara:    "Indonesia",
		Latitude:    "1235123423232sf",
		Longitude:	"165123.63645sf",
		Level:	1,
		Alamat:    "Jl. Kaliabang Bungur no.11 02/03 Desa Pejuang",
	},
}


var product = []models.Product{
	models.Product{
	Product : "Hasil Pertanian",
    Category : "Cabe Merah",
	Description : "Cabe merah asli pekalongan",
	Price : 15000,
	Qty : 20,
	Status : "Live",
	},
	models.Product{
	Product : "Hasil Peternakan",
    Category : "Sapi Bali",
	Description : "Sapi bali usia 12 bulan",
	Price : 10000000,
	Status : "Pause",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	models.Post{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}, &models.Product{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}, &models.Product{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}


	err = db.Debug().Model(&models.Product{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	
	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id","users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
	

	for i, _ := range users {

		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}

		product[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Product{}).Create(&product[i]).Error
		if err != nil {
			log.Fatalf("cannot seed product table: %v", err)
		}
		
	}
	
}