package core

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/wotlk/sim/core/stats"
)

// Callback for after a spell hits the target and after damage is calculated. Use it for proc effects
// or anything that comes from the final result of the spell.
type OnSpellHit func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect)
type EffectOnSpellHitDealt func(sim *Simulation, spell *Spell, spellEffect *SpellEffect)

// OnPeriodicDamage is called when dots tick, after damage is calculated. Use it for proc effects
// or anything that comes from the final result of a tick.
type OnPeriodicDamage func(aura *Aura, sim *Simulation, spell *Spell, spellEffect *SpellEffect)
type EffectOnPeriodicDamageDealt func(sim *Simulation, spell *Spell, spellEffect *SpellEffect)

type SpellEffect struct {
	// Target of the spell.
	Target *Unit

	BaseDamage     BaseDamageConfig // Callback for calculating base damage.
	OutcomeApplier OutcomeApplier   // Callback for determining outcome.

	// Used only for dot snapshotting. Internal-only.
	snapshotDamageMultiplier float64
	snapshotMeleeCritRating  float64
	snapshotSpellCritRating  float64

	// Use snapshotted values for crit/damage rather than recomputing them.
	isSnapshot bool

	// Used in determining snapshot based damage from effect details (e.g. snapshot crit and % damage modifiers)
	IsPeriodic bool

	// Indicates this is a healing spell, rather than a damage spell.
	IsHealing bool

	// Callbacks for providing additional custom behavior.
	OnSpellHitDealt       EffectOnSpellHitDealt
	OnPeriodicDamageDealt EffectOnPeriodicDamageDealt

	// Results
	Outcome HitOutcome
	Damage  float64 // Damage done by this cast.

	PreoutcomeDamage float64 // Damage done by this cast before Outcome is applied.
}

func (spellEffect *SpellEffect) Validate() {
	//if spell.ProcMask == ProcMaskUnknown {
	//	panic("SpellEffects must set a ProcMask!")
	//}
	//if spell.ProcMask.Matches(ProcMaskEmpty) && spell.ProcMask != ProcMaskEmpty {
	//	panic("ProcMaskEmpty must be exclusive!")
	//}
}

func (spellEffect *SpellEffect) Landed() bool {
	return spellEffect.Outcome.Matches(OutcomeLanded)
}

func (spellEffect *SpellEffect) DidCrit() bool {
	return spellEffect.Outcome.Matches(OutcomeCrit)
}

func (spellEffect *SpellEffect) calcThreat(spell *Spell) float64 {
	if spellEffect.Landed() {
		flatBonus := spell.FlatThreatBonus
		if spell.DynamicThreatBonus != nil {
			flatBonus += spell.DynamicThreatBonus(spellEffect, spell)
		}

		return (spellEffect.Damage*spell.ThreatMultiplier + flatBonus) * spell.Unit.PseudoStats.ThreatMultiplier
	} else {
		return 0
	}
}

func (spell *Spell) MeleeAttackPower() float64 {
	return spell.Unit.stats[stats.AttackPower] + spell.Unit.PseudoStats.MobTypeAttackPower
}

func (spell *Spell) RangedAttackPower(target *Unit) float64 {
	return spell.Unit.stats[stats.RangedAttackPower] +
		spell.Unit.PseudoStats.MobTypeAttackPower +
		target.PseudoStats.BonusRangedAttackPowerTaken
}

func (spell *Spell) BonusWeaponDamage() float64 {
	return spell.Unit.PseudoStats.BonusDamage
}

func (spell *Spell) ExpertisePercentage() float64 {
	expertiseRating := spell.Unit.stats[stats.Expertise] + spell.BonusExpertiseRating
	return math.Floor(expertiseRating/ExpertisePerQuarterPercentReduction) / 400
}

