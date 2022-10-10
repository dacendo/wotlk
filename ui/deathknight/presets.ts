import { Consumes, Explosive, PetFood, RaidTarget } from '../core/proto/common.js';
import { EquipmentSpec } from '../core/proto/common.js';
import { Flask } from '../core/proto/common.js';
import { Food } from '../core/proto/common.js';
import { Glyphs } from '../core/proto/common.js';
import { Potions } from '../core/proto/common.js';
import { SavedTalents } from '../core/proto/ui.js';
import { Spec } from '../core/proto/common.js';
import { Player } from '../core/player.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';

import {
	Deathknight_Rotation as DeathKnightRotation,
	Deathknight_Options as DeathKnightOptions,
	DeathknightMajorGlyph,
	DeathknightMinorGlyph,
} from '../core/proto/deathknight.js';

import * as Tooltips from '../core/constants/tooltips.js';
import { Deathknight_Rotation_StartingPresence, Deathknight_Rotation_BloodRuneFiller } from '../core/proto/deathknight.js';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.
export const FrostTalents = {
	name: 'Frost BL',
	data: SavedTalents.create({
		talentsString: '23050005-32005350352203012300033101351',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const FrostUnholyTalents = {
	name: 'Frost UH',
	data: SavedTalents.create({
		talentsString: '-32002350352203012300033101351-230200305003',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfObliterate,
			major2: DeathknightMajorGlyph.GlyphOfFrostStrike,
			major3: DeathknightMajorGlyph.GlyphOfDisease,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const UnholyDualWieldTalents = {
	name: 'Unholy DW',
	data: SavedTalents.create({
		talentsString: '-320043500002-2300303050032152000150013133051',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfDarkDeath,
			major3: DeathknightMajorGlyph.GlyphOfDeathAndDecay,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const Unholy2HTalents = {
	name: 'Unholy 2H',
	data: SavedTalents.create({
		talentsString: '-320050500002-2300303150032152000150013133151',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfTheGhoul,
			major2: DeathknightMajorGlyph.GlyphOfDarkDeath,
			major3: DeathknightMajorGlyph.GlyphOfIcyTouch,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const BloodTalents = {
	name: 'Blood DPS',
	data: SavedTalents.create({
		talentsString: '2305020530003303231023101351--230220305003',
		glyphs: Glyphs.create({
			major1: DeathknightMajorGlyph.GlyphOfDancingRuneWeapon,
			major2: DeathknightMajorGlyph.GlyphOfDeathStrike,
			major3: DeathknightMajorGlyph.GlyphOfDarkDeath,
			minor1: DeathknightMinorGlyph.GlyphOfHornOfWinter,
			minor2: DeathknightMinorGlyph.GlyphOfPestilence,
			minor3: DeathknightMinorGlyph.GlyphOfRaiseDead,
		}),
	}),
};

export const DefaultUnholyRotation = DeathKnightRotation.create({
	useDeathAndDecay: true,
	btGhoulFrenzy: true,
	refreshHornOfWinter: false,
	useEmpowerRuneWeapon: true,
	startingPresence: Deathknight_Rotation_StartingPresence.Unholy,
	bloodRuneFiller: Deathknight_Rotation_BloodRuneFiller.BloodBoil,
	useAms: false,
});

export const DefaultUnholyOptions = DeathKnightOptions.create({
	startingRunicPower: 0,
	petUptime: 1,
	precastGhoulFrenzy: true,
	precastHornOfWinter: true,
	unholyFrenzyTarget: RaidTarget.create({
		targetIndex: NO_TARGET, // In an individual sim the 0-indexed player is ourself.
	}),
});

export const DefaultFrostRotation = DeathKnightRotation.create({
	useDeathAndDecay: false,
	btGhoulFrenzy: false,
	refreshHornOfWinter: false,
	useEmpowerRuneWeapon: true,
	startingPresence: Deathknight_Rotation_StartingPresence.Blood,
	bloodRuneFiller: Deathknight_Rotation_BloodRuneFiller.BloodBoil,
	useAms: false,
	avgAmsSuccessRate: 1.0,
	avgAmsHit: 10000.0,
	oblitDelayDuration: 0,
});

export const DefaultFrostOptions = DeathKnightOptions.create({
	startingRunicPower: 0,
	petUptime: 1,
	precastHornOfWinter: true,
	unholyFrenzyTarget: RaidTarget.create({
		targetIndex: NO_TARGET, // In an individual sim the 0-indexed player is ourself.
	}),
});

export const DefaultBloodRotation = DeathKnightRotation.create({
	refreshHornOfWinter: false,
	useEmpowerRuneWeapon: true,
	startingPresence: Deathknight_Rotation_StartingPresence.Blood,
	bloodRuneFiller: Deathknight_Rotation_BloodRuneFiller.BloodStrike,
	useAms: false,
});

export const DefaultBloodOptions = DeathKnightOptions.create({
	startingRunicPower: 0,
	petUptime: 1,
	precastHornOfWinter: true,
	unholyFrenzyTarget: RaidTarget.create({
		targetIndex: NO_TARGET, // In an individual sim the 0-indexed player is ourself.
	}),
});

export const OtherDefaults = {
};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.PotionOfSpeed,
	petFood: PetFood.PetFoodSpicedMammothTreats,
	prepopPotion: Potions.PotionOfSpeed,
	thermalSapper: true,
	fillerExplosive: Explosive.ExplosiveSaroniteBomb,
});

export const P1_UNHOLY_2H_PRERAID_PRESET = {
	name: 'Pre-Raid 2H Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
    {
      "id": 41386,
      "enchant": 44879,
      "gems": [
        41400,
        49110
      ]
    },
    {
      "id": 37397
    },
    {
      "id": 37627,
      "enchant": 44871,
      "gems": [
        39996
      ]
    },
    {
      "id": 37647,
      "enchant": 44472
    },
    {
      "id": 39617,
      "enchant": 44489,
      "gems": [
        42142,
        39996
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
      "id": 39618,
      "enchant": 54999,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40688,
      "gems": [
        39996,
        42142
      ]
    },
    {
      "id": 37193,
      "enchant": 38374,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 44306,
      "enchant": 55016,
      "gems": [
        39996,
        39996
      ]
    },
    {
      "id": 37642
    },
    {
      "id": 44935
    },
    {
      "id": 40684
    },
    {
      "id": 42987
    },
    {
      "id": 41257,
      "enchant": 53344
    },
    {},
    {
      "id": 40867
    }
  ]}`),
};

export const P1_UNHOLY_2H_BIS_PRESET = {
	name: 'P1 2H Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
    {
      "id": 44006,
      "enchant": 44879,
      "gems": [
        41400,
        49110
      ]
    },
    {
      "id": 44664,
      "gems": [
        39996
      ]
    },
    {
      "id": 40557,
      "enchant": 44871,
      "gems": [
        39996
      ]
    },
    {
      "id": 40403,
      "enchant": 44472
    },
    {
      "id": 40550,
      "enchant": 44489,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 40330,
      "enchant": 44484,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40552,
      "enchant": 54999,
      "gems": [
        40038,
        0
      ]
    },
    {
      "id": 40278,
      "gems": [
        42142,
        42142
      ]
    },
    {
      "id": 40556,
      "enchant": 38374,
      "gems": [
        39996,
        39996
      ]
    },
    {
      "id": 40591,
      "enchant": 55016
    },
    {
      "id": 39401
    },
    {
      "id": 40075
    },
    {
      "id": 40256
    },
    {
      "id": 42987
    },
    {
      "id": 40384,
      "enchant": 53344
    },
    {},
    {
      "id": 40207
    }
  ]}`),
};

export const P1_UNHOLY_DW_PRERAID_PRESET = {
	name: 'Pre-Raid DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	gear: EquipmentSpec.fromJsonString(`{"items": [
    {
      "id": 41386,
      "enchant": 44879,
      "gems": [
        41400,
        49110
      ]
    },
    {
      "id": 37397
    },
    {
      "id": 37627,
      "enchant": 44871,
      "gems": [
        39996
      ]
    },
    {
      "id": 37647,
      "enchant": 44472
    },
    {
      "id": 39617,
      "enchant": 44489,
      "gems": [
        42142,
        39996
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
      "id": 39618,
      "enchant": 54999,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40688,
      "gems": [
        39996,
        42142
      ]
    },
    {
      "id": 37193,
      "enchant": 38374,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 44306,
      "enchant": 55016,
      "gems": [
        39996,
        39996
      ]
    },
    {
      "id": 37642
    },
    {
      "id": 44935
    },
    {
      "id": 40684
    },
    {
      "id": 42987
    },
    {
      "id": 41383,
      "enchant": 53344
    },
    {
      "id": 40703,
      "enchant": 44495
    },
    {
      "id": 40867
    }
  ]}`),
};

export const P1_UNHOLY_DW_BIS_PRESET = {
	name: 'P1 DW Unholy',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().summonGargoyle,
	gear: EquipmentSpec.fromJsonString(`{"items": [
    {
      "id": 44006,
      "enchant": 44879,
      "gems": [
        41398,
        42702
      ]
    },
    {
      "id": 39421
    },
    {
      "id": 40557,
      "enchant": 44871,
      "gems": [
        39996
      ]
    },
    {
      "id": 40403,
      "enchant": 44472
    },
    {
      "id": 40550,
      "enchant": 44489,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 40330,
      "enchant": 44484,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40347,
      "enchant": 54999,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40278,
      "gems": [
        42142,
        42142
      ]
    },
    {
      "id": 40294,
      "enchant": 38374
    },
    {
      "id": 39706,
      "enchant": 55016,
      "gems": [
        39996
      ]
    },
    {
      "id": 39401
    },
    {
      "id": 40075
    },
    {
      "id": 37390
    },
    {
      "id": 42987
    },
    {
      "id": 40402,
      "enchant": 53344
    },
    {
      "id": 40491,
      "enchant": 53344
    },
    {
      "id": 42620
    }
  ]}`),
};

export const P1_FROST_PRE_BIS_PRESET = {
	name: 'Pre-Raid Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().howlingBlast,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
    {
      "id": 41386,
      "enchant": 44879,
      "gems": [
        41398,
        49110
      ]
    },
    {
      "id": 42645,
      "gems": [
        42142
      ]
    },
    {
      "id": 34388,
      "enchant": 44871,
      "gems": [
        39996,
        39996
      ]
    },
    {
      "id": 37647,
      "enchant": 44472
    },
    {
      "id": 39617,
      "enchant": 44489,
      "gems": [
        42142,
        39996
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
      "id": 39618,
      "enchant": 54999,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 37171,
      "gems": [
        39996,
        39996
      ]
    },
    {
      "id": 37193,
      "enchant": 38374,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 44306,
      "enchant": 55016,
      "gems": [
        39996,
        39996
      ]
    },
    {
      "id": 42642,
      "gems": [
        39996
      ]
    },
    {
      "id": 44935
    },
    {
      "id": 40684
    },
    {
      "id": 42987
    },
    {
      "id": 41383,
      "enchant": 53343
    },
    {
      "id": 43611,
      "enchant": 53344
    },
    {
      "id": 40715
    }
  ]}`),
};

export const P1_FROST_BIS_PRESET = {
	name: 'P1 Frost',
	tooltip: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().howlingBlast,
	gear: EquipmentSpec.fromJsonString(`{   "items": [
    {
      "id": 44006,
      "enchant": 44879,
      "gems": [
        41398,
        42702
      ]
    },
    {
      "id": 44664,
      "gems": [
        39996
      ]
    },
    {
      "id": 40557,
      "enchant": 44871,
      "gems": [
        39996
      ]
    },
    {
      "id": 40403,
      "enchant": 44472
    },
    {
      "id": 40550,
      "enchant": 44489,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 40330,
      "enchant": 44484,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40552,
      "enchant": 54999,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40278,
      "gems": [
        39996,
        42142
      ]
    },
    {
      "id": 40556,
      "enchant": 38374,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 40591,
      "enchant": 55016
    },
    {
      "id": 39401
    },
    {
      "id": 40075
    },
    {
      "id": 40256
    },
    {
      "id": 42987
    },
    {
      "id": 40189,
      "enchant": 53343
    },
    {
      "id": 40189,
      "enchant": 53344
    },
    {
      "id": 40207
    }
  ]}`),
};

export const P1_BLOOD_BIS_PRESET = {
	name: 'P1 Blood',
	toolbar: Tooltips.BASIC_BIS_DISCLAIMER,
	enableWhen: (player: Player<Spec.SpecDeathknight>) => player.getTalents().dancingRuneWeapon,
	gear: EquipmentSpec.fromJsonString(`{ "items": [
    {
      "id": 44006,
      "enchant": 44879,
      "gems": [
        41398,
        49110
      ]
    },
    {
      "id": 44664,
      "gems": [
        39996
      ]
    },
    {
      "id": 40557,
      "enchant": 44871,
      "gems": [
        39996
      ]
    },
    {
      "id": 40250,
      "enchant": 44472
    },
    {
      "id": 40550,
      "enchant": 44489,
      "gems": [
        42142,
        39996
      ]
    },
    {
      "id": 40330,
      "enchant": 44484,
      "gems": [
        39996,
        0
      ]
    },
    {
      "id": 40552,
      "enchant": 54999,
      "gems": [
        40038,
        0
      ]
    },
    {
      "id": 40278,
      "gems": [
        42142,
        42142
      ]
    },
    {
      "id": 40556,
      "enchant": 38374,
      "gems": [
        39996,
        39996
      ]
    },
    {
      "id": 40591,
      "enchant": 55016
    },
    {
      "id": 43993,
      "gems": [
        39996
      ]
    },
    {
      "id": 40075
    },
    {
      "id": 40256
    },
    {
      "id": 42987
    },
    {
      "id": 40384,
      "enchant": 53344
    },
    {},
    {
      "id": 40207
    }
  ]}`),
};
