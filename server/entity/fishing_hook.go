package entity

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type FishingHookBehaviourConfig struct {
	ProjectileBehaviourConfig
}

func (conf FishingHookBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = &FishingHookBehaviour{
		ProjectileBehaviour: conf.ProjectileBehaviourConfig.New(),
	}
}

type FishingHookBehaviour struct {
	*ProjectileBehaviour
}

type fishingRodHolder interface {
	HeldItems() (mainHand, offHand item.Stack)
}

func (b *FishingHookBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	ownerHandle := b.Owner()
	if ownerHandle == nil {
		_ = e.Close()
		return nil
	}
	owner, ok := ownerHandle.Entity(tx)
	if !ok {
		_ = e.Close()
		return nil
	}
	if p, ok := owner.(fishingRodHolder); ok {
		held, _ := p.HeldItems()
		if r, ok := held.Item().(interface{ Rod() bool }); !ok || !r.Rod() {
			_ = e.Close()
			return nil
		}
	}

	m := b.ProjectileBehaviour.Tick(e, tx)
	if e.Position()[1] < float64(tx.World().Range()[0]) && e.Age()%10 == 0 {
		_ = e.Close()
	}
	return m
}

func NewFishingHook(opts world.EntitySpawnOpts, owner world.Entity) *world.EntityHandle {
	vel := opts.Velocity.Normalize().Add(mgl64.Vec3{
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
	}.Mul(0.007499999832361937)).Mul(1.3)
	vel[0] += opts.Velocity[0]
	vel[2] += opts.Velocity[2]
	opts.Velocity = vel

	conf := FishingHookBehaviourConfig{
		ProjectileBehaviourConfig: ProjectileBehaviourConfig{
			Gravity: 0.1,
			Drag:    0.02,
			Damage:  -1,
			Hit:     fishingHookHit,
		},
	}
	conf.Owner = owner.H()
	return opts.New(FishingHookType, conf)
}

func fishingHookHit(e *Ent, tx *world.Tx, result trace.Result) {
	b := e.Behaviour().(*FishingHookBehaviour)
	ownerHandle := b.Owner()
	if ownerHandle == nil {
		return
	}
	owner, _ := ownerHandle.Entity(tx)

	if res, ok := result.(trace.EntityResult); ok {
		if l, ok := res.Entity().(Living); ok {
			if _, vulnerable := l.Hurt(0.0, ProjectileDamageSource{Projectile: e, Owner: owner}); vulnerable {
				if owner != nil {
					ownerRot, ok1 := owner.(interface{ Rotation() cube.Rotation })
					lRot, ok2 := l.(interface{ Rotation() cube.Rotation })
					if ok1 && ok2 && ownerRot.Rotation().Vec3().Dot(lRot.Rotation().Vec3()) > 0 {
						// Pull back the target.
						l.KnockBack(l.Position().Add(e.Velocity()), 0.230, 0.372)
					} else {
						// Push back the target.
						l.KnockBack(l.Position().Sub(e.Velocity()), 0.374, 0.372)
					}
				}
			}
		}
	}
}

var FishingHookType fishingHookType

type fishingHookType struct{}

func (t fishingHookType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return Open(tx, handle, data)
}

func (fishingHookType) EncodeEntity() string {
	return "minecraft:fishing_hook"
}

func (fishingHookType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}

func (fishingHookType) DecodeNBT(_ map[string]any, data *world.EntityData) {
	conf := FishingHookBehaviourConfig{
		ProjectileBehaviourConfig: ProjectileBehaviourConfig{
			Gravity: 0.1,
			Drag:    0.02,
			Damage:  -1,
			Hit:     fishingHookHit,
		},
	}
	conf.Apply(data)
}
func (fishingHookType) EncodeNBT(*world.EntityData) map[string]any { return nil }
