# Juju Documentation Diagrams

This file is the contributor guide for diagrams in the Juju documentation. The short rules live in [AGENTS.documentation.rules.md](AGENTS.documentation.rules.md#diagrams); this file adds the process.

## When a diagram earns its place

A diagram earns its place where prose serialises a structure badly. Four kinds of fact do that, and each has one form:

| Kind | The fact | Form | Threshold |
|---|---|---|---|
| Entities | Four to twelve tables or entities joined by relations with multiplicities | ERD (a diagram of entities and their relations) | Fewer than four: prose. More than twelve: split by the sections of the page. |
| Operation | An ordered exchange among actors, started from an interface a client uses | Sequence | At least three actors and four ordered steps, with at least one actor outside the controller or a worker that runs on its own, and an order that decides what the reader sees (which error comes first, what is left after a failure). An operation that runs inside the controller (check, validate, store) gets no sequence: its authority, rules and errors say it all. |
| Life or status | Three or more states with allowed changes that the code enforces | State diagram | A state machine that allows nearly every change (the relation status allows 24 of 30) is a list of exceptions: write the list. |
| Topology | A deployed shape whose multiplicities differ per deployment | Diagram | Multiplicity that the reader would otherwise count in prose (three controller machines, two models). |

A subject with none of the four gets no diagram, whatever the page length. Check every statement of a figure against the code before drawing it; a figure inherited from an older page is weak evidence.

## Set up the tools

1. Create the docs virtualenv: run `make install` in `docs/`. It installs ggarch at the commit pinned in `docs/requirements.txt`.
2. Activate it (`. .venv/bin/activate`). `ggarch --help` lists the commands: `check`, `scaffold`, `dump`, `solve`, `route`, `render`, `render-seq`.
3. Preview: run `make run` in `docs/`. Edits to a page or to `docs/juju.ggarch` rebuild the affected pages.

The pinned renderer has two limits. The `ggarch render` command draws diagram views only, and the Sphinx build draws state views, so preview a state view through the build. The pin is a full commit SHA: bump it in the same PR as a diagram that needs a newer renderer, and build both formats afterwards.

## Add a diagram

All diagrams are views of the models in `docs/juju.ggarch`. A model declares what exists, and a view selects and arranges it.

1. **Choose the model.** `Juju` holds the architecture (software, clouds, entities). `Schema` holds the database tables behind the reference pages. `Lifecycles` holds state machines. Add a model only when node ids would collide with an existing one.
2. **Add the nodes and edges once.** Reuse an existing node when the story is the same (the client, the controller, Charmhub, the cloud). Use an existing node and edge type. Colour means ownership only: orange for Juju, grey for outside Juju.
3. **Add the view.** A `diagram` view has `select`, `positions`, and `annotations`. Stamp repeated machinery with `instances:` and show peer wiring with `pairing: mesh`. Keep one scenario to one view.
4. **For an ERD, scaffold.** Generate the nodes and edges from the DDL: `ggarch scaffold model.ddl --db model --tables a,b,c` (the DDL comes from `schema.ModelDDL().Patches()` in `domain/schema`). The scaffold adds the foreign-key parents of the tables you list. Select only the tables you mean to draw in the view. The scaffold drops quoted column names such as `"key"`, so add those fields by hand and say so in a comment.
5. **For a sequence, cite the code.** Every step of a behaviour needs a file and line in the code that performs it, written as a comment above the behaviour in `docs/juju.ggarch`. Keep to the simplest scenario, without alternatives or loops, and put the failures in the caption and the page's errors. A sequence sits under the operation it shows, never in a section of its own.
6. **For a state diagram, take the states and transitions from the code.** Declare each transition as an edge and as a step of a behaviour, and list the guard on the step. Put the steps that lead to a later state last, so the layout keeps the main flow in a row.
7. **Validate.** Run `ggarch check docs/juju.ggarch`.

## Embed a view

A figure is the first thing in its section, before the prose: it is a preview that stands on its own, so its caption and alt text must carry every fact it shows. Use the `{ggarch}` directive. `:caption:` carries the technical detail and defines each term in brackets where it first appears. Start it with the subject and the point. `:alt:` restates the figure for readers who cannot see it: every node and every labelled arrow, because AI readers meet the alt text and not the picture. The text around the figure stays light.

````text
```{ggarch}
:file: ../juju.ggarch
:view: Controller high availability
:caption: A highly available controller runs on several machines, and each machine runs an API server and a Dqlite node.
:alt: Three controller machines, each with an API server and a Dqlite node.
```
````

Use `:slides:` and `:slide-captions:` (captions separated by `|`) for several views in one place. Never commit rendered SVGs: the docs build renders them. When a diagram replaces an older figure, delete the old image and its `.excalidraw` source in the same PR.

## Review checklist

- The figure passes one of the four thresholds, and its source file and line (or DDL table) is recorded in the PR description.
- Every node, edge, label, and multiplicity matches the code, checked by reading it.
- `ggarch check docs/juju.ggarch` passes.
- Both builders pass: `sphinx-build -b dirhtml` and `sphinx-build -b latex` from `docs/` (Read the Docs builds a PDF, and a missing visitor fails the whole build). Only the baseline warnings appear (`unknown node type: <meta>` for each page, and the duplicate-target warnings in one how-to under latex).
- The rendered figure is readable in the Read the Docs preview in light and dark mode, and at phone width.
- The caption defines its terms and starts with the subject and the point. The alt text lists every node and arrow.
- No rendered SVG, no new node type without a new badge or border, and no second model for a story that an existing model tells.
