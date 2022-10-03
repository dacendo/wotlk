package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerUnstableAfflictionSpell() {
	baseCost := 0.15 * warlock.BaseMana
	actionID := core.ActionID{SpellID: 47843}
	spellSchool := core.SpellSchoolShadow

	warlock.UnstableAffliction = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskEmpty,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (1500 - 200*core.TernaryDuration(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfUnstableAffliction), 1, 0)),
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistShadowCrit() +
			3*core.CritRatingPerCritChance*float64(warlock.Talents.Malediction),
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true),
		CritMultiplier:           warlock.SpellCritMultiplier(1, 1),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier:  warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt: applyDotOnLanded(&warlock.UnstableAfflictionDot),
		}),
	})

	target := warlock.CurrentTarget
	spellCoefficient := 0.2 + 0.01*float64(warlock.Talents.EverlastingAffliction)
	applier := warlock.OutcomeFuncTick()
	if warlock.Talents.Pandemic {
		applier = warlock.OutcomeFuncMagicCrit()
	}

	warlock.UnstableAfflictionDot = core.NewDot(core.Dot{
		Spell: warlock.UnstableAffliction,
		Aura: target.RegisterAura(core.Aura{
			Label:    "UnstableAffliction-" + strconv.Itoa(int(warlock.Index)),
			ActionID: core.ActionID{SpellID: 47843},
		}),
		NumberOfTicks: 5,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(1150/5, spellCoefficient),
			OutcomeApplier: applier,
			IsPeriodic:     true,
		}),
	})
}
