package druid

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (druid *Druid) registerSwipeBearSpell() {
	cost := 20.0 - float64(druid.Talents.Ferocity)

	baseDamage := 108.0
	if druid.Equip[items.ItemSlotRanged].ID == 23198 { // Idol of Brutality
		baseDamage += 10
	} else if druid.Equip[items.ItemSlotRanged].ID == 38365 { // Idol of Perspicacious Attacks
		baseDamage += 24
	}

	lbdm := core.TernaryFloat64(druid.HasSetBonus(ItemSetLasherweaveBattlegear, 2), 1.2, 1.0)
	thdm := core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 4), 1.15, 1.0)
	fidm := 1.0 + 0.1*float64(druid.Talents.FeralInstinct)

	baseEffect := core.SpellEffect{
		BaseDamage: core.BaseDamageConfig{
			Calculator: core.BaseDamageFuncMelee(baseDamage, baseDamage, 0.07),
		},
		OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(),
	}

	druid.SwipeBear = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48562},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			ModifyCast:  druid.ApplyClearcasting,
			IgnoreHaste: true,
		},

		DamageMultiplier: lbdm * thdm * fidm,
		CritMultiplier:   druid.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncAOEDamageCapped(druid.Env, baseEffect),
	})
}

func (druid *Druid) registerSwipeCatSpell() {
	cost := 50.0 - float64(druid.Talents.Ferocity)

	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 0.0, false)
	weaponMulti := 2.5
	fidm := 1.0 + 0.1*float64(druid.Talents.FeralInstinct)

	baseEffect := core.SpellEffect{
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				return weaponBaseDamage(sim, hitEffect, spell)
			},
		},
		OutcomeApplier: druid.OutcomeFuncMeleeSpecialHitAndCrit(),
	}

	druid.SwipeCat = druid.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 62078},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
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

		DamageMultiplier: fidm * weaponMulti,
		CritMultiplier:   druid.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncAOEDamageCapped(druid.Env, baseEffect),
	})
}

func (druid *Druid) CanSwipeCat() bool {
	return druid.InForm(Cat) && (druid.CurrentEnergy() >= druid.CurrentSwipeCatCost() || druid.ClearcastingAura.IsActive())
}

func (druid *Druid) CurrentSwipeCatCost() float64 {
	return druid.SwipeCat.ApplyCostModifiers(druid.SwipeCat.BaseCost)
}

func (druid *Druid) CanSwipeBear() bool {
	return druid.InForm(Bear) && druid.CurrentRage() >= druid.SwipeBear.DefaultCast.Cost
}

func (druid *Druid) IsSwipeSpell(spell *core.Spell) bool {
	return spell == druid.SwipeBear || spell == druid.SwipeCat
}
