package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerRaptorStrikeSpell() {
	baseCost := 0.04 * hunter.BaseMana

	hunter.RaptorStrike = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 48996},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.2*float64(hunter.Talents.Resourcefulness)),
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusCritRating:  float64(hunter.Talents.SavageStrikes) * 10 * core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(false, false, hunter.CurrentTarget),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 335, true),
			OutcomeApplier: hunter.OutcomeFuncMeleeSpecialHitAndCrit(),
		}),
	})
}

// Returns true if the regular melee swing should be used, false otherwise.
func (hunter *Hunter) TryRaptorStrike(sim *core.Simulation) *core.Spell {
	//if hunter.Rotation.Weave == proto.Hunter_Rotation_WeaveAutosOnly || !hunter.RaptorStrike.IsReady(sim) || hunter.CurrentMana() < hunter.RaptorStrike.DefaultCast.Cost {
	//	return nil
	//}

	return hunter.RaptorStrike
}
