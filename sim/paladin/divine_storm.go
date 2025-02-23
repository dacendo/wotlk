package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerDivineStormSpell() {
	baseCost := paladin.BaseMana * 0.12

	baseEffectMH := core.SpellEffect{ // wait how will this work, something like whirlwind
		BaseDamage: core.BaseDamageConfigMeleeWeapon(
			core.MainHand,
			false, // ds is not subject to normalisation
			core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 45510, 235, 0)+ // Libram of Discord
				core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 38362, 81, 0), // Venture Co. Libram of Retribution
			// (much akin to what stuff like hs or ms have intrinsically)
			true,
		),
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(),
	}

	numHits := core.MinInt32(4, paladin.Env.GetNumTargets())
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		mhEffect := baseEffectMH
		mhEffect.Target = paladin.Env.GetTargetUnit(i)
		effects = append(effects, mhEffect)
	}

	paladin.DivineStorm = paladin.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 53385},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true, // ds is on phys gcd, which cannot be hasted
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		BonusCritRating: core.TernaryFloat64(paladin.HasSetBonus(ItemSetAegisBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,
		// base 1.1 multiplier, can be further improved by 10% via taow for a grand total of 1.21. NOTE: Unlike cs, ds tooltip IS NOT updated to reflect this.
		DamageMultiplierAdditive: 1 +
			paladin.getTalentTheArtOfWarBonus() +
			paladin.getItemSetRedemptionBattlegearBonus2(),
		DamageMultiplier: 1.1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})
}
