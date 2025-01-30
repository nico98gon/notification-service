# Weather Notification System

Este proyecto implementa un sistema de notificaciones climáticas utilizando una arquitectura de microservicios. Cada microservicio está diseñado para cumplir una función específica, siguiendo las mejores prácticas de escalabilidad, resiliencia y mantenimiento.

## Descripción General

El sistema consta de tres microservicios principales:

1. **Gestión de Usuarios**: CRUD para usuarios y sus preferencias de notificación.
2. **Gestión de Localidades**: Integración con la API del CPTEC para obtener información climática y de olas.
3. **Gestión de Notificaciones**: Envío de notificaciones basadas en eventos, con soporte para programación y múltiples canales de entrega.

Los microservicios están contenedorizados usando Docker, y se pueden ejecutar tanto de manera individual como con Docker Compose para facilitar el desarrollo y despliegue.

---

## Requisitos Funcionales y No Funcionales

### Funcionales

- Obtener información climática y de olas desde la API del CPTEC.
- Enviar notificaciones con previsiones para los próximos 4 días, incluyendo olas para localidades costeras.
- Permitir programar envíos y respetar el opt-out de usuarios.

### No Funcionales

- Alta escalabilidad y bajo nivel de latencia.
- Resiliencia ante fallos.
- Capacidad para añadir futuros canales de notificación como SMS y Push.
- Implementación siguiendo principios de arquitectura de microservicios.

---

## Arquitectura

### Microservicios y Patrones Utilizados

| **Microservicio**        | **Arquitectura**         | **Descripción**                                                                        |
|--------------------------|--------------------------|----------------------------------------------------------------------------------------|
| Gestión de Usuarios      | MVC                      | CRUD básico de usuarios y preferencias.                                                |
| Gestión de Localidades   | Domain-Driven Design     | Procesamiento complejo de datos climáticos y de olas desde la API externa.             |
| Gestión de Notificaciones| Event-Driven Architecture| Basado en eventos para programar y enviar notificaciones en tiempo real o programadas. |

### Archivo de documentación

En el repositorio del microservicio user-service existe un archivo con la documentación previa al desarrollo del challenge en formato .pdf

---

## Configuración y Ejecución

### Prerrequisitos

- **Docker** y **Docker Compose** instalados.
- [Golang Air](https://github.com/cosmtrek/air) para hot reload durante el desarrollo.
- Acceso a la API del CPTEC (credenciales necesarias).

### Estructura del Proyecto

nilus-challenge-backend/
    ├── docker-compose.yml
    ├── notification-service/
    ├── user-service/
    └── weather-service/

### Prueba de endpoints en postman

Se ha dejado el archivo .json disponible en el repositorio de cada microservicio.
En caso de no poder importar se puede ingresar al Workspace de Postman en el siguiente link: 
https://app.getpostman.com/join-team?invite_code=93887e071673589a5d26b909ee3d531511635ea794315f6f006d89a8b79b2351&target_code=fc47de7ca7fd71a7b02a4576bfa810e

### Docker Compose

Aquí está el código de docker compose que debe encontrarse a un nivel superior de los 3 microservicios.

services:
  postgres:
    image: postgres:15
    container_name: postgres
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: postgres101
      POSTGRES_DB: nilus_challenge_db
    ports:
      - "5432:5432"
    networks:
      - app-network

  user-service:
    build:
      context: ./user-service
      dockerfile: Dockerfile
    container_name: user-service
    restart: unless-stopped
    env_file:
      - ./user-service/.env
    ports:
      - "8081:8081"
    volumes:
      - ./user-service:/app
      - ~/.cache/go-build:/go
    command: ["air", "-c", ".air.toml"]
    networks:
      - app-network
    depends_on:
      - postgres
      - notification-service
    # environment:
    #   DATABASE_URL: postgres://user:password@postgres:5432/mydb

  notification-service:
    build:
      context: ./notification-service
      dockerfile: Dockerfile
    container_name: notification-service
    restart: unless-stopped
    env_file:
      - ./notification-service/.env
    ports:
      - "8082:8082"
    volumes:
      - ./notification-service:/app
      - ~/.cache/go-build:/go
    command: ["air", "-c", ".air.toml"]
    networks:
      - app-network
    depends_on:
      - weather-service
    # environment:
    #   USER_SERVICE_URL: http://user-service:8081
    #   WEATHER_SERVICE_URL: http://weather-service:8083

  weather-service:
    build:
      context: ./weather-service
      dockerfile: Dockerfile
    container_name: weather-service
    ports:
      - "8083:8083"
    volumes:
      - ./weather-service:/app
      - ~/.cache/go-build:/go
    networks:
      - app-network
    command: ["air", "-c", ".air.toml"]
    # environment:
    #   DATABASE_URL: postgres://user:password@postgres:5432/mydb

networks:
  app-network:
    driver: bridge
