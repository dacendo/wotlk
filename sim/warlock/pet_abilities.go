package warlock

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type PetAbilityType byte

const (
	Unknown PetAbilityType = iota
	Cleave
	Intercept
	LashOfPain
	ShadowBite
	Firebolt
)

// Returns whether the ability was successfully cast.
func (wp *WarlockPet) TryCast(sim *core.Simulation, target *core.Unit, spell *core.Spell) bool {
	if wp.CurrentMana() < spell.DefaultCast.Cost {
		return false
	}
	if !spell.IsReady(sim) {
		return false
	}

	spell.Cast(sim, target)
	return true
}

func (wp *WarlockPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) *core.Spell {
	switch abilityType {
	case Cleave:
		return wp.newCleave()
	case Intercept:
		return wp.newIntercept()
	case LashOfPain:
		return wp.newLashOfPain()
	case Firebolt:
		return wp.newFirebolt()
	case ShadowBite:
		return wp.newShadowBite()
	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

func (wp *WarlockPet) newIntercept() *core.Spell {
	return nil
}

func (wp *WarlockPet) newFirebolt() *core.Spell {
	baseCost := 180.0
	return wp.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47964},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (2500 - time.Duration(250*wp.owner.Talents.DemonicPower)),
			},
			IgnoreHaste: true,
		},

		BonusCritRating: wp.owner.masterDemonologistFireCrit(),
		DamageMultiplier: (1 + 0.1*float64(wp.owner.Talents.ImprovedImp)) *
			(1 + 0.2*core.TernaryFloat64(wp.owner.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImp), 1, 0)),
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(203, 227) + 0.571*spell.SpellPower()
			spell.CalcAndDealDamageMagicHitAndCrit(sim, target, baseDamage)
		},
	})
}

func (wp *WarlockPet) newCleave() *core.Spell {
	baseCost := 439.0 // 10% of base

	baseEffect := core.SpellEffect{
		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 124, true),
		OutcomeApplier: wp.OutcomeFuncMeleeSpecialHitAndCrit(),
	}

	numHits := core.MinInt32(2, wp.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		effects = append(effects, baseEffect)
		effects[i].Target = wp.Env.GetTargetUnit(i)
	}

	return wp.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47994},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}

func (wp *WarlockPet) newLashOfPain() *core.Spell {
	baseCost := 250.0
	return wp.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 47992},
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * (12 - time.Duration(3*wp.owner.Talents.DemonicPower)),
			},
		},

		BonusCritRating:  wp.owner.masterDemonologistShadowCrit(),
		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: the hidden 5% damage modifier succ currently gets also applies to this ...
			baseDamage := 237 + 0.429*spell.SpellPower()
			spell.CalcAndDealDamageMagicHitAndCrit(sim, target, baseDamage)
		},
	})
}

func (wp *WarlockPet) newShadowBite() *core.Spell {
	actionID := core.ActionID{SpellID: 54053}
	baseCost := 131.0 // TODO: should be 3% of BaseMana, but it's unclear what that actually refers to with pets

	var petManaMetrics *core.ResourceMetrics
	maxManaMult := 0.04 * float64(wp.owner.Talents.ImprovedFelhunter)
	impFelhunter := wp.owner.Talents.ImprovedFelhunter > 0
	if impFelhunter {
		petManaMetrics = wp.NewManaMetrics(actionID)
	}

	return wp.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * (6 - time.Duration(2*wp.owner.Talents.ImprovedFelhunter)),
			},
		},

		DamageMultiplier: 1 + 0.03*float64(wp.owner.Talents.ShadowMastery),
		CritMultiplier:   1.5 + 0.1*float64(wp.owner.Talents.Ruin),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(97+1, 97+41) + 0.429*spell.SpellPower()

			w := wp.owner
			dots := []*core.Dot{
				w.UnstableAfflictionDot, w.ImmolateDot, w.CurseOfAgonyDot,
				w.CurseOfDoomDot, w.CorruptionDot, w.ConflagrateDot,
				w.SeedDots[target.Index], w.DrainSoulDot,
				// missing: drain life, shadowflame
			}
			counter := 0
			for _, dot := range dots {
				if dot.IsActive() {
					counter++
				}
			}

			baseDamage *= 1 + 0.15*float64(counter)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if impFelhunter && result.Landed() {
				wp.AddMana(sim, wp.MaxMana()*maxManaMult, petManaMetrics, true)
			}
			spell.DealDamage(sim, &result)
		},
	})
}
