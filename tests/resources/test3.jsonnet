/* Semantic Evaluation Test
   Targeting: Thunk caching, Late binding, Closures, and Recursion.
*/

local StressTest = {

  // ---------------------------------------------------------
  // 1. LAZY EVALUATION & ERROR HANDLING
  // ---------------------------------------------------------
  // If your evaluator is eager, this will crash.
  // It should only crash if 'danger_zone' is actually accessed in the final JSON.
  lazy_check: {
    safe: 'I am safe',
    danger_zone:: error 'You evaluated me! I should have been lazy!',

    // Test: We access 'safe', so 'danger_zone' should remain unevaluated.
    result: self.safe,
  },

  // ---------------------------------------------------------
  // 2. THUNK CACHING (Memoization)
  // ---------------------------------------------------------
  // If you don't cache thunks, this calculation happens twice.
  // While hard to verify programmatically, logging inside your evaluator
  // should show "expensive_calc" computed only once.
  caching_check: {
    local heavy = std.makeArray(1000, function(x) x),
    len: std.length(heavy) + std.length(heavy),
  },

  // ---------------------------------------------------------
  // 3. OBJECT INHERITANCE & LATE BINDING (The Hardest Part)
  // ---------------------------------------------------------
  // When 'Base' is merged into 'Derived', 'self.name' in Base
  // must resolve to "Derived", not "Base".
  inheritance: {
    local Base = {
      name: 'Base',
      // 'self' here is dynamic. It depends on who calls it.
      who_am_i: 'I am ' + self.name,
    },

    local Derived = Base {
      name: 'Derived',
    },

    // Expectation: "I am Derived"
    result: Derived.who_am_i,
  },

  // ---------------------------------------------------------
  // 4. SUPER CHAINING
  // ---------------------------------------------------------
  // Testing access to the parent layer in a deep merge.
  super_chain: {
    local A = { val: 1 },
    local B = A { val: super.val + 1 },  // 1 + 1 = 2
    local C = B { val: super.val + 2 },  // 2 + 2 = 4

    final: C.val,
  },

  // ---------------------------------------------------------
  // 5. CLOSURES & SCOPE CAPTURE
  // ---------------------------------------------------------
  // The function 'makeAdder' returns a function that must remember 'x'.
  closures: {
    local makeAdder = function(x) function(y) x + y,
    local addFive = makeAdder(5),
    result: addFive(10),  // Should be 15
  },

  // ---------------------------------------------------------
  // 6. TAIL RECURSION / STACK SAFETY
  // ---------------------------------------------------------
  // Simple evaluators often blow the stack here.
  // Standard Jsonnet handles deep recursion reasonably well.
  recursion: {
    local sum(n, acc=0) =
      if n == 0 then acc else sum(n - 1, acc + n),

    // Sum of 0..100 = 5050
    result: sum(100),
  },
};

StressTest
