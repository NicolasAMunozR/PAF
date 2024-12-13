<template>
  <div class="login-page">
    <!-- Barra superior con imagotipo y título -->
    <div class="header">
      <img src="@/assets/logo.png" alt="Logo" class="logo" />
      <h1>Inicio de Sesión</h1>
    </div>

    <!-- Selección de rol -->
    <div class="role-selection">
      <button
        v-for="option in roleOptions"
        :key="option.value"
        :class="{ active: selectedRole === option.value }"
        @click="selectRole(option.value)"
      >
        {{ option.label }}
      </button>
    </div>

    <!-- Formulario de inicio de sesión -->
    <form @submit.prevent="handleLogin" v-if="selectedRole" class="login-form">
      <div class="form-group">
        <label for="run">RUN:</label>
        <input
          type="text"
          id="run"
          v-model="run"
          placeholder="Ingresa tu RUN"
          required
        />
      </div>
      <p class="role-info">
        Iniciarás sesión como: <strong>{{ selectedRole }}</strong>
      </p>
      <button type="submit">Ingresar</button>
    </form>

    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

export default {
  layout: false,
  setup() {
    const router = useRouter()
    const run = ref('')
    const selectedRole = ref('')
    const errorMessage = ref('')
    const roleOptions = ref([
      { value: 'personal-dei', label: 'Personal del Dei' },
      { value: 'encargado', label: 'Encargado' },
    ])

    const selectRole = (role) => {
      selectedRole.value = role
      errorMessage.value = ''
    }

    const handleLogin = async () => {
      if (run.value && selectedRole.value) {
        sessionStorage.setItem('rut', run.value) // Guardar en sesión
        try {
          const response = await axios.post('/api/login', {
            run: run.value,
            role: selectedRole.value,
          })
          if (response.data.success) {
            if (selectedRole.value === 'personal-dei') {
              router.push('/personas')
            } else if (selectedRole.value === 'encargado') {
              router.push('/encargado')
            }
          } else {
            errorMessage.value = 'El rol seleccionado no coincide con el usuario.'
          }
        } catch (error) {
          errorMessage.value = 'Hubo un error al iniciar sesión.'
          console.error(error)
        }
      } else {
        errorMessage.value = 'Por favor, completa todos los campos.'
      }
    }

    return {
      run,
      selectedRole,
      errorMessage,
      roleOptions,
      selectRole,
      handleLogin,
    }
  },
}
</script>

<style scoped>
.login-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.header {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 20px;
}

.logo {
  width: 100px;
  height: auto;
}

.role-selection {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.role-selection button {
  padding: 10px 20px;
  border: none;
  background-color: #394049;
  color: white;
  cursor: pointer;
}

.role-selection button.active {
  background-color: #ea7600;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  margin-bottom: 5px;
}

.form-group input {
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.role-info {
  margin-bottom: 10px;
}

.error {
  color: red;
  margin-top: 10px;
}
</style>