package database

import (
	"github.com/ling/muse/config"
	"github.com/ling/muse/entity"
)

type CharacterRepo struct {
}

func (c *CharacterRepo) InsertChara(character *entity.Character) error {
	db := config.GetDB()
	result := db.Select("WorldInfo").Create(character)
	return result.Error
}
