--------------------------------------------------------------------------------
-- Differential evolution algorithm module.
--
-- Copyright (C) 2011-2016 Stefano Peluchetti. All rights reserved.
--------------------------------------------------------------------------------

-- Here implemented is the differential evolution algorithm presented in the 
-- paper: "Self-adaptive Differential Evolution Algorithm for Constrained 
-- Real-Parameter Optimization", 2006.
-- The following modifications have been performed: 
-- + in the paper a moving window of LP generations is used to update the 
--   parameters; we simply re-update the parameters every LP generations using 
--   the last LP generations (no overlapping), provided that at least 100 
--   successful mutations have been obtained (i.e. 100 samples on which to base
--   the estimates)
-- + only one CR vector (common for all strategies), instead of 3 as in the 
--   paper
-- + different adaptation algorithm for CR, see sample_CR and initialization of 
--   CRmu, CRsigma

local xsys = require "xsys"
local alg  = require "sci.alg"
local prng = require "sci.prng"
local dist = require "sci.dist"
local math = require "sci.math"
local stat = require "sci.stat"

local min, max, abs, floor, ceil, step, sqrt = xsys.from(math,
     "min, max, abs, floor, ceil, step, sqrt")
     
local normald = dist.normal
local alg32 = alg.typeof("int32_t")

local function rand_1_bin(v, j, rj, xmin, x, F, K)
  local j1, j2, j3 = rj[{j,1}], rj[{j,2}], rj[{j,3}]
  local x1, x2, x3, F0 = x[j1], x[j2], x[j3], F[j]
  for i=1,#v do 
    v[i] = x1[i] + (x2[i] - x3[i])*F0 
  end
end

local function rand_2_bin(v, j, rj, xmin, x, F, K)
  local j1, j2, j3, j4, j5 = rj[{j,1}], rj[{j,2}], rj[{j,3}], rj[{j,4}], rj[{j,5}]
  local x1, x2, x3, x4, x5, F0 = x[j1], x[j2], x[j3], x[j4], x[j5], F[j]
  for i=1,#v do 
    v[i] = x1[i] + (x2[i] - x3[i])*F0 + (x4[i] - x5[i])*F0 
  end
end

local function randtobest_2_bin(v, j, rj, xmin, x, F, K)
  local j1, j2, j3, j4 = rj[{j,1}], rj[{j,2}], rj[{j,3}], rj[{j,4}]
  local x0, x1, x2, x3, x4, F0 = x[j], x[j1], x[j2], x[j3], x[j4], F[j]
  for i=1,#v do 
    v[i] = x0[i] + (xmin[i] - x0[i])*F0 + (x1[i] - x2[i])*F0 
                                        + (x3[i] - x4[i])*F0 
  end
end

local function currenttorand_1(v, j, rj, xmin, x, F, K)
  local j1, j2, j3 = rj[{j,1}], rj[{j,2}], rj[{j,3}]
  local x0, x1, x2, x3, F0, K0 = x[j], x[j1], x[j2], x[j3], F[j], K[j]
  for i=1,#v do 
    v[i] = x0[i] + (x1[i] - x0[i])*K0 + (x2[i] - x3[i])*F0 
  end
end

local strategies = { 
  rand_1_bin, 
  rand_2_bin,
  randtobest_2_bin,
  currenttorand_1,
}

-- Branch-free sampling of strategy according to the four probabilities in p.
local function sample_strategy_indices(rng, rs, p)
  for j=1,#rs do
    local u = rng:sample()
    rs[j] = 1 
          + step(u - (p[1]))               
          + step(u - (p[1] + p[2]))
          + step(u - (p[1] + p[2] + p[3]))
  end
end

-- Sample an integer uniformly distributed on the interval from, ... to.
local function sample_int(rng, from, to)
  return floor(from + (to + 1 - from)*rng:sample())
end

-- Branch-free version, each row of rj contains three distinct j uniformly 
-- distributed on 1, ..., NP.
local function sample_distinct_indices(rng, rj)
  local NP = rj:nrow()
  local j = 1
  while j <= NP do
    local j1, j2, j3, j4, j5 = sample_int(rng, 1, NP), sample_int(rng, 1, NP), 
       sample_int(rng, 1, NP), sample_int(rng, 1, NP), sample_int(rng, 1, NP)
    rj[{j,1}], rj[{j,2}], rj[{j,3}], rj[{j,4}], rj[{j,5}] = j1, j2, j3, j4, j5
    -- Zero if any pairwise match, integer otherwise.
    local m = (j1 - j)
             *(j2 - j1)*(j2 - j)
             *(j3 - j2)*(j3 - j1)*(j3 - j)
             *(j4 - j3)*(j4 - j2)*(j4 - j1)*(j4 - j)
             *(j5 - j4)*(j5 - j3)*(j5 - j2)*(j5 - j1)*(j5 - j) 
    j = j + min(abs(m), 1)
  end
