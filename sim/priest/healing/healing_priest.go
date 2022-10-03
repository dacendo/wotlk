package healing

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/priest"
)

func RegisterHealingPriest() {
	core.RegisterAgentFactory(
		proto.Player_HealingPriest{},
		proto.Spec_SpecHealingPriest,
		func(character core.Character, options proto.Player) core.Agent {
			return NewHealingPriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_HealingPriest)
			if !ok {
				panic("Invalid spec value for Healing Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

type HealingPriest struct {
	*priest.Priest

	rotation proto.HealingPriest_Rotation

	// Spells to rotate through for cyclic rotation.
	spellCycle     []*core.Spell
	nextCycleIndex int
}

func NewHealingPriest(character core.Character, options proto.Player) *HealingPriest {
	healingOptions := options.GetHealingPriest()

	selfBuffs := priest.SelfBuffs{
		UseInnerFire:   healingOptions.Options.UseInnerFire,
		UseShadowfiend: healingOptions.Options.UseShadowfiend,
	}

	if healingOptions.Talents.PowerInfusion && healingOptions.Options.PowerInfusionTarget != nil {
		selfBuffs.PowerInfusionTarget = *healingOptions.Options.PowerInfusionTarget
	} else {
		selfBuffs.PowerInfusionTarget.TargetIndex = -1
	}

	basePriest := priest.New(character, selfBuffs, *healingOptions.Talents)

	hpriest := &HealingPriest{
		Priest:   basePriest,
		rotation: *healingOptions.Rotation,
	}

	hpriest.EnableResumeAfterManaWait(hpriest.tryUseGCD)

	return hpriest
}

func (hpriest *HealingPriest) GetPriest() *priest.Priest {
	return hpriest.Priest
}

func (hpriest *HealingPriest) GetMainTarget() *core.Unit {
	target := hpriest.Env.Raid.GetFirstTargetDummy()
	if target == nil {
		return &hpriest.Unit
	} else {
		return &target.Unit
	}
}

func (hpriest *HealingPriest) Initialize() {
	hpriest.CurrentTarget = hpriest.GetMainTarget()
	hpriest.Priest.Initialize()
	hpriest.Priest.RegisterHealingSpells()

	hpriest.RegisterHymnOfHopeCD()

	hpriest.spellCycle = []*core.Spell{
		hpriest.GreaterHeal,
		hpriest.FlashHeal,
		hpriest.CircleOfHealing,
		hpriest.BindingHeal,
		hpriest.PrayerOfHealing,
		hpriest.PrayerOfMending,
		hpriest.PenanceHeal,
	}
}

func (hpriest *HealingPriest) Reset(sim *core.Simulation) {
	hpriest.Priest.Reset(sim)
	hpriest.nextCycleIndex = 0
}
