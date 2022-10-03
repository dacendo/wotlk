package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerChaosBoltSpell() {
	actionID := core.ActionID{SpellID: 59172}
	spellSchool := core.SpellSchoolFire

	baseCost := 0.07 * warlock.BaseMana
	spellCoeff := 0.7142 * (1 + 0.04*float64(warlock.Talents.ShadowAndFlame))

	warlock.ChaosBolt = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagIgnoreResists,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm]),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2500 - (time.Millisecond * 100 * time.Duration(warlock.Talents.Bane)),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.GCD = time.Duration(float64(cast.GCD) * warlock.backdraftModifier())
				cast.CastTime = time.Duration(float64(cast.CastTime) * warlock.backdraftModifier())
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * (12 - 2*core.TernaryDuration(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfChaosBolt), 1, 0)),
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistFireCrit() +
			core.TernaryFloat64(warlock.Talents.Devastation, 1, 0)*5*core.CritRatingPerCritChance,
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false),
		CritMultiplier:           warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(1429, 1813) + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamageMagicHitAndCrit(sim, target, baseDamage)
		},
	})
}
