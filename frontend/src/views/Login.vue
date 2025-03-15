<template>
  <div class="login-wrapper">
    <div class="login-container">
      <!-- Lado izquierdo: Frase motivacional -->
      <div class="investment-quote">
        <h1>Invierte en tu futuro</h1>
        <p>
          El mejor momento para empezar es ahora. ¡Haz crecer tu dinero con
          nosotros!
        </p>
      </div>

      <!-- Lado derecho: Formulario de login -->
      <div class="login-form">
        <h2 class="text-xl font-bold text-gray-900 mb-4 text-center">
          Iniciar Sesión
        </h2>
        <form @submit.prevent="handleLogin">
          <div class="form-group">
            <label for="username">Usuario</label>
            <input
              type="text"
              id="username"
              v-model="username"
              placeholder="Ingresa tu usuario"
              required
            />
          </div>
          <div class="form-group">
            <label for="password">Contraseña</label>
            <input
              type="password"
              id="password"
              v-model="password"
              placeholder="Ingresa tu contraseña"
              required
            />
          </div>
          <button type="submit">Iniciar Sesión</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useRouter } from "vue-router";

const username = ref("");
const password = ref("");
const router = useRouter();

function generateUserID(): string {
  return "user-" + Date.now() + "-" + Math.floor(Math.random() * 1000);
}

async function handleLogin() {
  try {
    const response = await fetch(
      "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/login",
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: username.value,
          password: password.value,
        }),
      }
    );

    if (response.ok) {
      const data = await response.json();
      localStorage.setItem("auth_token", data.auth_token);
      localStorage.setItem("current_user_id", generateUserID());

      router.push("/");
      setTimeout(() => {
        window.location.reload();
      }, 200); // 🔄 Forzar recarga
    } else {
      alert("Error en el login. Verifica tus credenciales.");
    }
  } catch (error) {
    alert("Ocurrió un error al intentar iniciar sesión.");
  }
}
</script>

<style scoped>
/* Fondo con degradado */
.login-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1f2937, #111827);
  padding: 20px;
}

/* Contenedor principal */
.login-container {
  display: flex;
  border-radius: 12px;
  overflow: hidden;
  width: 900px;
  max-width: 100%;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
}

/* Sección de inversión (lado izquierdo) */
.investment-quote {
  flex: 1;
  padding: 40px;
  background: #1e293b;
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
  text-align: center;
}

.investment-quote h1 {
  font-size: 2.5rem;
  margin-bottom: 15px;
  font-weight: bold;
}

.investment-quote p {
  font-size: 1.2rem;
  opacity: 0.9;
}

/* Sección del formulario (lado derecho) */
.login-form {
  flex: 1;
  padding: 40px;
  background: white;
  color: #333;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

/* Formulario */
.login-form form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Grupo de inputs */
.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  font-size: 1rem;
  margin-bottom: 5px;
  font-weight: 600;
  color: #374151;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #c3c3c3;
  border-radius: 8px;
  font-size: 1rem;
  transition: border 0.2s ease-in-out;
}

.form-group input:focus {
  border-color: #6366f1;
  outline: none;
}

/* Botón */
button {
  background-color: #6366f1;
  color: white;
  padding: 12px;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.3s ease-in-out, transform 0.1s ease-in-out;
}

button:hover {
  background-color: #4f46e5;
  transform: scale(1.05);
}
</style>
