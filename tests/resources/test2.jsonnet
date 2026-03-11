/* Comprehensive Jsonnet Parser Test
  Covers: Comments, Locals, Mixins, Comprehensions,
  Slicing, Hidden Fields, and String Literals.
*/

// 1. Top-level locals and arithmetic
local multiplier = 2;
local glue = ' - ';

local test(x) =
  x + ' test';

// 2. A base object (mixin) to test inheritance
local BaseMixin = {
  // Method definition with default arguments
  calculate(x, y=10):: (x * y) * multiplier,

  // Hidden field (double colon) - should not appear in final JSON
  internal_id:: 'hidden_base_01',

  common_field: 'base value',
};

// 3. The main object structure
{
  // mixin application
  base: BaseMixin,

  imported_field: import 'test.jsonnet',

  // 4. Fields and Sugar Syntax
  'quoted-key': true,
  unquoted_key: null,
  'single-quoted': 1.5e-2,  // Scientific notation

  // 5. Array Comprehensions & Slicing
  // Generate [1, 4, 9, 16] then slice to get [4, 9]
  squares_sliced: [x * x for x in [1, 2, 3, 4]][1:3],

  // 6. Object Comprehension
  // Dynamic key generation
  map_kv: {
    ['k_' + x]: x
    for x in ['a', 'b', 'c']
    if x != 'b'  // Conditional filter
  },

  // 7. String Block Literals (text blocks)
  description: |||
    This is a text block.
    It preserves newlines.
  |||,

  // 8. Self, Super, and Object merging (+)
  derived: BaseMixin {
    // Override standard field
    common_field: 'overridden value',

    asdf3: super.internal_id,

    // Accessing 'super'
    old_calc: super.calculate(5),


    // Accessing 'self' (late binding)
    // If we merge this object later, 'self.common_field' changes
    self_ref: self.common_field + glue + 'suffixed',
  },

  // 9. Assertions (Parsers must handle this syntax node)
  assert multiplier > 0 : 'Multiplier must be positive',

  // 10. Conditionals and logic
  status: if std.length('abc') == 3 then 'ok' else 'error',


  // 11. Error throwing (Parser should accept this keyword)
  // Commented out to prevent runtime failure, but parser should recognize:
  // crash: error "This is a fatal error",

  // 12: local func -> 'asdf test'
  asdf: test('asdf'),


  asdf2: BaseMixin.internal_id,
}
