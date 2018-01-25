__utf8 = require 'utf8'
function __decodeRune(s, i)
   return __utf8.sub(s, i+1, i+1)
end
