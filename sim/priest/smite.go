package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) RegisterSmiteSpell(memeDream bool) {
	baseCost := .15 * priest.BaseMana

	priest.Smite = priest.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48123},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2500 - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
		},

		BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + 0.05*float64(priest.Talents.SearingLight)) *
			core.TernaryFloat64(memeDream, 1.2, 1),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(713, 799) + 0.7143*spell.SpellPower()
			spell.CalcAndDealDamageMagicHitAndCrit(sim, target, baseDamage)
		},
	})
}
