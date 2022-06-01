# templates

Currently, only files named `*.gohtml` in this and in one layer subdirectories
are accessible from the template renderer.

- `file.gohtml`: Included
- `file.tmpl`: Not included
- `files/file.gohtml`: Included
- `files/more/file.gohtml`: Not included

## Naming
Templates should define their name using the `define` tag. They should include
the directory structure they are contained.

`base/head.gohtml`:
```gotemplate
{{define "base/head"}}
...
{{end}}
```
