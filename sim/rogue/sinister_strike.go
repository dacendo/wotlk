package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (rogue *Rogue) SinisterStrikeEnergyCost() float64 {
	return []float64{45, 42, 40}[rogue.Talents.ImprovedSinisterStrike]
}

func (rogue *Rogue) registerSinisterStrikeSpell() {
	energyCost := rogue.SinisterStrikeEnergyCost()
	refundAmount := energyCost * 0.8
	hasGlyphOfSinisterStrike := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfSinisterStrike)

	rogue.SinisterStrike = rogue.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48638},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder,
		ResourceType: stats.Energy,
		BaseCost:     energyCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: energyCost,
				GCD:  time.Second,
			},
			IgnoreHaste: true,
			ModifyCast:  rogue.CastModifier,
		},

		BonusCritRating: core.TernaryFloat64(rogue.HasSetBonus(ItemSetVanCleefs, 4), 5*core.CritRatingPerCritChance, 0) +
			[]float64{0, 2, 4, 6}[rogue.Talents.TurnTheTables]*core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.03*float64(rogue.Talents.Aggression) +
			0.05*float64(rogue.Talents.BladeTwisting) +
			core.TernaryFloat64(rogue.Talents.SurpriseAttacks, 0.1, 0) +
			core.TernaryFloat64(rogue.HasSetBonus(ItemSetSlayers, 4), 0.06, 0),
		CritMultiplier:   rogue.MeleeCritMultiplier(true, true),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 180, true),
			OutcomeApplier: rogue.OutcomeFuncMeleeSpecialHitAndCrit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					points := int32(1)
					if hasGlyphOfSinisterStrike && spellEffect.DidCrit() {
						if sim.RandomFloat("Glyph of Sinister Strike") < 0.5 {
							points += 1
						}
					}
					rogue.AddComboPoints(sim, points, spell.ComboPointMetrics())
				} else {
					rogue.AddEnergy(sim, refundAmount, rogue.EnergyRefundMetrics)
				}
			},
		}),
	})
}
