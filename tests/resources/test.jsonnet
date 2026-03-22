// local x = {
//   a: [],
//   b: {},
//   c: 'a',
// };

// {
//   json: std.manifestJson(x),
//   mini: std.manifestJsonMinified(x),
//   str: std.toString(x),
//   xx: x,
// }

// std.manifestJson(x)

// local asdf = function(selectors, extra_selectors=[], multiplier=3)
//   local x = selectors + extra_selectors;
//   std.length(x) * multiplier;


// asdf([1, 3, 4], multiplier=6)


// std.foldl(function(s, x) s + 7 + x, '123', '99')

// 'asdf %(asdf)s asdf' % { asdf: std.join('\n', ['a', 'b', 'c']) }

// 'asd %(asdf)s asd' % { asdf: 'banan' }
// local x = 3;
// local y = 4;

// std.startsWith(b='bb', a='bbaabb')

// std.splitLimit(str, c, maxsplits)
// std.splitLimit('1 2 3 4 5', maxsplits=2, c=' ')

// std.sort(arr=['b', 'a'])

// std.manifestYamlDoc({
//   from: '-Infinity',
// }, quote_keys=false)


// {
//   //   from: '-Infinity',
//   //   to: 'Infinity',
//   //   to: '+Infinity',
//   //   x: 'NaN',
//   a: '0x123',
// }

// local x = [
//   { name: 'cattle-fleet-system', system: true },
//   { name: 'cattle-impersonation-system', system: true },
//   { name: 'cattle-system', system: true },
//   { name: 'cattle-ui-plugin-system', system: true },
//   { name: 'falco', system: true },
//   { name: 'kube-node-lease', system: true },
//   { name: 'kube-public', system: true },
//   { name: 'kube-system', system: true },
//   { name: 'local', system: true },
//   { name: 'default', system: false },
//   { name: 'dependencytrack', system: false },
//   { name: 'falco', system: false },
//   { name: 'sonarqube', system: false },
// ];

// std.set(x, keyF=function(x) x.name)
// std.sort(x, keyF=function(x) x.name)

// [1, 2, 3, 4] <= [1, 2, 3]

// [1, 2, 3] <= [2, 3]

// std.assertEqual({ a: 'a' }, { a: 'b' })

// std.parseJson('{"a": "a"}')
// std.parseYaml('1234.23')

std.pi

// {
//   a: 'a',
//   b: 'b',
// }
// ==
// {
//   c:: 'a',
//   b: 'ba',
//   a: 'a',
// }
