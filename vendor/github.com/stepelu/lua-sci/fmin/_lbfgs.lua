-- LuaJIT implementation Limited memory BFGS (L-BFGS), based on libLBFGS
-- Copyright (c) 1990, Jorge Nocedal
-- Copyright (c) 2007-2010 Naoaki Okazaki
-- Copyright (c) 2014 Stefano peluchetti
-- All rights reserved.

-- TODO:
-- + method proposed by More and Thuente
-- + stop criteria (delta and past parameters)

local alg = require "sci.alg"

local vec = alg.vec
local sqrt, max = math.sqrt, math.max

-- Set to functions below:
local linesearches = {
  morethuente = true,
  armijo      = true,
  wolfe       = true,
  strongwolfe = true,
}

local function lbfgs_param(opt)
  opt = opt or { }
  local default = {
    ----------------------------------------------------------------------------
    -- The number of corrections to approximate the inverse Hessian matrix.
    -- The L-BFGS routine stores the computation results of previous \ref m
    -- iterations to approximate the inverse Hessian matrix of the current
    -- iteration. This parameter controls the size of the limited memories
    -- (corrections). The default value is \c 6. Values less than \c 3 are
    -- not recommended. Large values will result in excessive computing time.
    m = 6,

    ----------------------------------------------------------------------------
    -- Epsilon for convergence test.
    -- This parameter determines the accuracy with which the solution is to
    -- be found. A minimization terminates when
    --     ||g|| < \ref epsilon * max(1, ||x||),
    -- where ||.|| denotes the Euclidean (L2) norm. The default value is
    -- \c 1e-5.
    epsilon = 1e-6, -- TODO: Pass stopping criteria.

    -- The maximum number of iterations.
    -- Setting this parameter to zero continues an
    -- optimization process until a convergence or error. The default value
    -- is \c 0.
    max_iterations = 0, -- TODO: Pass stopping criteria.

    ----------------------------------------------------------------------------
    -- The line search algorithm.
    -- This parameter specifies a line search algorithm to be used by the
    -- L-BFGS routine.
    linesearch = "strongwolfe", -- TODO: Change default!

    -- The maximum number of trials for the line search.
    -- This parameter controls the number of function and gradients evaluations
    -- per iteration for the line search routine. The default value is \c 40.
    max_linesearch = 40,

    -- The minimum step of the line search routine.
    -- The default value is \c 1e-20. This value need not be modified unless
    -- the exponents are too large for the machine being used, or unless the
    -- problem is extremely badly scaled (in which case the exponents should
    -- be increased).
    min_step = 1e-20,

    -- The maximum step of the line search.
    -- The default value is \c 1e+20. This value need not be modified unless
    -- the exponents are too large for the machine being used, or unless the
    -- problem is extremely badly scaled (in which case the exponents should
    -- be increased).
    max_step = 1e20,

    -- A parameter to control the accuracy of the line search routine.
    -- The default value is \c 1e-4. This parameter should be greater
    -- than zero and smaller than \c 0.5.
    ftol = 1e-4, -- TODO: Add test!

    -- A coefficient for the Wolfe condition.
    -- This parameter is valid only when the backtracking line-search
    -- algorithm is used with the Wolfe condition,
    -- ::LBFGS_LINESEARCH_BACKTRACKING_STRONG_WOLFE or
    -- ::LBFGS_LINESEARCH_BACKTRACKING_WOLFE .
    -- The default value is \c 0.9. This parameter should be greater
    -- the \ref ftol parameter and smaller than \c 1.0.
    wolfe = 0.9, -- TODO: Add test!

    -- A parameter to control the accuracy of the line search routine.
    -- The default value is \c 0.9. If the function and gradient
    -- evaluations are inexpensive with respect to the cost of the
    -- iteration (which is sometimes the case when solving very large
    -- problems) it may be advantageous to set this parameter to a small
    -- value. A typical small value is \c 0.1. This parameter shuold be
    -- greater than the \ref ftol parameter (\c 1e-4) and smaller than
    -- \c 1.0.
    gtol = 0.9, -- TODO: Add test!

    -- The machine precision for floating-point values.
    -- This parameter must be a positive value set by a client program to
    -- estimate the machine precision. The line search routine will terminate
    -- with the status code (::LBFGSERR_ROUNDING_ERROR) if the relative width
    -- of the interval of uncertainty is less than this parameter.
    xtol = 1e-16,
  }
  local o = { }
  for k,v in pairs(default) do
    o[k] = opt[k] or v
  end
  assert(o.m        >= 1,      "opt.m must be strictly positive integer")
  assert(o.epsilon  >= 0,      "opt.epsilon must be positive")
  -- assert(o.past     >= 0,      "opt.past must be positive integer")
  -- assert(o.delta    >= 0,      "opt.delta must be positive")
  assert(o.min_step >= 0,      "opt.min_step must be positive")
  assert(o.max_step >= 0,      "opt.max_step must be positive")
  assert(o.min_step <= o.max_step, "opt.min_step <= opt.max_step required")
  assert(o.ftol     >= 0,      "opt.ftol must be positive")
  if o.linesearch == "wolfe" or o.linesearch == "strongwolfe" then
    assert(o.wolfe  >  o.ftol, "opt.wolfe > opt.ftol required")
    assert(o.wolfe  <  1,      "opt.wolfe < 1 required")
  end
  assert(o.gtol     >= 0,      "opt.gtol must be positive")
  assert(o.xtol     >= 0,      "opt.xtol must be positive")
  assert(o.max_linesearch > 0, "opt.max_linesearch must be strictly positive")
  assert(linesearches[o.linesearch], "invalid line search string")
  return o
