<script lang="ts">
  import type { ExperimentVariant } from "$lib/types.js";

  interface Props {
    variants: ExperimentVariant[];
    className?: string;
  }

  let { variants, className = "" }: Props = $props();

  const validVariants = $derived(
    variants.filter(v => v.results && (v.results.impressions > 0 || v.results.conversions > 0))
  );

  const maxImpressions = $derived(
    Math.max(...validVariants.map(v => v.results?.impressions || 0), 1)
  );

  const maxConversions = $derived(
    Math.max(...validVariants.map(v => v.results?.conversions || 0), 1)
  );

  function getConversionRate(variant: ExperimentVariant): number {
    if (!variant.results || variant.results.impressions === 0) return 0;
    return (variant.results.conversions / variant.results.impressions) * 100;
  }

  function formatNumber(num: number): string {
    if (num >= 1000000) {
      return (num / 1000000).toFixed(1) + "M";
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + "K";
    }
    return num.toLocaleString();
  }

  const chartHeight = 200;
  const barWidth = 60;
  const spacing = 40;
  const chartWidth = $derived(validVariants.length * (barWidth * 2 + spacing) + spacing);

  const totalImpressions = $derived(
    validVariants.reduce((sum, v) => sum + (v.results?.impressions || 0), 0)
  );

  const totalConversions = $derived(
    validVariants.reduce((sum, v) => sum + (v.results?.conversions || 0), 0)
  );

  const overallRate = $derived(
    totalImpressions > 0 ? (totalConversions / totalImpressions) * 100 : 0
  );
</script>

<div class="experiment-chart {className}">
  {#if validVariants.length === 0}
    <div class="text-center py-8 text-gray-500">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      </svg>
      <p class="mt-2 text-sm">No experiment data to display</p>
    </div>
  {:else}
    <div class="bg-white rounded-lg border border-gray-200 p-4">
      <h4 class="text-sm font-medium text-gray-900 mb-4">Variant Performance</h4>

      <!-- Chart -->
      <div class="overflow-x-auto">
        <svg width={chartWidth} height={chartHeight + 100} class="min-w-full">
          <!-- Background grid lines -->
          {#each [0.2, 0.4, 0.6, 0.8, 1.0] as ratio}
            <line
              x1="20"
              y1={20 + (chartHeight - 40) * (1 - ratio)}
              x2={chartWidth - 20}
              y2={20 + (chartHeight - 40) * (1 - ratio)}
              stroke="#f3f4f6"
              stroke-width="1"
            />
          {/each}

          {#each validVariants as variant, index}
            {@const xPos = spacing + index * (barWidth * 2 + spacing)}
            {@const impressionHeight = Math.max(((variant.results?.impressions || 0) / maxImpressions) * (chartHeight - 40), 2)}
            {@const conversionHeight = Math.max(((variant.results?.conversions || 0) / maxConversions) * (chartHeight - 40), 2)}
            {@const conversionRate = getConversionRate(variant)}

            <!-- Impressions bar -->
            <rect
              x={xPos}
              y={chartHeight - impressionHeight - 20}
              width={barWidth}
              height={impressionHeight}
              fill="#3b82f6"
              rx="2"
            />

            <!-- Conversions bar -->
            <rect
              x={xPos + barWidth + 10}
              y={chartHeight - conversionHeight - 20}
              width={barWidth}
              height={conversionHeight}
              fill="#10b981"
              rx="2"
            />

            <!-- Impressions value label -->
            <text
              x={xPos + barWidth / 2}
              y={chartHeight - impressionHeight - 25}
              text-anchor="middle"
              class="fill-gray-700 text-xs font-medium"
            >
              {formatNumber(variant.results?.impressions || 0)}
            </text>

            <!-- Conversions value label -->
            <text
              x={xPos + barWidth + 10 + barWidth / 2}
              y={chartHeight - conversionHeight - 25}
              text-anchor="middle"
              class="fill-gray-700 text-xs font-medium"
            >
              {formatNumber(variant.results?.conversions || 0)}
            </text>

            <!-- Variant name -->
            <text
              x={xPos + barWidth}
              y={chartHeight + 15}
              text-anchor="middle"
              class="fill-gray-900 text-xs font-medium"
            >
              {variant.name}
            </text>

            <!-- Conversion rate -->
            <text
              x={xPos + barWidth}
              y={chartHeight + 30}
              text-anchor="middle"
              class="fill-gray-600 text-xs"
            >
              {conversionRate.toFixed(1)}% CR
            </text>

            <!-- Weight -->
            <text
              x={xPos + barWidth}
              y={chartHeight + 45}
              text-anchor="middle"
              class="fill-gray-500 text-xs"
            >
              {variant.weight}% traffic
            </text>
          {/each}
        </svg>
      </div>

      <!-- Legend -->
      <div class="flex items-center justify-center space-x-6 mt-4 pt-4 border-t border-gray-200">
        <div class="flex items-center space-x-2">
          <div class="w-3 h-3 bg-blue-500 rounded-sm"></div>
          <span class="text-xs text-gray-600">Impressions</span>
        </div>
        <div class="flex items-center space-x-2">
          <div class="w-3 h-3 bg-green-500 rounded-sm"></div>
          <span class="text-xs text-gray-600">Conversions</span>
        </div>
      </div>

      <!-- Summary stats -->
      <div class="grid grid-cols-1 gap-4 mt-4 pt-4 border-t border-gray-200 sm:grid-cols-3">
        <div class="text-center">
          <div class="text-lg font-semibold text-gray-900">
            {formatNumber(totalImpressions)}
          </div>
          <div class="text-xs text-gray-500">Total Impressions</div>
        </div>
        <div class="text-center">
          <div class="text-lg font-semibold text-gray-900">
            {formatNumber(totalConversions)}
          </div>
          <div class="text-xs text-gray-500">Total Conversions</div>
        </div>
        <div class="text-center">
          <div class="text-lg font-semibold text-gray-900">
            {overallRate.toFixed(1)}%
          </div>
          <div class="text-xs text-gray-500">Overall CR</div>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .experiment-chart {
    @apply w-full;
  }
</style>