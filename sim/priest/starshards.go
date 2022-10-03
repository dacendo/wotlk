package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (priest *Priest) registerStarshardsSpell() {
	actionID := core.ActionID{SpellID: 25446}
	priest.Starshards = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			OutcomeApplier: priest.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					priest.StarshardsDot.Apply(sim)
				}
			},
		}),
	})

	target := priest.CurrentTarget
	priest.StarshardsDot = core.NewDot(core.Dot{
		Spell: priest.Starshards,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Starshards-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks: 5,
		TickLength:    time.Second * 3,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			IsPeriodic:     true,
			BaseDamage:     core.BaseDamageConfigMagicNoRoll(785/5, 0.167),
			OutcomeApplier: priest.OutcomeFuncTick(),
		}),
	})
}
