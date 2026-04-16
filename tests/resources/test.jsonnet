// std.parseYaml('1\n---')

// std.objectRemoveKey({ foo: 1, bar: 2, baz: 3 }, 'foo')
// std.objectRemoveKey({ foo: 1, bar: 2, baz:: 3 }, 'foo').baz
// { x: 1 } + std.objectRemoveKey({ a: 1 } + { b: super.a }, 'a')

// { a: 1 } + std.objectRemoveKey({ b: super.a }, 'a')

std.assertEqual({ a: 1 } + std.objectRemoveKey({ b: super.a }, 'a'), { a: 1, b: 1 })