end

local function sample_F(rng, F)
  for i=1,#F do 
    F[i] = rng:sample()
  end
end

local function sample_K(rng, K)
  for i=1,#K do 
    K[i] = rng:sample() 
  end
end

-- MODIFICATION: Our algorithm performs flooring at 0 and 1.
-- CRm has column 1 set at 100: always CR = 1.
local function sample_CR(rng, CR, CRmu, CRsigma)
  for j=1,#CR do
    local v = normald(CRmu, CRsigma):sample(rng)
    CR[j] = max(0, min(v, 1))
  end
end

-- Element equal to one means mutation.
local function sample_mutations(rng, rz, rs, CR)
  local NP, dim = #rz, #rz[1]
  for j=1,NP do
    if rs[j] == 4 then
      for d=1,dim do
        rz[j][d] = 1
      end
    else
      for d=1,dim do
        rz[j][d] = step(CR[j] - rng:sample())
      end
      -- Always move at least among one dimension:
      rz[j][sample_int(rng, 1, dim)] = 1
    end
  end
end

local function nan_to_inf(x)
  return x == x and x or 1/0
end

-- Notice that to avoid stagnation we favor f1 and v1 over f2 and v2.
local function compare(scale, f1, v1, f2, v2)
  local vmin, vmax = min(v1, v2), max(v1, v2)
  if vmax == 0 then -- Both satisfy the constraints.    
    return nan_to_inf(scale*f1) <= nan_to_inf(scale*f2)
  elseif vmin > 0 then -- Neither satisfy the constraints.
    return v1 <= v2
  else
    return v1 == 0 -- Who satisfy the constraint wins, and one of them do.
  end
end

local function updatemin(scale, xmin, fmin, vmin, xnew, fnew, vnew)
  if compare(scale, fnew, vnew, fmin, vmin) then
    -- Note copy is made only when new particle is a new global minimum. 
    -- But the copy NEEDS to be made.
    return xnew:copy(), fnew, vnew
  else
    return xmin, fmin, vmin
  end
end

local function range(x, col)
  local v = x[{1,col}]
  local l, u = v, v
  for r=2,x:nrow() do
    v = x[{r,col}]
    l = min(l, v)
    u = max(u, v)
  end
  return u - l
end

-- TODO: Maybe a joint absolute and relative stopping criteria?
local function stop_x_range_no_violation(xrange)
  xrange = xrange or 1e-4
  if not (xrange > 0) then
    error('strictly positive maximum x-range is required, is: '..xrange)
  end
  return function(_, _, vmin, xval, _, _)
    if vmin > 0 then
      return false
    end
    local v = 0
    for c=1,xval:ncol() do
      v = max(v, range(xval, c))
    end
    return v < xrange
  end
end

local function linear_lt(l, u)
  return max(0, l - u)
end

local function hyper_cube_constrain(xl, xu)
  return function(x, lt)
    local v = 0
    for i=1,#x do
      v = v + lt(xl[i], x[i]) + lt(x[i], xu[i])
    end
    return v
  end
end

local function rows_as_vec(x)
  local o =  {}
  for i=1,x:nrow() do
    o[i] = alg.vec(x:ncol())
    for j=1,x:ncol() do
      o[i][j] = x[{i,j}]
    end
  end
  return o
end

