# basic-info.md

## Language And Tone
1. Write in Russian.
2. Use a formal business style.
3. Keep wording concise and unambiguous, without conversational filler.

## Obsidian Compatibility
1. All content must be compatible with Obsidian Markdown.
2. Use standard markdown headings and Obsidian wikilinks.

## Page Structure
1. Each page must start with an introductory section containing a brief page summary.
2. Use a flat structure: primarily `##` headings.
3. Avoid `###` and lower-level headings whenever possible.

## Glossary Rules
1. All key terms are stored in `Глоссарий.md`.
2. Each page must include a block with links to the terms used on that page.
3. Use this pattern to insert a term:
`![[Глоссарий#{heading-name}]]`
4. Every important term in the text must link to the corresponding heading in `Глоссарий.md`.

## Minimal Page Template
```md
Краткое резюме страницы.

## Глоссарий
![[Глоссарий#Термин-1]]
![[Глоссарий#Термин-2]]

## Раздел 1
Содержимое страницы.

## Раздел 2
```
