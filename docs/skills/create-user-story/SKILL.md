---
name: create-user-story
description: Create or update FFBC business-analysis User Story sections in Obsidian Markdown. Use when Codex is asked to write a User Story for a given filepath and story name in the FFBC docs, especially under `docs/бизнес анализ/*.md`; the skill instructions are in English, but generated story content must be in Russian.
---

# Create User Story

## Overview

Use this skill to write a page-local User Story in the FFBC documentation style. The user should provide two inputs: `filepath` and `storyname`.

The skill workflow is written in English. The User Story output itself must be written in Russian, in a concise formal business style.

## Required Inputs

- `filepath`: path to the target Markdown file, usually relative to the current repository or docs root.
- `storyname`: exact Russian section title to create or update.

If either input is missing, ask for the missing value before editing.

## Workflow

1. Read `agent_docs/basic-info.md` before reading or editing the target file.
2. Read `filepath` and identify the local context around `storyname`.
3. If the story already exists, preserve its intent and improve it only as needed.
4. If the story does not exist, add a new `## {storyname}` section in the most appropriate local User Story area.
5. Write the story in Russian, using the structure below.
6. Keep existing page terminology and Obsidian links. Add glossary links only when the target page pattern already makes the required term clear.
7. Do not rewrite unrelated sections.

## Story Structure

Follow the pattern from the `Привязывание модификатора` story:

```md
## {storyname}

**Роль:** [user role].

**Цель:** [business goal].

**Предусловие:** [state before the scenario starts].

**Основной сценарий:**

1. [Actor action].
2. [System response].
3. [Actor action].
    1. Alt. `[Alternative scenario name]`

**[Alternative scenario name]:**
3. [Alternative step].
4. Возврат к пункту N **основного сценария**.

**Результат:** [observable final business state].
```

## Writing Rules

- Use Russian for story content.
- Use formal business wording.
- Keep scenario steps short and observable.
- Separate actor actions from system responses where possible.
- Use alternatives only when the normal flow can branch.
- Reference alternative scenarios from the main scenario with `Alt. \`Name\`` and define each alternative after the main scenario.
- End alternatives with a return point when the flow rejoins the main scenario.
- Use the same numbering style as the surrounding page.
- Do not add implementation details, UI component names, API names, or database details unless the source page already treats them as business terms.

## Output Contract

After editing, report the changed file and the story section name. Mention if the story could not be written because the source context was insufficient.
