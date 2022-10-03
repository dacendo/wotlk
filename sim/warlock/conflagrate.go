package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerConflagrateSpell() {
	actionID := core.ActionID{SpellID: 17962}
	spellSchool := core.SpellSchoolFire
	target := warlock.CurrentTarget
	hasGlyphOfConflag := warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)

	baseCost := 0.16 * warlock.BaseMana
	directFlatDamage := 0.6 * 785 / 5 * float64(warlock.ImmolateDot.NumberOfTicks)
	directSpellCoeff := 0.6 * 0.2 * float64(warlock.ImmolateDot.NumberOfTicks)
	dotFlatDamage := 0.4 / 3 * 785 / 5 * float64(warlock.ImmolateDot.NumberOfTicks)
	dotSpellCoeff := 0.4 / 3 * 0.2 * float64(warlock.ImmolateDot.NumberOfTicks)

	warlock.Conflagrate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskSpellDamage,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm]),
				GCD:  core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 10,
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				if !warlock.ImmolateDot.IsActive() {
					panic("Conflagrate spell is cast while Immolate is not active.")
				}
				if !hasGlyphOfConflag {
					warlock.ImmolateDot.Deactivate(sim)
					//warlock.ShadowflameDot.Deactivate(sim)
				}
			},
		},

		BonusCritRating: 0 +
			warlock.masterDemonologistFireCrit() +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0) +
			5*float64(warlock.Talents.FireAndBrimstone)*core.CritRatingPerCritChance,
		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, false),
		CritMultiplier:           warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier:         1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := directFlatDamage + directSpellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				warlock.ConflagrateDot.Apply(sim)
			}
			spell.DealDamage(sim, &result)
		},
	})

	warlock.ConflagrateDot = core.NewDot(core.Dot{
		Spell: warlock.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: spellSchool,
			ProcMask:    core.ProcMaskSpellDamage,

			BonusCritRating:          warlock.Conflagrate.BonusCritRating,
			DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true),
			CritMultiplier:           warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
			ThreatMultiplier:         warlock.Conflagrate.ThreatMultiplier,
		}),
		Aura: target.RegisterAura(core.Aura{
			Label:    "conflagrate-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 3,
		TickLength:    time.Second * 2,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			IsPeriodic:     true,
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(dotFlatDamage, dotSpellCoeff),
			OutcomeApplier: warlock.OutcomeFuncMagicCrit(),
		}),
	})
}
