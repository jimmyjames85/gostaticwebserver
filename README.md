# gostaticwebserver

Execute `run.sh` and try curling the following:

```bash
curl localhost:8080/parse/json -d '{"ID":123,"NAME":"jimbo", "Info":"extra info"}'
curl localhost:8080/parse/json -d '{"USERID":123,"NAME":"jimbo", "Info":"extra info"}'
curl localhost:8080/parse/json -d '{"userid":123,"name":"jimbo"}'

curl localhost:8080/parse/parameters -d foo="one" -d bar="fish" -d biz="two" -d baz="fish"
curl localhost:8080/parse/parameters -d foo="foo1" -d foo="foo2" -d foo=3 -d bar="do ampersands have a speacial meaning"
curl localhost:8080/parse/parameters -d bar="do & ampersands & have & a & speacial & meaning"
curl localhost:8080/parse/parameters -d bar="how do you escape ampersands \& hmmmm?"
curl localhost:8080/parse/parameters -d "`ls -lah`"
```
