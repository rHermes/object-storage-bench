# Object storage benchmark

[![Build Status](https://rhermes.semaphoreci.com/badges/object-storage-bench/branches/main.svg?style=shields)](https://rhermes.semaphoreci.com/projects/object-storage-bench)

This repo serves as a set of tools to answer some questions about the
object storages.

It currently has support for the following stores:

- Local directory

## Tests


### Prefix speed

This test see if the prefixing capabilities of the system grow slower
as more files are added or if the underlying storage uses some sort of
**trie** implementation to make these lookups constant time.

It generates various directory trees and creates these in the underlying
stores. Some are deep, some are wide and there is also a mix. We explore
how the variation of these parameters affect the performance of the
object listings.
