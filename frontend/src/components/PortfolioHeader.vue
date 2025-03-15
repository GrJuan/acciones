<template>
  <div>
    <!-- Header del Balance -->
    <div
      v-motion="portfolioMotion"
      style="margin-bottom: 1.5rem"
      class="flex items-center justify-between p-4 bg-gray-900 rounded-lg shadow-lg border border-gray-700 relative"
    >
      <div
        class="absolute -top-3 left-4 bg-green-500 p-2 rounded-full shadow-md"
      >
        <Icon name="mdi:cash-multiple" class="text-white text-2xl" />
      </div>

      <div class="text-gray-100 flex-1 pl-10">
        <h2 class="text-2xl font-bold text-green-400 flex items-center">
          ${{ balance.toFixed(2) }}
          <Icon name="mdi:trending-up" class="ml-2 text-lg text-green-400" />
        </h2>
        <p class="text-gray-400 text-sm">Disponible para operaciones</p>
      </div>

      <!-- Botón para abrir el modal de recarga -->
      <button
        @click="showModal = true"
        class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md flex items-center text-sm transition-all"
      >
        <Icon name="mdi:plus-circle-outline" class="mr-2 text-lg" />
        Recargar
      </button>
    </div>

    <!-- Modal de Recarga -->
    <div
      v-if="showModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div
        class="bg-gray-800 p-6 rounded-lg shadow-lg w-[400px] border border-gray-700 relative"
      >
        <button
          @click="showModal = false"
          class="absolute top-2 right-2 text-gray-400 hover:text-gray-200 text-2xl"
        >
          <Icon name="mdi:close-circle-outline" />
        </button>

        <h2 class="text-xl font-bold text-gray-100 mb-4 flex items-center">
          <Icon
            name="mdi:credit-card-outline"
            class="mr-2 text-blue-400 text-xl"
          />
          Recargar Balance
        </h2>

        <p class="text-gray-400 text-sm mb-4">
          Selecciona una cantidad para recargar tu cuenta.
        </p>

        <!-- Botones de montos rápidos -->
        <div class="grid grid-cols-3 gap-3 mb-4">
          <button
            v-for="amount in quickAmounts"
            :key="amount"
            @click="selectedAmount = amount"
            :class="[
              'py-2 rounded-md text-white text-sm transition-all',
              selectedAmount === amount
                ? 'bg-green-600'
                : 'bg-gray-700 hover:bg-gray-600',
            ]"
          >
            ${{ amount }}
          </button>
        </div>

        <!-- Input para monto personalizado -->
        <input
          type="number"
          v-model="customAmount"
          placeholder="Otro monto"
          class="w-full p-2 bg-gray-700 text-white border border-gray-600 rounded-md text-center"
          @focus="selectedAmount = null"
        />

        <!-- Botón de Confirmar Recarga -->
        <button
          @click="confirmRecharge"
          class="mt-4 w-full bg-green-600 hover:bg-green-700 text-white py-2 rounded-md transition-all"
        >
          Confirmar Recarga
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useDashboardData } from "@/composables/dashboardData";

const { balance, updateBalance } = useDashboardData();
const showModal = ref(false);
const selectedAmount = ref<number | null>(null);
const customAmount = ref<string>("");

// Montos rápidos
const quickAmounts = [10, 50, 100, 500, 1000, 5000];

const portfolioMotion = {
  initial: { opacity: 0, y: -20 },
  enter: { opacity: 1, y: 0, transition: { duration: 0.5 } },
};

function confirmRecharge() {
  const amount = selectedAmount.value || parseFloat(customAmount.value);
  if (!amount || amount <= 0) return;

  // Simula la recarga del saldo
  updateBalance(amount);

  // Cierra el modal después de un breve delay
  setTimeout(() => {
    showModal.value = false;
    selectedAmount.value = null;
    customAmount.value = "";
  }, 800);
}
</script>
