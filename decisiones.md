# Decisiones — Proyecto IngSoft3

## TP1 — Control de versiones

### Por qué Git no pudo resolver el conflicto solo
Las ramas `feature/titulo-a` y `feature/titulo-b` partieron las dos del mismo commit de `main`, 
y ambas modificaron la misma línea del `README.md`, cada una con un texto distinto
("versión A" vs "versión B"). El merge de tres vías que usa Git compara el commit común (la base),
la rama A y la rama B — y solo puede fusionar automático cuando los cambios no se pisan. Como las
dos tocaron exactamente la misma línea con contenido distinto, Git no tiene forma de decidir y necesita que una persona lo resuelva.

La rama B tendría que haber traído los cambios de `main` (ya con A mergeada) antes de abrir su
 propio PR — así el conflicto aparecía localmente, antes de llegar a GitHub, sabiendo de antemano qué había en A.

### Problemas encontrados y cómo los resolví
Para que el conflicto entre las ramas A y B apareciera de verdad (y no se mergeara solo), las dos
  tenían que partir del mismo commit de `main` y modificar la misma línea del README. Si una rama
  partía de la otra en vez de `main`, Git no generaba ningún conflicto — hubo que prestar atención
  al orden: crear ambas ramas desde `main`, y recién mergear la primera después de tener las dos
  abiertas.

### Uso de IA
Usé Claude como guía durante todo el TP1: para entender la protección de rama, interpretar los
marcadores de conflicto antes de resolverlos, redactar la descripción de la release, y resolver
dudas puntuales. Verifiqué cada paso ejecutándolo yo misma y mirando el resultado real en GitHub.
