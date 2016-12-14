
tab = {}
--tab ['']=function(ast)
--	print("ast",ast)
--end

local temp = template[[
Here is Here!!! {{.a}}
]]
local td = tdata()
td.a = "Homogenous"

P.select(tab,call.source)


local B = buffer()
temp:render(B,td)

--print(B)
call.result = B:str()

