<template>
    <div class="login-page">
      <h1>Inicio de Sesión</h1>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="email">Correo Electrónico:</label>
          <input type="email" id="email" v-model="email" placeholder="Ingresa tu correo" required />
        </div>
        <div class="form-group">
          <label for="password">Contraseña:</label>
          <input type="password" id="password" v-model="password" placeholder="Ingresa tu contraseña" required />
        </div>
        <div class="form-group">
          <label for="role">Selecciona tu Rol:</label>
          <select id="role" v-model="role" required>
            <option value="" disabled>Seleccione una opción</option>
            <option value="profesor">Profesor</option>
            <option value="personal-dei">Personal del DEI</option>
            <option value="encargado">Encargado</option>
          </select>
        </div>
        <button type="submit">Ingresar</button>
      </form>
      <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
    </div>
  </template>
  
  <script>
  export default {
    layout: false, // Oculta el layout para esta página
    data() {
      return {
        email: "",
        password: "",
        role: "",
        errorMessage: "",
      };
    },
    methods: {
      handleLogin() {
        if (this.email && this.password && this.role) {
          // Lógica para manejar el inicio de sesión según el rol seleccionado
          alert(`Iniciando sesión como ${this.role}`);
          // Aquí puedes redirigir o hacer la lógica necesaria según el rol
        } else {
          this.errorMessage = "Por favor, completa todos los campos.";
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .login-page {
    max-width: 400px;
    margin: 0 auto;
    padding: 20px;
    border: 1px solid #ddd;
    border-radius: 8px;
    background-color: #f9f9f9;
  }
  
  h1 {
    text-align: center;
    margin-bottom: 20px;
  }
  
  .form-group {
    margin-bottom: 15px;
  }
  
  label {
    display: block;
    font-weight: bold;
    margin-bottom: 5px;
  }
  
  input,
  select {
    width: 100%;
    padding: 8px;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  
  button {
    width: 100%;
    padding: 10px;
    background-color: #007bff;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button:hover {
    background-color: #0056b3;
  }
  
  .error {
    color: red;
    margin-top: 10px;
    text-align: center;
  }
  </style>
  