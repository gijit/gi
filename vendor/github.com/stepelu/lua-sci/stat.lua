--------------------------------------------------------------------------------
-- Statistical functions module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

-- Variances and covariances are computed according to the unbiased version of
-- the algorithm.
-- Welford-type algorithms are used for superior numerical stability, see:
-- http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance
-- http://www.johndcook.com/standard_deviation.html

-- TODO: BIC, AIC.
-- TODO: Function join(...) to join results for parallel computing.

-- TODO: Speed-up via OpenBLAS.

local ffi = require "ffi" 
local alg = require "sci.alg"

local sqrt, abs, max, log = math.sqrt, math.abs, math.max, math.log
local vec, mat, arrayct = alg.vec, alg.mat, alg.arrayct
local typeof, metatype = ffi.typeof, ffi.metatype

local function clear(x)
  for i=1,#x do
    x[i] = 0
  end
end

local function mean(x)
  if #x < 1 then
    error("#x >=1 required: #x="..#x)
  end
  local mu = 0
  for i=1,#x do
    mu = mu + (x[i] - mu)/i
  end
  return mu
end

local function var(x)
  if #x < 2 then
    error("#x >= 2 required: #x"..#x)
  end
  local mu, s2 = 0, 0
  for i=1,#x do
    local delta = x[i] - mu
    mu = mu + delta/i
    s2 = s2 + delta*(x[i] - mu)
  end
  return s2/(#x - 1)
end

local function cov(x, y)
  local mux, muy, s2c = 0, 0, 0
  if not #x == #y then
    error("#x ~= #y: #x="..#x..", #y="..#y)
  end
  if #x < 2 then
    error("#x >= 2 required: #x="..#x)
  end
  for i=1,#x do
    local deltax = x[i] - mux
    local deltay = y[i] - muy
    local r = 1/i
    mux = mux + deltax*r
    muy = muy + deltay*r
    s2c = s2c + deltax*(y[i] - muy)
  end
  return s2c/(#x - 1)
end

local function cor(x, y)
  return cov(x, y)/sqrt(var(x)*var(y))
end

local function chk_dim(self, x)
  local d = self._d
  if d ~= #x then
    error("argument with #="..#x.." passed to statistics of dimension="..d)
  end
  return d
end

local function chk_eq_square(X, Y)
  local n, m = X:nrow(), X:ncol()
  local ny, my = Y:nrow(), Y:ncol()
  if not (n == m and n == ny and n == my) then
    error("arguments must be square matrices of equal size, passed: "..n.."x"..m
        ..", "..ny.."x"..my)
  end
  return n, m
end

local function dim(self)
  return self._d
end

local function len(self)
  return self._n
end

local function tos_mean(self)
  local m = vec(self._d)
  self:mean(m)
  return "mean:\n"..m:width()
end

local meand_mt = {
  dim = dim,
  len = len,
  clear = function(self)
    self._n = 0
    clear(self._mu)
  end,
  push = function(self, x)
    local d = chk_dim(self, x)
    self._n = self._n + 1
    for i=1,d do self._mu[i] = self._mu[i] + (x[i] - self._mu[i])/self._n end   
  end,
  mean = function(self, mean)
    local d = chk_dim(self, mean)
    if self._n < 1 then
      error("n >= 1 required: n="..self._n)
    end
    for i=1,d do mean[i] = self._mu[i] end
  end,
  __tostring = tos_mean,
}
meand_mt.__index = meand_mt

local mean0_mt = {
  dim = dim,
  len = len,
  clear = function(self)
    self._n = 0
    self._mu = 0
  end,
  push = function(self, x)
    self._n = self._n + 1
    self._mu = self._mu + (x - self._mu)/self._n
  end,
  mean = function(self)
    return self._mu
  end,
  __tostring = tos_mean,
}
mean0_mt.__index = mean0_mt

local meand_ct = typeof("struct { int32_t _d, _n; $& _mu; }", arrayct)
local mean0_ct = typeof("struct { int32_t _d, _n; double _mu; }")
meand_ct = metatype(meand_ct, meand_mt)
mean0_ct = metatype(mean0_ct, mean0_mt)

local function tos_var(self)
  local v = mat(self._d, self._d)
  self:var(v)
  return tos_mean(self).."\nvar:\n"..v:width()
end

local vard_mt = {
  dim = dim,
  len = len,
  clear = function(self)
    self._n = 0
    clear(self._mu)
    clear(self._s2)
  end,
  push = function(self, x)
    local d = chk_dim(self, x)
    self._n = self._n + 1
    local r = 1/self._n
    for i=1,d do
      self._delta[i] = x[i] - self._mu[i]
      self._mu[i] = self._mu[i] + self._delta[i]*r
      self._s2[i] = self._s2[i] + self._delta[i]*(x[i] - self._mu[i])
    end
  end,
  mean = meand_mt.mean,
  var = function(self, var)
    local d = chk_dim(self, var)
    if self._n < 2 then
      error("n >= 2 required: n="..self._n)
    end
    for i=1,d do var[i] = self._s2[i]/(self._n - 1) end
  end,
  __tostring = tos_var,
}
vard_mt.__index = vard_mt

local var0_mt = {
  dim = dim,
  len = len,
  clear = function(self)
    self._n = 0
    self._mu = 0
    self._s2 = 0
  end,
  push = function(self, x)
    self._n = self._n + 1
    local r = 1/self._n
    self._delta = x - self._mu
    self._mu = self._mu + self._delta*r
    self._s2 = self._s2 + self._delta*(x - self._mu)
  end,
  mean = mean0_mt.mean,
  var = function(self)
    if self._n < 2 then
      error("n >= 2 required: n="..self._n)
    end
    return self._s2/(self._n - 1)
  end,
  __tostring = tos_var,
}
var0_mt.__index = var0_mt

local vard_ct = typeof("struct { int32_t _d, _n; $& _mu; $& _delta; $& _s2; }", 
  arrayct, arrayct, arrayct)
local var0_ct = typeof("struct { int32_t _d, _n; double _mu, _delta, _s2; }")
vard_ct = metatype(vard_ct, vard_mt)
var0_ct = metatype(var0_ct, var0_mt)

-- Y *can* alias X.
local function covtocor(X, Y)
  local n, m = chk_eq_square(X, Y)
  for r=1,n do
    for c=1,m do
      if r ~= c then
        Y[{r,c}] = X[{r,c}]/sqrt(X[{r,r}]*X[{c,c}])
      end
    end
  end
  for i=1,n do Y[{i,i}] = 1 end
end

local function tos_cor(self)
  local c, r = mat(self._d, self._d), mat(self._d, self._d)
  self:cov(c)
  self:cor(r)
  return tos_mean(self).."\ncov:\n"..c:width().."\ncor:\n"..r:width()
end

local covd_mt = {
  dim = dim,
  len = len,
  clear = function(self)
    self._n = 0
    clear(self._mu)
    clear(self._s2)
  end,
  push = function(self, x)
    local d = chk_dim(self, x)
    self._n = self._n + 1
    local r = 1/self._n
    for i=1,d do
      self._delta[i] = x[i] - self._mu[i]
      self._mu[i] = self._mu[i] + self._delta[i]*r
    end
    for i=1,d do for j=1,d do
      self._s2[{i,j}] = self._s2[{i,j}] + self._delta[i]*(x[j] - self._mu[j])
    end end
  end,
  mean = meand_mt.mean,
  var = function(self, var)
    local d = chk_dim(self, var)
    if self._n < 2 then
      error("n >= 2 required: n="..self._n)
    end
    for i=1,d do var[i] = self._s2[{i,i}]/(self._n - 1) end
  end,
  cov = function(self, cov)
    local n, m = chk_eq_square(self._s2, cov)
    if self._n < 2 then
      error("n >= 2 required: n="..self._n)
    end
    for i=1,n do for j=1,m do 
      cov[{i,j}] = self._s2[{i,j}]/(self._n - 1)
    end end
  end,
  cor = function(self, cor)
    self:cov(cor)
    covtocor(cor, cor)
  end,
  __tostring = tos_cor,
}
covd_mt.__index = covd_mt

local covd_ct = typeof("struct { int32_t _d, _n; $& _mu; $& _delta; $& _s2; }", 
  arrayct, arrayct, arrayct)
covd_ct = metatype(covd_ct, covd_mt)

local samples_mt = {
  dim = dim,
  len = len,
  clear = function(self)
    self._n = 0
    self._x = { }
  end,
  push = function(self, x)
    local d = chk_dim(self, x)
    self._n = self._n + 1
    self._x[n] = x:copy()
  end,
  samples = function(self, samples)
    local n, m = samples:nrow(), samples:ncol()
    if n ~= self._n or m ~= self._d then
      error("output matrix has wrong dimensions")
    end
    for i=1,n do for j=1,m do 
      samples[{i,j}] = self._x[{i,j}]
    end end
  end,
}
samples_mt.__index = samples_mt

local anchor = setmetatable({ }, { __mode = "k" })

local function olmean(dim)
  if dim == 0 then
    return mean0_ct(dim)
  else
    local mu = vec(dim)
    local v = meand_ct(dim, 0, mu)
    anchor[v] = { mu }
    return v
  end
end
local function olvar(dim)
  if dim == 0 then
    return var0_ct(dim)
  else
    local mu, delta, s2 = vec(dim), vec(dim), vec(dim)
    local v = vard_ct(dim, 0, mu, delta, s2)
    anchor[v] = { mu, delta, s2 }
    return v
  end
end
local function olcov(dim)
    local mu, delta, s2 = vec(dim), vec(dim), mat(dim, dim)
    local v = covd_ct(dim, 0, mu, delta, s2)
    anchor[v] = { mu, delta, s2 }
    return v
end
local function olsamples(dim)
  return setmetatable({ _d = dim, _n = 0, _x = { } }, samples_mt)
end

return {
  mean      = mean,
  var       = var,
  cov       = cov,
  cor       = cor,
  covtocor  = covtocor,
  olmean    = olmean,
  olvar     = olvar, 
  olcov     = olcov,
  olsamples = olsamples, 
}
