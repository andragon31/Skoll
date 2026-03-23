# PRD — Skoll
## Product Requirements Document v1.0

**Producto:** Skoll  
**Tipo:** MCP Plugin — RSAW Orchestration Layer  
**Lenguaje:** Go 1.22+  
**Licencia:** MIT  
**Versión del documento:** 1.0  
**Fecha:** Marzo 2026  
**Estado:** Draft

---

## Tabla de contenidos

1. [Visión del producto](#1-visión-del-producto)
2. [Problema que resuelve](#2-problema-que-resuelve)
3. [Usuarios objetivo](#3-usuarios-objetivo)
4. [Objetivos y no-objetivos](#4-objetivos-y-no-objetivos)
5. [El sistema RSAW](#5-el-sistema-rsaw)
6. [Requisitos funcionales](#6-requisitos-funcionales)
7. [Requisitos no funcionales](#7-requisitos-no-funcionales)
8. [Arquitectura técnica](#8-arquitectura-técnica)
9. [Modelo de datos y estructura de archivos](#9-modelo-de-datos-y-estructura-de-archivos)
10. [MCP Tools — Especificación completa](#10-mcp-tools--especificación-completa)
11. [CLI — Comandos](#11-cli--comandos)
12. [Configuración](#12-configuración)
13. [Integración con Fenrir](#13-integración-con-fenrir)
14. [Adapters por herramienta](#14-adapters-por-herramienta)
15. [Distribución](#15-distribución)
16. [Métricas de éxito](#16-métricas-de-éxito)
17. [Roadmap y milestones](#17-roadmap-y-milestones)
18. [Riesgos y mitigaciones](#18-riesgos-y-mitigaciones)

---

## 1. Visión del producto

Skoll es la capa de **orquestación estructurada** para equipos que desarrollan software con agentes de IA.

Define quién hace qué, cómo se hace y en qué orden — sin importar qué herramienta AI esté activa. Convierte el caos del desarrollo con IA en un proceso reproducible, con roles claros, conocimiento reutilizable y flujos de trabajo explícitos.

> **Misión:** Que ningún agente de IA trabaje sin saber quién es, qué sabe y cuál es el proceso que debe seguir.

### Propuesta de valor única

Un equipo humano bien organizado no necesita que cada persona reinvente cómo hacer un commit, cómo diseñar un endpoint o cómo estructurar un módulo. Ese conocimiento está en la cultura del equipo. Skoll da esa cultura a cada agente, desde el primero hasta el último.

### Relación con Fenrir

Skoll y Fenrir son productos independientes. Funcionan solos. Pero cuando coexisten, se complementan:

- **Skoll** responde a: *"¿Quién hace qué, cómo y en qué orden?"*
- **Fenrir** responde a: *"¿Qué pasó, qué sabemos y qué está permitido?"*

Skoll puede llamar a Fenrir en puntos específicos del ciclo de vida (inicio de workflow, handoff entre agentes, cierre de workflow) para enriquecer la ejecución con contexto histórico y gobernanza. Esta integración es opcional — no es una dependencia.

---

## 2. Problema que resuelve

Los agentes de IA son potentes pero desestructurados. Sin un framework de roles y procesos, generan:

### Problema A — Contexto indefinido
El agente no sabe cuál es su rol en el proyecto. "¿Soy backend? ¿Frontend? ¿Me encargo de tests también?" El resultado es código que cruza responsabilidades, archivos modificados fuera del scope y confusión en los handoffs.

### Problema B — Conocimiento no reutilizable
Cada vez que el agente necesita hacer un commit, un test o diseñar un endpoint, improvisa. El conocimiento de "cómo se hace aquí" no existe de forma estructurada. El mismo patrón se reinventa en cada sesión, con variaciones.

### Problema C — Procesos implícitos
"Desarrollar un feature" no tiene pasos definidos. El agente interpreta libremente: a veces hace tests, a veces no; a veces consulta antes de implementar, a veces no. La reproducibilidad del proceso es cero.

### Problema D — Handoffs rotos
Cuando el backend agent termina y el frontend agent necesita consumir la API, no hay un protocolo de transferencia. El frontend agent adivina el contrato. Los bugs de integración son consecuencia directa de esto.

### Problema E — Reglas inconsistentes
Cada developer tiene su propio CLAUDE.md o .cursorrules con reglas distintas. El mismo proyecto tiene restricciones contradictorias según qué developer configuró el agente. No hay una fuente de verdad única.

---

## 3. Usuarios objetivo

### Persona primaria — El Developer con contexto de equipo

- Trabaja en equipo de 2–15 personas usando IA
- Quiere que el agente respete los roles y estándares del equipo
- Ha visto cómo el agente modifica archivos que no le corresponden o genera código inconsistente con el resto del proyecto
- **Pain point:** "El agente no sabe dónde termina su responsabilidad"

### Persona secundaria — El Tech Lead / Arquitecto

- Define la arquitectura y los estándares del equipo
- Quiere que esos estándares se apliquen automáticamente en cada sesión de IA, sin necesidad de repetirlos
- **Pain point:** "Cada vez que alguien usa el agente, tengo que recordarle las convenciones del proyecto"

### Persona terciaria — El Developer que trabaja solo en proyectos grandes

- Proyecto con múltiples dominios (frontend, backend, infra, QA)
- Usa diferentes "sombreros" según la tarea
- Quiere que el agente cambie de modo según el contexto
- **Pain point:** "El agente mezcla responsabilidades y genera código inconsistente entre sesiones"

---

## 4. Objetivos y no-objetivos

### Objetivos

- Proveer un sistema de cuatro capas (Rules, Skills, Agents, Workflows) como archivos Markdown estructurados en el proyecto
- Permitir activar un agente específico con su scope, skills y comportamiento definidos
- Ejecutar workflows con pasos explícitos, agentes asignados y Definition of Done
- Gestionar handoffs entre agentes con contratos explícitos
- Ser completamente independiente de herramientas AI específicas
- Funcionar con o sin Fenrir instalado
- Instalarse y configurarse en menos de 60 segundos

### No-objetivos

- **No** es un sistema de memoria — eso es responsabilidad de Fenrir
- **No** genera código ni hace sugerencias de implementación
- **No** ejecuta código directamente — orquesta agentes que lo hacen
- **No** reemplaza el CLAUDE.md del developer — lo genera y mantiene
- **No** requiere servidor — es stateless por diseño (el estado vive en archivos Markdown)
- **No** depende de Fenrir para funcionar
- **No** gestiona permisos de sistema operativo ni acceso a archivos

---

## 5. El sistema RSAW

Skoll implementa el framework **RSAW** (Rules, Skills, Agents, Workflows) como una jerarquía de cuatro capas.

### La regla de oro para clasificar

```
¿Esto aplica SIEMPRE sin importar nada?          →  Rules
¿Esto describe CÓMO hacer algo técnico?           →  Skill
¿Esto define QUIÉN hace algo y en qué contexto?  →  Agent
¿Esto define el ORDEN en que se hacen las cosas? →  Workflow
```

---

### Capa 1 — Rules (La base inamovible)

Son restricciones no negociables que aplican en todo momento, sin importar qué agente esté activo o qué workflow se esté ejecutando.

**Características:**
- Siempre activas, sin excepción
- No describen cómo hacer algo, solo qué está permitido/prohibido
- Son las más cortas y específicas
- No tienen pasos ni procesos
- No saben de tecnologías específicas

**Formato de archivo:**

```markdown
# Rules: [Categoría]

## Regla 1 — [Nombre corto]
[Una línea que describe la restricción]

## Regla 2 — [Nombre corto]
[Una línea que describe la restricción]
```

**Ejemplo — `global.md`:**
```markdown
# Rules: Global

## No secrets hardcodeados
Nunca escribir API keys, passwords o tokens directamente en código.
Usar siempre variables de entorno o secret managers.

## No any en TypeScript
Nunca usar el tipo `any` en TypeScript. Usar `unknown` con type guards o tipos específicos.

## Confirmar antes de destruir
Siempre confirmar con el usuario antes de eliminar archivos, bases de datos o datos de producción.

## Idioma del usuario
Responder siempre en el idioma del usuario, independientemente del idioma del código.

## Scope del proyecto
Nunca modificar archivos fuera del directorio raíz del proyecto sin confirmación explícita.
```

---

### Capa 2 — Skills (Conocimiento aplicado)

Son bloques de conocimiento reutilizable sobre cómo hacer algo específico. Un skill no sabe quién lo va a usar — solo describe el "cómo". Son autocontenidos, técnicos y tienen ejemplos concretos.

**Características:**
- Autocontenidos — no dependen de ningún agente
- Reutilizables — múltiples agentes aplican el mismo skill
- Técnicos y específicos — tienen ejemplos, comandos, patrones
- Definen cuándo aplicarse (trigger)
- No saben del proyecto en particular

**Formato de archivo:**

```markdown
# Skill: [Nombre]

## Cuándo aplicar
[Trigger: en qué situación debe usarse este skill]

## Contexto
[Por qué importa y qué problema resuelve]

## Proceso
[Pasos concretos con ejemplos de código si aplica]

## Checklist
- [ ] [Verificación 1]
- [ ] [Verificación 2]

## Anti-patrones
[Lo que NO se debe hacer y por qué]
```

**Ejemplo — `git-workflow.md`:**
```markdown
# Skill: Git Workflow

## Cuándo aplicar
Al crear commits, branches o pull requests.

## Proceso

### Branch naming
feat/[descripción-corta]
fix/[descripción-corta]
chore/[descripción-corta]

### Commit message
Formato: `tipo(scope): descripción en imperativo`
Ejemplos:
  feat(auth): add JWT refresh token rotation
  fix(payments): handle stripe webhook timeout
  chore(deps): upgrade axios to 1.7.0

### Antes de commitear
- Revisar el diff completo
- Verificar que no hay secrets en el diff
- Asegurar que los tests pasan

## Anti-patrones
- No usar mensajes como "fix", "changes", "wip"
- No commitear archivos .env o secrets
- No hacer commits de más de 400 líneas cambiadas sin justificación
```

---

### Capa 3 — Agents (El rol con personalidad)

Un agent es un rol con identidad, scope y skills asignados. Es quien recibe la tarea, sabe qué herramientas (skills) usar, y sabe cuándo derivar a otro agente.

**Características:**
- Tiene un rol claro con nombre
- Define su scope: qué archivos/tareas le pertenecen
- Lista los skills que aplica
- Sabe cuándo hacer handoff a otro agente
- Conoce el stack tecnológico específico del proyecto
- No define procesos paso a paso (eso es el workflow)
- No repite reglas (ya están en rules/)

**Formato de archivo:**

```markdown
# Agent: [Nombre del rol]

## Rol
[Quién eres en una o dos líneas — identidad y especialización]

## Scope
### Archivos que me pertenecen
- [path/pattern]

### Lo que NO toco
- [path/pattern]

## Stack
[Tecnologías específicas que uso en este proyecto]

## Skills que aplico
- [skill-nombre.md]

## Cuándo hacer handoff
### → [Otro agente]
[Cuándo: condición que dispara el handoff]
[Contrato: qué se transfiere exactamente]
```

> ⚠️ **Regla crítica del Agent:** No incluir pasos de proceso ni instrucciones de "cómo hacer" las cosas. Si necesitas describir un proceso, es un **Workflow**. Si necesitas describir cómo se hace algo técnico, es un **Skill**. El agent solo responde: *"¿Quién eres, qué te pertenece y qué sabes?"*

**Ejemplo — `backend.md`:**
```markdown
# Agent: Backend Engineer

## Rol
Soy un Senior Backend Engineer especializado en NestJS y Prisma.
Me especializo en la lógica de negocio, el dominio y la capa de persistencia.

## Scope
### Archivos que me pertenecen
- src/modules/**
- src/services/**
- src/domain/**
- prisma/**
- test/unit/**

### Lo que NO toco
- src/components/**  (Frontend agent)
- src/pages/**       (Frontend agent)
- .github/workflows/ (DevOps agent)
- infrastructure/**  (DevOps agent)

## Stack
- NestJS 10 + TypeScript strict
- Prisma ORM + PostgreSQL
- Jest para unit tests
- Zod para validación

## Skills que aplico
- clean-architecture.md
- error-handling.md
- testing.md
- api-design.md
- git-workflow.md

## Cuándo hacer handoff
### → Frontend Agent
Cuando: al definir o modificar cualquier endpoint de API
Contrato: DTO de request, DTO de response, códigos de error posibles

### → QA Agent
Cuando: al completar la implementación de un módulo
Contrato: lista de casos de prueba esperados, happy path y edge cases

### → DevOps Agent
Cuando: se necesitan cambios en infraestructura o variables de entorno
Contrato: qué variables nuevas se requieren y en qué entornos
```

---

### Capa 4 — Workflows (El director de orquesta)

Un workflow coordina múltiples agentes y skills en un proceso con pasos definidos. Es el único componente que puede "llamar" a otros agentes.

**Características:**
- Tiene pasos ordenados y explícitos
- Indica qué agent ejecuta cada paso
- Define el Definition of Done
- Cubre los casos de error del proceso
- No repite lógica técnica (está en agents/skills)
- No tiene detalles de implementación

**Formato de archivo:**

```markdown
# Workflow: [Nombre]

## Propósito
[Qué proceso orquesta este workflow]

## Cuándo usarlo
[Trigger: en qué situaciones activar este workflow]

## Prerequisitos
- [ ] [Condición necesaria]

## Pasos

### Paso N — [Nombre] ([Agent responsable] | ninguno)
[Qué hace este paso en una línea — sin detalles de implementación]
Skills: [skill-a.md, skill-b.md] ← referencias, no repetición
**Output esperado:** [Qué artefacto o resultado concreto debe existir al terminar]

## Definition of Done
- [ ] [Criterio verificable de completitud]

## Casos de error
### Si [condición de error]
[Qué hacer — sin implementación, solo decisión de flujo]
```

> ⚠️ **Regla crítica del Workflow:** Los pasos describen *qué* hace cada agente y *qué skill aplica*, nunca *cómo* lo hace. El *cómo* vive en el skill referenciado. Si un paso tiene más de 2 líneas de descripción, hay detalles de implementación que deben moverse al skill correspondiente.

**Ejemplo — `feature.md`:**
```markdown
# Workflow: Feature Development

## Propósito
Proceso completo para desarrollar un feature desde el requerimiento hasta el commit.

## Cuándo usarlo
Cuando el usuario pide desarrollar cualquier funcionalidad nueva no trivial.

## Prerequisitos
- [ ] El requerimiento está claro o se han hecho las preguntas necesarias
- [ ] El scope del feature está definido

## Pasos

### Paso 1 — Entender (ninguno)
Leer el requerimiento. Preguntar si hay ambigüedad. No escribir código hasta tener claridad.
**Output esperado:** Confirmación de entendimiento del usuario.

### Paso 2 — Domain (Backend Agent)
Modelar el dominio del feature.
Skills: clean-architecture.md
**Output esperado:** Entidades y tipos del dominio definidos.

### Paso 3 — Application (Backend Agent)
Implementar la lógica de aplicación.
Skills: clean-architecture.md, error-handling.md, testing.md
**Output esperado:** Handler implementado con tests unitarios pasando.

### Paso 4 — API (Backend Agent → handoff Frontend Agent)
Exponer el dominio como API y transferir contrato al frontend.
Skills: api-design.md
**Output esperado:** Endpoint funcional. Contrato documentado y entregado al Frontend Agent.

### Paso 5 — UI (Frontend Agent)
Consumir el contrato recibido del Backend Agent.
Skills: [skill de componentes del proyecto]
**Output esperado:** Componente conectado al endpoint con manejo de todos los estados.

### Paso 6 — Tests (QA Agent)
Validar el flujo completo.
Skills: testing.md
**Output esperado:** Tests E2E del flujo completo pasando en CI.

### Paso 7 — Commit (todos)
Commitear los cambios de cada capa.
Skills: git-workflow.md
**Output esperado:** Commits con mensajes convencionales por cada agente.

## Definition of Done
- [ ] Todos los tests pasan en CI
- [ ] No hay warnings de TypeScript
- [ ] El endpoint maneja happy path y casos de error
- [ ] El UI maneja loading, error, empty y success
- [ ] Los commits siguen el formato convencional

## Casos de error
### Si el requerimiento es ambiguo
Volver al Paso 1. No continuar sin claridad total.

### Si el Backend Agent no puede completar el dominio
Escalar al usuario antes de avanzar al Paso 3.
```

---

## 6. Requisitos funcionales

### RF-01 — Gestión de Rules

| ID | Requisito | Prioridad |
|---|---|---|
| RF-01-01 | El sistema debe cargar y parsear todos los archivos en `.skoll/rules/` al iniciar | MUST |
| RF-01-02 | Las rules deben estar disponibles en todo momento, sin necesidad de cargarlas explícitamente | MUST |
| RF-01-03 | El tool `rule_check` debe verificar si una acción viola alguna rule activa | MUST |
| RF-01-04 | Las rules deben poderse listar vía CLI y TUI | MUST |
| RF-01-05 | `skoll init --rules` debe generar un set de rules base sensatas para cualquier proyecto | SHOULD |
| RF-01-06 | Las rules deben poder marcarse como `active` o `disabled` sin eliminarlas | SHOULD |

### RF-02 — Gestión de Skills

| ID | Requisito | Prioridad |
|---|---|---|
| RF-02-01 | El sistema debe indexar todos los archivos en `.skoll/skills/` | MUST |
| RF-02-02 | `skill_load` debe retornar el contenido completo de un skill por nombre | MUST |
| RF-02-03 | `skill_search` debe buscar skills por nombre o descripción | MUST |
| RF-02-04 | Los skills deben tener un campo `trigger` que permita encontrarlos por contexto | MUST |
| RF-02-05 | `skoll skill add <nombre>` debe crear un archivo de skill con la plantilla correcta | MUST |
| RF-02-06 | `skoll skill import <fuente>` debe importar skills desde Gentleman-Skills u otras fuentes | SHOULD |
| RF-02-07 | Los skills deben poder validarse para verificar que tienen los campos obligatorios | SHOULD |

### RF-03 — Gestión de Agents

| ID | Requisito | Prioridad |
|---|---|---|
| RF-03-01 | El sistema debe indexar todos los archivos en `.skoll/agents/` | MUST |
| RF-03-02 | `agent_activate` debe activar un agente inyectando su identidad, scope y skills en el contexto | MUST |
| RF-03-03 | `agent_list` debe retornar todos los agentes disponibles en el proyecto | MUST |
| RF-03-04 | `agent_context` debe retornar el contexto completo de un agente (identidad + scope + skills cargados) | MUST |
| RF-03-05 | `agent_handoff` debe formalizar la transferencia entre agentes con un contrato explícito | MUST |
| RF-03-06 | El sistema debe validar que el scope de un agente no se solapa con el de otro | SHOULD |
| RF-03-07 | `skoll agent add <nombre>` debe crear un archivo de agente con la plantilla correcta | MUST |

### RF-04 — Gestión de Workflows

| ID | Requisito | Prioridad |
|---|---|---|
| RF-04-01 | El sistema debe indexar todos los archivos en `.skoll/workflows/` | MUST |
| RF-04-02 | `workflow_start` debe iniciar un workflow, cargando el primer paso y el agente asignado | MUST |
| RF-04-03 | `workflow_step` debe avanzar al siguiente paso retornando: descripción, agente, output esperado | MUST |
| RF-04-04 | `workflow_status` debe retornar el estado actual del workflow activo | MUST |
| RF-04-05 | `workflow_complete` debe cerrar el workflow verificando el Definition of Done | MUST |
| RF-04-06 | Solo puede haber un workflow activo por sesión | MUST |
| RF-04-07 | El estado del workflow activo debe persistir entre mensajes del agente | MUST |
| RF-04-08 | `skoll workflow add <nombre>` debe crear un archivo de workflow con la plantilla correcta | MUST |

### RF-05 — Init y Setup

| ID | Requisito | Prioridad |
|---|---|---|
| RF-05-01 | `skoll init` debe crear la estructura `.skoll/` completa con ejemplos funcionales | MUST |
| RF-05-02 | `skoll init` debe detectar el stack del proyecto y generar skills relevantes | SHOULD |
| RF-05-03 | `skoll init` debe generar los archivos de configuración para herramientas AI detectadas | MUST |
| RF-05-04 | `skoll init --rules-only` debe generar solo rules base | SHOULD |
| RF-05-05 | `skoll init --from <url>` debe inicializar desde un template de repositorio | COULD |

### RF-06 — Validación

| ID | Requisito | Prioridad |
|---|---|---|
| RF-06-01 | `skoll validate` debe verificar que todos los archivos RSAW tienen los campos obligatorios | MUST |
| RF-06-02 | La validación debe detectar referencias a skills inexistentes en los agents | MUST |
| RF-06-03 | La validación debe detectar referencias a agentes inexistentes en los workflows | MUST |
| RF-06-04 | La validación debe detectar pasos de workflow con más de 3 líneas de descripción (señal de detalles de implementación) | SHOULD |
| RF-06-05 | La validación debe detectar secciones de proceso paso a paso dentro de archivos Agent (violación de RSAW) | SHOULD |
| RF-06-06 | La validación debe correr automáticamente en CI vía `skoll validate --ci` | SHOULD |

### RF-07 — TUI

| ID | Requisito | Prioridad |
|---|---|---|
| RF-07-01 | Dashboard: resumen de rules activas, agents disponibles, skills indexados, workflows | MUST |
| RF-07-02 | Browser de Rules: listar y ver contenido de cada rule | MUST |
| RF-07-03 | Browser de Skills: buscar y ver skills con su trigger | MUST |
| RF-07-04 | Browser de Agents: ver agentes con scope y skills | MUST |
| RF-07-05 | Workflow Runner: iniciar y avanzar un workflow interactivamente | SHOULD |
| RF-07-06 | Navegación vim-style | MUST |

---

## 7. Requisitos no funcionales

### RNF-01 — Performance

| ID | Requisito | Métrica |
|---|---|---|
| RNF-01-01 | Tiempo de startup del servidor MCP | < 100ms |
| RNF-01-02 | Tiempo de `skoll init` | < 30 segundos |
| RNF-01-03 | Tiempo de indexado de .skoll/ completo | < 500ms para hasta 100 archivos |
| RNF-01-04 | Tiempo de respuesta de cualquier tool MCP | < 100ms |
| RNF-01-05 | Tamaño del binario | < 15MB |

### RNF-02 — Confiabilidad

| ID | Requisito |
|---|---|
| RNF-02-01 | Skoll es stateless — no tiene base de datos propia, el estado vive en archivos Markdown |
| RNF-02-02 | Si `.skoll/` no existe, Skoll degrada gracefully y sugiere correr `skoll init` |
| RNF-02-03 | Archivos RSAW con errores de formato no deben crashear el servidor — deben loguearse y omitirse |

### RNF-03 — Compatibilidad

| ID | Requisito |
|---|---|
| RNF-03-01 | Compatible con MCP protocol spec 1.0+ |
| RNF-03-02 | Funciona en Linux, macOS (Intel y Apple Silicon), Windows |
| RNF-03-03 | Compatible con Claude Code, Cursor, Windsurf, GitHub Copilot, Gemini CLI, OpenCode |
| RNF-03-04 | Funciona sin Fenrir instalado |
| RNF-03-05 | Funciona con Fenrir instalado, enriqueciéndose con su contexto |

### RNF-04 — Extensibilidad

| ID | Requisito |
|---|---|
| RNF-04-01 | Cualquier developer puede agregar skills, agents y workflows propios sin tocar código |
| RNF-04-02 | El formato RSAW debe ser estable — cambios breaking requieren major version bump |
| RNF-04-03 | Skills pueden importarse desde repositorios externos (Gentleman-Skills y otros) |

---

## 8. Arquitectura técnica

### Decisión de diseño: stateless

Skoll no tiene base de datos. Su "estado" son los archivos Markdown en `.skoll/`. Esta decisión es intencional:

- Los archivos RSAW son versionables con git — forman parte del proyecto
- El equipo puede revisar, modificar y hacer PR a cualquier rule, skill, agent o workflow
- No hay migración de datos, no hay corrupción de DB
- El único estado transitorio (workflow activo en una sesión) vive en memoria del proceso MCP

### Stack

| Componente | Librería | Razón |
|---|---|---|
| Lenguaje | Go 1.22+ | Binario único, sin runtime |
| MCP | github.com/mark3labs/mcp-go | Consistencia con Fenrir |
| Markdown parsing | github.com/yuin/goldmark | Parser Markdown más completo en Go |
| CLI | github.com/spf13/cobra | Estándar de facto |
| Config | github.com/spf13/viper | JSON/YAML/ENV |
| TUI | github.com/charmbracelet/bubbletea + lipgloss | Consistencia con Fenrir |
| File watching | github.com/fsnotify/fsnotify | Reindexar .skoll/ en cambios |
| Testing | testing + testify | Tests unitarios |
| Release | GoReleaser | Distribución |

### Estructura de directorios del proyecto

```
skoll/
├── cmd/
│   └── skoll/
│       └── main.go
├── internal/
│   ├── rsaw/
│   │   ├── loader.go        # Carga y parsea archivos Markdown de .skoll/
│   │   ├── index.go         # Índice en memoria de Rules, Skills, Agents, Workflows
│   │   ├── validator.go     # Validación de estructura y referencias
│   │   ├── watcher.go       # File watching para reindexar en cambios
│   │   └── types.go         # Rule, Skill, Agent, Workflow, Step, Handoff...
│   ├── mcp/
│   │   ├── server.go        # MCP server
│   │   └── tools.go         # Registro de los 16 tools
│   ├── engine/
│   │   ├── workflow.go      # Estado de workflow activo en sesión
│   │   ├── handoff.go       # Protocolo de handoff entre agentes
│   │   └── fenrir.go     # Cliente opcional de Fenrir MCP
│   ├── generator/
│   │   ├── init.go          # skoll init — genera estructura .skoll/
│   │   ├── detect.go        # Detección de stack del proyecto
│   │   ├── templates.go     # Templates de Rules, Skills, Agents, Workflows
│   │   └── adapters.go      # Genera configs para herramientas AI
│   ├── adapters/
│   │   ├── claude_code.go
│   │   ├── cursor.go
│   │   ├── windsurf.go
│   │   ├── copilot.go
│   │   ├── gemini.go
│   │   └── opencode.go
│   └── tui/
│       ├── model.go
│       ├── styles.go
│       ├── update.go
│       └── view.go
├── .goreleaser.yaml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

### Estructura `.skoll/` en el proyecto del usuario

```
.skoll/
├── rules/
│   ├── global.md          # Reglas universales del proyecto
│   ├── security.md        # Reglas de seguridad
│   └── quality.md         # Reglas de calidad de código
├── skills/
│   ├── git-workflow.md
│   ├── clean-architecture.md
│   ├── testing.md
│   ├── api-design.md
│   ├── error-handling.md
│   └── [framework-específico].md
├── agents/
│   ├── backend.md
│   ├── frontend.md
│   ├── qa.md
│   └── devops.md
├── workflows/
│   ├── feature.md
│   ├── bugfix.md
│   ├── refactor.md
│   └── release.md
└── skoll.json          # Configuración de Skoll para el proyecto
```

---

## 9. Modelo de datos y estructura de archivos

### Tipos internos (en memoria, no persistidos)

```go
type Rule struct {
    ID          string
    Name        string
    Category    string
    Description string
    Status      string    // active | disabled
    FilePath    string
    LoadedAt    time.Time
}

type Skill struct {
    ID          string
    Name        string
    Trigger     string    // Cuándo aplicar este skill
    Context     string    // Por qué importa
    Content     string    // Contenido completo del skill
    Tags        []string
    FilePath    string
    LoadedAt    time.Time
}

type Agent struct {
    ID           string
    Name         string
    Role         string        // Identidad en una o dos líneas
    Scope        AgentScope
    Stack        []string
    Skills       []string      // Referencias a skill IDs
    Handoffs     []HandoffDef
    FilePath     string
    LoadedAt     time.Time
}

type AgentScope struct {
    Owns    []string    // Paths/patterns que le pertenecen
    Avoids  []string    // Paths/patterns que no toca
}

type HandoffDef struct {
    ToAgent  string
    When     string
    Contract string
}

type Workflow struct {
    ID           string
    Name         string
    Purpose      string
    Trigger      string
    Prerequisites []string
    Steps        []WorkflowStep
    DoD          []string    // Definition of Done
    ErrorCases   []ErrorCase
    FilePath     string
    LoadedAt     time.Time
}

type WorkflowStep struct {
    Number       int
    Name         string
    Agent        string    // Agent ID o "none"
    Description  string
    Output       string
    IsHandoff    bool
    HandoffTo    string
}

type ActiveWorkflow struct {
    WorkflowID   string
    CurrentStep  int
    StartedAt    time.Time
    SessionID    string
    Context      map[string]string
}
```

### `skoll.json` (configuración del proyecto)

```json
{
  "project": "mi-proyecto",
  "version": "1.0",
  "wolf_dir": ".wolf",
  "default_agent": "backend",
  "active_workflows_limit": 1,
  "validation": {
    "strict": false,
    "warn_on_scope_overlap": true
  },
  "fenrir": {
    "enabled": false,
    "integration_points": {
      "workflow_start": true,
      "agent_activate": false,
      "workflow_complete": true,
      "handoff": false
    }
  },
  "skill_sources": [
    {
      "name": "gentleman-skills",
      "url": "https://github.com/Gentleman-Programming/Gentleman-Skills",
      "auto_update": false
    }
  ]
}
```

---

## 10. MCP Tools — Especificación completa

Skoll expone **16 herramientas MCP** organizadas en 4 módulos RSAW más un módulo de sistema.

---

### Módulo Rules (3 tools)

#### `rule_list`
Lista todas las rules activas del proyecto.

```json
{
  "name": "rule_list",
  "description": "List all active rules for this project",
  "inputSchema": {
    "type": "object",
    "properties": {
      "category": { "type": "string", "description": "Filter by category" }
    }
  }
}
```

**Retorna:**
```json
{
  "rules": [
    { "id": "no-secrets", "name": "No secrets hardcodeados", "category": "security", "description": "..." }
  ],
  "total": 8
}
```

---

#### `rule_check`
Verifica si una acción propuesta viola alguna rule activa.

```json
{
  "name": "rule_check",
  "description": "Check if a proposed action violates any active rule",
  "inputSchema": {
    "type": "object",
    "required": ["action"],
    "properties": {
      "action":  { "type": "string", "description": "The action you're about to take" },
      "context": { "type": "string", "description": "Additional context" }
    }
  }
}
```

**Retorna:**
```json
{
  "violations": [
    {
      "rule_id": "no-any-typescript",
      "rule_name": "No any en TypeScript",
      "severity": "hard",
      "suggestion": "Use 'unknown' with type guards instead"
    }
  ],
  "approved": false
}
```

---

#### `rule_get`
Retorna el contenido completo de una rule específica.

```json
{
  "name": "rule_get",
  "description": "Get full content of a specific rule",
  "inputSchema": {
    "type": "object",
    "required": ["rule_id"],
    "properties": {
      "rule_id": { "type": "string" }
    }
  }
}
```

---

### Módulo Skills (3 tools)

#### `skill_load`
Carga el contenido completo de un skill por nombre o ID.

```json
{
  "name": "skill_load",
  "description": "Load a skill's full content by name",
  "inputSchema": {
    "type": "object",
    "required": ["name"],
    "properties": {
      "name": { "type": "string", "description": "Skill name or ID" }
    }
  }
}
```

**Retorna:** El contenido Markdown completo del skill con toda su guía técnica.

---

#### `skill_search`
Busca skills relevantes para un contexto o tarea.

```json
{
  "name": "skill_search",
  "description": "Find relevant skills for a given context or task",
  "inputSchema": {
    "type": "object",
    "required": ["query"],
    "properties": {
      "query": { "type": "string", "description": "What you're trying to do" },
      "limit": { "type": "integer", "default": 5 }
    }
  }
}
```

---

#### `skill_list`
Lista todos los skills disponibles con su trigger.

```json
{
  "name": "skill_list",
  "description": "List all available skills with their triggers",
  "inputSchema": {
    "type": "object",
    "properties": {
      "tag": { "type": "string", "description": "Filter by tag" }
    }
  }
}
```

---

### Módulo Agents (4 tools)

#### `agent_list`
Lista todos los agentes disponibles en el proyecto.

```json
{
  "name": "agent_list",
  "description": "List all available agents in this project",
  "inputSchema": {
    "type": "object",
    "properties": {}
  }
}
```

---

#### `agent_activate`
Activa un agente cargando su identidad, scope y skills en el contexto.

```json
{
  "name": "agent_activate",
  "description": "Activate an agent role, loading identity, scope and skills into context",
  "inputSchema": {
    "type": "object",
    "required": ["agent_id"],
    "properties": {
      "agent_id": { "type": "string" },
      "load_skills": { "type": "boolean", "default": true, "description": "Whether to load full skill content" }
    }
  }
}
```

**Retorna:**
```json
{
  "agent": {
    "name": "Backend Engineer",
    "identity": "...",
    "scope": { "owns": ["src/modules/**"], "avoids": ["src/components/**"] },
    "skills_loaded": ["clean-architecture", "error-handling", "testing"],
    "handoffs": [...]
  },
  "rules_reminder": ["No secrets hardcodeados", "No any en TypeScript"],
  "fenrir_context": null
}
```

---

#### `agent_context`
Retorna el contexto completo del agente actualmente activo.

```json
{
  "name": "agent_context",
  "description": "Get full context of the currently active agent",
  "inputSchema": {
    "type": "object",
    "properties": {}
  }
}
```

---

#### `agent_handoff`
Formaliza la transferencia entre agentes con un contrato explícito.

```json
{
  "name": "agent_handoff",
  "description": "Formalize handoff from current agent to another with explicit contract",
  "inputSchema": {
    "type": "object",
    "required": ["to_agent", "contract"],
    "properties": {
      "to_agent":  { "type": "string", "description": "Target agent ID" },
      "contract":  { "type": "string", "description": "What you're handing off: artifacts, decisions, context" },
      "completed": { "type": "string", "description": "What the current agent completed" },
      "pending":   { "type": "string", "description": "What remains for the next agent" }
    }
  }
}
```

**Retorna:** Contexto completo del agente receptor + contrato registrado.

---

### Módulo Workflows (4 tools)

#### `workflow_start`
Inicia un workflow, cargando el primer paso y el agente asignado.

```json
{
  "name": "workflow_start",
  "description": "Start a workflow and load first step with assigned agent",
  "inputSchema": {
    "type": "object",
    "required": ["workflow_id"],
    "properties": {
      "workflow_id": { "type": "string" },
      "context":     { "type": "string", "description": "Goal or context for this workflow run" }
    }
  }
}
```

**Retorna:**
```json
{
  "workflow": "Feature Development",
  "step": 1,
  "step_name": "Entender",
  "agent": "none",
  "description": "Leer el requerimiento completo...",
  "output_expected": "Confirmación de entendimiento",
  "total_steps": 7
}
```

---

#### `workflow_step`
Avanza al siguiente paso del workflow activo.

```json
{
  "name": "workflow_step",
  "description": "Advance to next step in active workflow",
  "inputSchema": {
    "type": "object",
    "properties": {
      "step_output": { "type": "string", "description": "What was accomplished in the current step" },
      "skip_reason": { "type": "string", "description": "If skipping this step, why" }
    }
  }
}
```

---

#### `workflow_status`
Retorna el estado actual del workflow activo.

```json
{
  "name": "workflow_status",
  "description": "Get current status of active workflow",
  "inputSchema": {
    "type": "object",
    "properties": {}
  }
}
```

---

#### `workflow_complete`
Cierra el workflow verificando el Definition of Done.

```json
{
  "name": "workflow_complete",
  "description": "Complete active workflow and verify Definition of Done",
  "inputSchema": {
    "type": "object",
    "properties": {
      "dod_status": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "criterion": { "type": "string" },
            "completed": { "type": "boolean" },
            "notes":     { "type": "string" }
          }
        }
      }
    }
  }
}
```

---

### Módulo Sistema (2 tools)

#### `skoll_status`
Estado del sistema Skoll: índice RSAW y conectividad con Fenrir.

```json
{
  "name": "skoll_status",
  "description": "Skoll system status: RSAW index health and Fenrir connectivity",
  "inputSchema": {
    "type": "object",
    "properties": {}
  }
}
```

**Retorna:**
```json
{
  "rsaw": {
    "rules": 8,
    "skills": 12,
    "agents": 4,
    "workflows": 4,
    "validation_errors": []
  },
  "fenrir": {
    "available": true,
    "version": "1.2.0"
  },
  "active_workflow": null,
  "active_agent": "backend"
}
```

---

#### `skoll_validate`
Valida la estructura RSAW del proyecto detectando referencias rotas y campos faltantes.

```json
{
  "name": "skoll_validate",
  "description": "Validate RSAW structure: broken references, missing fields, scope overlaps",
  "inputSchema": {
    "type": "object",
    "properties": {
      "strict": { "type": "boolean", "default": false }
    }
  }
}
```

---

## 11. CLI — Comandos

```
skoll init [--dry-run] [--rules-only] [--from <url>]
    Crear estructura .skoll/ con ejemplos, detectar stack, generar configs

skoll mcp
    Iniciar servidor MCP en stdio

skoll tui
    Lanzar TUI interactivo

skoll validate [--ci] [--strict]
    Validar estructura RSAW completa

skoll status
    Estado del sistema y del índice RSAW

# Rules
skoll rules list [--category <cat>]
skoll rules check "<acción>"

# Skills
skoll skills list [--tag <tag>]
skoll skills show <nombre>
skoll skills add <nombre>
skoll skills import <fuente> [<nombre>]

# Agents
skoll agents list
skoll agents show <nombre>
skoll agents add <nombre>

# Workflows
skoll workflows list
skoll workflows show <nombre>
skoll workflows add <nombre>
skoll workflows run <nombre>

skoll version
```

---

## 12. Configuración

### `skoll.json`

Ver sección 9 — Modelo de datos.

### Variables de entorno

```
SKOLL_DIR            Directorio RSAW (default: .skoll)
SKOLL_LOG_LEVEL      debug|info|warn|error (default: info)
SKOLL_FENRIR_MCP     URL del servidor MCP de Fenrir (default: vía stdio)
```

---

## 13. Integración con Fenrir

La integración es completamente opcional. Skoll detecta si Fenrir está disponible al iniciar y activa la integración automáticamente si `fenrir.enabled: true` en `skoll.json`.

### Puntos de integración

| Evento Skoll | Tool Fenrir llamado | Por qué |
|---|---|---|
| `workflow_start` | `mem_context` + `predict` | Cargar contexto histórico del módulo y alertas predictivas |
| `agent_activate` | `arch_verify` (opcional) | Verificar que el scope del agente no viola decisiones existentes |
| `agent_handoff` | `mem_save` | Persistir el contrato de handoff como observación |
| `workflow_complete` | `mem_session_end` | Cerrar sesión con el DNA del workflow completo |
| `rule_check` (violation) | `audit_log` | Persistir violaciones en el audit trail |

### Cómo se comunican

Cuando Fenrir está instalado como binario en el mismo sistema, Skoll lo llama vía su HTTP API (puerto 7438) o directamente como proceso MCP en stdio. No hay dependencia de importación entre ambos — son dos binarios independientes.

```go
// internal/engine/fenrir.go
type FenrirClient struct {
    baseURL    string
    available  bool
}

func (c *FenrirClient) Ping() bool {
    resp, err := http.Get(c.baseURL + "/health")
    return err == nil && resp.StatusCode == 200
}

func (c *FenrirClient) MemContext(module string) (*ContextResult, error) {
    // HTTP call a Fenrir si está disponible
    // Retorna nil gracefully si no está disponible
}
```

---

## 14. Adapters por herramienta

`skoll init` genera los archivos de configuración para cada herramienta detectada, similar a Fenrir pero con el protocolo FENRIR.md de Skoll.

### Claude Code → `FENRIR.md` (append si ya existe)

```markdown
## Skoll RSAW Protocol

You have access to Skoll via 16 MCP tools for role-based structured development.

### At session start:
- Call `skoll_status` to verify RSAW index is loaded
- Call `agent_activate` with your role for this task

### Before any task:
- Call `rule_check` to verify you're not violating project rules
- Call `skill_search` to find relevant skills for what you're about to do

### For complex features:
- Call `workflow_start` with the appropriate workflow
- Follow steps strictly — do NOT skip steps without `skip_reason`

### On handoff:
- Call `agent_handoff` with explicit contract before switching roles
```

---

## 15. Distribución

### Instalación

```bash
# Homebrew
brew install skoll-dev/tap/skoll

# Desde source
git clone https://github.com/tu-org/skoll
cd skoll && go install ./cmd/skoll

# Instalar ambos juntos (recomendado)
brew install fenrir/tap/fenrir skoll-dev/tap/skoll
```

### Setup en un proyecto

```bash
skoll init
# ✅ Creado: .skoll/rules/global.md
# ✅ Creado: .skoll/rules/security.md
# ✅ Creado: .skoll/skills/git-workflow.md
# ✅ Creado: .skoll/skills/clean-architecture.md (detectado: NestJS)
# ✅ Creado: .skoll/agents/backend.md (detectado: NestJS)
# ✅ Creado: .skoll/agents/frontend.md (detectado: React)
# ✅ Creado: .skoll/workflows/feature.md
# ✅ Creado: .skoll/skoll.json
# ✅ Actualizado: FENRIR.md (protocolo del agente)
# 🐺 Skoll listo. Reinicia tu agente para activar.
```

---

## 16. Métricas de éxito

| Métrica | Objetivo 3 meses | Objetivo 6 meses |
|---|---|---|
| Instalaciones via Homebrew | 300 | 1,500 |
| Proyectos con .skoll/ en GitHub (detectable vía búsqueda) | 50 | 300 |
| Estrellas en GitHub | 150 | 600 |
| Tasa de workflow_complete / workflow_start | > 70% | > 80% |
| Skills creados por la comunidad | 20 | 80 |

---

## 17. Roadmap y milestones

| Fase | Semanas | Deliverable | Milestone |
|---|---|---|---|
| 0 — Setup | 1 | Repo, CI, tipos Go, Makefile | Compilación limpia |
| 1 — RSAW Loader | 2–3 | Parseo de Markdown, índice en memoria, file watcher | `skoll status` muestra índice RSAW |
| 2 — Rules Module | 4 | rule_list, rule_check, rule_get | Agente puede verificar rules antes de actuar |
| 3 — Skills Module | 5 | skill_load, skill_search, skill_list | Agente puede cargar skills on-demand |
| 4 — Agents Module | 6–7 | agent_activate, agent_handoff, agent_list, agent_context | Handoff entre agentes con contrato |
| 5 — Workflows Module | 8–9 | workflow_start, workflow_step, workflow_status, workflow_complete | `skoll workflows run feature` funciona end-to-end |
| 6 — Init & Detect | 10 | skoll init con detección de stack, templates base | Setup en < 30s en proyecto vacío |
| 7 — Fenrir Integration | 11 | fenrir.go client, puntos de integración | Workflow enriquecido con contexto de Fenrir |
| 8 — Adapters | 12 | 6 adapters de herramientas AI | Config generada para Claude Code, Cursor, Windsurf |
| 9 — Validate & TUI | 13–14 | skoll validate, TUI navegable | CI puede validar .skoll/ automáticamente |
| 10 — Release | 15 | Homebrew tap, docs, GitHub release | `brew install skoll` funciona |

---

## 18. Riesgos y mitigaciones

| Riesgo | Probabilidad | Impacto | Mitigación |
|---|---|---|---|
| Formato Markdown es demasiado libre para parsear de forma confiable | Alta | Alto | Definir schema estricto con campos obligatorios; validación en `skoll validate` |
| Developers mezclan implementación en workflows o procesos en agents | Alta | Alto | `skoll validate` detecta estas violaciones con mensajes claros; ejemplos base de `skoll init` sirven como referencia |
| Los agentes ignoran el protocolo FENRIR.md y no llaman los tools | Media | Alto | Hacer el protocolo lo más corto y claro posible; `skoll_status` al inicio hace evidente el estado |
| Scope overlap entre agents genera confusión | Media | Medio | `skoll validate` detecta overlap y advierte; `agent_activate` lo informa |
| Adopción lenta porque requiere crear archivos Markdown | Alta | Medio | `skoll init` genera todo automáticamente en < 30 segundos |
| Integración con Fenrir no es fluida en práctica | Media | Bajo | Integración completamente opcional; Skoll es 100% funcional sin Fenrir |
| Skills quedan desactualizados respecto a versiones de frameworks | Media | Medio | `skoll skills import` + skill_sources en config permite actualización controlada |

---

*Skoll PRD v1.0 — Marzo 2026*
*Go 1.22+ · Stateless (Markdown files) · MCP Protocol · MIT License*
*Funciona solo o con Fenrir · Sin base de datos · Sin dependencias de runtime*
