-- The only point of this file is to generate an interesting stack trace

function call3(n)
	if n == 3 then
		error("Some error")
	end
end

function call2()
	for i = 1, 4, 1 do
		call3(i)
	end
end

function call1()
	call2()
end

call1()
