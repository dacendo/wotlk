package shaman

import (
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var TotemOfTheAstralWinds int32 = 27815
var TotemOfSplintering int32 = 40710

func (shaman *Shaman) newWindfuryImbueSpell(isMH bool) *core.Spell {
	apBonus := 1250.0
	if shaman.Equip[proto.ItemSlot_ItemSlotRanged].ID == TotemOfTheAstralWinds {
		apBonus += 80
	} else if shaman.Equip[proto.ItemSlot_ItemSlotRanged].ID == TotemOfSplintering {
		apBonus += 212
	}

	actionID := core.ActionID{SpellID: 58804}

	baseEffect := core.SpellEffect{
		OutcomeApplier: shaman.OutcomeFuncMeleeSpecialHitAndCrit(),
	}
	if isMH {
		bonusApDmg := shaman.AutoAttacks.MH.SwingSpeed * apBonus / core.MeleeAttackRatingPerDamage
		actionID.Tag = 1
		baseEffect.BaseDamage = core.BaseDamageConfigMeleeWeapon(core.MainHand, false, bonusApDmg, true)
	} else {
		bonusApDmg := shaman.AutoAttacks.OH.SwingSpeed * apBonus / core.MeleeAttackRatingPerDamage
		actionID.Tag = 2
		baseEffect.BaseDamage = core.BaseDamageConfigMeleeWeapon(core.OffHand, false, bonusApDmg, true)
	}

	effects := []core.SpellEffect{
		baseEffect,
		baseEffect,
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMelee,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		DamageMultiplier: 1 + math.Round(float64(shaman.Talents.ElementalWeapons)*13.33)/100,
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDamageMultipleTargeted(effects),
	})
}

func (shaman *Shaman) ApplyWindfuryImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	var proc = 0.2
	if mh && oh {
		proc = 0.36
	}
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfWindfuryWeapon) {
		proc += 0.02 //TODO: confirm how this actually works
	}

	mhSpell := shaman.newWindfuryImbueSpell(true)
	ohSpell := shaman.newWindfuryImbueSpell(false)

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Second * 3,
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Windfury Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// ProcMask: 20
			if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			isMHHit := spell.IsMH()
			if (!mh && isMHHit) || (!oh && !isMHHit) {
				return // cant proc if not enchanted
			}
			if !icd.IsReady(sim) {
				return
			}
			if sim.RandomFloat("Windfury Imbue") > proc {
				return
			}
			icd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, spellEffect.Target)
			} else {
				ohSpell.Cast(sim, spellEffect.Target)
			}
		},
	})
}

func (shaman *Shaman) newFlametongueImbueSpell(isMH bool) *core.Spell {
	var baseDamage float64
	var spellCoeff float64
	if isMH {
		if weapon := shaman.GetMHWeapon(); weapon != nil {
			baseDamage = weapon.SwingSpeed * 68.5
			spellCoeff = 0.1 * weapon.SwingSpeed / 2.6
		}
	} else {
		if weapon := shaman.GetOHWeapon(); weapon != nil {
			baseDamage = weapon.SwingSpeed * 68.5
			spellCoeff = 0.1 * weapon.SwingSpeed / 2.6
		}
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58790},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * 1 * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseDamage + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamageMagicHitAndCrit(sim, target, baseDamage)
		},
	})
}

func (shaman *Shaman) ApplyFlametongueImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	imbueCount := 1.0
	spBonus := 211.0
	spMod := 1.0 + 0.1*float64(shaman.Talents.ElementalWeapons)
	if mh && oh { // grant double SP+Crit bonuses for ft/ft (possible bug, but currently working on beta, its unclear)
		imbueCount += 1.0
	}
	shaman.AddStat(stats.SpellPower, spBonus*spMod*imbueCount)
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon) {
		shaman.AddStat(stats.SpellCrit, 2*core.CritRatingPerCritChance*imbueCount)
	}

	ftIcd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueImbueSpell(true)
	ohSpell := shaman.newFlametongueImbueSpell(false)

	shaman.RegisterAura(core.Aura{
		Label:    "Flametongue Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			isMHHit := spell.IsMH()
			if (isMHHit && !mh) || (!isMHHit && !oh) {
				return // cant proc if not enchanted
			}
			if !ftIcd.IsReady(sim) {
				return
			}
			ftIcd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, spellEffect.Target)
			} else {
				ohSpell.Cast(sim, spellEffect.Target)
			}
		},
	})
}

