<template>
  <div class="relative pl-6 border-l-2 border-gray-700">
    <div
      v-for="item in sortedItemsByTime"
      :key="item.ticker"
      class="mb-6 relative"
    >
      <span
        class="absolute -left-3 top-3 w-3 h-3 rounded-full bg-gray-800 border-2 border-white"
      ></span>
      <div class="bg-gray-800 p-4 rounded-lg shadow border border-gray-700">
        <h3 class="text-lg font-semibold text-gray-100 mb-2">
          {{ item.ticker }} - {{ item.company }}
        </h3>
        <p class="text-xs text-gray-400 mb-2">{{ formatDate(item.time) }}</p>
        <p class="text-sm text-gray-200">
          <strong>Rating:</strong> {{ item.rating_from }} → {{ item.rating_to }}
        </p>
        <p class="text-sm text-gray-200">
          <strong>Objetivo:</strong> {{ item.target_from }} -
          {{ item.target_to }}
        </p>
        <p class="text-sm text-gray-200">
          <strong>Acción:</strong> {{ item.action }}
        </p>
        <p class="text-sm text-gray-200">
          <strong>Corretaje:</strong> {{ item.brokerage }}
        </p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useDashboardData } from "@/composables/dashboardData";
const { sortedItemsByTime, formatDate } = useDashboardData();
</script>
