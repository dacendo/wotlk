package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) channelCheck(sim *core.Simulation, dot *core.Dot, maxTicks int) *core.Spell {
	if dot.IsActive() && dot.TickCount+1 < maxTicks {
		return warlock.DrainSoulChannelling
	} else {
		return warlock.DrainSoul
	}
}

func (warlock *Warlock) registerDrainSoulSpell() {
	actionID := core.ActionID{SpellID: 47855}
	spellSchool := core.SpellSchoolShadow
	soulSiphonMultiplier := 0.03 * float64(warlock.Talents.SoulSiphon)
	baseCost := warlock.BaseMana * 0.14
	channelTime := 3 * time.Second
	epsilon := 1 * time.Millisecond

	warlock.DrainSoul = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  spellSchool,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        core.SpellFlagChanneled,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
		},

		DamageMultiplierAdditive: warlock.staticAdditiveDamageMultiplier(actionID, spellSchool, true),
		// For performance optimization, the execute modifier is basekit since we never use it before execute
		DamageMultiplier: (4.0 + 0.04*float64(warlock.Talents.DeathsEmbrace)) / (1 + 0.04*float64(warlock.Talents.DeathsEmbrace)),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: warlock.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				warlock.DrainSoulDot.Apply(sim)
				warlock.DrainSoulDot.Aura.UpdateExpires(warlock.DrainSoulDot.Aura.ExpiresAt() + epsilon)
			},
		}),
	})

	target := warlock.CurrentTarget

	effect := core.SpellEffect{
		IsPeriodic:     true,
		OutcomeApplier: warlock.OutcomeFuncTick(),
		BaseDamage: core.WrapBaseDamageConfig(core.BaseDamageConfigMagicNoRoll(142, 0.429),
			func(oldCalc core.BaseDamageCalculator) core.BaseDamageCalculator {
				return func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					w := warlock
					auras := []*core.Aura{
						w.HauntDebuffAura(target),
						w.UnstableAfflictionDot.Aura, w.CorruptionDot.Aura,
						w.SeedDots[hitEffect.Target.Index].Aura, w.CurseOfAgonyDot.Aura,
						w.CurseOfDoomDot.Aura, w.CurseOfElementsAura, w.CurseOfWeaknessAura,
						w.CurseOfTonguesAura,
						// missing: death coil
					}
					counter := 0
					for _, aura := range auras {
						if aura.IsActive() {
							counter++
						}
					}
					return oldCalc(sim, hitEffect, spell) *
						(1.0 + float64(core.MinInt(3, counter))*soulSiphonMultiplier)
				}
			}),
	}

	warlock.DrainSoulDot = core.NewDot(core.Dot{
		Spell: warlock.DrainSoul,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Drain Soul-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       1,
		TickLength:          3 * time.Second,
		AffectedByCastSpeed: true,

		TickEffects: core.TickFuncSnapshot(target, effect),
	})

	warlock.DrainSoulChannelling = warlock.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagNoLogs | core.SpellFlagNoMetrics,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
				CastTime:    0,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.DrainSoulDot.Apply(sim) // TODO: do we want to just refresh and continue ticking with same snapshot or update snapshot?
			warlock.DrainSoulDot.Aura.UpdateExpires(warlock.DrainSoulDot.Aura.ExpiresAt() + epsilon)
		},
	})
}
