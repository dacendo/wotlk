package rogue

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (rogue *Rogue) registerHackAndSlash(mask core.ProcMask) {
	if rogue.Talents.HackAndSlash < 1 || mask == core.ProcMaskUnknown {
		return
	}
	var hackAndSlashSpell *core.Spell
	icd := core.Cooldown{
		Timer:    rogue.NewTimer(),
		Duration: time.Millisecond * 500,
	}
	procChance := 0.01 * float64(rogue.Talents.HackAndSlash)
	rogue.RegisterAura(core.Aura{
		Label:    "Hack and Slash",
		Duration: core.NeverExpires,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			hackAndSlashSpell = rogue.GetOrRegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 13964},
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskMeleeMHAuto,
				Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

				DamageMultiplier: rogue.AutoAttacks.MHConfig.DamageMultiplier,
				CritMultiplier:   rogue.MeleeCritMultiplier(true, false),
				ThreatMultiplier: rogue.AutoAttacks.MHConfig.ThreatMultiplier,

				ApplyEffects: rogue.AutoAttacks.MHConfig.ApplyEffects,
			})
		},
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			if !spell.ProcMask.Matches(mask) {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			if sim.RandomFloat("Sword Specialization") > procChance {
				return
			}
			icd.Use(sim)
			hackAndSlashSpell.Cast(sim, spellEffect.Target)
		},
	})
}
