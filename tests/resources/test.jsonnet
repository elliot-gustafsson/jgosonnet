// std.parseYaml('1\n---')

// std.objectRemoveKey({ foo: 1, bar: 2, baz: 3 }, 'foo')
// std.objectRemoveKey({ foo: 1, bar: 2, baz:: 3 }, 'foo').baz
// { x: 1 } + std.objectRemoveKey({ a: 1 } + { b: super.a }, 'a')

// { a: 1 } + std.objectRemoveKey({ b: super.a }, 'a')

// std.assertEqual({ a: 1 } + std.objectRemoveKey({ b: super.a }, 'a'), { a: 1, b: 1 })

// std.assertEqual(std.trace('', null), null) &&
// std.assertEqual(std.trace('', true), true) &&
// std.assertEqual(std.trace('', 77), 77) &&
// std.assertEqual(std.trace('', 77.88), 77.88) &&
// std.assertEqual(std.trace('', 'word'), 'word') &&
// std.assertEqual(std.foldl(std.trace('', function(acc, i) acc + i), [1, 2, 3], 0), 6) &&
// std.assertEqual(std.trace('', {}), {}) &&
// std.assertEqual(std.trace('', { a: {} }), { a: {} }) &&
// std.assertEqual(std.trace('', []), []) &&
// std.assertEqual(std.trace('', [{ a: 'b' }, { a: 'b' }]), [{ a: 'b' }, { a: 'b' }]) &&
// std.assertEqual(std.trace('Some Trace Message', { a: {} }), { a: {} }) &&

// true

// std.trace('asdf', 1)

// Test values at boundary of safe integer range
// local max_safe = 9007199254740991;  // 2^53 - 1
// local min_safe = -9007199254740991;  // -(2^53 - 1)

// std.assertEqual(max_safe & 1, 1) &&
// std.assertEqual(min_safe & 1, 1)


// ~(max_safe - 1)
importstr '/home/elliot.gustafsson@fnox.it/Projects/jgosonnet/tests/resources/jsonnet-cpp/test_suite/unicode_bmp1.jsonnet.in'
