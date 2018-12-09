---
title: "Intro"

# Set type to 'docs' if you want to render page outside of configured section or if you render section other than 'docs'
type: docs

# Set page weight to re-arrange items in file-tree menu (if BookMenuBundle not set)
weight: 1

# (Optional) Set to mark page as flat section in file-tree menu (if BookMenuBundle not set)
bookFlatSection: true

# (Optional) Set to hide table of contents, overrides global value
bookShowToC: true
---

# Welcome to the dProc docs!

So you read the homepage and clicked this link. Now you're probably asking yourself what the hell is dProc and why do I need it? So without further ado...

## What is dProc?

dProc is a small, generic data processing library written in Go. It is modeled after [Actors](https://en.wikipedia.org/wiki/Actor_model) and a bit of [Flow-based programming](https://en.wikipedia.org/wiki/Flow-based_programming). It is a small library (< 250 lines of code) with very few interfaces.

Actors, if you are not familiar, are components which are often used in concurrent computation. They feature an inbox, essentially a queue, in which they receive messages for processing. Afterwards, they may send a new message to another actor and the process continues until the input data is fully processed.

In dProc, a pipeline is established first then it is initialized. Flow-based programming is similar in nature in that there is an outline or diagram of processes which handle the data. Often times, flow-based programming is built with a network of black boxes. Each with their own methods of processing the input data and forwarding output data. In this library, this black-box components are called processes.

dProc merges these models into a single package where the data pipeline is made up of components which happen to be similar to actors. After the pipeline is built, an engine starts the pipeline and runs until completion. Generally, there is a source process which generates the data stream and is then sent to the child processes.

## Why do I need it?

Well, to be honest. You might not. But it depends on what you need. If you have data that needs processing in various ways, perhaps using streams, then dProc might be useful. The ideas behind dProc were initially written in Python around 2010. It has since been used for many projects and has proven to be a good model for multi-stage, stream processing.

It is particularly useful for processing data using various means to collate and reduce large streams of data. For instance, if you have a large file or network stream that needs to be analyzed in several ways and then output to a secondary file, then dProc was built for this. However, most data can be manipulated in a similar fashion. This framework allows for a stable development path which allows for growth and variability.

Please, checkout the examples and in-depth articles to better understand the design and principles of the package.

You can find the godocs [here](https://godoc.org/github.com/eliquious/dproc).
