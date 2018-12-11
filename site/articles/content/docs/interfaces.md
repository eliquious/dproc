---
title: "Interfaces"

# Set type to 'docs' if you want to render page outside of configured section or if you render section other than 'docs'
type: docs

# Set page weight to re-arrange items in file-tree menu (if BookMenuBundle not set)
weight: 1

# (Optional) Set to mark page as flat section in file-tree menu (if BookMenuBundle not set)
bookFlatSection: true

# (Optional) Set to hide table of contents, overrides global value
bookShowToC: false
---

# Interfaces and structs

dProc is a very small package and as such there are very few interfaces. The package is fairly opinionated but the interfaces allows for custom processes if the need arises.

For the high-level overview, there are 4 main interfaces or types: `Message`, `Process`, `Handler` and the `Engine`. Messages are used to pass data through the pipeline, processes handle the pipeline life-cycle, handlers handle the messages and the engine manages all the processes.

MessageTypes and Handlers are written for each specific pipeline, while the processes and the engine are more generic. There may be a need to write a specific processes if the life-cycle needs to change outside the norm.

There is a separate page for each of the main types. Check them out on the left (or in the menu).

