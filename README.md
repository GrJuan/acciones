# 📌 Proyecto de Gestión de Acciones con Vue 3 y AWS Terraform

Este proyecto permite gestionar acciones bursátiles de manera visual e interactiva utilizando Vue 3 y AWS. Se incluyen gráficos dinámicos, manejo de autenticación y despliegue en AWS con Terraform.

---

## 🚀 Tecnologías utilizadas

- **Frontend:** Vue 3, Tailwind CSS, Heroicons
- **Backend:** API Gateway AWS, Lambda Functions (Node.js)
- **Base de Datos:** cockroach
- **Infraestructura:** Terraform, AWS EC2
- **Autenticación:** API de autenticación personalizada

---

## 🔧 Instalación y Configuración

### 📌 1. Clonar el repositorio
```bash
git clone https://github.com/tu-repo.git
cd frontend
```

### 📌 2. Instalar dependencias
```bash
npm install
```

### 📌 4. Ejecutar el proyecto localmente
```bash
npm run dev
```

---

## 📦 Despliegue en AWS con Terraform

Terraform se utiliza para desplegar la infraestructura necesaria en AWS, incluyendo una instancia EC2 para el frontend.

### 📌 1. Configurar Terraform
Asegúrate de tener instalado Terraform y configura AWS CLI:
```bash
aws configure
```

### 📌 2. Inicializar Terraform
```bash
cd terraform
terraform init
```

### 📌 3. Aplicar cambios en AWS
```bash
terraform apply -auto-approve
```
Esto creará la infraestructura y mostrará la IP pública de la instancia.

### 📌 4. Acceder a la aplicación
Usa la IP pública de la instancia EC2 en tu navegador:
```bash
http://<IP-PUBLICA>
```

---

## 🔥 Funcionalidades Clave

✅ **Autenticación con login**
✅ **Gestión de Acciones Disponibles y Portafolio**
✅ **Gráficos Interactivos con Chart.js**
✅ **Lazy Loading para mejorar rendimiento**
✅ **Kanban y Línea de Tiempo**
✅ **Recomendaciones de compra con IA**
✅ **Despliegue automático con Terraform**

---

