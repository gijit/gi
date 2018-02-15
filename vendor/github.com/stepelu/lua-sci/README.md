SciLua: Scientific Computing with LuaJIT
========================================

A complete framework for numerical computing based on LuaJIT which combines the ease of use of scripting languages (MATLAB, R, ...) with the high performance of compiled languages (C/C++, Fortran, ...).

## Modules

<table>
<tr><th>Sub-Module</th><th>Description</th></tr>
<tr><td><code><a href="http://www.scilua.org/sci_math.html">sci.math</a></code></td><td>special mathematical functions</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_diff.html">sci.diff</a></code></td><td>automatic differentiation</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_alg.html">sci.alg</a></code></td><td>vector and matrix algebra</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_quad.html">sci.quad</a></code></td><td>quadrature algorithms</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_root.html">sci.root</a></code></td><td>root-finding algorithms</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_fminfmax.html">sci.fmin</a></code></td><td>function minimization algorithms</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_fminfmax.html">sci.fmax</a></code></td><td>function maximization algorithms</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_prng.html">sci.prng</a></code></td><td>pseudo random number generators</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_qrng.html">sci.qrng</a></code></td><td>quasi random number generators</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_stat.html">sci.stat</a></code></td><td>statistical functions</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_dist.html">sci.dist</a></code></td><td>statistical distributions</td></tr>
<tr><td><code><a href="http://www.scilua.org/sci_mcmc.html">sci.mcmc</a></code></td><td>MCMC algorithms</td></tr>
</table>

## Install

This module is included in the [ULua](http://ulua.io) distribution, to install it use:
```
upkg add sci
```

Alternatively, manually install this module making sure that all dependencies listed in the `require` section of [`__meta.lua`](__meta.lua) are installed as well (dependencies starting with `clib_` are standard C dynamic libraries).

## Documentation

Refer to the [official documentation](http://scilua.org).