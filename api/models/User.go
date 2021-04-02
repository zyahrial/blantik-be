package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key; auto_increment" json:"id"`
	Nickname  string    `gorm:"size:255;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Phone     string    `gorm:"size:20;not null;unique" json:"phone"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Kecamatan  string    `gorm:"size:255;not null;" json:"kecamatan"`
	Kota  string    `gorm:"size:255;not null;" json:"kota"`
	Provinsi  string    `gorm:"size:255;not null;" json:"provinsi"`
	Negara  string    `gorm:"size:255;not null;" json:"negara"`
	Alamat  string    `gorm:"size:255;not null;" json:"alamat"`
	Latitude  string    `gorm:"size:255;not null;" json:"latitude"`
	Longitude  string   `gorm:"size:255;not null;" json:"longitude"`
	Level  int64   `gorm:"size:11;not null;" json:"level"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Phone = html.EscapeString(strings.TrimSpace(u.Phone))
	u.Kecamatan = html.EscapeString(strings.TrimSpace(u.Kecamatan))
	u.Kota = html.EscapeString(strings.TrimSpace(u.Kota))
	u.Provinsi = html.EscapeString(strings.TrimSpace(u.Provinsi))
	u.Negara = html.EscapeString(strings.TrimSpace(u.Negara))
	u.Alamat = html.EscapeString(strings.TrimSpace(u.Alamat))
	u.Latitude = html.EscapeString(strings.TrimSpace(u.Latitude))
	u.Longitude = html.EscapeString(strings.TrimSpace(u.Longitude))
	u.Level = 0
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if u.Phone == "" {
			return errors.New("Required Phone")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		// if u.Kecamatan == "" {
		// 	return errors.New("Required Kecamatan")
		// }
		// if u.Kota == "" {
		// 	return errors.New("Required Kota")
		// }
		// if u.Provinsi == "" {
		// 	return errors.New("Required Provinsi")
		// }
		// if u.Negara == "" {
		// 	return errors.New("Required Negara")
		// }
		// if u.Alamat == "" {
		// 	return errors.New("Required Alamat")
		// }
		}

		return nil
		case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
// 	var err error
// 	users := []User{}

// 	err = db.Debug().Model(&User{}).Where("level = 0").Limit(100).Find(&users).Error

// 	if err != nil {
// 		return &[]User{}, err
// 	}
// 	return &users, err
// }

func (p *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	if len(users) > 0 {
		for i, _ := range users {
			err := db.Debug().Model(&User{}).Where("level = ?", users[i].Level).Take(&users[i]).Error
			if err != nil {
				return &[]User{}, err
			}
		}
	}
	return &users, nil
}


func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) FindMe(db *gorm.DB, email string) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("email = ?", email).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"nickname":  u.Nickname,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}