end

-- Algebra ---------------------------------------------------------------------

local function vecadd(y, x, c)
  for i=1,#x do y[i] = y[i] + c*x[i] end
end

local function vecdiff(z, x, y)
  for i=1,#x do z[i] = x[i] - y[i] end
end

local function vecmul(y, x)
  for i=1,#x do y[i] = y[i]*x[i] end
end

---------------------------------------
local function vecset(x, c)
  for i=1,#x do x[i] = c end
end

local function veccpy(y, x)
  for i=1,#x do y[i] = x[i] end
end

local function vecdot(x, y)
  local s = 0
  for i=1,#x do s = s + x[i]*y[i] end
  return s
end

local function vecscale(x, scale)
  for i=1,#x do x[i] = x[i]*scale end
end

-- Line searches ---------------------------------------------------------------

-- [Backtracking method with the Armijo condition]
-- The backtracking method finds the step length such that it satisfies
-- the sufficient decrease (Armijo) condition:
--     - f(x + a * d) <= f(x) + opt.ftol * a * g(x)^T d
-- where x is the current point, d is the current search direction, and
-- a is the step length.
--
-- [Backtracking method with regular Wolfe condition]
-- The backtracking method finds the step length such that it satisfies
-- both the Armijo condition and the curvature condition:
--     - g(x + a * d)^T d >= opt.wolfe * g(x)^T d
-- where x is the current point, d is the current search direction, and
-- a is the step length.
--
-- [Backtracking method with strong Wolfe condition]
-- The backtracking method finds the step length such that it satisfies
-- both the Armijo condition and the following condition:
--     - |g(x + a * d)^T d| <= opt.wolfe * |g(x)^T d|
-- where x is the current point, d is the current search direction, and
-- a is the step length.
--
-- All these cases are covered in:
local function backtracking(x, finit, g, s, step, xp, gp, wp, gradf, param, 
    scale)
  local dec, inc = 0.5, 2.1

  if step <= 0 then
    return nil, "step size must be positive"
  end
  -- Compute the initial gradient in the search direction:
  local dginit = vecdot(g, s)
  -- Make sure that s points to a descent direction:
  if 0 < dginit then
      return nil, "s is not pointing to a descent direction"
  end

  local dgtest = param.ftol*dginit
  for count=1,1/0 do
    veccpy(x, xp)
    vecadd(x, s, step)
    -- Evaluate the function and gradient values:
    local f = scale*gradf(x, g)
    vecscale(g, scale)

    local width
    if f > finit + step*dgtest then
      width = dec
    else 
      if param.linesearch == "armijo" then -- OK Armijo condition.
        return count, f
      end
      local dg = vecdot(g, s)
      if dg < param.wolfe*dginit then
        width = inc
      else
        if param.linesearch == "wolfe" then -- OK Wolfe condition.
          return count, f
        end
        if dg > -param.wolfe*dginit  then
          width = dec
        else  -- OK strong Wolfe condition.
          return count, f
        end
      end
    end

    if step < param.min_step then
      return nil, "step smaller than opt.min_step"
    end
    if step > param.max_step then
      return nil, "step larger than opt.max_step"
    end
    if param.max_linesearch <= count then
      return nil, "opt.max_linesearch line search iterations exceeded"
    end

    step = step*width
  end
