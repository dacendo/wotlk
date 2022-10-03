package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerMultiShotSpell(timer *core.Timer) {
	baseCost := 0.09 * hunter.BaseMana
	numHits := core.MinInt32(3, hunter.Env.GetNumTargets())

	hunter.MultiShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 49048},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskRangedSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.03*float64(hunter.Talents.Efficiency)),

				GCD:      core.GCDDefault,
				CastTime: 1, // Dummy value so core doesn't optimize the cast away
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.CastTime = hunter.MultiShotCastTime()
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*10 - core.TernaryDuration(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfMultiShot), time.Second*1, 0),
			},
		},

		BonusCritRating: 0 +
			4*core.CritRatingPerCritChance*float64(hunter.Talents.ImprovedBarrage),
		DamageMultiplierAdditive: 1 +
			.04*float64(hunter.Talents.Barrage),
		DamageMultiplier: 1 *
			hunter.markedForDeathMultiplier(),
		CritMultiplier:   hunter.critMultiplier(true, false, hunter.CurrentTarget),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			sharedDmg := hunter.AutoAttacks.Ranged.BaseDamage(sim) +
				hunter.AmmoDamageBonus +
				spell.BonusWeaponDamage() +
				408

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := sharedDmg + 0.2*spell.RangedAttackPower(curTarget)
				spell.CalcAndDealDamageRangedHitAndCrit(sim, curTarget, baseDamage)

				curTarget = sim.Environment.NextTargetUnit(target)
			}
		},
	})
}

func (hunter *Hunter) MultiShotCastTime() time.Duration {
	return time.Duration(float64(time.Millisecond*500) / hunter.RangedSwingSpeed())
}
