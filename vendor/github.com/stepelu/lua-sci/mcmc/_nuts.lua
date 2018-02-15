--------------------------------------------------------------------------------
-- NUTS MCMC sampler.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

-- TODO: At the moment single cache, algorithm, ecc ecc so nested calls are not
-- TODO: possible nor are use of NUTS for sampling subset of the parameters, 
-- TODO: change when improving MCMC framework.

local alg  = require "sci.alg"
local dist = require "sci.dist"
local stat = require "sci.stat"
local math = require "sci.math"
local xsys = require "xsys"

local normald = dist.normal(0, 1)
local unifm1p1d = dist.uniform(-1, 1)
local exp, log, sqrt, sign, step, min, abs, floor = xsys.from(math, 
     "exp, log, sqrt, sign, step, min, abs, floor")
local vec, mat = alg.vec, alg.mat
local width = xsys.string.width

-- Only leapfrog (for th1, r1) and evalgrad (for grad) modify elements of some
-- vector (the newly initialized ones), no vector is modified after
-- initialization => it's safe to cache based on reference instead of value for
-- the parameters in evalgrad.
local cache = { }

local invert, factor, correlate, sigestimate, kinectic

-- Caching is done based on theta by-ref. Fine as thetas are never mutated.
local function evalfgrad(fgrad, theta)
  local found = cache[theta] 
  if found then
    return found[1], found[2]
  else
    local grad = vec(#theta)
    -- Newly created grad vector is initialized:
    local val = fgrad(theta, grad)
    if not (abs(val) < 1/0) then -- If val is nan or not finite.
      val = - 1/0   -- No nan allowed.
      -- Gradient is almost surely nan or not finite:
      for i=1,#grad do grad[i] = 0 end 
    end
    cache[theta] = { val, grad }
    return val, grad
  end
end

local function stop_iter(n)
  local c = 0 
  return function()
    c = c + 1
    return c >= n
  end
end

local function invert(m, sig, mass)
  if mass == "diagonal" then
    for i=1,#m do
      m[i] = 1/sig[i]
    end
  else
    alg.invert(m, sig, "posdef")
  end
end

local function factor(mfac, m, mass)
  if mass == "diagonal" then
    for i=1,#mfac do mfac[i] = sqrt(m[i]) end
  else
    alg.factor(mfac, m, "posdef")
  end
end

local function correlate(r, mfac, mass)
  if mass == "diagonal" then
    for i=1,#r do r[i] = mfac[i]*r[i] end
  else
    alg.mulmv(r, mfac, r)
  end
end

local function sigestimate(sig, sigestimator, mass)
  if mass == "diagonal" then
    sigestimator:var(sig)
  else
    sigestimator:cov(sig)
  end
end

local function logpdf(fgrad, theta, r, sig, mass)
  local n = #r
  -- Kinetic energy: p^t*M^1*p/2 == p^t*sig*p/2.
  local sum = 0; 
  if mass == "diagonal" then
    for i=1,n do sum = sum + r[i]^2*sig[i] end
  else
    local t = r:stack().vec(#r)
    alg.mulmv(t, sig, r)
    for i=1,n do sum = sum + t[i]*r[i] end
    r:stack().clear() 
  end
  local v = evalfgrad(fgrad, theta) - 0.5*sum
  return v
end

-- Newly created vectors th1, r1 vectors are initialized:
local function leapfrog(fgrad, eps, th0, r0, sig, mass)
  local n = #r0
  local th1, r1 = vec(n), vec(n)
  -- Leapfrog step: eps*M^-1*esp == eps*sig.
  if mass == "diagonal" then
    local _, grad = evalfgrad(fgrad, th0)
    for i=1,n do r1[i]  = r0[i]  + 0.5*eps*grad[i] end
    for i=1,n do th1[i] = th0[i] + eps*sig[i]*r1[i] end
    local _, grad = evalfgrad(fgrad, th1)
    for i=1,n do r1[i]  = r1[i]  + 0.5*eps*grad[i] end
  else
    local _, grad = evalfgrad(fgrad, th0)
    for i=1,n do r1[i]  = r0[i]  + 0.5*eps*grad[i] end
    local t = th1:stack().vec(#th1)
    alg.mulmv(t, sig, r1)
    for i=1,n do th1[i] = th0[i] + eps*t[i] end
    th1:stack().clear()  
    local _, grad = evalfgrad(fgrad, th1)
    for i=1,n do r1[i]  = r1[i]  + 0.5*eps*grad[i] end 
  end
  return th1, r1
end

local function heuristiceps(rng, fgrad, th0, sig, mass)
  local dim = #th0
  local r0 = vec(dim)
  local eps = 1
  for i=1,#r0 do r0[i] = normald:sample(rng) end
  local logpdf0 = logpdf(fgrad, th0, r0, sig, mass)
  local tht, rt = leapfrog(fgrad, eps, th0, r0, sig, mass)
  local alpha = sign((logpdf(fgrad, tht, rt, sig, mass) - logpdf0) - log(0.5))
  while alpha*(logpdf(fgrad, tht, rt, sig, mass) - logpdf0) > alpha*log(0.5) do
    eps = eps*2^alpha
    tht, rt = leapfrog(fgrad, eps, th0, r0, sig, mass)    
  end
  return eps
end

local function evalsr(sl, thp, thm, rp, rm)
  local n = #thp
  local sum1 = 0; for i=1,n do sum1 = sum1 + (thp[i] - thm[i])*rm[i] end; 
  local sum2 = 0; for i=1,n do sum2 = sum2 + (thp[i] - thm[i])*rp[i] end;
  return sl*step(sum1)*step(sum2)
end

-- Shorts: m = minus, p = plus, 1 = 1 prime, 2 = 2 primes.
-- This function does not modify any of its arguments.
-- Also all of the returned arguments are not modified in the recursion.
-- As long as the input or returned vectors are *not* modified it's fine to
-- work with references that may alias each other (see also caching of fgrad).
-- This means that all vectors here are effectively immutable after 
-- "initialization", which in this context means after being passed to the 
-- leapfrog function which modifies them.
local function buildtree(rng, fgrad, th, r, logu, v, j, eps, th0, r0, sig, 
    dlmax, mass)
  local dim = #th
  if j == 0 then
    local th1, r1 = leapfrog(fgrad, v*eps, th, r, sig, mass)
    local logpdfv1 = logpdf(fgrad, th1, r1, sig, mass)
    local n1 = step(logpdfv1 - logu)
    local s1 = step(dlmax + logpdfv1 - logu)
    return th1, r1, th1, r1, th1, n1, s1,
      min(exp(logpdfv1 - logpdf(fgrad, th0, r0, sig, mass)), 1), 1
  else
    local thm, rm, thp, rp, th1, n1, s1, a1, na1 =
    buildtree(rng, fgrad, th, r, logu, v, j - 1, eps, th0, r0, sig, dlmax, mass)
    local _, th2, n2, s2, a2, na2
    if s1 == 1 then
      if v == -1 then
        thm, rm, _, _, th2, n2, s2, a2, na2 =
        buildtree(rng, fgrad, thm, rm, logu, v, j - 1, eps, th0, r0, sig, dlmax,
          mass)
      else
        _, _, thp, rp, th2, n2, s2, a2, na2 =
        buildtree(rng, fgrad, thp, rp, logu, v, j - 1, eps, th0, r0, sig, dlmax,
          mass)
      end
      if rng:sample() < n2/(n1 + n2) then
        th1 = th2
      end
      a1 = a1 + a2
      na1 = na1 + na2
      s1 = evalsr(s2, thp, thm, rp, rm)
      n1 = n1 + n2
    end
    return thm, rm, thp, rp, th1, n1, s1, a1, na1
  end
end

-- Work by references on the vectors, only newly created one are modified in
-- the leapforg function. Use vectors as read-only.
local function nuts(rng, fgrad, theta0, o)
  local stop = o.stop
  local stopadapt = o.stopadapt or 1024
  if type(stop) == "number" then
    stop = stop_iter(stop)
  end
  if stopadapt ~= 0 and (log(stopadapt)/log(2) ~= floor(log(stopadapt)/log(2) or 
      stopadapt < 64)) then
    error("stopadapt must be 0 or a power of 2 which is >= 64")
  end
  local olstat = o.olstat
  local delta    = o.delta    or 0.8
  local gamma    = o.gamma    or 0.05
  local t0       = o.t0       or 10
  local k        = o.k        or 0.75
  local dlmax    = o.deltamax or 1000
  local mass     = o.mass     or "diagonal"

  local dim = #theta0
  local thr0, thr1 = theta0:copy()
  local r0 = vec(dim)
  -- m must be precision matrix (Cov^-1) or Var^-1 for diagonal mass case:
  local sigestimator, sig, mfac, m 
  if mass == "diagonal" then
    sigestimator, sig, mfac, m = stat.olvar(dim), vec(dim), vec(dim), vec(dim)
    if o.var then
      sig = o.var:copy()
    elseif o.cov then
      for i=1,dim do sig[i] = o.cov[i][i] end
    else
      for i=1,dim do sig[i] = 1/1e3 end
    end
  elseif mass == "dense" then
    sigestimator, sig, mfac, m = stat.olcov(dim), mat(dim, dim), mat(dim, dim),
      mat(dim, dim)
    if o.cov then
      sig = o.cov:copy()
    elseif o.var then
      for i=1,dim do sig[i][i] = o.var[i] end
    else
      for i=1,dim do for j=1,dim do sig[i][j] = 1/1e3 end end
    end
  else
    error("mass can only be 'diagonal' or 'dense', passed: "..mass)
  end
  local varu = 32
  invert(m, sig, mass)
  factor(mfac, m, mass)
  local eps = o.eps or heuristiceps(rng, fgrad, thr0, sig, mass)
  local mu = log(10*eps)
  local madapt, Ht, lepst, leps = 0, 0, 0
  local totadapt = 0
  while true do
    for i=1,dim do r0[i] = normald:sample(rng) end
    correlate(r0, mfac, mass)
    local logpdfr0 = logpdf(fgrad, thr0, r0, sig, mass)
    assert(logpdfr0 > -1/0)
    local logu = log(rng:sample()) + logpdfr0
    local thm, thp = thr0, thr0
    local rm, rp = r0, r0
    local j, n, s = 0, 1, 1
    local a, na
    thr1 = thr0
    while s == 1 do
      local v = sign(unifm1p1d:sample(rng))
      local _, th1, n1, s1
      if v == -1 then
        thm, rm, _, _, th1, n1, s1, a, na =
        buildtree(rng, fgrad, thm, rm, logu, v, j, eps, thr0, r0, sig, dlmax, 
          mass)
      else
        _, _, thp, rp, th1, n1, s1, a, na =
        buildtree(rng, fgrad, thp, rp, logu, v, j, eps, thr0, r0, sig, dlmax,
          mass)
      end
      if s1 == 1 then
        if rng:sample() < min(n1/n, 1) then
          thr1 = th1
        end
      end
      n = n + n1
      s = evalsr(s1, thp, thm, rp, rm)
      j = j + 1
      -- Limit tree depth during initial adaptation phase:
      if totadapt <= stopadapt/2 and j >= 10 then break end
    end
    local alpha = a/na
    if totadapt < stopadapt then
      totadapt = totadapt + 1 -- Last possible totadapt is stopadapt.
      madapt = madapt + 1
      Ht = (1 - 1/(madapt + t0))*Ht + 1/(madapt + t0)*(delta - alpha)
      leps = mu - sqrt(madapt)/gamma*Ht -- Log-eps used in adaptation.
      lepst = (madapt^-k)*leps + (1 - madapt^-k)*lepst -- Optimal log-eps.
      if totadapt == stopadapt then
        eps = exp(lepst)
      else
        eps = exp(leps)
      end
      sigestimator:push(thr1)
      if totadapt == varu then
        -- Set mass using current var estimates.
        sigestimate(sig, sigestimator, mass)
        invert(m, sig, mass)
        factor(mfac, m, mass)
        -- Re-initialize var estimators.
        sigestimator:clear()
        varu = varu*2
        mu = log(10*eps)
        madapt, Ht, lepst = 1, 0, 0
      end
    else
      olstat:push(thr1)
      if stop(thr1) then break end
    end
    thr0 = thr1
    cache = { thr0 = cache[thr0] }
  end
  return thr1:copy(), eps, sig:copy()
end

return {
  mcmc = nuts
}