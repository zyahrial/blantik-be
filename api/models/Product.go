package models

import (
	"errors"
	"html"
	"strings"
	"time"

	guuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

	type Product struct {
	ID           	string		`gorm:"primary_key;not null;unique" json:"id"`
	Product        string		`gorm:"size:255;not null;unique" json:"product"`
	Category       string		`gorm:"size:255;not null;unique" json:"category"`
	Price          int64		`gorm:"size:255;not null;unique" json:"price"`
	Qty			   uint32		`gorm:"size:255;not null;unique" json:"qty"`
	Author    		User      `json:"author"`
	AuthorID  		uint32    `gorm:"not null" json:"author_id"`
	Description     string		`gorm:"size:255;not null;unique" json:"description"`
	CreatedAt    	time.Time		`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    	time.Time		`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
  }

  func (p *Product) Prepare() {
	p.ID = guuid.New()
	p.Product = html.EscapeString(strings.TrimSpace(p.Product))
	p.Category = html.EscapeString(strings.TrimSpace(p.Category))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Price = html.EscapeString(strings.TrimSpace(p.Price))
	p.Qty = html.EscapeString(strings.TrimSpace(p.Qty))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	}

	func (p *Product) Validate() error {

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
		if p.Price < 1 {
			return errors.New("Required Price")
		}
		if p.Qty < 1 {
			return errors.New("Required Qty")
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