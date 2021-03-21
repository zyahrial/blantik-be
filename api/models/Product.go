package models

import (
	"errors"
	"html"
	"strings"
	"time"
	"log"


	"github.com/sony/sonyflake"
	"github.com/jinzhu/gorm"
)

	type Product struct {
	UUID           uint64		`gorm:"primary_key;not null;unique" json:"uuid"`
	Product        string		`gorm:"size:255;not null" json:"product"`
	Category       string		`gorm:"size:255;not null" json:"category"`
	Price          int64		`gorm:"not null" json:"price"`
	Qty			   uint32		`gorm:"not null" json:"qty"`
	Author    		User      `json:"author"`
	AuthorID  		uint32    `gorm:"not null" json:"author_id"`
	Description     string		`gorm:"size:500;not null" json:"description"`
	Status     		string		`gorm:"size:255;not null" json:"status"`
	CreatedAt    	time.Time		`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    	time.Time		`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
  }

  func (p *Product) Prepare() {

	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}

	p.UUID = id
	p.Product = html.EscapeString(strings.TrimSpace(p.Product))
	p.Category = html.EscapeString(strings.TrimSpace(p.Category))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Status = html.EscapeString(strings.TrimSpace(p.Status))
	p.Author = User{}
	p.Price = int64(p.Price)
	p.Qty = uint32(p.Qty)
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	}

	func (p *Product) Validate() error {

		if p.UUID < 1 {
			return errors.New("Required UUID")
		}
		if p.Product == "" {
			return errors.New("Required Product")
		}
		if p.Category == "" {
			return errors.New("Required Category")
		}
		if p.Description == "" {
			return errors.New("Required Description")
		}
		if p.AuthorID < 1 {
			return errors.New("Required Author")
		}
		if p.Status == "" {
			return errors.New("Required Status")
		}
		return nil
	}

	func (p *Product) SaveProduct(db *gorm.DB) (*Product, error) {
		var err error
		err = db.Debug().Model(&Product{}).Create(&p).Error
		if err != nil {
			return &Product{}, err
		}
		
		return p, nil
	}

	func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
		var err error
		products := []Product{}
		err = db.Debug().Model(&Product{}).Limit(100).Find(&products).Error
		if err != nil {
			return &[]Product{}, err
		}
		if len(products) > 0 {
			for i, _ := range products {
				err := db.Debug().Model(&User{}).Where("id = ?", products[i].AuthorID).Take(&products[i].Author).Error
				if err != nil {
					return &[]Product{}, err
				}
			}
		}
		return &products, nil
	}

	func (p *Product) FindProductByID(db *gorm.DB, pid uint64) (*Product, error) {
		var err error
		err = db.Debug().Model(&Product{}).Where("uuid = ?", pid).Take(&p).Error
		if err != nil {
			return &Product{}, err
		}
		if p.UUID != 0 {
			err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
			if err != nil {
				return &Product{}, err
			}
		}
		return p, nil
	}

	func (p *Product) UpdateAProduct(db *gorm.DB) (*Product, error) {

		var err error
	
		err = db.Debug().Model(&Product{}).Where("uuid = ?", p.UUID).Updates(Product{Product: p.Product, Category: p.Category, Description: p.Description, Status: p.Status, Price: p.Price, Qty: p.Qty, UpdatedAt: time.Now()}).Error
		if err != nil {
			return &Product{}, err
		}
		if p.UUID != 0 {
			err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
			if err != nil {
				return &Product{}, err
			}
		}
		return p, nil
	}

	func (p *Product) DeleteAProduct(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

		db = db.Debug().Model(&Product{}).Where("uuid = ? and author_id = ?", pid, uid).Take(&Product{}).Delete(&Product{})
	
		if db.Error != nil {
			if gorm.IsRecordNotFoundError(db.Error) {
				return 0, errors.New("Product not found")
			}
			return 0, db.Error
		}
		return db.RowsAffected, nil
	}