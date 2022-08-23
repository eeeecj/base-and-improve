package sql

import (
	"errors"
	"github.com/balance/boolfilter"
	"github.com/balance/boolfilter/global"
	"github.com/balance/boolfilter/memeory"
	"gorm.io/gorm"
)

func NewSqlFilter(db *gorm.DB, tableName, key string, blen int, hashes ...global.HashFunc) (boolfilter.Adapter, error) {
	f := &SqlFilter{
		Filter: &memeory.Filter{
			Bits:      make([]byte, blen, blen),
			Hashes:    hashes,
			IsChanged: false,
		},
		key:       key,
		db:        db,
		tableName: tableName,
	}
	err := f.Init(f.db, f.tableName, f.key, f.Bits)
	if err != nil {
		err := f.CreateTable()
		if err != nil {
			return nil, err
		}
		return NewSqlFilter(f.db, f.tableName, f.key, blen, f.Hashes...)
	}
	return f, err
}

type SqlFilter struct {
	*memeory.Filter
	key       string
	db        *gorm.DB
	tableName string
}
type Bloom struct {
	Id  string `gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Val []byte `gorm:"column:value"`
}

func (s *SqlFilter) Clear() error {
	bloom := Bloom{
		Id:  s.key,
		Val: s.Bits,
	}
	err := s.db.Table(s.tableName).Delete(&bloom).Error
	if err != nil {
		return err
	}
	return s.Filter.Clear()
}

func (s *SqlFilter) Init(db *gorm.DB, tableName, key string, bytes []byte) error {
	if key == "" {
		return errors.New("empty key")
	}
	s.db = db
	s.tableName = tableName
	s.key = key
	var bloom Bloom
	err := s.db.Table(s.tableName).Where("id=?", s.key).Limit(1).Scan(&bloom).Error
	if err != nil {
		return err
	}
	if bloom.Id == "" {
		bloom = Bloom{Id: s.key, Val: s.Bits}
		err := s.db.Table(s.tableName).Create(&bloom).Error
		if err != nil {
			return err
		}
	} else {
		if len(s.Bits) != len(bloom.Val) {
			return errors.New("the length is not consistent")
		} else {
			s.Bits = bloom.Val
		}
	}
	return nil
}

func (s *SqlFilter) CreateTable() error {
	return s.db.Table(s.tableName).AutoMigrate(&Bloom{})
}

func (s *SqlFilter) Write() error {
	return s.db.Table(s.tableName).Where("id=?", s.key).Update("value", s.Bits).Error
}

func (s *SqlFilter) Close() error {
	db, err := s.db.DB()
	if db == nil {
		return err
	}
	return db.Close()
}
