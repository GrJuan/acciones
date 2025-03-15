<template>
  <div class="flex space-x-2 mb-6">
    <button
      @click="$emit('update:activeTab', 'cards')"
      :class="[
        'px-6 py-2 rounded-lg transition-colors',
        activeTab === 'cards'
          ? 'bg-gray-700 text-white'
          : 'bg-gray-600 text-gray-300 hover:bg-gray-500',
      ]"
    >
      Tarjetas
    </button>

    <button
      @click="
        portfolioStocks.length > 0 ? $emit('update:activeTab', 'charts') : null
      "
      :class="[
        'px-6 py-2 rounded-lg transition-colors',
        portfolioStocks.length > 0
          ? activeTab === 'charts'
            ? 'bg-gray-700 text-white'
            : 'bg-gray-600 text-gray-300 hover:bg-gray-500'
          : 'bg-gray-800 text-gray-500 cursor-not-allowed opacity-50',
      ]"
      :disabled="portfolioStocks.length === 0"
    >
      Gráficos
    </button>
  </div>
</template>

<script lang="ts" setup>
import { defineProps } from "vue";
import { useDashboardData } from "@/composables/dashboardData";

const props = defineProps({
  activeTab: {
    type: String,
    required: true,
  },
});

const { portfolioStocks } = useDashboardData(); // Accedemos a las acciones compradas
</script>
