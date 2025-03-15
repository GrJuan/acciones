<template>
  <div
    class="stock-item relative bg-gray-800 rounded-lg shadow-lg p-5 border transition-all duration-300 transform hover:scale-105"
    :class="getBorderClass()"
  >
    <!-- Encabezado con el ticker y botones -->
    <div class="flex items-center justify-between mb-3">
      <h2
        :class="{ 'line-through': action === 'sell' && sold }"
        class="text-xl font-semibold text-white flex items-center"
      >
        <component :is="getActionIcon()" class="w-5 h-5 mr-2" />
        {{ item.ticker }}
      </h2>
      <div class="flex space-x-2">
        <button
          v-if="action === 'buy'"
          @click="handleBuy"
          class="btn-action bg-green-600 hover:bg-green-700"
        >
          Buy
        </button>
        <button
          v-if="action === 'sell'"
          @click="handleSell"
          class="btn-action bg-red-600 hover:bg-red-700"
        >
          Sell
        </button>
      </div>
    </div>

    <!-- Información de la empresa -->
    <p class="text-sm text-gray-400 mb-2">{{ item.company }}</p>

    <!-- Detalles de la acción -->
    <div class="text-sm text-gray-300 space-y-1">
      <p>
        <strong>📊 Rating:</strong>
        <span class="text-yellow-400">{{ item.rating_to }}</span>
      </p>
      <p>
        <strong>🎯 Objetivo:</strong>
        <span class="text-green-400">{{ item.target_from }}</span> -
        <span class="text-red-400">{{ item.target_to }}</span>
      </p>
    </div>

    <!-- Animación de compra/venta -->
    <div
      v-if="animationVisible"
      class="absolute inset-0 flex items-center justify-center pointer-events-none"
    >
      <div
        v-motion="animationMotion"
        class="bg-white bg-opacity-80 px-5 py-2 rounded-lg text-gray-800 font-bold shadow-lg"
      >
        {{ animationMessage }}
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import {
  ArrowUpIcon,
  ArrowDownIcon,
  CheckCircleIcon,
} from "@heroicons/vue/24/solid";
import confetti from "canvas-confetti";

// Props y eventos
const props = defineProps<{ item: any; action: string }>();
const emit = defineEmits<{ (e: "buy"): void; (e: "sell"): void }>();

// Estados
const sold = ref(false);
const animationVisible = ref(false);
const animationMessage = ref("");

const animationMotion = {
  initial: { opacity: 0, scale: 0.8 },
  enter: { opacity: 1, scale: 1, transition: { duration: 0.5 } },
  leave: { opacity: 0, scale: 0.8, transition: { duration: 0.5 } },
};

// Determina el color del borde según la acción
function getBorderClass() {
  return props.action === "buy" ? "border-green-500" : "border-red-500";
}

// Obtiene el ícono de la acción
function getActionIcon() {
  return props.action === "buy" ? ArrowUpIcon : ArrowDownIcon;
}

// Maneja la compra de una acción
function handleBuy() {
  sold.value = false;
  animationMessage.value = "✅ Bought!";
  showAnimation();
  emit("buy");
}

// Maneja la venta de una acción
function handleSell() {
  sold.value = true;
  animationMessage.value = "❌ Sold!";
  showAnimation();

  // Confetti para ventas altas
  const salePrice = parseFloat(props.item.target_to.replace("$", ""));
  if (salePrice > 200) {
    confetti({ particleCount: 100, spread: 70, origin: { y: 0.6 } });
  }
  emit("sell");
}

// Muestra la animación de compra/venta
function showAnimation() {
  animationVisible.value = true;
  setTimeout(() => {
    animationVisible.value = false;
  }, 1500);
}
</script>

<style scoped>
/* Botón de acción */
.btn-action {
  color: white;
  padding: 0.5rem 0.75rem;
  border-radius: 0.375rem;
  transition: background-color 0.2s ease-in-out;
  font-size: 0.875rem;
  font-weight: 500;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
