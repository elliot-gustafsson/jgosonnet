// std.format('thing-%5.3d', [10.3])
// std.format('thing-%#5.3o', [10.3])
// std.format('thing-%6.4x', [910.3])
// std.format('%#.3g', [1000000001])  // '1.00e+09'

// local asdf(x, y=2, z=3) = [x, y, z];
// asdf(1, y=9, y=2)

local asdf(x, x) = x;
asdf(1, 2)
