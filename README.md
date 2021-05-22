# gostaticwebserver

## Deploy

Export environment variables

 - PORT (default is 8080)
 - ROUTE_FILE (should be json formated)

## Route File Example

``` json
{
    "basedir": "/home/jim/tmp",
    "routes": {
	    "": "{{.basedir}}",
	    "dir1": "{{.basedir}}/dir1",
	    "dir2": "{{.basedir}}/dir2"
    }
}
```
 - `""` empty route is base route
 - `dir1`  will route dir1.domain.com to basedir/dir1
 - `dir2`  will route dir1.domain.com to basedir/dir2
