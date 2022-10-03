package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerShieldOfRighteousnessSpell() {
	baseCost := paladin.BaseMana * 0.06

	paladin.ShieldOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 61411},
		SpellSchool:  core.SpellSchoolHoly,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					// TODO: Derive or find accurate source for DR curve
					bv := paladin.GetStat(stats.BlockValue)
					if bv <= 2400.0 {
						return 520.0 + bv
					} else {
						bv = 2400.0 + (bv-2400.0)/2
						return 520.0 + core.TernaryFloat64(bv > 2760.0, 2760.0, bv)
					}
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(),
		}),
	})
}
