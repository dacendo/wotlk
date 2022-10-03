import { Consumes } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { ItemSpec } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Faction } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Player } from '../core/player.js';

import {
	WarriorShout,
	WarriorTalents as WarriorTalents,
	Warrior,
	Warrior_Rotation as WarriorRotation,
	Warrior_Rotation_SunderArmor as SunderArmor,
	Warrior_Options as WarriorOptions,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
	Warrior_Rotation_StanceOption as StanceOption,
} from '../core/proto/warrior.js';

import * as Gems from '../core/proto_utils/gems.js';
import * as Tooltips from '../core/constants/tooltips.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const ArmsTalents = {
	name: 'Arms',
	data: SavedTalents.create({
		talentsString: '3022032023335100102012213231251-305-2033',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfRending,
			major2: WarriorMajorGlyph.GlyphOfMortalStrike,
			major3: WarriorMajorGlyph.GlyphOfExecution,
			minor1: WarriorMinorGlyph.GlyphOfThunderClap,
			minor2: WarriorMinorGlyph.GlyphOfCommand,
			minor3: WarriorMinorGlyph.GlyphOfCharge,
		}),
	}),
};

export const FuryTalents = {
	name: 'Fury',
	data: SavedTalents.create({
		talentsString: '32002301233-305053000520310053120500351',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfWhirlwind,
			major2: WarriorMajorGlyph.GlyphOfHeroicStrike,
			major3: WarriorMajorGlyph.GlyphOfExecution,
			minor1: WarriorMinorGlyph.GlyphOfBattle,
			minor2: WarriorMinorGlyph.GlyphOfBloodrage,
			minor3: WarriorMinorGlyph.GlyphOfCharge,
		}),
	}),
};

export const DefaultRotation = WarriorRotation.create({
	useRend: false,
	useMs: true,
	useCleave: false,

	prioritizeWw: true,
	sunderArmor: SunderArmor.SunderArmorNone,

	msRageThreshold: 35,
	hsRageThreshold: 30,
	rendRageThresholdBelow: 70,
	slamRageThreshold: 25,
	rendCdThreshold: 0,
	useHsDuringExecute: true,
	useBtDuringExecute: true,
	useWwDuringExecute: true,
	useSlamOverExecute: true,
	spamExecute: true,
	stanceOption: StanceOption.DefaultStance,
});

export const ArmsRotation = WarriorRotation.create({
	useRend: true,
	useMs: true,
	useCleave: false,
	sunderArmor: SunderArmor.SunderArmorNone,
	msRageThreshold: 355,
	slamRageThreshold: 25,
	hsRageThreshold: 50,
	rendCdThreshold: 0,
	useHsDuringExecute: true,
	spamExecute: true,
	stanceOption: StanceOption.DefaultStance,
});

export const DefaultOptions = WarriorOptions.create({
	startingRage: 0,
	useRecklessness: true,
	shout: WarriorShout.WarriorShoutCommanding,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodSpicedWormBurger,
	defaultPotion: Potions.IndestructiblePotion,
	prepopPotion: Potions.PotionOfSpeed,
});

