    -- configuration file for program `pp'
width = 300
height = 300

background = {r=30, g=20, b=10}

function f (x, y)
	return x+y
end

function print_foreground() 
	print(foreground.red)
	print(foreground.green)
	print(foreground.blue)
end

function print_background() 
	print(background.r)
	print(background.g)
	print(background.b)
end

function print_summator()
	print(summator(5,5))
end
