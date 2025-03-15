<template>
  <div class="w-full overflow-hidden bg-gray-800 py-2">
    <div
      class="whitespace-nowrap animate-marquee text-gray-100 text-sm font-medium"
    >
      <span v-for="(stock, index) in stocks" :key="index" class="mx-8">
        {{ stock.symbol }}: ${{ stock.price.toFixed(2) }}
        <span v-if="stock.change >= 0" class="text-green-500">
          ▲ {{ stock.changesPercentage.toFixed(2) }}%
        </span>
        <span v-else class="text-red-500">
          ▼ {{ stock.changesPercentage.toFixed(2) }}%
        </span>
      </span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import axios from "axios";

interface Stock {
  symbol: string;
  price: number;
  change: number;
  changesPercentage: number;
}

const stocks = ref<Stock[]>([]);

async function fetchStockData() {
  // Lista de símbolos a consultar
  const symbols = "AAPL,GOOGL,AMZN,MSFT,TSLA";
  const apiKey = "SIRt7hBKcxaOztaxynxHf9AhYzVUCJ0F";
  try {
    const response = await axios.get(
      `https://financialmodelingprep.com/api/v3/quote/${symbols}?apikey=${apiKey}`
    );
    stocks.value = response.data;
  } catch (error) {
    console.error("Error fetching stock data:", error);
  }
}

onMounted(() => {
  fetchStockData();
});
</script>

<style scoped>
@keyframes marquee {
  0% {
    transform: translateX(0%);
  }
  100% {
    transform: translateX(-100%);
  }
}
.animate-marquee {
  display: inline-block;
  animation: marquee 15s linear infinite;
}
</style>