end

linesearches.armijo      = backtracking
linesearches.wolfe       = backtracking
linesearches.strongwolfe = backtracking

-- Method proposed by More and Thuente:
-- linesearches.morethuente (TODO: make default)
-- TODO: Implement.

local function stop(x, g, epsilon)
  -- Compute x and g norms.
  local xnorm = sqrt(vecdot(x, x))
  local gnorm = sqrt(vecdot(g, g))
  -- Converged if:
  --     |g(x)| / \max(1, |x|) < \epsilon
  if gnorm/max(1, xnorm) <= epsilon then
    return true
  end
end

-- LBFGS -----------------------------------------------------------------------
local function lbfgs(scale, gradf, param)
  local x = assert(param.x0, "x0 is required")
  param = lbfgs_param(param)
  local epsilon = param.epsilon
  local n = #x
  assert(n > 0, "problem size must be positive")
  local m = param.m
  local linesearch = linesearches[param.linesearch]

  local xp, g, gp, d, w = vec(n), vec(n), vec(n), vec(n), vec(n)
  local lm = { }
  for i=0,m-1 do
    lm[i] = { alpha = 0, ys = 0, y = vec(n), s = vec(n) }
  end

  local fx = scale*gradf(x, g) -- Initial function value and gradient.
  vecscale(g, scale)
  -- Initial direction, assume initial Hessian as identity matrix:
  for i=1,n do d[i] = -g[i] end

  if stop(x, g, epsilon) then
    return x, fx/scale
  end

  local step = 1/sqrt(vecdot(d, d)) -- Initial step.

  local iend = 0
  for k=1,1/0 do
    -- Store the current position and gradient vectors:
    --print(x:width())
    veccpy(xp, x)
    veccpy(gp, g)

    -- Search for an optimal step:
    local ls, fx_or_err = linesearch(x, fx, g, d, step, xp, gp, w, gradf, param, 
      scale)
    if not ls then
      return nil, fx_or_err
    end

    if stop(x, g, epsilon) then
      return x, fx_or_err/scale
    end

    -- TODO: Use stop().
    if param.max_iterations ~= 0 and param.max_iterations < k+1 then
      return nil, "opt.max_iterations LBFGS iterations exceeded"
    end

    --[[Update vectors s and y:
            s_{k+1} = x_{k+1} - x_{k} = \step * d_{k}.
            y_{k+1} = g_{k+1} - g_{k}.]]--
    local it = lm[iend]
    vecdiff(it.s, x, xp)
    vecdiff(it.y, g, gp)

    --[[Compute scalars ys and yy:
            ys = y^t \cdot s = 1 / \rho.
            yy = y^t \cdot y.
        Notice that yy is used for scaling the Hessian matrix H_0 (Cholesky
        factor).]]
    local ys = vecdot(it.y, it.s)
    local yy = vecdot(it.y, it.y)
    it.ys = ys

    --[[Recursive formula to compute dir = -(H \cdot g).
        This is described in page 779 of:
        Jorge Nocedal.
        Updating Quasi-Newton Matrices with Limited Storage.
        Mathematics of Computation, Vol. 35, No. 151,
        pp. 773--782, 1980.]]
    local bound = (m <= k) and m or k
    iend = (iend + 1) % m

    -- Compute the steepest direction:
    for i=1,n do d[i] = -g[i] end

    local j = iend
    for i=0,bound-1 do
      j = (j + m - 1) % m
      it = lm[j]
      -- \alpha_{j} = \rho_{j} s^{t}_{j} \cdot q_{k+1}.
      it.alpha = vecdot(it.s, d)
      it.alpha = it.alpha/it.ys
      -- q_{i} = q_{i+1} - \alpha_{i} y_{i}.
      vecadd(d, it.y, -it.alpha)
    end

    for i=1,n do d[i] = d[i]*(ys/yy) end

    for i=0,bound-1 do
        it = lm[j]
        -- \beta_{j} = \rho_{j} y^t_{j} \cdot \gamma_{i}.
        local beta = vecdot(it.y, d, n);
        beta = beta/it.ys
        -- \gamma_{i+1} = \gamma_{i} + (\alpha_{j} - \beta_{j}) s_{j}.
        vecadd(d, it.s, it.alpha - beta)
        j = (j + 1) % m
    end

    -- Now the search direction d is ready. We try step = 1 first:
    step = 1
  end
end

return {
  optim = lbfgs
}




