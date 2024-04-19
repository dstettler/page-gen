# Usage

This document will just provide the command line arguments that can/must be set
to run page-gen. For the input files, see [format.md](format.md).

## Required Positional Arguments
- `title` - The title the generated page will have (`[title].html`).

## Optional Arugments
| Argument                  | Default                    | Description |
| ---                       | ---                        | --- |
| `template`                | `reldir/template.html`     | The template that will be used to generate the page. |
| `rel-dir`                 | Current working dir        | Directory for page-gen to read files relative from. |
| `output-dir`              | `reldir/output/`           | Directory to output generated page (s) |
| `content_file`            | `reldir/[title]-conent.yml`|  
| `default-variable-val`    | N/A                        | Global default for ***ALL*** unset `V_` values. |