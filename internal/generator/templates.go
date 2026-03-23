package generator

const GlobalRulesTemplate = `# Rules: Global

## No secrets hardcodeados
Nunca escribir API keys, passwords o tokens directamente en código.
Usar siempre variables de entorno o secret managers.

## No any en TypeScript
Nunca usar el tipo any en TypeScript. Usar unknown con type guards o tipos específicos.

## Idioma del usuario
Responder siempre en el idioma del usuario, independientemente del idioma del código.
`

const GitSkillTemplate = `# Skill: Git Workflow

## Cuándo aplicar
Al crear commits, branches o pull requests.

## Proceso
Formato de commit: type(scope): description
Ejemplo: feat(auth): add JWT refresh token rotation

## Checklist
- [ ] Revisar el diff
- [ ] No hay secrets
- [ ] tests pasan
`

const BackendAgentTemplate = `# Agent: Backend Engineer

## Rol
Arquitecto y desarrollador senior de backend.

## Scope
### Archivos que me pertenecen
- src/modules/**
- src/services/**
- internal/**

### Lo que NO toco
- src/components/**
- .github/**

## Skills que aplico
- git-workflow.md
`

const FeatureWorkflowTemplate = `# Workflow: Feature Development

## Propósito
Proceso completo para desarrollar un feature.

## Pasos
### Paso 1 — Entender (none)
Leer el requerimiento.
**Output esperado:** Confirmación del usuario.

### Paso 2 — Implementar (Backend Engineer)
Escribir el código sugerido.
Skills: git-workflow.md
**Output esperado:** Funcionalidad implementada.
`
