package warrior

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerExecuteSpell() {
	cost := 15.0 - float64(warrior.Talents.FocusedRage)
	if warrior.Talents.ImprovedExecute == 1 {
		cost -= 2
	} else if warrior.Talents.ImprovedExecute == 2 {
		cost -= 5
	}
	if warrior.HasSetBonus(ItemSetOnslaughtBattlegear, 2) {
		cost -= 3
	}
	refundAmount := cost * 0.8
	gcd := core.GCDDefault - core.TernaryDuration(warrior.HasSetBonus(ItemSetYmirjarLordsBattlegear, 4), 500, 0)*time.Millisecond

	var extraRage float64
	extraRageBonus := core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfExecution), 10, 0)

	warrior.Execute = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47471},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  gcd,
			},
			IgnoreHaste: true,
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.Cost = math.Min(spell.Unit.CurrentRage(), 30)
				extraRage = cast.Cost - spell.BaseCost
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1.25,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return 1456 + 0.2*spell.MeleeAttackPower() + 38*(extraRage+extraRageBonus)
				},
			},
			OutcomeApplier: warrior.OutcomeFuncMeleeSpecialHitAndCrit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) SpamExecute(spam bool) bool {
	return warrior.CurrentRage() >= warrior.Execute.BaseCost && spam && warrior.Talents.MortalStrike
}

func (warrior *Warrior) CanExecute() bool {
	return warrior.CurrentRage() >= warrior.Execute.BaseCost
}

func (warrior *Warrior) CanSuddenDeathExecute() bool {
	return warrior.CurrentRage() >= warrior.Execute.BaseCost && warrior.SuddenDeathAura.IsActive()
}

func (warrior *Warrior) CastExecute(sim *core.Simulation, target *core.Unit) bool {
	if warrior.Ymirjar4pcProcAura.IsActive() && warrior.SuddenDeathAura.IsActive() {
		warrior.Execute.DefaultCast.GCD = time.Second * 1
		warrior.Ymirjar4pcProcAura.RemoveStack(sim)
	} else {
		warrior.Execute.DefaultCast.GCD = core.GCDDefault
	}

	return warrior.Execute.Cast(sim, target)
}