func (spell *Spell) PhysicalHitChance(target *Unit) float64 {
	hitRating := spell.Unit.stats[stats.MeleeHit] +
		spell.BonusHitRating +
		target.PseudoStats.BonusMeleeHitRatingTaken
	return hitRating / (MeleeHitRatingPerHitChance * 100)
}

func (spellEffect *SpellEffect) PhysicalCritChance(unit *Unit, spell *Spell, attackTable *AttackTable) float64 {
	critRating := 0.0
	if spellEffect.isSnapshot {
		// periodic spells apply crit from snapshot at time of initial cast if capable of a crit
		// ignoring units real time crit in this case
		critRating = spellEffect.snapshotMeleeCritRating
	} else {
		critRating = spell.physicalCritRating(spellEffect.Target)
	}

	return (critRating / (CritRatingPerCritChance * 100)) - attackTable.CritSuppression
}
func (spell *Spell) physicalCritRating(target *Unit) float64 {
	return spell.Unit.stats[stats.MeleeCrit] +
		spell.BonusCritRating +
		target.PseudoStats.BonusCritRatingTaken
}

func (spell *Spell) SpellPower() float64 {
	return spell.Unit.GetStat(stats.SpellPower) +
		spell.BonusSpellPower +
		spell.Unit.PseudoStats.MobTypeSpellPower
}

func (spell *Spell) SpellHitChance(target *Unit) float64 {
	hitRating := spell.Unit.stats[stats.SpellHit] +
		spell.BonusHitRating +
		target.PseudoStats.BonusSpellHitRatingTaken

	return hitRating / (SpellHitRatingPerHitChance * 100)
}

func (spellEffect *SpellEffect) SpellCritChance(unit *Unit, spell *Spell) float64 {
	critRating := 0.0
	if spellEffect.isSnapshot {
		// periodic spells apply crit from snapshot at time of initial cast if capable of a crit
		// ignoring units real time crit in this case
		critRating = spellEffect.snapshotSpellCritRating
	} else {
		critRating = spell.spellCritRating(spellEffect.Target)
	}

	return critRating / (CritRatingPerCritChance * 100)
}
func (spell *Spell) spellCritRating(target *Unit) float64 {
	return spell.Unit.stats[stats.SpellCrit] +
		spell.BonusCritRating +
		target.PseudoStats.BonusCritRatingTaken +
		target.PseudoStats.BonusSpellCritRatingTaken
}

func (spell *Spell) HealingPower() float64 {
	return spell.SpellPower()
}
func (spellEffect *SpellEffect) HealingCritChance(unit *Unit, spell *Spell) float64 {
	critRating := 0.0
	if spellEffect.isSnapshot {
		// periodic spells apply crit from snapshot at time of initial cast if capable of a crit
		// ignoring units real time crit in this case
		critRating = spellEffect.snapshotSpellCritRating
	} else {
		critRating = spell.healingCritRating()
	}

	return critRating / (CritRatingPerCritChance * 100)
}
func (spell *Spell) healingCritRating() float64 {
	return spell.Unit.GetStat(stats.SpellCrit) + spell.BonusCritRating
}

func (spell *Spell) ApplyPostOutcomeDamageModifiers(sim *Simulation, result *SpellEffect) {
	for i := range result.Target.DynamicDamageTakenModifiers {
		result.Target.DynamicDamageTakenModifiers[i](sim, spell, result)
	}
	result.Damage = MaxFloat(0, result.Damage)
}

// For spells that do no damage but still have a hit/miss check.
func (spell *Spell) CalcOutcome(sim *Simulation, target *Unit, outcomeApplier NewOutcomeApplier) SpellEffect {
	attackTable := spell.Unit.AttackTables[target.UnitIndex]
	result := SpellEffect{
		Target: target,
		Damage: 0,
	}
	outcomeApplier(sim, &result, attackTable)
	return result
}

