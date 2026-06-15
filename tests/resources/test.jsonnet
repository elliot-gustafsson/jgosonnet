local myObj = {
  // CONTEXT BOUNDARY: object field
  crashField: [
    'safe_element',

    // CONTEXT BOUNDARY: array element
    // PASS-THROUGHS: The 'if' condition and the '+' binary operator
    // should NOT add frames to your stack.
    if true then 100 + error 'Intentional Boom!' else 0,
  ],
};

local myFunc(x) =
  // CONTEXT BOUNDARY: thunk
  local myThunk = myObj.crashField[1];

  // PASS-THROUGH: The '+' binary operator
  x + myThunk;

// We use std.map to force an anonymous function and a builtin function boundary
local trigger = std.map(
  // CONTEXT BOUNDARY: anonymous function
  function(y)
    // CONTEXT BOUNDARY: named function
    myFunc(y),
  [42]
);

// CONTEXT BOUNDARY: $ (Root execution)
trigger[0]
