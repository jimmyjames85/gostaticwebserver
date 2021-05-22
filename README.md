# gostaticwebserver

## Deploy

Export environment variables

 - PORT (default is 8080)
 - ROUTE_FILE (See json example below)

## Route File Example

``` json
{
    "basedir": "/www",
    "routes": {
	    "": "{{.basedir}}",
	    "dir1": "{{.basedir}}/dir1",
	    "dir2": "{{.basedir}}/dir2",
	    "dir2.subdir": "{{.basedir}}/dir2/subdir"
    }
}
```

 - The route file is parsed as a go template populating instances of `{{.basedir}}` with the value in
   root level `basedir` (e.g. /www)
 - The empty route (`""`) routes root `/` level requests
   - `domain.com/file1`->`/www/file1`
   - `domain.com/file2`->`/www/file2`
 - All other routes are considered domain prefixes e.g.
   - `dir1.domain.com` -> `/www/dir1`
   - `dir2.domain.com` -> `/www/dir2`
   - `dir2.domain.com/foo` -> `/www/dir2/foo`
   - `dir2.subdir.com/foo` -> `/www/dir2/subdir/foo`
