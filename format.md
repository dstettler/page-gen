# Format

This document will explain the general format of the files that will be read by
page-gen to generate output webpages.

## Templates

Within a template file, variables are written using `${var}` syntax.
When generating the page, each variable will be filled in using preset definitions
from a separate `contents.yml` file.
By default, page-gen will use the format of `[pagename]-contents.yml` relative to
the given `--reldir` directory. If this cannot be found, 
and an alternative is not specified by way of a `--contentfile` argument, execution will be halted.

There are some default variables that page-gen keeps internally:
- `PG_TITLE` - Automatically generated based on the given title, but can also be overridden
in the contents file using the `title` key at the top level.
- `PG_DATE_SAVED` - Will be automatically generated when page-gen is initially run.
This cannot be overwritten.

Additionally, top-level variables in the contents file can be referenced by appending
`V_` to their name. For example, if one declares a variable `new_header_name` in their
contents file, this would be referenced in a template with `V_new_header_name`.
This allows a more granular approach to templating, but page-gen will *require* that variable to be
found in the contents file, unless the `--default-variable-val` parameter is set.

The `for` tag, explained in the [next section](#contents), is also usable here, and uses the exact same
types of attributes as in content body. Note that if an iterator `var` is not found, even with a
`--default-variable-var` set, page-gen will immediately halt execution.

### `template.xml` Example
```xml
<html>
	<head>
		<!-- header stuff -->
	</head>

	<body>
		<h1>${PG_TITLE}</h1>
		<p id="date">${PG_DATE_SAVED}</p>
		<p id="author">${V_author}</p>
		${V_body}
		<for var="footeritems" refname="item">
			<div class="footer-content"> </div>
		</for>
	</body>
</html>
```

## Contents

Variables can be declared in the contents file with either a global name, which can be accessed
both in the body and template, or in the form of arrays. For the sake of simplicity, an array structure
may only be one level deep, and may *only* contain keys, 
so `app.appname` would be possible, but `app.purposes.true_purpose` would not be.

In addition, iteration - provided with the fancy new `<for>` tag - cannot be done manually with ranges.
The `for` tag requires two attributes:
- `var` - This is the name of the array to read from in the contents file.
- `refname` - This is the name with which the current iteration will be referred to.

Within a for block, a given value can be selected using `${refname.value}` syntax, much like with templates.

It is also possible to declare a default value for a given `iterable.value` in an array using the `defaults` structure.
Here, simply use the name of your array, followed by the name of the value you wish to set a default for.
The parser will check for a default value if one isn't found in a given array structure.

### `contents.yml` Example
```yaml
---
  title: exceptionally important title!
  body: |
    <for var="apps" refname="app"> 
    <div class="box">
	<p>Appname: ${app.appname}, Purpose: ${app.purpose}</p>
    </div>
    </for>
  
  defaults:
    - apps.purpose: Wasn't sure what to put here!
    - apps.appname: wowie what

  apps:
    - appname: atomizer3000
      purpose: Destroy the world!

    - appname: foobar2000
      purpose: Groovin

    - appname: app-inator1000
      purpose: Pad my github stats

    # Purpose will be the default value!
    - appname: out-of-ideas-0000

  footeritems: [nothing of note number 1, nothing of note number 2]

  author: Me!
```

## Custom Tags

For all your fancy operation needs!

- `<if-exists>` - Only renders content if the given key is defined in contents.
	- `varname` - Key to check contents for
- `<if>` - Only renders content if the value of `var` is `condition` to `val`.
	- `var` - Name of the variable defined in contents
	- `val` - Value to compare `var` to
	- `condition` - Type of comparator to use (note that only the top two are available for strings):
		- `eq` - Equal to (==)
		- `ne` - Not equal to (!=)
		- `gt` - Greater than (>)
		- `lt` - Less than (<)
		- `gte` - Greater than or equal to (>=)
		- `lte` - Less than or equal to (<=)
- `<for>` - Renders the content for each iteration of
the array referenced by name `var`
	- `var` - Key to array defined in contents
	- `refname` - Name to refer to the current iteration as


## (Possible) Future Expansion

I'd like to add relative file structures in the future to this.
For example, I'd like for the user to be able to specify something to the effect of `pg-gen:filename`
like shown here:
```html
<img src="pg-gen:img/img.jpg" />
<a href="pg-gen:pages/page1.html">Another generated page!</a>
```
using a prepended value like `pg-gen` to be able to access files in the specified `--reldir`.
The final generated page would include and absolute path to the reldir + whatever was specified.

This would require a lot of file verification when parsing templates and bodies, but I think
it could make the tool especially valuable. 
