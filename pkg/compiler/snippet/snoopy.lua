
        Dog = __reg:RegisterInterface("Dog");
        Beagle = __reg:RegisterStruct("Beagle");

	    function Beagle:Write(with)
            b = self;
  		    return b.word .. ":it was a dark and stormy night, " .. with;
     	end;

        snoopy = __reg:NewInstance("Beagle",{["word"]="hiya"});

  	    _r = snoopy:Write("with a pen");
  	    book = _r;