func (spell *Spell) CalcDamage(sim *Simulation, target *Unit, baseDamage float64, outcomeApplier NewOutcomeApplier) SpellEffect {
	attackTable := spell.Unit.AttackTables[target.UnitIndex]
	result := SpellEffect{
		Target: target,
		Damage: baseDamage,
	}

	result.Damage *= spell.CasterDamageMultiplier()
	result.Damage *= spell.ResistanceMultiplier(sim, false, attackTable)
	result.Damage = spell.applyTargetModifiers(result.Damage, attackTable, false)
	result.PreoutcomeDamage = result.Damage
	outcomeApplier(sim, &result, attackTable)
	spell.ApplyPostOutcomeDamageModifiers(sim, &result)

	return result
}

func (spell *Spell) DealOutcome(sim *Simulation, result *SpellEffect) {
	spell.DealDamage(sim, result)
}

// Applies the fully computed spell result to the sim.
func (spell *Spell) DealDamage(sim *Simulation, result *SpellEffect) {
	spell.SpellMetrics[result.Target.UnitIndex].TotalDamage += result.Damage
	spell.SpellMetrics[result.Target.UnitIndex].TotalThreat += result.calcThreat(spell)

	// Mark total damage done in raid so far for health based fights.
	// Don't include damage done by EnemyUnits to Players
	if result.Target.Type == EnemyUnit {
		sim.Encounter.DamageTaken += result.Damage
	}

	if sim.Log != nil {
		spell.Unit.Log(sim, "%s %s %s. (Threat: %0.3f)", result.Target.LogLabel(), spell.ActionID, result, result.calcThreat(spell))
	}

	spell.Unit.OnSpellHitDealt(sim, spell, result)
	result.Target.OnSpellHitTaken(sim, spell, result)
}

func (spell *Spell) CalcHealing(sim *Simulation, target *Unit, baseHealing float64, outcomeApplier NewOutcomeApplier) SpellEffect {
	attackTable := spell.Unit.AttackTables[target.UnitIndex]
	result := SpellEffect{
		Target: target,
		Damage: baseHealing,
	}

	result.Damage *= spell.CasterHealingMultiplier()
	result.Damage = spell.applyTargetHealingModifiers(result.Damage, attackTable)
	outcomeApplier(sim, &result, attackTable)

	return result
}

// Applies the fully computed spell result to the sim.
func (spell *Spell) DealHealing(sim *Simulation, result *SpellEffect) {
	spell.SpellMetrics[result.Target.UnitIndex].TotalHealing += result.Damage
	spell.SpellMetrics[result.Target.UnitIndex].TotalThreat += result.calcThreat(spell)
	if result.Target.HasHealthBar() {
		result.Target.GainHealth(sim, result.Damage, spell.HealthMetrics(result.Target))
	}

	if sim.Log != nil {
		spell.Unit.Log(sim, "%s %s %s. (Threat: %0.3f)", result.Target.LogLabel(), spell.ActionID, result, result.calcThreat(spell))
	}

	spell.Unit.OnHealDealt(sim, spell, result)
	result.Target.OnHealTaken(sim, spell, result)
}

func (spell *Spell) WaitTravelTime(sim *Simulation, callback func(*Simulation)) {
	travelTime := time.Duration(float64(time.Second) * spell.Unit.DistanceFromTarget / spell.MissileSpeed)
	StartDelayedAction(sim, DelayedActionOptions{
		DoAt:     sim.CurrentTime + travelTime,
		OnAction: callback,
	})
}

func (spellEffect *SpellEffect) calculateBaseDamage(sim *Simulation, spell *Spell) float64 {
	if spellEffect.BaseDamage.Calculator == nil {
		return 0
	} else {
		return spellEffect.BaseDamage.Calculator(sim, spellEffect, spell)
	}
}

