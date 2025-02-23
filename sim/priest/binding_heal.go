package priest

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerBindingHealSpell() {
	baseCost := .27 * priest.BaseMana
	spellCoeff := 0.8057 + 0.04*float64(priest.Talents.EmpoweredHealing)

	priest.BindingHeal = priest.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48120},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskSpellHealing,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,

				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCritRating: float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.DivineProvidence)),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 0.5 * (1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve]),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healFromSP := spellCoeff * spell.HealingPower()

			selfHealing := sim.Roll(1959, 2516) + healFromSP
			spell.CalcAndDealHealingCrit(sim, &priest.Unit, selfHealing)

			targetHealing := sim.Roll(1959, 2516) + healFromSP
			spell.CalcAndDealHealingCrit(sim, target, targetHealing)
		},
	})
}
