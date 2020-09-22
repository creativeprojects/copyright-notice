# Copyright notice

This is a quick CLI tool to manage your copyright notices in your source code.
It can:
- add a copyright header automatically on files selected by extension
- exclude some folders or files (`node_modules` anyone?)
- change the year of an existing copyright
- detect a different copyright header and not touch it
- detect auto-generated files
- keep the Windows BOM on UTF-8 files

## TODO:

The tool is actually fully working, but:
- The documentation is very much work in progress.
- I want to do some refactoring.
- I also want to be able to load a configuration file instead of typing a long command line every time