func (spellEffect *SpellEffect) calcDamageSingle(sim *Simulation, spell *Spell, attackTable *AttackTable) {
	if sim.Log != nil {
		baseDmg := spellEffect.Damage
		spellEffect.applyAttackerModifiers(sim, spell)
		afterAttackMods := spellEffect.Damage
		spellEffect.applyResistances(sim, spell, attackTable)
		afterResistances := spellEffect.Damage
		spellEffect.applyTargetModifiers(spell, attackTable)
		afterTargetMods := spellEffect.Damage
		spellEffect.PreoutcomeDamage = spellEffect.Damage
		spellEffect.OutcomeApplier(sim, spell, spellEffect, attackTable)
		afterOutcome := spellEffect.Damage
		spell.Unit.Log(
			sim,
			"%s %s [DEBUG] MAP: %0.01f, RAP: %0.01f, SP: %0.01f, BaseDamage:%0.01f, AfterAttackerMods:%0.01f, AfterResistances:%0.01f, AfterTargetMods:%0.01f, AfterOutcome:%0.01f",
			spellEffect.Target.LogLabel(), spell.ActionID, spell.Unit.GetStat(stats.AttackPower), spell.Unit.GetStat(stats.RangedAttackPower), spell.Unit.GetStat(stats.SpellPower), baseDmg, afterAttackMods, afterResistances, afterTargetMods, afterOutcome)
	} else {
		spellEffect.applyAttackerModifiers(sim, spell)
		spellEffect.applyResistances(sim, spell, attackTable)
		spellEffect.applyTargetModifiers(spell, attackTable)
		spellEffect.PreoutcomeDamage = spellEffect.Damage
		spellEffect.OutcomeApplier(sim, spell, spellEffect, attackTable)
	}
}
func (spellEffect *SpellEffect) calcDamageTargetOnly(sim *Simulation, spell *Spell, attackTable *AttackTable) {
	spellEffect.applyResistances(sim, spell, attackTable)
	spellEffect.applyTargetModifiers(spell, attackTable)
	spellEffect.OutcomeApplier(sim, spell, spellEffect, attackTable)
}

func (spellEffect *SpellEffect) finalize(sim *Simulation, spell *Spell) {
	if spell.MissileSpeed == 0 {
		if spellEffect.IsHealing {
			spellEffect.finalizeHealingInternal(sim, spell)
		} else {
			spellEffect.finalizeInternal(sim, spell)
		}
	} else {
		travelTime := time.Duration(float64(time.Second) * spell.Unit.DistanceFromTarget / spell.MissileSpeed)

		// We need to make a copy of this SpellEffect because some spells re-use the effect objects.
		effectCopy := *spellEffect

		StartDelayedAction(sim, DelayedActionOptions{
			DoAt: sim.CurrentTime + travelTime,
			OnAction: func(sim *Simulation) {
				if spellEffect.IsHealing {
					effectCopy.finalizeHealingInternal(sim, spell)
				} else {
					effectCopy.finalizeInternal(sim, spell)
				}
			},
		})
	}
}

// Applies the fully computed results from this SpellEffect to the sim.
func (spellEffect *SpellEffect) finalizeInternal(sim *Simulation, spell *Spell) {
	spell.ApplyPostOutcomeDamageModifiers(sim, spellEffect)

	spell.SpellMetrics[spellEffect.Target.UnitIndex].TotalDamage += spellEffect.Damage
	spell.SpellMetrics[spellEffect.Target.UnitIndex].TotalThreat += spellEffect.calcThreat(spell)

	// Mark total damage done in raid so far for health based fights.
	// Don't include damage done by EnemyUnits to Players
	if spellEffect.Target.Type == EnemyUnit {
		sim.Encounter.DamageTaken += spellEffect.Damage
	}

	if sim.Log != nil {
		if spellEffect.IsPeriodic {
			spell.Unit.Log(sim, "%s %s tick %s. (Threat: %0.3f)", spellEffect.Target.LogLabel(), spell.ActionID, spellEffect, spellEffect.calcThreat(spell))
		} else {
			spell.Unit.Log(sim, "%s %s %s. (Threat: %0.3f)", spellEffect.Target.LogLabel(), spell.ActionID, spellEffect, spellEffect.calcThreat(spell))
		}
	}

	if !spellEffect.IsPeriodic {
		if spellEffect.OnSpellHitDealt != nil {
			spellEffect.OnSpellHitDealt(sim, spell, spellEffect)
		}
		spell.Unit.OnSpellHitDealt(sim, spell, spellEffect)
		spellEffect.Target.OnSpellHitTaken(sim, spell, spellEffect)
	} else {
		if spellEffect.OnPeriodicDamageDealt != nil {
			spellEffect.OnPeriodicDamageDealt(sim, spell, spellEffect)
		}
		spell.Unit.OnPeriodicDamageDealt(sim, spell, spellEffect)
		spellEffect.Target.OnPeriodicDamageTaken(sim, spell, spellEffect)
	}
}

