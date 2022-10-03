package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerMortalStrikeSpell(cdTimer *core.Timer) {
	cost := 30.0
	refundAmount := cost * 0.8

	warrior.MortalStrike = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47486},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Second*6 - time.Millisecond*200*time.Duration(warrior.Talents.ImprovedMortalStrike),
			},
		},

		BonusCritRating: core.TernaryFloat64(warrior.HasSetBonus(ItemSetSiegebreakerBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 +
			[]float64{0, 0.03, 0.06, 0.1}[warrior.Talents.ImprovedMortalStrike] +
			core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfMortalStrike), 0.1, 0) +
			core.TernaryFloat64(warrior.HasSetBonus(ItemSetOnslaughtBattlegear, 4), 0.05, 0),
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 380, true),
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) CanMortalStrike(sim *core.Simulation) bool {
	if warrior.Talents.MortalStrike {
		return warrior.CurrentRage() >= warrior.MortalStrike.DefaultCast.Cost && warrior.MortalStrike.IsReady(sim)
	} else {
		return false
	}
}