local function optim(scale, f, o)
  local rng = o.rng or prng.std()

  local stop = o.stop
  if type(stop) == 'nil' or type(stop) == "number" then
    stop = stop_x_range_no_violation(stop)
  end
  
  local constraint = o.constraint
  local xl, xu = o.xl and o.xl:copy(), o.xu and o.xu:copy()
  if xl then
    if #xl ~= #xu then
      error("xl and xu must have the same size")
    end
    for i=1,#xl do
      if not (xl[i] < xu[i]) then
        error("xl < xu required")
      end
    end
    constraint = constraint or hyper_cube_constrain(xl, xu)
  end
  if not constraint then
    error('constraint must be provided if xl and xu are not provided')
  end

  local NP, dim, xval, x -- Current population.
  if o.x0 then
    xval = o.x0:copy()
    NP, dim = xval:nrow(), xval:ncol()
    x = rows_as_vec(xval)
  else
    dim = #xl
    NP = o.np or max(10, 8*dim)
    xval = alg.mat(NP, dim)
    x = rows_as_vec(xval)
    local popd = dist.mvuniform(xl, xu)
    for j=1,NP do
      popd:sample(rng, x[j])
      for i=1,dim do
        xval[{j,i}] = x[j][i] 
      end
    end
  end
  if NP < 10 then
    error("NP >= 10 required")
  end

  local fval, vval = alg.vec(NP), alg.vec(NP)
  local xmin, fmin, vmin = nil, 1/0, 1/0
  for j=1,NP do
    vval[j] = constraint(x[j], linear_lt)
    fval[j] = vval[j] == 0 and f(x[j]) or 0/0
    xmin, fmin, vmin = updatemin(scale, xmin, fmin, vmin, x[j], fval[j], 
      vval[j])
  end
     
  -- Equal probability to each strategy:
  local Pmu = alg.vec(4, 1/4)
  -- Centered around 0.5 with good dispersion: 68.2% of mass in (0, 1).
  local CRmu, CRsigma = 0.5, 0.5
  
  local Pstat = { } 
    for i=1,4 do Pstat[i] = stat.olmean(0) 
  end
  local CRstat = stat.olvar(0)
  local nsuccess = 0

  local F  = alg.vec(NP)
  local K  = alg.vec(NP)
  local CR = alg.vec(NP)
  local rs = alg32.vec(NP)    -- Random strategies.
  local rj = alg32.mat(NP, 5) -- Random j-indices.
  local rz = { } -- Dimensions which mutate.
  local v  = alg.vec(dim) -- Potential mutation particle.
  local u  = { } -- Mutated particles.
  for j=1,NP do
    rz[j], u[j] = alg.vec(dim), alg.vec(dim)
  end
      
  local generation = 0
  while not stop(xmin, fmin, vmin, xval, fval, vval) do
    generation = generation + 1
    -- Update meta-parameters.
    if generation % 20 == 0 and nsuccess >= 100 then      
      for i=1,4 do
        Pmu[i] = Pstat[i]:mean() + 0.01
        Pstat[i]:clear()
      end
      local sum = 0 
      for i=1,4 do sum = sum + Pmu[i] end
      for i=1,4 do Pmu[i] = Pmu[i]/sum end
      CRmu = max(0, min(CRmu + 2*(CRstat:mean() - CRmu), 1))
      CRsigma = max(0.1, sqrt(CRstat:var()))
      CRstat:clear()
      nsuccess = 0
    end
    -- Sample required quantities:
    sample_strategy_indices(rng, rs, Pmu)
    sample_distinct_indices(rng, rj)
    sample_F(rng, F)
    sample_K(rng, K)
    sample_CR(rng, CR, CRmu, CRsigma)
    sample_mutations(rng, rz, rs, CR)
    -- Evolve the population:
    for j=1,NP do
      -- Potential mutation:
      strategies[rs[j]](v, j, rj, xmin, x, F, K)
      -- Mutated particle:
      for i=1,dim do 
        u[j][i] = rz[j][i]*v[i] + (1 - rz[j][i])*x[j][i] 
      end
      -- No hyper-cube style bounds are applied.     
    end    
    -- Selection:
    for j=1,NP do
      local s = rs[j]
      local vuj = constraint(u[j], linear_lt)
      local fuj = vuj == 0 and f(u[j]) or 0/0
      if compare(scale, fuj, vuj, fval[j], vval[j]) then
        -- It's an improvement --> select.
        nsuccess = nsuccess + 1
        Pstat[s]:push(1)
        CRstat:push(CR[j])
        for i=1,dim do 
          x[j][i] = u[j][i] 
        end
        fval[j], vval[j] = fuj, vuj        
        xmin, fmin, vmin = updatemin(scale, xmin, fmin, vmin, u[j], fuj, vuj)
      else
        Pstat[s]:push(0)
      end
      -- Update xval for next generation (used only in stop()):
      for i=1,dim do 
        xval[{j,i}] = x[j][i] 
      end
    end
  end
  
  return xmin:copy(), fmin, vmin, xval, fval:copy(), vval:copy()
end

return { 
  optim = optim,
}