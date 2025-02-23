package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerShadowBoltSpell() {
	baseCost := 0.17 * warlock.BaseMana
	actionID := core.ActionID{SpellID: 47809}
	spellSchool := core.SpellSchoolShadow
	spellCoeff := 0.857 * (1 + 0.04*float64(warlock.Talents.ShadowAndFlame))
	ISBProcChance := 0.2 * float64(warlock.Talents.ImprovedShadowBolt)

	warlock.ShadowBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskSpellDamage,
		MissileSpeed: 20,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm] - core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowBolt), 0.1, 0)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (3000 - 100*time.Duration(warlock.Talents.Bane)),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.GCD = time.Duration(float64(cast.GCD) * warlock.backdraftModifier())
				cast.CastTime = time.Duration(float64(cast.CastTime) * warlock.backdraftModifier())
				if warlock.Talents.Nightfall > 0 {
					warlock.applyNightfall(cast)
				}
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistShadowCrit() +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 4), 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false),
		CritMultiplier:           warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(694, 775) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					// ISB debuff
					if ISBProcChance > 0 && (ISBProcChance == 1 || sim.RandomFloat("ISB") <= ISBProcChance) {
						// TODO precalculate aura
						core.ShadowMasteryAura(target).Activate(sim) // calls refresh if already active.
					}
				}
				spell.DealDamage(sim, &result)
			})
		},
	})
}
