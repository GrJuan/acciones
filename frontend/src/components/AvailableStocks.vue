<template>
  <div>
    <h2 class="text-xl font-bold text-gray-100 mb-4">Acciones Disponibles</h2>
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      <StockItem
        v-for="item in sortedAvailableStocks"
        :key="item.ticker"
        :item="item"
        action="buy"
        @buy="buyStock(item)"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useDashboardData } from "@/composables/dashboardData";
import StockItem from "./StockItem.vue";

const { sortedAvailableStocks, buyStock } = useDashboardData();
</script>
