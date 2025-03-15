// src/composables/chatData.ts
import { ref } from 'vue';

export interface ChatMessage {
    sender: string;
    text: string;
    timestamp: number;
}

export interface Conversation {
    messages: ChatMessage[];
    unread: number;
}

// Estado global de conversaciones: clave = ID del remitente
export const chatConversations = ref<Record<string, Conversation>>({});

// Función para agregar un mensaje a una conversación
export function addMessage(fromUserId: string, message: ChatMessage, isRecipient = true) {
    if (!chatConversations.value[fromUserId]) {
        chatConversations.value[fromUserId] = { messages: [], unread: 0 };
    }
    chatConversations.value[fromUserId].messages.push(message);
    if (isRecipient) {
        chatConversations.value[fromUserId].unread++;
    }
}

export function conversationKey(id1: string, id2: string): string {
    // Ordena los IDs alfabéticamente para que la clave sea la misma, independientemente del orden.
    const ids = [id1, id2].sort();
    return ids.join("-");
}


// Función para marcar los mensajes como leídos
export function markAsRead(fromUserId: string) {
    if (chatConversations.value[fromUserId]) {
        chatConversations.value[fromUserId].unread = 0;
    }
}
