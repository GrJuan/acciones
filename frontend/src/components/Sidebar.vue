<template>
  <aside
    class="bg-gray-800 text-gray-200 p-4 shadow-md transition-all duration-300 relative flex flex-col h-screen"
    :class="isExpanded ? 'w-56' : 'w-20'"
    @mouseenter="expandSidebar"
    @mouseleave="collapseSidebar"
  >
    <!-- Botón del Asistente IA en la parte superior -->
    <div class="flex-none mb-4">
      <button
        @click="openAssistant"
        :class="
          isExpanded
            ? 'w-full px-2 py-1 bg-[#16ac86] text-white rounded text-sm hover:bg-[#149a75] transition-colors focus:outline-none focus:ring-0'
            : 'w-10 h-10 flex items-center justify-center bg-[#16ac86] rounded-full hover:bg-[#149a75] transition-colors focus:outline-none focus:ring-0'
        "
      >
        <template v-if="isExpanded"> Asistente IA </template>
        <template v-else>
          <img
            src="https://static.vecteezy.com/system/resources/previews/024/558/807/non_2x/openai-chatgpt-logo-icon-free-png.png"
            alt="ChatGPT"
            class="w-8 h-8 rounded-full border-0 focus:outline-none focus:ring-0"
          />
        </template>
      </button>

      <AIChatAssistant ref="assistantRef" />
    </div>

    <!-- Espacio flexible para empujar la lista de usuarios al fondo -->
    <div class="flex-1"></div>

    <!-- Lista de usuarios en la parte inferior -->
    <div class="space-y-4">
      <div
        v-for="user in sortedUsers"
        :key="user.id"
        class="flex items-center justify-start"
        style="height: 40px"
      >
        <!-- Contenedor del avatar con indicador de conexión -->
        <div
          class="relative w-10 h-10 shrink-0 cursor-pointer"
          @click="handleAvatarClick(user)"
        >
          <img
            :src="user.avatar"
            alt="Avatar"
            class="w-10 h-10 rounded-full border-[3px] border-green-500"
          />

          <!-- Indicador de conexión: un pequeño círculo verde -->
          <span
            class="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full border border-gray-800"
          ></span>
          <span
            v-if="user.id !== currentUserId && getUnreadCount(user.id) > 0"
            class="absolute -top-1 -right-1 bg-red-600 text-xs w-5 h-5 rounded-full flex items-center justify-center"
          >
            {{ getUnreadCount(user.id) }}
          </span>
        </div>
        <!-- Botón "Enviar" -->
        <div v-if="isExpanded && user.id !== currentUserId" class="ml-2">
          <button
            @click.stop="openChatForUser(user)"
            class="px-2 py-2 bg-[#25D366] rounded text-white text-sm font-semibold hover:bg-[#1ebe59] transition-colors flex items-center justify-center"
          >
            <span>Enviar Mensaje</span>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 ml-1"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M17 8h2a2 2 0 012 2v8a2 2 0 01-2 2h-6l-4 4v-4H5a2 2 0 01-2-2V10a2 2 0 012-2h2"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Modal de confirmación sofisticado -->
    <ConfirmModal
      v-if="showConfirmModal"
      title="Iniciar Chat Privado"
      :message="confirmMessage"
      @confirm="confirmChat"
      @cancel="cancelChat"
    />
  </aside>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue";
import { useDashboardData } from "@/composables/dashboardData";
import { chatConversations, conversationKey } from "@/composables/chatData";
import AIChatAssistant from "./AIChatAssistant.vue";
import ConfirmModal from "./ConfirmModal.vue";

const { connectedUsers } = useDashboardData();
const currentUserId = localStorage.getItem("current_user_id") || "";
const isExpanded = ref(false);

// Control del modal de confirmación
const showConfirmModal = ref(false);
const selectedUser = ref<any>(null);
const confirmMessage = ref(
  "Antes de iniciar el chat, recuerda que es importante mantener siempre el respeto y la cortesía en la comunicación. ¿Deseas continuar?"
);

const sortedUsers = computed(() => {
  return [...connectedUsers.value];
});

function expandSidebar() {
  if (connectedUsers.value.length > 1) {
    isExpanded.value = true;
  }
}
function collapseSidebar() {
  isExpanded.value = false;
}

function handleAvatarClick(user: any) {
  if (user.id !== currentUserId) {
    selectedUser.value = user;
  }
}

function openChatForUser(user: any) {
  selectedUser.value = user;
  showConfirmModal.value = true;
}

function confirmChat() {
  window.dispatchEvent(
    new CustomEvent("open-chat", { detail: selectedUser.value })
  );
  showConfirmModal.value = false;
}
function cancelChat() {
  showConfirmModal.value = false;
}

function getUnreadCount(otherUserId: string): number {
  const key = conversationKey(otherUserId, currentUserId);
  return chatConversations.value[key]?.unread || 0;
}

const assistantRef = ref<InstanceType<typeof AIChatAssistant> | null>(null);
function openAssistant() {
  assistantRef.value?.open();
}

// Reproducción de sonido para mensajes entrantes
const newMessageSound = new Audio("/src/assets/notification-sound.mp3");
onMounted(() => {
  window.addEventListener("notification-sound", () => {
    newMessageSound.play();
  });
});
</script>

<style scoped>
/* Ajusta el indicador de conexión si deseas modificar su tamaño o estilo */
</style>
