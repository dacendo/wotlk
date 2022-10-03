package druid

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerRakeSpell() {
	actionID := core.ActionID{SpellID: 48574}

	cost := 40.0 - float64(druid.Talents.Ferocity)

	mangleAura := core.MangleAura(druid.CurrentTarget)

	t9bonus := core.TernaryInt(druid.HasT9FeralSetBonus(2), 1, 0)

	druid.Rake = druid.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIgnoreResists,
		ResourceType: stats.Energy,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  time.Second,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 + 0.1*float64(druid.Talents.SavageFury),
		CritMultiplier:   druid.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					damage := 176 + 0.01*spell.MeleeAttackPower()
					if mangleAura.IsActive() {
						return damage * 1.3
					} else {
						return damage
					}
				},
			},
			OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())
					druid.RakeDot.Apply(sim)
				} else {
					druid.AddEnergy(sim, spell.CurCast.Cost*0.8, druid.EnergyRefundMetrics)
				}
			},
		}),
	})

	dotAura := druid.CurrentTarget.RegisterAura(druid.applyRendAndTear(core.Aura{
		Label:    "Rake-" + strconv.Itoa(int(druid.Index)),
		ActionID: actionID,
		Duration: time.Second * 9,
	}))
	druid.RakeDot = core.NewDot(core.Dot{
		Spell:         druid.Rake,
		Aura:          dotAura,
		NumberOfTicks: 3 + t9bonus,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			IsPeriodic: true,
			BaseDamage: core.BaseDamageConfig{
				Calculator: core.BaseDamageFuncMelee(358, 358, 0.06),
			},
			OutcomeApplier: core.Ternary(druid.HasSetBonus(ItemSetLasherweaveBattlegear, 4), druid.OutcomeFuncTickHitAndCrit(), druid.OutcomeFuncTick()),
		})),
	})
}

func (druid *Druid) CanRake() bool {
	return druid.InForm(Cat) && ((druid.CurrentEnergy() >= druid.CurrentRakeCost()) || druid.ClearcastingAura.IsActive())
}

func (druid *Druid) CurrentRakeCost() float64 {
	return druid.Rake.ApplyCostModifiers(druid.Rake.BaseCost)
}