func (spellEffect *SpellEffect) finalizeHealingInternal(sim *Simulation, spell *Spell) {
	spell.SpellMetrics[spellEffect.Target.UnitIndex].TotalThreat += spellEffect.calcThreat(spell)
	spell.SpellMetrics[spellEffect.Target.UnitIndex].TotalHealing += spellEffect.Damage
	if spellEffect.Target.HasHealthBar() {
		spellEffect.Target.GainHealth(sim, spellEffect.Damage, spell.HealthMetrics(spellEffect.Target))
	}

	if sim.Log != nil {
		if spellEffect.IsPeriodic {
			spell.Unit.Log(sim, "%s %s tick %s. (Threat: %0.3f)", spellEffect.Target.LogLabel(), spell.ActionID, spellEffect, spellEffect.calcThreat(spell))
		} else {
			spell.Unit.Log(sim, "%s %s %s. (Threat: %0.3f)", spellEffect.Target.LogLabel(), spell.ActionID, spellEffect, spellEffect.calcThreat(spell))
		}
	}

	if !spellEffect.IsPeriodic {
		if spellEffect.OnSpellHitDealt != nil {
			spellEffect.OnSpellHitDealt(sim, spell, spellEffect)
		}
		spell.Unit.OnHealDealt(sim, spell, spellEffect)
		spellEffect.Target.OnHealTaken(sim, spell, spellEffect)
	} else {
		if spellEffect.OnPeriodicDamageDealt != nil {
			spellEffect.OnPeriodicDamageDealt(sim, spell, spellEffect)
		}
		spell.Unit.OnPeriodicHealDealt(sim, spell, spellEffect)
		spellEffect.Target.OnPeriodicHealTaken(sim, spell, spellEffect)
	}
}

func (spellEffect *SpellEffect) String() string {
	outcomeStr := spellEffect.Outcome.String()
	if !spellEffect.Landed() {
		return outcomeStr
	}
	if spellEffect.IsHealing {
		return fmt.Sprintf("%s for %0.3f healing", outcomeStr, spellEffect.Damage)
	} else {
		return fmt.Sprintf("%s for %0.3f damage", outcomeStr, spellEffect.Damage)
	}
}

func (spellEffect *SpellEffect) applyAttackerModifiers(sim *Simulation, spell *Spell) {
	if spell.Flags.Matches(SpellFlagIgnoreAttackerModifiers) {
		// Even when ignoring attacker multipliers we still apply this one, because its
		// specific to the spell.
		spellEffect.Damage *= spell.DamageMultiplier * spell.DamageMultiplierAdditive
		return
	}

	// For dot snapshots, everything has already been stored in spellEffect.snapshotDamageMultiplier.
	if spellEffect.isSnapshot {
		spellEffect.Damage *= spellEffect.snapshotDamageMultiplier
		return
	}

	attacker := spell.Unit

	if spellEffect.IsHealing {
		spellEffect.Damage *= attacker.PseudoStats.HealingDealtMultiplier
		return
	}

	spellEffect.Damage *= spell.CasterDamageMultiplier()
}