export const P1_PRERAID_FURY_PRESET = {
	name: 'P1 Pre-Raid Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 41386,
          "enchant": 44879,
          "gems": [
            41398,
            42702
          ]
        },
        {
          "id": 42645,
          "gems": [
            40003
          ]
        },
        {
          "id": 44195,
          "enchant": 44871
        },
        {
          "id": 37647,
          "enchant": 55002
        },
        {
          "id": 39606,
          "enchant": 44489,
          "gems": [
            40003,
            40003
          ]
        },
        {
          "id": 44203,
          "enchant": 44484,
          "gems": [
            0
          ]
        },
        {
          "id": 39609,
          "enchant": 54999,
          "gems": [
            40037,
            0
          ]
        },
        {
          "id": 40694,
          "gems": [
            42149,
            42149
          ]
        },
        {
          "id": 44205,
          "enchant": 38374
        },
        {
          "id": 44306,
          "enchant": 55016,
          "gems": [
            40037,
            40037
          ]
        },
        {
          "id": 42642,
          "gems": [
            42149
          ]
        },
        {
          "id": 37642
        },
        {
          "id": 42987
        },
        {
          "id": 40684
        },
        {
          "id": 37852,
          "enchant": 44492
        },
        {
          "id": 37852,
          "enchant": 44492
        },
        {
          "id": 37191
        }
      ]}`),
};

export const P1_FURY_PRESET = {
	name: 'P1 Fury Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().bloodthirst,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {
          "id": 44006,
          "enchant": 44879,
          "gems": [
            41285,
            42702
          ]
        },
        {
          "id": 44664,
          "gems": [
            40037
          ]
        },
        {
          "id": 40530,
          "enchant": 44871,
          "gems": [
            42149
          ]
        },
        {
          "id": 40403,
          "enchant": 55002
        },
        {
          "id": 40539,
          "enchant": 44489,
          "gems": [
            40037
          ]
        },
        {
          "id": 39765,
          "enchant": 44484,
          "gems": [
            40037
          ]
        },
        {
          "id": 40541,
          "enchant": 54999,
          "gems": [
            0
          ]
        },
        {
          "id": 40205,
          "gems": [
            42149
          ]
        },
        {
          "id": 40529,
          "enchant": 38374,
          "gems": [
            40037,
            40086
          ]
        },
        {
          "id": 40591,
          "enchant": 55016
        },
        {
          "id": 43993,
          "gems": [
            42149
          ]
        },
        {
          "id": 40717
        },
        {
          "id": 42987
        },
        {
          "id": 40256
        },
        {
          "id": 40384,
          "enchant": 44492
        },
        {
          "id": 40384,
          "enchant": 44492
        },
        {
          "id": 40385
        }
      ]}`),
};

export const P1_PRERAID_ARMS_PRESET = {
	name: 'P1 Pre-Raid Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
        {
          "id": 41386,
          "enchant": 44879,
          "gems": [
            41285,
            42702
          ]
        },
        {
          "id": 42645,
          "gems": [
            40002
          ]
        },
        {
          "id": 44195,
          "enchant": 44871
        },
        {
          "id": 37647,
          "enchant": 55002
        },
        {
          "id": 39606,
          "enchant": 44489,
          "gems": [
            40002,
            40002
          ]
        },
        {
          "id": 41355,
          "enchant": 44484,
          "gems": [
            0
          ]
        },
        {
          "id": 39609,
          "enchant": 54999,
          "gems": [
            40037,
            0
          ]
        },
        {
          "id": 40694,
          "gems": [
            42149,
            42149
          ]
        },
        {
          "id": 37193,
          "enchant": 38374,
          "gems": [
            40002,
            40037
          ]
        },
        {
          "id": 44306,
          "enchant": 55016,
          "gems": [
            40086,
            40002
          ]
        },
        {
          "id": 42642,
          "gems": [
            42149
          ]
        },
        {
          "id": 37642
        },
        {
          "id": 42987
        },
        {
          "id": 40684
        },
        {
          "id": 37852,
          "enchant": 44492
        },
        {},
        {
          "id": 37191
        }
      ]}`),
};

export const P1_ARMS_PRESET = {
	name: 'P1 Arms Preset',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecWarrior>) => player.getTalents().mortalStrike,
	gear: EquipmentSpec.fromJsonString(`{"items": [
        {
          "id": 40528,
          "enchant": 44879,
          "gems": [
            41398,
            42153
          ]
        },
        {
          "id": 44664,
          "gems": [
            39996
          ]
        },
        {
          "id": 40530,
          "enchant": 44871,
          "gems": [
            40038
          ]
        },
        {
          "id": 40403,
          "enchant": 55002
        },
        {
          "id": 40539,
          "enchant": 44489,
          "gems": [
            42153
          ]
        },
        {
          "id": 40330,
          "enchant": 44484,
          "gems": [
            40002
          ]
        },
        {
          "id": 40541,
          "enchant": 54999,
          "gems": [
            0
          ]
        },
        {
          "id": 40205,
          "gems": [
            42153
          ]
        },
        {
          "id": 40318,
          "enchant": 38374,
          "gems": [
            49110,
            40038
          ]
        },
        {
          "id": 40591,
          "enchant": 55016
        },
        {
          "id": 43993,
          "gems": [
            40002
          ]
        },
        {
          "id": 40474
        },
        {
          "id": 42987
        },
        {
          "id": 40256
        },
        {
          "id": 40384,
          "enchant": 44492
        },
        {},
        {
          "id": 40385
        }
      ]}`),
};
