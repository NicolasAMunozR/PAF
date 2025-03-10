<template>
    <div class="login-page">
      <!-- Barra superior con imagotipo y título -->
      <header class="header-bar">
        <h1>PAF</h1>
      </header>
  
      <!-- Contenido principal -->
      <div class="content">
        <h2>Inicio de Sesión</h2>
  
        <!-- Formulario de inicio de sesión -->
        <form @submit.prevent="handleLogin" class="login-form">
          <div class="form-group">
            <label for="username">username:</label>
            <input
              type="text"
              id="username"
              v-model="username"
              placeholder="Ingresa tu username"
              required
            />
          </div>
          <div class="form-group">
            <label for="password">password:</label>
            <input
              type="password"
              id="password"
              v-model="password"
              placeholder="Ingresa tu contraseña"
              required
            />
          </div>
          <button type="submit">Ingresar</button>
        </form>
  
        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
      </div>
      <!-- Modal para seleccionar rol -->
    <div v-if="showRoleModal" class="modal">
      <div class="modal-content">
        <h3>Selecciona un Rol</h3>
        <div class="form-group">
          <label for="role">Rol:</label>
          <select v-model="roleElegido" id="role">
            <option value="personal-dei">DIPRE</option>
            <option value="encargado">Encargado</option>
            <option value="profesor">Profesor</option>
          </select>
        </div>
        <div v-if="roleElegido === 'profesor'" class="form-group">
          <label for="newRut">Nuevo RUT:</label>
          <input
            type="text"
            id="newRut"
            v-model="newRut"
            placeholder="Ingresa el nuevo RUT"
          />
        </div>
        <button @click="confirmRole">Confirmar</button>
        <button @click="closeRoleModal">Cancelar</button>
      </div>
    </div>
    </div>
  </template>
  
  <script>
  export default {
    layout: false,
    data() {
      return {
        run: "",
        selectedRole: "",
        errorMessage: "",
        showRoleModal: false,
        roleElegido: "",
        username: "",
        password: "",
        newRut: "",
        roleOptions: [
          { value: "profesor", label: "Docente" },
          { value: "personal-dei", label: "DIPRE" },
          { value: "encargado", label: "Encargado" },
        ],
      };
    },
    methods: {
      selectRole(role) {
        this.selectedRole = role;
        this.errorMessage = "";
      },
      async handleLogin() {
        sessionStorage.clear();
        if (this.password && this.username) {
        const data = {
          User: this.username,
          Password: this.password
        };

        try {
          const response = await Promise.race([
            this.$axios.post("/api/paf-en-linea/login", data),
            new Promise((_, reject) => setTimeout(() => reject(new Error("timeout")), 5000))
          ]);
          const rut = response.data.data.rut;
          sessionStorage.setItem("rut", rut);
          console.log(response.data);
          const response2 = await this.$axios.get(`/api/paf-en-linea/pipelsoft/contratos-run/${rut}`);
          console.log(response2.data);
            if (!response2.data || response2.data.length === 0) {
              const response1 = await this.$axios.get(`/api/paf-en-linea/usuario/rut/${rut}`);
          console.log(response1.data);
          if (response1.data[0].Rol === "admin") {
            this.showRoleModal = true;
        }
          if(response1.data[0].Rol === "personal-dei") {
            this.$router.push("paf-en-linea/personas");
          }
          if(response1.data[0].Rol === "encargado") {
            if (response1.data[0].Acceso === 0) {
              this.errorMessage = "Este usuario no posee accesos.";
              return;
            }
            if (response1.data[0].Vista_universidad === 1) {
              this.$router.push("paf-en-linea/seguimientoPAF");
            } 
            if(response1.data[0].Vista_facultad === 1) {
              sessionStorage.setItem("unidadMayor", response1.data[0].UnidadMayor);
              this.$router.push(`paf-en-linea/unidadMayorPAF?UnidadMayor=${response1.data[0].UnidadMayor}`);
            }
            if (response1.data[0].Acceso === 1 && response1.data[0].Vista_facultad === 0 && response1.data[0].Vista_universidad === 0) {
              sessionStorage.setItem("unidadMayor", response1.data[0].UnidadMayor);
              sessionStorage.setItem("unidadMenor", response1.data[0].UnidadMenor);
            this.$router.push(`paf-en-linea/unidadMayorPAF?UnidadMayor=${response1.data[0].UnidadMayor}&UnidadMenor=${response1.data[0].UnidadMenor}`)
          }
        } 
            }

            else {
              this.$router.push(`paf-en-linea/profesorPAF?run=${rut}`);
              
            }
        } catch (error) {
          if (error.message === "timeout") {
            this.errorMessage = "La solicitud ha tardado demasiado. Por favor, inténtalo de nuevo.";
          } else {
            this.errorMessage = "Usuario o contraseña no válidos.";
          }
          return;
        }
      } 
  /*   
  if (this.run && this.selectedRole) {
    sessionStorage.setItem("rut", this.run); // Guardar en sesión
    try {
      let response;
      // Accede a $axios desde el contexto del componente
      if (this.selectedRole === "profesor") {
        response = await this.$axios.get(`/api/paf-en-linea/pipelsoft/contratos-run/${this.run}`);
      }
      else if (this.selectedRole === "personal-dei" || this.selectedRole === "encargado") {
        response = await this.$axios.get(`/api/paf-en-linea/usuario/rut/${this.run}`);
      }
      if (!response.data || response.data.length === 0) {
        this.errorMessage = "Usuario no encontrado.";
        return;
      }
      if (this.selectedRole === "profesor") {
        // Redirigir a la página de seguimiento de contratos
        this.$router.push(`paf-en-linea/profesorPAF?run=${this.run}`);
      } 
      // Verificar cada caso en el array
      const userMatch = response.data.find(user =>
        user.Rol === this.selectedRole &&
        (user.Rol === "encargado"  ||
         user.Rol === "personal-dei")
      );
        if (userMatch) {
          // Redirigir según el caso encontrado
        if (userMatch.Rol === "personal-dei") {
          this.$router.push("paf-en-linea/personas");
        }
        else if (userMatch.Rol === "encargado") {
          if (userMatch.Acceso === 0) {
            this.errorMessage = "Este usuario no posee accesos.";
            return;
          }
          else if (userMatch.Vista_universidad === 1) {
            this.$router.push("paf-en-linea/seguimientoPAF");
          } 
          else if(userMatch.Vista_facultad === 1) {
            sessionStorage.setItem("unidadMayor", userMatch.UnidadMayor);
            this.$router.push(`paf-en-linea/unidadMayorPAF?UnidadMayor=${userMatch.UnidadMayor}`);
          }
          else if (userMatch.Acceso === 1 && userMatch.Vista_facultad === 0 && userMatch.Vista_universidad === 0) {
            sessionStorage.setItem("unidadMayor", userMatch.UnidadMayor);
            sessionStorage.setItem("unidadMenor", userMatch.UnidadMenor);
          this.$router.push(`paf-en-linea/unidadMayorPAF?UnidadMayor=${userMatch.UnidadMayor}&UnidadMenor=${userMatch.UnidadMenor}`)
        }
        } 
      } 
       else {
        this.errorMessage = "El rol seleccionado no coincide con los usuarios disponibles.";
      }
    } catch (error) {
      this.errorMessage = "Hubo un error al iniciar sesión.";
      console.error(error);
    }
  }*/
   else {
    this.errorMessage = "Por favor, completa todos los campos.";
  }
},
async confirmRole() {
  if (this.roleElegido === "personal-dei") {
            this.$router.push("paf-en-linea/personas");
  } 
  else if (this.roleElegido === "encargado") {
    const response1 = await this.$axios.get(`/api/paf-en-linea/usuario/rut/${sessionStorage.getItem("rut")}`);
            if (response1.data[0].Acceso === 0) {
              this.errorMessage = "Este usuario no posee accesos.";
              return;
            } else if (response1.data[0].Vista_universidad === 1) {
              this.$router.push("paf-en-linea/seguimientoPAF");
            } else if (response1.data[0].Vista_facultad === 1) {
              sessionStorage.setItem("unidadMayor", response1.data[0].UnidadMayor);
              this.$router.push(`paf-en-linea/unidadMayorPAF?UnidadMayor=${response1.data[0].UnidadMayor}`);
            } else if (response1.data[0].Acceso === 1 && response1.data[0].Vista_facultad === 0 && response1.data[0].Vista_universidad === 0) {
              sessionStorage.setItem("unidadMayor", response1.data[0].UnidadMayor);
              sessionStorage.setItem("unidadMenor", response1.data[0].UnidadMenor);
              this.$router.push(`paf-en-linea/unidadMayorPAF?UnidadMayor=${response1.data[0].UnidadMayor}&UnidadMenor=${response1.data[0].UnidadMenor}`);
      }
    }
    else if (this.roleElegido === "profesor") {
      this.$router.push(`paf-en-linea/profesorPAF?run=${this.newRut}`);
    } 
      this.closeRoleModal();
    },
    closeRoleModal() {
      this.showRoleModal = false;
      this.roleElegido = "";
      this.newRut = "";
    }
    },
  };
