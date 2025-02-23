import { BooleanPicker } from '../core/components/boolean_picker.js';
import { EnumPicker } from '../core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '../core/components/icon_enum_picker.js';
import { IconPickerConfig } from '../core/components/icon_picker.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	EnhancementShaman_Options as ShamanOptions,
	ShamanTotems,
	ShamanShield,
	ShamanImbue,
	ShamanSyncType,
	EnhancementShaman_Rotation_PrimaryShock as PrimaryShock,
	EnhancementShaman_Rotation_RotationType as RotationType,
	EnhancementShaman_Rotation_CustomRotationSpell as CustomRotationSpell
} from '../core/proto/shaman.js';
import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Target } from '../core/target.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Bloodlust = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecEnhancementShaman>({
	fieldName: 'bloodlust',
	id: ActionId.fromSpellId(2825),
});
export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ color: 'grey', value: ShamanShield.NoShield },
		{ actionId: ActionId.fromSpellId(57960), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
	],
});

export const ShamanImbueMH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
	fieldName: 'imbueMh',
	values: [
		{ color: 'grey', value: ShamanImbue.NoImbue },
		{ actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
		{ actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon, text: 'R10'},
		{ actionId: ActionId.fromSpellId(58789), value: ShamanImbue.FlametongueWeaponDownrank, text: 'R9'},
		{ actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
	],
});

export const ShamanImbueOH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
	fieldName: 'imbueOh',
	values: [
		{ color: 'grey', value: ShamanImbue.NoImbue },
		{ actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
		{ actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon, text: 'R10'},
		{ actionId: ActionId.fromSpellId(58789), value: ShamanImbue.FlametongueWeaponDownrank, text: 'R9'},
		{ actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
	],
});

export const SyncTypeInput = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecEnhancementShaman, ShamanSyncType>({
	fieldName: 'syncType',
	label: 'Sync/Stagger Setting',
	labelTooltip: 'Choose your sync or stagger option, Perfect Sync makes your weapons always attack at the same time, which is ideal for mixed imbues. Delayed Offhand is similar but additionally adds a slight delay to the offhand attacks while staying within the 0.5s flurry ICD window, ideal for matched imbues.',
	values: [
		{ name: 'None', value: ShamanSyncType.NoSync },
		{ name: 'Perfect Sync', value: ShamanSyncType.SyncMainhandOffhandSwings },
		{ name: 'Delayed Offhand', value: ShamanSyncType.DelayOffhandSwings },
	],
});

export const EnhancementShamanRotationConfig = {
	inputs:
		[
			InputHelpers.makeRotationEnumInput<Spec.SpecEnhancementShaman, RotationType>({
				fieldName: 'rotationType',
				label: 'Type',
				labelTooltip:
					`<ul>
					<li>
						<div>Standard: Priority Rotation</div>
					</li>
					<li>
						<div>Custom: Highest spell that is ready will be cast.</div>
					</li>
				</ul>`,
				values: [
					{ name: 'Standard', value: RotationType.Priority },
					{ name: 'Custom', value: RotationType.Custom },
				],
			}),
			InputHelpers.makeCustomRotationInput<Spec.SpecEnhancementShaman, CustomRotationSpell>({
				fieldName: 'customRotation',
				numColumns: 2,
				values: [
					{ actionId: ActionId.fromSpellId(49238), value: CustomRotationSpell.LightningBolt},
					{ actionId: ActionId.fromSpellId(49238), value: CustomRotationSpell.LightningBoltWeave, text: "Weave" },
					{ actionId: ActionId.fromSpellId(17364), value: CustomRotationSpell.StormstrikeDebuffMissing, text: "Debuff"  },
					{ actionId: ActionId.fromSpellId(17364), value: CustomRotationSpell.Stormstrike },
					{ actionId: ActionId.fromSpellId(49233), value: CustomRotationSpell.FlameShock },
					{ actionId: ActionId.fromSpellId(49231), value: CustomRotationSpell.EarthShock },
					{ actionId: ActionId.fromSpellId(61657), value: CustomRotationSpell.FireNova },
					{ actionId: ActionId.fromSpellId(60103), value: CustomRotationSpell.LavaLash },
					{ actionId: ActionId.fromSpellId(49281), value: CustomRotationSpell.LightningShield },
					{ actionId: ActionId.fromSpellId(60043), value: CustomRotationSpell.LavaBurst, text: "Weave" },
					{ actionId: ActionId.fromSpellId(49236), value: CustomRotationSpell.FrostShock},
				],
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType == RotationType.Custom,
			}),
			InputHelpers.makeRotationEnumInput<Spec.SpecEnhancementShaman, PrimaryShock>({
				fieldName: 'primaryShock',
				label: 'Primary Shock',
				values: [
					{ name: 'None', value: PrimaryShock.None },
					{ name: 'Earth Shock', value: PrimaryShock.Earth },
					{ name: 'Frost Shock', value: PrimaryShock.Frost },
				],
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType != RotationType.Custom
			}),
			InputHelpers.makeRotationBooleanInput<Spec.SpecEnhancementShaman>({
				fieldName: 'weaveFlameShock',
				label: 'Weave Flame Shock',
				labelTooltip: 'Use Flame Shock whenever the target does not already have the DoT.',
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType != RotationType.Custom
			}),
            InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
                fieldName:  'flameShockClipTicks',
                label:  'Refresh Flame Shock at ticks remaining',
                labelTooltip: 'Set to 0 to require the debuff be missing. A tick is 3s, affected by spell haste',
                enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().weaveFlameShock,
                showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().weaveFlameShock
            }),
			InputHelpers.makeRotationBooleanInput<Spec.SpecEnhancementShaman>({
				fieldName: 'lightningboltWeave',
				label: 'Enable Weaving Lightning Bolt',
				labelTooltip: 'Will provide a DPS increase, but is harder to execute',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getTalents().maelstromWeapon > 0,
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType != RotationType.Custom
			}),
			InputHelpers.makeRotationEnumInput<Spec.SpecEnhancementShaman, number>({
				fieldName: 'maelstromweaponMinStack',
				label: 'Minimum Maelstrom Stacks to Weave',
				labelTooltip: '3 stacks is the most realistic, however there are cases where lower might be possible, just much harder to do in practice',
				values: [
					{ name: '1', value: 1 },
					{ name: '2', value: 2 },
					{ name: '3', value: 3 },
					{ name: '4', value: 4 },
				],
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().lightningboltWeave || player.getRotation().rotationType == RotationType.Custom,
			}),
			InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
				fieldName: 'weaveReactionTime',
				label: 'Weaving Reaction Time',
				labelTooltip: 'The Reaction time to gaining maelstrom stacks after a auto attack in milliseconds',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().lightningboltWeave || player.getRotation().rotationType == RotationType.Custom,
			}),
			InputHelpers.makeRotationBooleanInput<Spec.SpecEnhancementShaman>({
				fieldName: 'lavaburstWeave',
				label: 'Enable Weaving Lava Burst',
				labelTooltip: 'Not particularily useful for dual wield, mostly a 2h option',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().lightningboltWeave,
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType != RotationType.Custom
			}),
			InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
				fieldName: 'firenovaManaThreshold',
				label: 'Minimum mana to cast Fire Nova',
				labelTooltip: 'Fire Nova will not be cast when mana is below this value. Set this medium-low, it has a bad mana-to-damage ratio',
			}),
			InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
				fieldName: 'shamanisticRageManaThreshold',
				label: 'Mana % to use Shamanistic Rage',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getTalents().shamanisticRage,
			}),
		],
};
