package mage

import (
	"testing"

	_ "github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func init() {
	RegisterMage()
}

func TestArcane(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P1Arcane", GearSet: P1ArcaneGear},
		SpecOptions: core.SpecOptionsCombo{Label: "ArcaneRotation", SpecOptions: PlayerOptionsArcane},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "AOE", SpecOptions: PlayerOptionsArcaneAOE},
		},

		Consumes: FullArcaneConsumes,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},
	}))
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P1Fire", GearSet: P1FireGear},
		SpecOptions: core.SpecOptionsCombo{Label: "FireRotation", SpecOptions: PlayerOptionsFire},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "AOE", SpecOptions: PlayerOptionsFireAOE},
		},

		Consumes: FullFireConsumes,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},
	}))
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GearSetCombo{Label: "P1Frost", GearSet: P1FrostGear},
		SpecOptions: core.SpecOptionsCombo{Label: "FrostRotation", SpecOptions: PlayerOptionsFrost},
		OtherSpecOptions: []core.SpecOptionsCombo{
			core.SpecOptionsCombo{Label: "AOE", SpecOptions: PlayerOptionsFrostAOE},
		},

		Consumes: FullFrostConsumes,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},
	}))
}
