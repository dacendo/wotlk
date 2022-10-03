package deathknight

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO: Cleanup death strike the same way we did for plague strike
var DeathStrikeActionID = core.ActionID{SpellID: 49924}

func (dk *Deathknight) newDeathStrikeSpell(isMH bool, onhit func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect)) *RuneSpell {
	bonusBaseDamage := dk.sigilOfAwarenessBonus()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 297.0+bonusBaseDamage, true)
	if !isMH {
		// SpellID 66953
		weaponBaseDamage = core.BaseDamageFuncMeleeWeapon(core.OffHand, true, 148.0+bonusBaseDamage, true)
	}

	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDeathStrike)

	effect := core.SpellEffect{
		BaseDamage: core.BaseDamageConfig{
			Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
				glyphDmgMultiplier := core.TernaryFloat64(hasGlyph, 1.0+0.01*core.MinFloat(dk.CurrentRunicPower(), 25.0), 1.0)

				return weaponBaseDamage(sim, hitEffect, spell) * dk.RoRTSBonus(hitEffect.Target) * glyphDmgMultiplier
			},
		},

		OnSpellHitDealt: onhit,
	}

	procMask := dk.threatOfThassarianProcMasks(isMH, &effect)

	conf := core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating: (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .75 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			(1.0 + 0.15*float64(dk.Talents.ImprovedDeathStrike)),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	}

	rs := &RuneSpell{}
	if isMH {
		conf.ResourceType = stats.RunicPower
		conf.BaseCost = float64(core.NewRuneCost(uint8(15.0+2.5*float64(dk.Talents.Dirge)), 0, 1, 1, 0))
		conf.Cast = core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: conf.BaseCost,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			IgnoreHaste: true,
		}
		conf.ApplyEffects = dk.withRuneRefund(rs, effect, false)
		if dk.Talents.DeathRuneMastery == 3 {
			rs.DeathConvertChance = 1.0
		} else {
			rs.DeathConvertChance = float64(dk.Talents.DeathRuneMastery) * 0.33
		}
		rs.ConvertType = RuneTypeFrost | RuneTypeUnholy
	}

	if isMH {
		return dk.RegisterSpell(rs, conf, func(sim *core.Simulation) bool {
			return dk.CastCostPossible(sim, 0.0, 0, 1, 1) && dk.DeathStrike.IsReady(sim)
		}, nil)
	} else {
		return dk.RegisterSpell(rs, conf, nil, nil)
	}
}

func (dk *Deathknight) registerDeathStrikeSpell() {
	healthMetrics := dk.NewHealthMetrics(DeathStrikeActionID)

	dk.DeathStrikeOhHit = dk.newDeathStrikeSpell(false, nil)
	dk.DeathStrikeMhHit = dk.newDeathStrikeSpell(true, func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
		dk.LastOutcome = spellEffect.Outcome

		if spellEffect.Landed() {
			healingAmount := 0.05 * dk.dkCountActiveDiseases(dk.CurrentTarget) * dk.MaxHealth() * (1.0 + 0.5*float64(dk.Talents.ImprovedDeathStrike)) * (1.0 + core.TernaryFloat64(dk.VampiricBloodAura.IsActive(), 0.35, 0.0))
			dk.GainHealth(sim, healingAmount, healthMetrics)
			dk.DeathStrikeHeals = append(dk.DeathStrikeHeals, healingAmount)
		}

		dk.threatOfThassarianProc(sim, spellEffect, dk.DeathStrikeOhHit)
	})
	dk.DeathStrike = dk.DeathStrikeMhHit
}

func (dk *Deathknight) registerDrwDeathStrikeSpell() {
	bonusBaseDamage := dk.sigilOfAwarenessBonus()
	weaponBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, true, 297.0+bonusBaseDamage, true)

	hasGlyph := dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfDeathStrike)

	dk.RuneWeapon.DeathStrike = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    DeathStrikeActionID.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		BonusCritRating:  (dk.annihilationCritBonus() + dk.improvedDeathStrikeCritBonus()) * core.CritRatingPerCritChance,
		DamageMultiplier: .75 * dk.improvedDeathStrikeDamageBonus(),
		CritMultiplier:   dk.RuneWeapon.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: dk.RuneWeapon.OutcomeFuncMeleeWeaponSpecialHitAndCrit(),
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					bonusDamage := core.TernaryFloat64(hasGlyph, 1.0+core.MinFloat(0.25, dk.CurrentRunicPower()/100.0), 1.0)
					return weaponBaseDamage(sim, hitEffect, spell) * bonusDamage
				},
			},
		}),
	})
}
