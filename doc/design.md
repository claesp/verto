# Design
This document describes the overall design for the Verto system. At its heart is
the _object model_. The system takes an `Importer` that can produce a struct
based on the object model from a source, and then gives that struct to an
`Exporter` that can produce some type of output from said model struct.

This should allow for arbitrary combinations of _importers_ and _exporters_.

##  Object Model
### VertoDevice
A `VertoDevice` is the root object for the model.