</script>

<style scoped>
/* Página principal */
.login-page {
  max-width: 500px;
  margin: 20px auto;
  border: 1px solid #394049;
  border-radius: 8px;
  background-color: #f9f9f9;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  text-align: center;
}

/* Barra superior */
.header-bar {
  background-color: #00a499;
  color: #fff;
  padding: 10px;
  font-family: "Bebas Neue Pro", sans-serif;
  font-size: 1.5rem;
  font-weight: bold;
}

/* Contenido principal */
.content {
  padding: 20px;
}

/* Formulario */
.login-form {
  margin-top: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  font-family: "Helvetica Neue LT", sans-serif;
  font-size: 0.9rem;
  text-align: left;
  margin-bottom: 5px;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 0.9rem;
}

button[type="submit"] {
  background-color: #EA7600;
  border: none;
  color: white;
  font-family: "Bebas Neue Pro", sans-serif;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button[type="submit"]:hover {
  background-color: #C8102E;
}

/* Modal */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.modal-content .form-group {
  margin-bottom: 15px;
}

button {
  background-color: #EA7600;
  border: none;
  color: white;
  font-family: "Bebas Neue Pro", sans-serif;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #C8102E;
}

/* Errores */
.error {
  color: red;
  margin-top: 10px;
  font-size: 0.9rem;
}
</style>