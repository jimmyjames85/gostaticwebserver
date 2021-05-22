# gostaticwebserver

## Deploy

Export environment variables

 - PORT (default is 8080)
 - ROUTE_FILE (should be json formated)

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

 - empty route (`""`) routes root `/` level requests
   - `domain.com/file1`->`/www/file1`
   - `domain.com/file2`->`/www/file2`
 - all other routes are considered domain prefixes e.g.
   - `dir1.domain.com` -> `/www/dir1`
   - `dir2.domain.com` -> `/www/dir2`
   - `dir2.domain.com/foo` -> `/www/dir2/foo`
   - `dir2.subdir.com/foo` -> `/www/dir2/subdir/foo`
