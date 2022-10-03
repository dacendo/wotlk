package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerSweepingStrikesCD() {
	if !warrior.Talents.SweepingStrikes {
		return
	}

	actionID := core.ActionID{SpellID: 12328}

	var curDmg float64
	ssHit := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		// No proc mask, so it won't proc itself.
		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Cancel out caster multiplier so it doesn't double dip.
			spell.CalcAndDealDamageAlwaysHit(sim, target, curDmg/spell.CasterDamageMultiplier())
		},
	})

	ssAura := warrior.RegisterAura(core.Aura{
		Label:     "Sweeping Strikes",
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.SetStacks(sim, 10)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if aura.GetStacks() == 0 || spellEffect.Damage == 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			// TODO: If the triggering spell is Execute and 2nd target health > 20%, do a normalized MH hit instead.

			// Undo armor reduction to get the raw damage value.
			curDmg = spellEffect.Damage / warrior.AttackTables[spellEffect.Target.Index].GetArmorDamageModifier(spell)

			ssHit.Cast(sim, warrior.Env.NextTargetUnit(spellEffect.Target))
			ssHit.SpellMetrics[spellEffect.Target.UnitIndex].Casts--
			if aura.GetStacks() > 0 {
				aura.RemoveStack(sim)
			}
		},
	})

	hasGlyph := warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfSweepingStrikes)
	cost := 30.0
	ssCD := warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if hasGlyph {
					cast.Cost = 0
				}
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			ssAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: ssCD,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return sim.GetNumTargets() > 1 && warrior.CurrentRage() >= ssCD.DefaultCast.Cost
		},
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return true
		},
	})
}
