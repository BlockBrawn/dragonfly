package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// NewFallingText creates and returns a new Text entity with the text and position provided with gravity.
func NewFallingText(text string, pos mgl64.Vec3) *world.EntityHandle {
	return world.EntitySpawnOpts{Position: pos, NameTag: text}.New(FallingTextType, textConf)
}

var fallingTextConf = FallingBlockBehaviourConfig{
	Gravity: 0.04,
	Drag:    0.02,
}

// FallingTextType is a world.EntityType implementation for Text.
var FallingTextType fallingTextType

type fallingTextType struct{}

func (t fallingTextType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}
func (fallingTextType) EncodeEntity() string        { return "dragonfly:falling_text" }
func (fallingTextType) BBox(world.Entity) cube.BBox { return cube.BBox{} }
func (fallingTextType) NetworkEncodeEntity() string { return "minecraft:falling_block" }

func (fallingTextType) DecodeNBT(_ map[string]any, data *world.EntityData) {
	data.Data = fallingTextConf.New()
}
func (fallingTextType) EncodeNBT(_ *world.EntityData) map[string]any { return nil }
