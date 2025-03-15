<template>
  <div
    v-if="visible"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
  >
    <div class="bg-gray-900 rounded-lg shadow-lg p-8 max-w-lg w-full">
      <h3 class="text-3xl font-bold text-gray-100 mb-6">{{ title }}</h3>
      <p class="text-lg text-gray-300 mb-8">{{ message }}</p>
      <div class="flex justify-end space-x-4">
        <button
          @click="cancel"
          class="px-6 py-3 bg-gray-600 text-white rounded hover:bg-gray-500"
        >
          Cancelar
        </button>
        <button
          @click="confirm"
          class="px-6 py-3 bg-indigo-600 text-white rounded hover:bg-indigo-700"
        >
          Aceptar
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, defineProps, defineEmits } from "vue";

const props = defineProps({
  title: { type: String, default: "Confirmación" },
  message: { type: String, default: "Por favor, confirma para continuar." },
});
const emits = defineEmits(["confirm", "cancel"]);

const visible = ref(true);

function confirm() {
  emits("confirm");
  visible.value = false;
}

function cancel() {
  emits("cancel");
  visible.value = false;
}
</script>

<style scoped>
/* Puedes ajustar o agregar estilos adicionales si lo deseas */
</style>
