package balance

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/druid"
)

func RegisterBalanceDruid() {
	core.RegisterAgentFactory(
		proto.Player_BalanceDruid{},
		proto.Spec_SpecBalanceDruid,
		func(character core.Character, options proto.Player) core.Agent {
			return NewBalanceDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_BalanceDruid)
			if !ok {
				panic("Invalid spec value for Balance Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewBalanceDruid(character core.Character, options proto.Player) *BalanceDruid {
	balanceOptions := options.GetBalanceDruid()
	selfBuffs := druid.SelfBuffs{}

	if balanceOptions.Options.InnervateTarget != nil {
		selfBuffs.InnervateTarget = *balanceOptions.Options.InnervateTarget
	} else {
		selfBuffs.InnervateTarget.TargetIndex = -1
	}

	moonkin := &BalanceDruid{
		Druid:    druid.New(character, druid.Moonkin, selfBuffs, *balanceOptions.Talents),
		Rotation: *balanceOptions.Rotation,
	}

	moonkin.ResetTalentsBonuses()
	moonkin.RegisterTalentsBonuses()
	moonkin.EnableResumeAfterManaWait(moonkin.tryUseGCD)
	return moonkin
}

type BalanceCooldown struct {
	Name     string
	Stat     stats.Stat
	Priority float64
	ID       core.ActionID
}

type BalanceDruid struct {
	*druid.Druid

	Rotation           proto.BalanceDruid_Rotation
	CooldownsAvailable []*core.MajorCooldown

	// CDS
	hyperSpeedMCD      *core.MajorCooldown
	potionSpeedMCD     *core.MajorCooldown
	potionWildMagicMCD *core.MajorCooldown
	onUseTrinket1      *core.MajorCooldown
	onUseTrinket2      *core.MajorCooldown
	potionUsed         bool
}

func (moonkin *BalanceDruid) GetDruid() *druid.Druid {
	return moonkin.Druid
}

func (moonkin *BalanceDruid) Initialize() {
	moonkin.Druid.Initialize()
	moonkin.RegisterBalanceSpells()
}

func (moonkin *BalanceDruid) Reset(sim *core.Simulation) {
	moonkin.Druid.Reset(sim)
	moonkin.RebirthTiming = moonkin.Env.BaseDuration.Seconds() * sim.RandomFloat("Rebirth Timing")

	if moonkin.Rotation.UseSmartCooldowns == true {
		moonkin.potionUsed = false
		consumes := &moonkin.Consumes

		if consumes.DefaultPotion == proto.Potions_PotionOfSpeed {
			moonkin.potionSpeedMCD = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: 40211})
		}
		if consumes.DefaultPotion == proto.Potions_PotionOfWildMagic {
			moonkin.potionWildMagicMCD = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: 40212})
		}
		if moonkin.HasProfession(proto.Profession_Engineering) {
			moonkin.hyperSpeedMCD = moonkin.getBalanceMajorCooldown(core.ActionID{SpellID: 54758})
		}
		moonkin.onUseTrinket1 = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: moonkin.Equip[items.ItemSlotTrinket1].ID})
		moonkin.onUseTrinket2 = moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: moonkin.Equip[items.ItemSlotTrinket2].ID})
	}
}

// Takes out a Cooldown from the generic MajorCooldownManager and adds it to a custom Slice of Cooldowns
func (moonkin *BalanceDruid) getBalanceMajorCooldown(actionID core.ActionID) *core.MajorCooldown {
	if moonkin.Character.HasMajorCooldown(actionID) {
		majorCd := moonkin.Character.GetMajorCooldown(actionID)
		majorCd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			return false
		}
		return majorCd
	}
	return nil
}
