<template>
  <div
    v-if="visible"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
  >
    <div class="bg-gray-900 rounded-lg shadow-lg w-full max-w-lg p-6 relative">
      <button
        @click="close"
        class="absolute top-4 right-4 text-gray-400 hover:text-gray-200"
      >
        &times;
      </button>
      <h2 class="text-xl font-bold text-gray-100 mb-4">
        Chat con {{ chatPartner?.id }}
      </h2>
      <div
        class="messages h-64 overflow-y-auto bg-gray-800 rounded p-4 text-gray-200 mb-4"
      >
        <div
          v-for="(msg, index) in conversation.messages"
          :key="index"
          class="mb-2"
        >
          <span class="font-bold">{{ msg.sender }}:</span>
          <span>{{ msg.text }}</span>
        </div>
      </div>
      <form @submit.prevent="sendMessage" class="flex">
        <input
          v-model="newMessage"
          type="text"
          placeholder="Escribe un mensaje..."
          class="flex-1 p-2 rounded-l-md bg-gray-700 border border-gray-600 text-gray-200"
          required
        />
        <button
          type="submit"
          class="px-4 bg-indigo-600 text-white rounded-r-md hover:bg-indigo-700"
        >
          Enviar
        </button>
      </form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from "vue";
import {
  chatConversations,
  markAsRead,
  addMessage,
  ChatMessage,
  conversationKey,
} from "@/composables/chatData";
import { useDashboardData } from "@/composables/dashboardData";

const visible = ref(false);
const newMessage = ref("");
const chatPartner = ref<any>(null);
const conversation = ref<{ messages: ChatMessage[]; unread: number }>({
  messages: [],
  unread: 0,
});
const currentUserId = localStorage.getItem("current_user_id") || "Anónimo";
const { socket } = useDashboardData();

// Escucha el evento global para abrir el chat
// En ChatWindow.vue, en el listener del evento 'open-chat'
window.addEventListener("open-chat", (event: any) => {
  // Primero, asignamos el usuario seleccionado
  chatPartner.value = event.detail;

  // Generamos la clave única para la conversación entre el usuario actual y el chatPartner
  const key = conversationKey(chatPartner.value.id, currentUserId);

  // Marcamos la conversación como leída usando la clave
  markAsRead(key);

  // Si no existe la conversación para ese key, la inicializamos
  if (!chatConversations.value[key]) {
    chatConversations.value[key] = { messages: [], unread: 0 };
  }

  // Asignamos la conversación global a la variable local
  conversation.value = chatConversations.value[key];

  // Mostramos el chat
  visible.value = true;
});

// Watcher para actualizar la conversación cuando cambie en el estado global
watch(
  () => chatConversations.value[chatPartner.value?.id],
  (newConv) => {
    console.log(
      "Conversación actualizada para:",
      chatPartner.value?.id,
      newConv
    );
    if (newConv) {
      conversation.value = newConv;
    }
  },
  { deep: true }
);

function sendMessage() {
  if (!newMessage.value.trim()) return;
  const messagePayload = {
    event: "private-message",
    sender: currentUserId,
    recipient: chatPartner.value.id,
    text: newMessage.value,
  };
  if (socket.value && socket.value.readyState === WebSocket.OPEN) {
    socket.value.send(JSON.stringify(messagePayload));
    // Agregar mensaje localmente
    conversation.value.messages.push({
      sender: currentUserId,
      text: newMessage.value,
      timestamp: Date.now(),
    });
    newMessage.value = "";
  } else {
    console.error("Socket no está disponible o no está abierto");
  }
}

function close() {
  visible.value = false;
}
</script>

<style scoped>
/* Estilos adicionales para ChatWindow */
</style>
