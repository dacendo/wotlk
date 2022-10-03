package paladin

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (paladin *Paladin) registerSealOfCommandSpellAndAura() {
	/*
	 * Seal of Command is an Spell/Aura that when active makes the paladin capable of procing
	 * 2 different SpellIDs depending on a paladin's casted spell or melee swing.
	 *
	 * SpellID 20467 (Judgement of Command):
	 *   - Procs off of any "Primary" Judgement (JoL, JoW, JoJ).
	 *   - Cannot miss or be dodged/parried.
	 *   - Deals hybrid AP/SP damage.
	 *   - Crits off of a melee modifier.
	 *
	 * SpellID 20424 (Seal of Command):
	 *   - Procs off of any melee special ability, or white hit.
	 *   - If the ability is SINGLE TARGET, it hits up to 2 extra targets.
	 *   - Deals hybrid AP/SP damage * current weapon speed.
	 *   - Crits off of a melee modifier.
	 *   - CAN MISS, BE DODGED/PARRIED/BLOCKED.
	 */

	baseEffect := core.SpellEffect{
		BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 0, true),
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialHitAndCrit(),
	}

	numHits := core.MinInt32(3, paladin.Env.GetNumTargets()) // primary target + 2 others
	effects := make([]core.SpellEffect, 0, numHits)
	for i := int32(0); i < numHits; i++ {
		mhEffect := baseEffect
		mhEffect.Target = paladin.Env.GetTargetUnit(i)
		effects = append(effects, mhEffect)
	}

	baseMultiplierAdditive := 1 +
		paladin.getItemSetLightswornBattlegearBonus4() +
		paladin.getTalentTwoHandedWeaponSpecializationBonus()

	onSpecialOrSwingActionID := core.ActionID{SpellID: 20424}
	onSpecialOrSwingProcCleave := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    onSpecialOrSwingActionID, // Seal of Command damage bonus for single target spells.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplierAdditive: baseMultiplierAdditive,
		DamageMultiplier:         0.36,
		CritMultiplier:           paladin.MeleeCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: core.ApplyEffectFuncDamageMultiple(effects),
	})

	onSpecialOrSwingProc := paladin.RegisterSpell(core.SpellConfig{
		ActionID:    onSpecialOrSwingActionID, // Seal of Command damage bonus for cleaves.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplierAdditive: baseMultiplierAdditive,
		DamageMultiplier:         0.36,
		CritMultiplier:           paladin.MeleeCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(baseEffect),
	})

	var glyphManaMetrics *core.ResourceMetrics
	glyphManaGain := .08 * paladin.BaseMana
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfSealOfCommand) {
		glyphManaMetrics = paladin.NewManaMetrics(core.ActionID{ItemID: 41094})
	}

	var onJudgementProc *core.Spell

	// Seal of Command aura.
	auraActionID := core.ActionID{SpellID: 20375}
	paladin.SealOfCommandAura = paladin.RegisterAura(core.Aura{
		Label:    "Seal of Command",
		Tag:      "Seal",
		ActionID: auraActionID,
		Duration: SealDuration,

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if glyphManaMetrics != nil && spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				paladin.AddMana(sim, glyphManaGain, glyphManaMetrics, false)
			}

			// Don't proc on misses or our own procs.
			if !spellEffect.Landed() || spell == onJudgementProc || spell.SameAction(onSpecialOrSwingActionID) {
				return
			}

			// Differ between judgements and other melee abilities.
			if spell.Flags.Matches(SpellFlagPrimaryJudgement) {
				onJudgementProc.Cast(sim, spellEffect.Target)
				if paladin.Talents.JudgementsOfTheJust > 0 {
					// Special JoJ talent behavior, procs swing seal on judgements
					// For SoC this is a cleave.
					onSpecialOrSwingProcCleave.Cast(sim, spellEffect.Target)
				}
			} else {
				if spell.IsMelee() {
					// Temporary check to avoid AOE double procing.
					if spell.SpellID == paladin.HammerOfTheRighteous.SpellID || spell.SpellID == paladin.DivineStorm.SpellID {
						onSpecialOrSwingProc.Cast(sim, spellEffect.Target)
					} else {
						onSpecialOrSwingProcCleave.Cast(sim, spellEffect.Target)
					}
				}
			}
		},
	})

	aura := paladin.SealOfCommandAura
	baseCost := paladin.BaseMana * 0.14
	paladin.SealOfCommand = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    auraActionID, // Seal of Command self buff.
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskEmpty,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(paladin.Talents.Benediction)),
				GCD:  core.GCDDefault,
			},
		},

		BonusCritRating: (6 * float64(paladin.Talents.Fanaticism) * core.CritRatingPerCritChance) +
			(core.TernaryFloat64(paladin.HasSetBonus(ItemSetTuralyonsBattlegear, 4) || paladin.HasSetBonus(ItemSetLiadrinsBattlegear, 4), 5, 0) * core.CritRatingPerCritChance),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			if paladin.CurrentSeal != nil {
				paladin.CurrentSeal.Deactivate(sim)
			}
			paladin.CurrentSeal = aura
			paladin.CurrentSeal.Activate(sim)
		},
	})

	onJudgementProc = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20467}, // Judgement of Command
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeOrRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagSecondaryJudgement,

		DamageMultiplierAdditive: baseMultiplierAdditive +
			paladin.getMajorGlyphOfJudgementBonus() +
			paladin.getTalentTheArtOfWarBonus(),
		DamageMultiplier: 1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return 0.19*core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 0, true).Calculator(sim, hitEffect, spell) +
						0.08*spell.MeleeAttackPower() +
						0.13*spell.SpellPower()
				},
			},
			OutcomeApplier: paladin.OutcomeFuncMeleeSpecialCritOnly(), // Secondary Judgements cannot miss if the Primary Judgement hit, only roll for crit.
		}),
	})
}
