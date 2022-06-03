# templates

As templates are used for html rendering as well as for mail bodies,
they are located in `html/` and `mail/`.

## Naming
Templates should define their name using the `define` tag. They should include
the directory structure they are contained (without `html/` or `mail/`).

`html/base/head.gohtml`:
```gotemplate
{{define "base/head"}}
...
{{end}}
```

## html templates

Currently, only files named `*.gohtml` in `html/` and in direct 
subdirectories of `html/` are included.

- `html/file.gohtml`: Included
- `html/file.tmpl`: Not included
- `html/files/file.gohtml`: Included
- `html/files/more/file.gohtml`: Not included

## mail templates

Currently, only files named `*.gotxt` in `mail/` are included.

- `mail/file.gotxt`: Included
- `mail/file.tmpl`: Not included
- `mail/files/file.gotxt`: Not Included
