gxs
===

G(o) X(cross) S(titch) is a simple web-based pattern designer utilizing JSON to define
patterns

## build

```
make
```

## web-based design

run
```
./bin/gxs
```

open the bound address (e.g. http://localhost:10987) in the browser and begin editing the JSON

## cli design

edit a file and use `gxs` to show the results

```
cat design.json | ./bin/gxs --file --
# or
./bin/gxs --file design.json
```

### advanced

#### size

Increase the pattern size via `-size <int>` (it will always be a square)

#### development

Change the internal padding identifier definitions (especially `-size` > 9999)
