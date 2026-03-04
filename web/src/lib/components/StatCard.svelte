<script lang="ts">
interface Props {
  label: string;
  value: number | string;
  trend?: {
    value: number;
    direction: "up" | "down";
  };
  icon?: string;
  loading?: boolean;
}

let { label, value, trend, icon, loading = false }: Props = $props();
</script>

<div class="bg-white rounded-lg border border-gray-200 p-6 shadow-sm">
  <div class="flex items-center justify-between">
    <div class="flex-1">
      <div class="flex items-center space-x-2">
        {#if icon}
          <span class="text-2xl">{icon}</span>
        {/if}
        <p class="text-sm font-medium text-gray-600">{label}</p>
      </div>

      {#if loading}
        <div class="mt-2">
          <div class="h-8 bg-gray-200 rounded animate-pulse w-20"></div>
        </div>
      {:else}
        <p class="text-3xl font-bold text-gray-900 mt-2">{value}</p>
      {/if}

      {#if trend && !loading}
        <div class="flex items-center mt-2">
          <div class="flex items-center space-x-1 text-sm {trend.direction === 'up' ? 'text-green-600' : 'text-red-600'}">
            {#if trend.direction === "up"}
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M3.293 9.707a1 1 0 010-1.414l6-6a1 1 0 011.414 0l6 6a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L4.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
              </svg>
            {:else}
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 10.293a1 1 0 010 1.414l-6 6a1 1 0 01-1.414 0l-6-6a1 1 0 111.414-1.414L9 14.586V3a1 1 0 012 0v11.586l4.293-4.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            {/if}
            <span>{Math.abs(trend.value)}%</span>
          </div>
          <span class="text-gray-500 text-sm ml-2">vs last month</span>
        </div>
      {/if}
    </div>
  </div>
</div>