package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const fireTotemDuration time.Duration = time.Second * 120

func (shaman *Shaman) registerFireElementalTotem() {
	if !shaman.Totems.UseFireElemental {
		return
	}

	actionID := core.ActionID{SpellID: 2894}
	manaCost := 0.23 * shaman.BaseMana

	fireElementalAura := shaman.RegisterAura(core.Aura{
		Label:    "Fire Elemental Totem",
		ActionID: actionID,
		Duration: fireTotemDuration,
	})

	shaman.FireElementalTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.Mana,
		BaseCost:     manaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: manaCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * time.Duration(core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFireElementalTotem), 5, 10)),
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// TODO Need to handle ToW
			shaman.MagmaTotemDot.Cancel(sim)
			shaman.SearingTotemDot.Cancel(sim)
			shaman.NextTotemDrops[FireTotem] = sim.CurrentTime + fireTotemDuration

			shaman.FireElemental.EnableWithTimeout(sim, shaman.FireElemental, fireTotemDuration)

			// Add a dummy aura to show in metrics
			fireElementalAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell:    shaman.FireElementalTotem,
		Priority: core.CooldownPriorityDefault, // TODO needs to be altered due to snap shotting.
		Type:     core.CooldownTypeDPS,
	})
}
