Outputting paths to a terminal is nice, one can leverage all the *nix tools they have (grep, awk, sed, etc) to filter and modify the exported paths.

```
❯ yangpath export -y ~/projects/openconfig/public/ \
                  -m ~/projects/openconfig/public/release/models/bgp/openconfig-bgp.yang \
                  | grep neighbor | grep as-path
[rw]  /bgp/neighbors/neighbor[neighbor-address=*]/as-path-options/config/allow-own-as  uint8
[rw]  /bgp/neighbors/neighbor[neighbor-address=*]/as-path-options/config/disable-peer-as-filter  boolean
[rw]  /bgp/neighbors/neighbor[neighbor-address=*]/as-path-options/config/replace-peer-as  boolean
[ro]  /bgp/neighbors/neighbor[neighbor-address=*]/as-path-options/state/allow-own-as  uint8
[ro]  /bgp/neighbors/neighbor[neighbor-address=*]/as-path-options/state/disable-peer-as-filter  boolean
[ro]  /bgp/neighbors/neighbor[neighbor-address=*]/as-path-options/state/replace-peer-as  boolean
```

But we also thought that creating an HTML service for filtering the once exported paths is a nice and powerful idea. So we added an `html` output format option.

## HTML templating
With HTML templating its possible to nicely style your paths by creating an HTML file with them. The HTML file then can be served by a server or stored somewhere with shared access.

As usual with templates, the logic for the template generation resides within the template itself. For `yangpath` the template must be done in Go templating syntax, but don't worry, its very easy to create one.

#### Creating a basic HTML template
Lets create an HTML boilerplate that will store our paths in a table format with some bootstrap styling applied.

```html
<head>
	<!-- Font Awesome -->
	<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.8.2/css/all.css">
	<!-- Google Fonts -->
	<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto:300,400,500,700&display=swap">
	<!-- Bootstrap core CSS -->
	<link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.5.0/css/bootstrap.min.css" rel="stylesheet">
	<!-- Material Design Bootstrap -->
	<link href="https://cdnjs.cloudflare.com/ajax/libs/mdbootstrap/4.19.1/css/mdb.min.css" rel="stylesheet">
</head>

<div class="container-lg">
	<h1>yangpath paths</h1>
	<table class="table table-striped">
		<thead>
			<tr>
				<th>#</th>
                <th>Path</th>
                <th>Type</th>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>1</td>
                <td>/c1/l1</td>
                <td>string</td>
			</tr>
		</tbody>
	</table>
</div>
```

Right now its not a template, but just an HTML page which we populated with static data. Lets make it a template.

#### Path variables
To create a template out of the HTML base we need to embed variables inside the file. By default `yangpath` passes the list of exported Paths into the templating engine which are accessible via `.Paths` notation. The `Paths` list contains Path structures which export the [following fields](https://github.com/hellt/yangpath/blob/master/pkg/path/path.go#L31-L36).

This means that inside our template we can range over that list and access each of the Path fields. Lets implement this:

```html
<div class="container-lg">
	<h1>yangpath paths</h1>
	<table class="table table-striped">
		<thead>
			<tr>
				<th>#</th>
				<th>Path</th>
				<th>Type</th>
			</tr>
		</thead>
		<tbody>
			{{range $i, $p  := .Paths}}
			<tr>
				<td>{{$i}}</td>
				<td>{{$p.XPath}}</td>
				<td>{{$p.Type.Name}}</td>
			</tr>
			{{end}}
		</tbody>
	</table>
</div>
```

Now lets generate the HTML using this template saved as `myTemplate.html`:

```
yangpath export -m pkg/path/testdata/test3/test3.yang \
                -f html \
                --template myTemplate.html > paths.html
```

And here we have our paths in the nice looking HTML page:
![templated](https://gitlab.com/rdodin/pics/-/wikis/uploads/620446b83db793b9de6186a4e0a3ab54/image.png)

#### Custom variables
Sometimes you need to augment your template with user-defined variables. For such scenarios we made it possible to pass variables over the CLI in a format of `key:::value` pairs.

Lets say we want to make the header of the table to reflect the original YANG file name from which the paths were exported. To do that we need to pass the value over the CLI with the user-defined key:

```
yangpath export -m pkg/path/testdata/test3/test3.yang \
                -f html \
                --template-vars fname:::pkg/path/testdata/test3/test3.yang \
                --template myTemplate.html > paths.html
```

By now, the variable `fname` will hold the string value of `pkg/path/testdata/test3/test3.yang`, which enables us to change the header in our template as follows:

```html
	<h1>yangpath paths exported from <code>{{.Vars.fname}}</code></h1>
```

Here is the result:
![custom-vars](https://gitlab.com/rdodin/pics/-/wikis/uploads/e97d52f13ee07172bad77b3fc263045e/image.png)

You can find the final template in the [repository](https://github.com/hellt/yangpath/blob/master/docs/myTemplate.html).

## Examples
With the above explained techniques its possible to create a full featured **Path Browser** which will have filtering capabilities. That is what we did for Nokia SR OS YANG models at [hellt/nokia-yangtree](https://github.com/hellt/nokia-yangtree).

![path-broser](https://gitlab.com/rdodin/pics/-/wikis/uploads/14fcedd12a511674a1ff8829ab28390f/CleanShot_2020-05-19_at_17.37.55.gif)