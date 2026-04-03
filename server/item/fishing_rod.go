package item

import (
	"time"

	"github.com/df-mc/dragonfly/server/world"
)

type fishingRodUser interface {
	User
	FishingHook() *world.EntityHandle
	SetFishingHook(*world.EntityHandle)
}

type FishingRod struct{}

func (FishingRod) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: 384,
		BrokenItem:    func() Stack { return Stack{} },
	}
}

func (FishingRod) MaxCount() int {
	return 1
}

func (FishingRod) Rod() bool {
	return true
}

func (FishingRod) Cooldown() time.Duration {
	return 100 * time.Millisecond
}

func (FishingRod) Use(tx *world.Tx, user User, ctx *UseContext) bool {
	p, ok := user.(fishingRodUser)
	if !ok {
		return false
	}

	if hook := p.FishingHook(); hook != nil {
		if ent, ok := hook.Entity(tx); ok {
			_ = ent.Close()
		}
		p.SetFishingHook(nil)
		ctx.DamageItem(1)

		return true
	}

	create := tx.World().EntityRegistry().Config().FishingHook
	opts := world.EntitySpawnOpts{Position: eyePosition(user), Velocity: user.Rotation().Vec3().Mul(1.3)}
	handle := create(opts, user)
	tx.AddEntity(handle)
	p.SetFishingHook(handle)

	ctx.DamageItem(1)
	return true
}

func (FishingRod) EncodeItem() (name string, meta int16) {
	return "minecraft:fishing_rod", 0
}
