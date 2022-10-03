package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerCrusaderStrikeSpell() {
	baseCost := paladin.BaseMana * 0.05

	paladin.CrusaderStrike = paladin.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 35395},
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost *
					(1 - 0.02*float64(paladin.Talents.Benediction)) *
					core.TernaryFloat64(paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfCrusaderStrike), 0.8, 1),
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // cs is on phys gcd, which cannot be hasted
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 4, // the cd is 4 seconds in 3.3
			},
		},

		BonusCritRating: core.TernaryFloat64(paladin.HasSetBonus(ItemSetAegisBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1 +
			paladin.getTalentSanctityOfBattleBonus() +
			paladin.getTalentTheArtOfWarBonus() +
			paladin.getItemSetGladiatorsVindicationBonusGloves(),
		DamageMultiplier: 0.75,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfigMeleeWeapon(
				core.MainHand,
				true, // cs is subject to normalisation
				core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 31033, 36, 0)+ // Libram of Righteous Power
					core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 40191, 79, 0), // Libram of Radiance
				true,
			),
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(),
		}),
	})
}
