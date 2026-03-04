<script lang="ts">
	import { api } from '../api.js';
	import type { Flag } from '../types.js';

	interface Props {
		flag: Flag;
		environment: string;
		onUpdate?: () => void;
	}

	let { flag, environment, onUpdate }: Props = $props();

	let isLoading = $state(false);
	let localValue = $state(0);

	// Initialize local value from flag data
	$effect(() => {
		localValue = flag.environments[environment]?.rollout_percent ?? 0;
	});

	async function handleChange(event: Event) {
		if (isLoading) return;

		const target = event.target as HTMLInputElement;
		const newValue = parseFloat(target.value);

		try {
			isLoading = true;

			// Update the flag with new rollout percentage
			const updatedFlag = {
				...flag,
				environments: {
					...flag.environments,
					[environment]: {
						...flag.environments[environment],
						rollout_percent: newValue
					}
				}
			};

			await api.updateFlag(flag.key, updatedFlag);
			localValue = newValue;
			onUpdate?.();
		} catch (error) {
			console.error('Failed to update rollout percentage:', error);
			// Reset to original value on error
			localValue = flag.environments[environment]?.rollout_percent ?? 0;
		} finally {
			isLoading = false;
		}
	}

	function handleInput(event: Event) {
		const target = event.target as HTMLInputElement;
		localValue = parseFloat(target.value);
	}
</script>

<div class="space-y-2">
	<div class="flex items-center justify-between">
		<label for="rollout-slider-{flag.key}-{environment}" class="text-sm font-medium text-gray-700">
			Rollout Percentage
		</label>
		<span class="text-sm text-gray-500">
			{Math.round(localValue)}%
		</span>
	</div>
	<div class="relative">
		<input
			id="rollout-slider-{flag.key}-{environment}"
			type="range"
			min="0"
			max="100"
			step="1"
			value={localValue}
			oninput={handleInput}
			onchange={handleChange}
			disabled={isLoading}
			class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider {isLoading
				? 'opacity-50 cursor-not-allowed'
				: ''}"
		/>
		<div class="flex justify-between text-xs text-gray-500 mt-1">
			<span>0%</span>
			<span>50%</span>
			<span>100%</span>
		</div>
	</div>
</div>

<style>
	.slider::-webkit-slider-thumb {
		appearance: none;
		height: 20px;
		width: 20px;
		border-radius: 50%;
		background: #3b82f6;
		cursor: pointer;
		border: 2px solid white;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.slider::-moz-range-thumb {
		height: 20px;
		width: 20px;
		border-radius: 50%;
		background: #3b82f6;
		cursor: pointer;
		border: 2px solid white;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.slider::-webkit-slider-track {
		height: 8px;
		border-radius: 4px;
		background: linear-gradient(
			to right,
			#3b82f6 0%,
			#3b82f6 var(--value),
			#e5e7eb var(--value),
			#e5e7eb 100%
		);
	}
</style>