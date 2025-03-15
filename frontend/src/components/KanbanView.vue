<template>
  <div class="flex space-x-4 overflow-x-auto pb-4">
    <div
      v-for="(items, rating) in kanbanGroups"
      :key="rating"
      class="w-80 flex-shrink-0 bg-gray-800 rounded-lg shadow p-4 border border-gray-700"
    >
      <h2
        :class="[
          'text-center text-lg font-bold py-2 mb-4 rounded-md',
          getKanbanHeaderClass(rating),
        ]"
      >
        {{ rating }}
      </h2>
      <div
        v-for="item in items"
        :key="item.ticker"
        class="bg-gray-700 rounded-md p-3 mb-4 shadow-sm border border-gray-600"
      >
        <h3 class="font-semibold text-gray-100 mb-1">{{ item.ticker }}</h3>
        <p class="text-sm text-gray-300 mb-1">{{ item.company }}</p>
        <p class="text-sm text-gray-300">
          <strong>Objetivo:</strong> {{ item.target_from }} -
          {{ item.target_to }}
        </p>
        <p class="text-sm text-gray-300">
          <strong>Fecha:</strong> {{ formatDate(item.time) }}
        </p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useDashboardData } from "@/composables/dashboardData";
const { kanbanGroups, getKanbanHeaderClass, formatDate } = useDashboardData();
</script>
