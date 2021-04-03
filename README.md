gxs
===

G(o) X(cross) S(titch) is a simple ascii-driven cross stitch pattern builder

## build

```
make
```

## patterns

`gxs` uses a declaration of patterns which is based on building 1 to N layers
and committing each layer once it is "complete" (allowing alterations on top of previous layers).
This allows for an "ASCII-like" representation of the pattern to exist for editing via text editor
(e.g. `vim`) while ending up with a full resulting pattern

### example

#### mode

mode is the type of stitch that will be built

| type | explanation |
| ---  | ---         |
| xstitch   | full X cross stitch |
| topedge   | back stitch at the top edge |
| bottomedge   | back stitch at the bottom edge |
| leftedge   | back stitch on the left edge |
| rightedge   | back stitch on the right edge |
| hline | back stitch horizontally through |
| vline | back stitch vertically through |

```
mode => {
    xstitch
}
```

#### palette

configures the colors to use within the ascii pattern

```
palette => {
    x => red
    y => #333333
    z => NONE
}
```

in the above the "x" symbol will be red, "y" will be the css "#333333" and "z"
uses the "NONE" special keyword to allow for spacing ascii symbols but not producing a
color.

_When a named color is given that matches a known DMC floss, it will result in the
RGB floss for that color, in the above example the 'red' value will match a floss_

#### pattern

define the ascii pattern to draw onto a resulting grid

```
pattern => {
    xyzyz
    zyxyz
    zzxzz
}
```

which will produce a simple pattern

#### action

finally tell `gxs` to commit the stitching layer

```
action => {
    commit
}
```

#### example

```
# comments start with '#' but can NOT be within a '{' and '}' block
palette => {
    x => red
    y => #333333
    z => NONE
}
mode => {
    xstitch
}
pattern => {
    xyzyz
    zyxyz
    zzxzz
}
action => {
    commit
}
# reuse the palette, single-line block (only for single-line commands)
mode => {topedge}
# patterns must always be redefined
pattern => {
    xxxxx
    zzzzz
    xxxxx
}
action => {commit}
```

see `examples/` for more.