// Returns the combined attacker modifiers. For snapshot dots, these are precomputed and stored.
func (spell *Spell) CasterDamageMultiplier() float64 {
	if spell.Flags.Matches(SpellFlagIgnoreAttackerModifiers) {
		return 1
	}

	attacker := spell.Unit
	multiplier := attacker.PseudoStats.DamageDealtMultiplier
	multiplier *= spell.DamageMultiplier
	multiplier *= spell.DamageMultiplierAdditive

	if spell.Flags.Matches(SpellFlagDisease) {
		multiplier *= attacker.PseudoStats.DiseaseDamageDealtMultiplier
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		multiplier *= attacker.PseudoStats.PhysicalDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolArcane) {
		multiplier *= attacker.PseudoStats.ArcaneDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFire) {
		multiplier *= attacker.PseudoStats.FireDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFrost) {
		multiplier *= attacker.PseudoStats.FrostDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolHoly) {
		multiplier *= attacker.PseudoStats.HolyDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolNature) {
		multiplier *= attacker.PseudoStats.NatureDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolShadow) {
		multiplier *= attacker.PseudoStats.ShadowDamageDealtMultiplier
	}

	return multiplier
}

func (spellEffect *SpellEffect) applyTargetModifiers(spell *Spell, attackTable *AttackTable) {
	if spellEffect.IsHealing && !spell.Flags.Matches(SpellFlagIgnoreTargetModifiers) {
		spellEffect.Damage *= attackTable.Defender.PseudoStats.HealingTakenMultiplier * attackTable.HealingDealtMultiplier
		return
	}

	spellEffect.Damage = spell.applyTargetModifiers(spellEffect.Damage, attackTable, spellEffect.IsPeriodic)
}
func (spell *Spell) applyTargetModifiers(damage float64, attackTable *AttackTable, isPeriodic bool) float64 {
	if spell.Flags.Matches(SpellFlagIgnoreTargetModifiers) {
		return damage
	}

	target := attackTable.Defender
	damage *= attackTable.DamageDealtMultiplier
	damage *= target.PseudoStats.DamageTakenMultiplier

	if spell.Flags.Matches(SpellFlagDisease) {
		damage *= target.PseudoStats.DiseaseDamageTakenMultiplier
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		if spell.Flags.Matches(SpellFlagIncludeTargetBonusDamage) {
			damage += target.PseudoStats.BonusPhysicalDamageTaken
		}
		if isPeriodic {
			damage *= target.PseudoStats.PeriodicPhysicalDamageTakenMultiplier
		}
		damage *= target.PseudoStats.PhysicalDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolArcane) {
		damage *= target.PseudoStats.ArcaneDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFire) {
		damage *= target.PseudoStats.FireDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolFrost) {
		damage *= target.PseudoStats.FrostDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolHoly) {
		damage *= target.PseudoStats.HolyDamageTakenMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolNature) {
		damage *= target.PseudoStats.NatureDamageTakenMultiplier
		damage *= attackTable.NatureDamageDealtMultiplier
	} else if spell.SpellSchool.Matches(SpellSchoolShadow) {
		damage *= target.PseudoStats.ShadowDamageTakenMultiplier
		if isPeriodic {
			damage *= attackTable.PeriodicShadowDamageDealtMultiplier
		}
	}

	return damage
}

func (spell *Spell) CasterHealingMultiplier() float64 {
	if spell.Flags.Matches(SpellFlagIgnoreAttackerModifiers) {
		return 1
	}

	return spell.Unit.PseudoStats.HealingDealtMultiplier *
		spell.DamageMultiplier *
		spell.DamageMultiplierAdditive
}
func (spell *Spell) applyTargetHealingModifiers(damage float64, attackTable *AttackTable) float64 {
	if spell.Flags.Matches(SpellFlagIgnoreTargetModifiers) {
		return damage
	}

	return damage *
		attackTable.Defender.PseudoStats.HealingTakenMultiplier *
		attackTable.HealingDealtMultiplier
}