func (shaman *Shaman) newFlametongueDownrankImbueSpell(isMH bool) *core.Spell {
	var baseDamage float64
	var spellCoeff float64
	if isMH {
		if weapon := shaman.GetMHWeapon(); weapon != nil {
			baseDamage = weapon.SwingSpeed * 64
			spellCoeff = 0.1 * weapon.SwingSpeed / 2.6
		}
	} else {
		if weapon := shaman.GetOHWeapon(); weapon != nil {
			baseDamage = weapon.SwingSpeed * 64
			spellCoeff = 0.1 * weapon.SwingSpeed / 2.6
		}
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58789},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * 1 * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := baseDamage + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamageMagicHitAndCrit(sim, target, baseDamage)
		},
	})
}

func (shaman *Shaman) ApplyFlametongueDownrankImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	imbueCount := 1.0
	spBonus := 186.0
	spMod := 1.0 + 0.1*float64(shaman.Talents.ElementalWeapons)
	if mh && oh { // grant double SP+Crit bonuses for ft/ft (possible bug, but currently working on beta, its unclear)
		imbueCount += 1.0
	}
	shaman.AddStat(stats.SpellPower, spBonus*spMod*imbueCount)
	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFlametongueWeapon) {
		shaman.AddStat(stats.SpellCrit, 2*core.CritRatingPerCritChance*imbueCount)
	}

	ftDownrankIcd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond,
	}

	mhSpell := shaman.newFlametongueDownrankImbueSpell(true)
	ohSpell := shaman.newFlametongueDownrankImbueSpell(false)

	shaman.RegisterAura(core.Aura{
		Label:    "Flametongue Imbue (downranked)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			isMHHit := spell.IsMH()
			if (isMHHit && !mh) || (!isMHHit && !oh) {
				return // cant proc if not enchanted
			}
			if !ftDownrankIcd.IsReady(sim) {
				return
			}
			ftDownrankIcd.Use(sim)

			if isMHHit {
				mhSpell.Cast(sim, spellEffect.Target)
			} else {
				ohSpell.Cast(sim, spellEffect.Target)
			}
		},
	})
}

func (shaman *Shaman) FrostbrandDebuffAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Frostbrand Attack-" + shaman.Label,
		ActionID: core.ActionID{SpellID: 58799},
		Duration: time.Second * 8,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell.Unit != &shaman.Unit {
				return
			}
			//if spell != shaman.LightningBolt || shaman.ChainLightning || shaman.FlameShock || shaman.FrostShock || shaman.EarthShock || shaman.LavaLash {
			//	return
			//}
			//modifier goes here, maybe?? might need to go in the specific spell files i guess
		},
		//TODO: figure out how to implement frozen power (might not be here)
	})
}

func (shaman *Shaman) newFrostbrandImbueSpell(isMH bool) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 58796},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskEmpty,

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * 1 * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 530 + 0.1*spell.SpellPower()
			spell.CalcAndDealDamageMagicHitAndCrit(sim, target, baseDamage)
		},
	})
}

func (shaman *Shaman) ApplyFrostbrandImbue(mh bool, oh bool) {
	if !mh && !oh {
		return
	}

	mhSpell := shaman.newFrostbrandImbueSpell(true)
	ohSpell := shaman.newFrostbrandImbueSpell(false)
	procMask := core.GetMeleeProcMaskForHands(mh, oh)
	ppmm := shaman.AutoAttacks.NewPPMManager(9.0, procMask)

	shaman.RegisterAura(core.Aura{
		Label:    "Frostbrand Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spell.ProcMask.Matches(procMask) {
				return
			}

			if !ppmm.Proc(sim, spell.ProcMask, "Frostbrand Weapon") {
				return
			}

			if spell.IsMH() {
				mhSpell.Cast(sim, spellEffect.Target)
			} else {
				ohSpell.Cast(sim, spellEffect.Target)
			}
			shaman.FrostbrandDebuffAura(shaman.CurrentTarget).Activate(sim)
		},
	})
}

//earthliving? not important for dps sims though
